package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	"github.com/spf13/cobra"
)

type option struct {
	port             int
	filterPath       []string
	saveResponseBody bool
	output           string
	upstreamProxy    string
	verbose          bool
	username         string
	password         string
}

// NewRootCmd creates the root command
func NewRootCmd() (c *cobra.Command) {
	opt := &option{}
	c = &cobra.Command{
		Use:   "atest-collector",
		Short: "A collector for API testing, it will start a HTTP proxy server",
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 8080, "The port for the proxy")
	flags.StringSliceVarP(&opt.filterPath, "filter-path", "", []string{}, "The path prefix for filtering")
	flags.BoolVarP(&opt.saveResponseBody, "save-response-body", "", false, "Save the response body")
	flags.StringVarP(&opt.output, "output", "o", "sample.yaml", "The output file")
	flags.StringVarP(&opt.upstreamProxy, "upstream-proxy", "", "", "The upstream proxy")
	flags.StringVarP(&opt.username, "username", "", "", "The username for basic auth")
	flags.StringVarP(&opt.password, "password", "", "", "The password for basic auth")
	flags.BoolVarP(&opt.verbose, "verbose", "", false, "Verbose mode")

	_ = cobra.MarkFlagRequired(flags, "filter-path")
	return
}

type responseFilter struct {
	urlFilter *filter.URLPathFilter
	collects  *pkg.Collects
	ctx       context.Context
}

func (f *responseFilter) filter(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return resp
	}

	req := resp.Request
	if f.urlFilter.Filter(req.URL) {
		simpleResp := &pkg.SimpleResponse{StatusCode: resp.StatusCode}

		if resp.Body != nil {
			buf := new(bytes.Buffer)
			io.Copy(buf, resp.Body)
			simpleResp.Body = buf.String()
			resp.Body = io.NopCloser(buf)
		}

		f.collects.Add(req.Clone(f.ctx), simpleResp)
	}
	return resp
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	urlFilter := &filter.URLPathFilter{PathPrefix: o.filterPath}
	collects := pkg.NewCollects()
	responseFilter := &responseFilter{urlFilter: urlFilter, collects: collects, ctx: cmd.Context()}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = o.verbose
	if o.upstreamProxy != "" {
		proxy.Tr.Proxy = func(r *http.Request) (*url.URL, error) {
			return url.Parse(o.upstreamProxy)
		}
		proxy.ConnectDial = proxy.NewConnectDialToProxy(o.upstreamProxy)
		cmd.Println("Using upstream proxy", o.upstreamProxy)
	}
	if o.username != "" && o.password != "" {
		auth.ProxyBasic(proxy, "my_realm", func(user, pwd string) bool {
			return user == o.username && o.password == pwd
		})
	}
	proxy.OnResponse().DoFunc(responseFilter.filter)

	exporter := pkg.NewSampleExporter(o.saveResponseBody)
	collects.AddEvent(exporter.Add)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", o.port),
		Handler: proxy,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		collects.Stop()
		_ = srv.Shutdown(context.Background())
	}()

	cmd.Println("Starting the proxy server with port", o.port)
	_ = srv.ListenAndServe()
	var data string
	if data, err = exporter.Export(); err == nil {
		err = os.WriteFile(o.output, []byte(data), 0644)
	}
	return
}
