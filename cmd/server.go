/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd provides all the commands
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	_ "embed"
	pprof "net/http/pprof"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/linuxsuren/api-testing/pkg/oauth"
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
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func createServerCmd(execer fakeruntime.Execer, httpServer server.HTTPServer) (c *cobra.Command) {
	opt := &serverOption{
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
	flags.StringVarP(&opt.auth, "auth", "", os.Getenv("AUTH_MODE"), "The auth mode, supported: oauth. Keep it empty to disable auth")
	flags.StringVarP(&opt.oauthProvider, "oauth-provider", "", "github", "The oauth provider, supported: github")
	flags.StringVarP(&opt.oauthServer, "oauth-server", "", "", "The oAuth server address, required if it is a private server")
	flags.BoolVarP(&opt.oauthSkipTls, "oauth-skip-tls", "", false, "Skip TLS verify when connect to oauth server")
	flags.StringArrayVarP(&opt.oauthGroup, "oauth-group", "", []string{}, "Alow specific groups, all groups is ok if it is empty")
	flags.StringVarP(&opt.clientID, "client-id", "", os.Getenv("OAUTH_CLIENT_ID"), "ClientID is the application's ID")
	flags.StringVarP(&opt.clientSecret, "client-secret", "", os.Getenv("OAUTH_CLIENT_SECRET"), "ClientSecret is the application's secret")
	flags.BoolVarP(&opt.dryRun, "dry-run", "", false, "Do not really start a gRPC server")

	// gc releated flags
	flags.IntVarP(&opt.gcPercent, "gc-percent", "", 100, "The GC percent of Go")

	c.Flags().MarkHidden("dry-run")
	c.Flags().MarkHidden("gc-percent")
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

	auth          string
	oauthProvider string
	// ClientID is the application's ID.
	clientID string
	// ClientSecret is the application's secret.
	clientSecret string
	oauthServer  string
	oauthSkipTls bool
	oauthGroup   []string

	gcPercent int

	dryRun bool

	// inner fields, not as command flags
	provider oauth.OAuthProvider
}

func (o *serverOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	o.execer.WithContext(cmd.Context())
	var grpcOpts []grpc.ServerOption

	if o.auth == "oauth" {
		if o.provider = oauth.GetOAuthProvider(o.oauthProvider); o.provider == nil {
			err = fmt.Errorf("not support: %q", o.oauthProvider)
			return
		}

		if o.provider.GetServer() != "" {
			// returns empty string if it's a private server
			o.oauthServer = o.provider.GetServer()
		} else {
			o.provider.SetServer(o.oauthServer)
		}

		if o.clientID == "" || o.clientSecret == "" {
			err = errors.New("--client-id and --client-secret flags are required when auth enabled")
			return
		}

		if o.oauthServer == "" {
			err = errors.New("oAuth server address is required")
			return
		}

		grpcOpts = append(grpcOpts, oauth.NewAuthInterceptor(o.oauthGroup))
	}

	if o.dryRun {
		o.gRPCServer = &fakeGRPCServer{}
	} else {
		o.gRPCServer = grpc.NewServer(grpcOpts...)
	}

	o.configDir = os.ExpandEnv(o.configDir)
	err = o.execer.MkdirAll(o.configDir, 0755)
	return
}

func (o *serverOption) runE(cmd *cobra.Command, args []string) (err error) {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	debug.SetGCPercent(o.gcPercent)

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
	defer storeExtMgr.StopAll()

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
	signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

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
		log.Println("stopping the server")
		_ = lis.Close()
		_ = o.httpServer.Shutdown(ctx)
	}()

	mux := runtime.NewServeMux(runtime.WithMetadata(server.MetadataStoreFunc))
	err = server.RegisterRunnerHandlerFromEndpoint(ctx, mux, "127.0.0.1:7070", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err == nil {
		mux.HandlePath(http.MethodGet, "/", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/assets/{asset}", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/healthz", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/favicon.ico", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath(http.MethodGet, "/get", o.getAtestBinary)

		postRequestProxyFunc := postRequestProxy(o.skyWalking)
		mux.HandlePath(http.MethodPost, "/browser/{app}", postRequestProxyFunc)
		mux.HandlePath(http.MethodPost, "/v3/segments", postRequestProxyFunc)

		// Create non-global registry.
		reg := prometheus.NewRegistry()

		// register oauth endpoint
		if o.auth == "oauth" {
			authHandler := oauth.NewAuth(o.provider, oauth2.Config{
				ClientID:     o.clientID,
				ClientSecret: o.clientSecret,
			}, o.oauthSkipTls)
			mux.HandlePath(http.MethodGet, "/oauth2/token", authHandler.RequestCode)
			mux.HandlePath(http.MethodGet, "/oauth2/getLocalCode", authHandler.RequestLocalCode)
			mux.HandlePath(http.MethodGet, "/oauth2/getUserInfoFromLocalCode", authHandler.RequestLocalToken)
			mux.HandlePath(http.MethodGet, "/oauth2/callback", authHandler.Callback)
		}

		// Add go runtime metrics and process collectors.
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
		mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)
		})

		debugHandler(mux, remoteServer)
		o.httpServer.WithHandler(mux)
		log.Printf("HTTP server listening at %v", httplis.Addr())
		log.Printf("Server is running.")
		err = o.httpServer.Serve(httplis)
		err = util.IgnoreErrServerClosed(err)
	}
	return
}

func postRequestProxy(proxy string) func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if proxy == "" {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {}
	}

	proxy = strings.TrimSuffix(proxy, "/")
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
		case strings.HasSuffix(target, ".ico"):
			content = uiResourceIcon
			customHeader[util.ContentType] = "image/x-icon"
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

func debugHandler(mux *runtime.ServeMux, remoteServer server.RunnerServer) {
	mux.HandlePath(http.MethodGet, "/debug/pprof/{sub}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		sub := pathParams["sub"]
		extName := r.URL.Query().Get("name")
		if extName != "" && remoteServer != nil {
			log.Println("get pprof of extension:", extName)

			ctx := metadata.NewIncomingContext(r.Context(), metadata.New(map[string]string{
				server.HeaderKeyStoreName: extName,
			}))

			data, err := remoteServer.PProf(ctx, &server.PProfRequest{
				Name: sub,
			})
			if err == nil {
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write(data.Data)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			}

			return
		}

		switch sub {
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

func (o *serverOption) getAtestBinary(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	name := util.EmptyThenDefault(r.URL.Query().Get("name"), "atest")

	w.Header().Set(util.ContentDisposition, fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Set(util.ContentType, "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-cache")

	var data []byte
	if atestPath, err := o.execer.LookPath(name); err == nil {
		if data, err = os.ReadFile(atestPath); err != nil {
			data = []byte(fmt.Sprintf("failed to read %q: %v", name, err))
		}
	} else {
		data = []byte(fmt.Sprintf("not found %q", name))
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

//go:embed data/favicon.ico
var uiResourceIcon string
