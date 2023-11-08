// Package cmd provides all the commands
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	_ "embed"
	pprof "net/http/pprof"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	template "github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func createServerCmd(execer fakeruntime.Execer, gRPCServer gRPCServer, httpServer server.HTTPServer) (c *cobra.Command) {
	opt := &serverOption{
		gRPCServer: gRPCServer,
		httpServer: httpServer,
		execer:     execer,
	}
	c = &cobra.Command{
		Use:     "server",
		Short:   "Run as a server mode",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7070, "The RPC server port")
	flags.IntVarP(&opt.httpPort, "http-port", "", 8080, "The HTTP server port")
	flags.BoolVarP(&opt.printProto, "print-proto", "", false, "Print the proto content and exit")
	flags.StringArrayVarP(&opt.localStorage, "local-storage", "", []string{"*.yaml"}, "The local storage path")
	flags.StringVarP(&opt.consolePath, "console-path", "", "", "The path of the console")
	flags.StringVarP(&opt.configDir, "config-dir", "", os.ExpandEnv("$HOME/.config/atest"), "The config directory")
	flags.StringVarP(&opt.secretServer, "secret-server", "", "", "The secret server URL")
	flags.StringVarP(&opt.skyWalking, "skywalking", "", "", "Push the browser tracing data to the Apache SkyWalking HTTP URL")
	return
}

type serverOption struct {
	gRPCServer gRPCServer
	httpServer server.HTTPServer
	execer     fakeruntime.Execer

	port         int
	httpPort     int
	printProto   bool
	localStorage []string
	consolePath  string
	secretServer string
	configDir    string
	skyWalking   string
}

func (o *serverOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	o.configDir = os.ExpandEnv(o.configDir)
	err = o.execer.MkdirAll(o.configDir, 0755)
	return
}

func (o *serverOption) runE(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	if o.printProto {
		for _, val := range server.GetProtos() {
			cmd.Println(val)
		}
		return
	}

	var (
		lis     net.Listener
		httplis net.Listener
	)
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}
	httplis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.httpPort))
	if err != nil {
		return
	}

	loader := testing.NewFileWriter("")
	for _, storage := range o.localStorage {
		if loadErr := loader.Put(storage); loadErr != nil {
			cmd.PrintErrf("failed to load %s, error: %v\n", storage, loadErr)
			continue
		}
	}

	var secretServer remote.SecretServiceServer
	if o.secretServer != "" {
		if secretServer, err = remote.NewGRPCSecretFrom(o.secretServer); err != nil {
			return
		}

		template.SetSecretGetter(remote.NewGRPCSecretGetter(secretServer))
	}

	storeExtMgr := server.NewStoreExtManager(o.execer)

	remoteServer := server.NewRemoteServer(loader, remote.NewGRPCloaderFromStore(), secretServer, storeExtMgr, o.configDir)
	kinds, storeKindsErr := remoteServer.GetStoreKinds(ctx, nil)
	if storeKindsErr != nil {
		cmd.PrintErrf("failed to get store kinds, error: %p\n", storeKindsErr)
	} else {
		if err = startPlugins(storeExtMgr, kinds); err != nil {
			return
		}
	}

	clean := make(chan os.Signal, 1)
	signal.Notify(clean, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s := o.gRPCServer
	go func() {
		if gRPCServer, ok := s.(reflection.GRPCServer); ok {
			reflection.Register(gRPCServer)
		}
		server.RegisterRunnerServer(s, remoteServer)
		log.Printf("gRPC server listening at %v", lis.Addr())
		s.Serve(lis)
	}()

	go func() {
		<-clean
		_ = lis.Close()
		_ = o.httpServer.Shutdown(ctx)
		_ = storeExtMgr.StopAll()
	}()

	mux := runtime.NewServeMux(runtime.WithMetadata(server.MetadataStoreFunc)) //  runtime.WithIncomingHeaderMatcher(func(key string) (s string, b bool) {
	err = server.RegisterRunnerHandlerServer(ctx, mux, remoteServer)
	if err == nil {
		mux.HandlePath(http.MethodGet, "/", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/assets/{asset}", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/healthz", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/get", o.getAtestBinary)

		postRequestProxyFunc := postRequestProxy(o.skyWalking)
		mux.HandlePath(http.MethodPost, "/browser/{app}", postRequestProxyFunc)
		mux.HandlePath(http.MethodPost, "/v3/segments", postRequestProxyFunc)

		// Create non-global registry.
		reg := prometheus.NewRegistry()

		// Add go runtime metrics and process collectors.
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
		mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)
		})

		debugHandler(mux)
		o.httpServer.WithHandler(mux)
		log.Printf("HTTP server listening at %v", httplis.Addr())
		err = o.httpServer.Serve(httplis)
		err = util.IgnoreErrServerClosed(err)
	}
	return
}

func postRequestProxy(proxy string) func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if proxy == "" {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {}
	}

	if strings.HasSuffix(proxy, "/") {
		proxy = strings.TrimSuffix(proxy, "/")
	}

	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		http.Post(fmt.Sprintf("%s%s", proxy, r.URL.Path), "application/json", r.Body)
	}
}

func startPlugins(storeExtMgr server.ExtManager, kinds *server.StoreKinds) (err error) {
	const socketPrefix = "unix://"

	for _, kind := range kinds.Data {
		if kind.Enabled && strings.HasPrefix(kind.Url, socketPrefix) {
			if err = storeExtMgr.Start(kind.Name, kind.Url); err != nil {
				break
			}
		}
	}
	return
}

func frontEndHandlerWithLocation(consolePath string) func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		target := r.URL.Path
		if target == "/" {
			target = "/index.html"
		} else if target == "/healthz" {
			w.Write([]byte("ok"))
			return
		}

		var content string
		customHeader := map[string]string{}
		switch {
		case strings.HasSuffix(target, ".html"):
			content = uiResourceIndex
		case strings.HasSuffix(target, ".js"):
			content = uiResourceJS
			customHeader[util.ContentType] = "text/javascript; charset=utf-8"
		case strings.HasSuffix(target, ".css"):
			content = uiResourceCSS
			customHeader[util.ContentType] = "text/css"
		}

		if content != "" {
			for k, v := range customHeader {
				w.Header().Set(k, v)
			}
			http.ServeContent(w, r, "", time.Now(), bytes.NewReader([]byte(content)))
		} else {
			http.ServeFile(w, r, path.Join(consolePath, target))
		}
	}
}

func debugHandler(mux *runtime.ServeMux) {
	mux.HandlePath(http.MethodGet, "/debug/pprof/{sub}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		switch sub := pathParams["sub"]; sub {
		case "cmdline":
			pprof.Cmdline(w, r)
		case "profile":
			pprof.Profile(w, r)
		case "symbol":
			pprof.Symbol(w, r)
		case "trace":
			pprof.Trace(w, r)
		case "allocs", "block", "goroutine", "heap", "mutex", "threadcreate":
			pprof.Index(w, r)
		case "":
			pprof.Index(w, r)
		}
	})
}

var filesNeedToBeRemoved = []string{}

func (o *serverOption) getAtestBinary(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	w.Header().Set(util.ContentDisposition, "attachment; filename=atest")
	w.Header().Set(util.ContentType, "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache")

	var data []byte
	if atestPath, err := o.execer.LookPath("atest"); err == nil {
		if data, err = os.ReadFile(atestPath); err != nil {
			data = []byte(fmt.Sprintf("failed to read atest: %v", err))
		}
	} else {
		data = []byte("not found atest")
	}
	w.Write(data)
}

type gRPCServer interface {
	Serve(lis net.Listener) error
	grpc.ServiceRegistrar
}

type fakeGRPCServer struct {
}

// NewFakeGRPCServer creates a fake gRPC server
func NewFakeGRPCServer() gRPCServer {
	return &fakeGRPCServer{}
}

// Serve is a fake method
func (s *fakeGRPCServer) Serve(net.Listener) error {
	return nil
}

// RegisterService is a fake method
func (s *fakeGRPCServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	// Do nothing due to this is a fake method
}

//go:embed data/index.js
var uiResourceJS string

//go:embed data/index.css
var uiResourceCSS string

//go:embed data/index.html
var uiResourceIndex string
