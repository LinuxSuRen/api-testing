package main

import (
	"os"

	_ "github.com/apache/skywalking-go"
	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	exec "github.com/linuxsuren/go-fake-runtime"
	"google.golang.org/grpc"
)

func main() {
	gRPCServer := grpc.NewServer()
	c := cmd.NewRootCmd(exec.DefaultExecer{}, gRPCServer,
		server.NewDefaultHTTPServer())
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
