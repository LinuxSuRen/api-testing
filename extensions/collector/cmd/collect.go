package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/elazarl/goproxy"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	atestpkg "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type option struct {
	port       int
	filterPath string
	output     string
}

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

	cobra.MarkFlagRequired(flags, "filter-path")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	urlFilter := &filter.URLPathFilter{PathPrefix: o.filterPath}
	collects := pkg.NewCollects()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if urlFilter.Filter(r.URL) {
				collects.Add(r.Clone(context.TODO()))
			}
			return r, nil
		})

	exporter := &sampleExporter{
		testSuite: atestpkg.TestSuite{
			Name: "sample",
		},
	}
	collects.AddEvent(exporter.add)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", o.port),
		Handler: proxy,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		collects.Stop()
		srv.Shutdown(context.Background())
	}()

	cmd.Println("Starting the proxy server with port", o.port)
	srv.ListenAndServe()
	var data string
	if data, err = exporter.export(); err == nil {
		err = os.WriteFile(o.output, []byte(data), 0644)
	}
	return
}

type sampleExporter struct {
	testSuite atestpkg.TestSuite
}

func (e *sampleExporter) add(r *http.Request) {
	body := r.Body
	data, _ := io.ReadAll(body)

	fmt.Println("receive", r.URL.Path)
	req := atestpkg.Request{
		API:    r.URL.String(),
		Method: r.Method,
		Header: map[string]string{},
		Body:   string(data),
	}

	testCase := atestpkg.TestCase{
		Request: req,
		Expect: atestpkg.Response{
			StatusCode: 200,
		},
	}

	specs := strings.Split(r.URL.Path, "/")
	if len(specs) > 0 {
		testCase.Name = specs[len(specs)-1]
	}

	if val := r.Header.Get("Content-Type"); val != "" {
		req.Header["Content-Type"] = val
	}

	e.testSuite.Items = append(e.testSuite.Items, testCase)
}

var prefix = `#!api-testing
# yaml-language-server: $schema=https://gitee.com/linuxsuren/api-testing/raw/master/sample/api-testing-schema.json
`

func (e *sampleExporter) export() (string, error) {
	marker := map[string]int{}

	for i, item := range e.testSuite.Items {
		if _, ok := marker[item.Name]; ok {
			marker[item.Name]++
			e.testSuite.Items[i].Name = fmt.Sprintf("%s-%d", item.Name, marker[item.Name])
		} else {
			marker[item.Name] = 0
		}
	}

	data, err := yaml.Marshal(e.testSuite)
	return prefix + string(data), err
}
