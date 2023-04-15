package main

import (
	"os"

	"github.com/linuxsuren/api-testing/cmd"
	exec "github.com/linuxsuren/go-fake-runtime"
	"google.golang.org/grpc"
)

func main() {
	gRPCServer := grpc.NewServer()
	c := cmd.NewRootCmd(exec.DefaultExecer{}, gRPCServer)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
