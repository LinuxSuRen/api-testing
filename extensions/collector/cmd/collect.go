package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/elazarl/goproxy"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	"github.com/spf13/cobra"
)

type option struct {
	port       int
	filterPath string
	output     string
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
	flags.StringVarP(&opt.filterPath, "filter-path", "", "", "The path prefix for filtering")
	flags.StringVarP(&opt.output, "output", "o", "sample.yaml", "The output file")

	_ = cobra.MarkFlagRequired(flags, "filter-path")
	return
}

type responseFilter struct {
	urlFilter *filter.URLPathFilter
	collects  *pkg.Collects
}

func (f *responseFilter) filter(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return resp
	}

	req := resp.Request
	if f.urlFilter.Filter(req.URL) {
		f.collects.Add(req.Clone(context.TODO()))
	}
	return resp
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	urlFilter := &filter.URLPathFilter{PathPrefix: o.filterPath}
	collects := pkg.NewCollects()
	responseFilter := &responseFilter{urlFilter: urlFilter, collects: collects}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnResponse().DoFunc(responseFilter.filter)

	exporter := pkg.NewSampleExporter()
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
