// Package cmd provides all the commands
package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"path"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func createServerCmd(gRPCServer gRPCServer, httpServer server.HTTPServer) (c *cobra.Command) {
	opt := &serverOption{
		gRPCServer: gRPCServer,
		httpServer: httpServer,
	}
	c = &cobra.Command{
		Use:   "server",
		Short: "Run as a server mode",
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7070, "The RPC server port")
	flags.IntVarP(&opt.httpPort, "http-port", "", 8080, "The HTTP server port")
	flags.BoolVarP(&opt.printProto, "print-proto", "", false, "Print the proto content and exit")
	flags.StringVarP(&opt.localStorage, "local-storage", "", "", "The local storage path")
	flags.StringVarP(&opt.consolePath, "console-path", "", "", "The path of the console")
	return
}

type serverOption struct {
	gRPCServer gRPCServer
	httpServer server.HTTPServer

	port         int
	httpPort     int
	printProto   bool
	localStorage string
	consolePath  string
}

func (o *serverOption) runE(cmd *cobra.Command, args []string) (err error) {
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
	if o.localStorage != "" {
		if err = loader.Put(o.localStorage); err != nil {
			return
		}
	}

	removeServer := server.NewRemoteServer(loader)
	s := o.gRPCServer
	go func() {
		server.RegisterRunnerServer(s, removeServer)
		log.Printf("gRPC server listening at %v", lis.Addr())
		s.Serve(lis)
	}()

	mux := runtime.NewServeMux()
	err = server.RegisterRunnerHandlerServer(cmd.Context(), mux, removeServer)
	if err == nil {
		mux.HandlePath("GET", "/", frontEndHandlerWithLocation(o.consolePath))
		mux.HandlePath("GET", "/assets/{asset}", frontEndHandlerWithLocation(o.consolePath))
		o.httpServer.WithHandler(mux)
		err = o.httpServer.Serve(httplis)
	}
	return
}

func frontEndHandlerWithLocation(consolePath string) func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		target := r.URL.Path
		if target == "/" {
			target = "/index.html"
		}

		http.ServeFile(w, r, path.Join(consolePath, target))
	}
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
