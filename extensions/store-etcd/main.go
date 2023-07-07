package main

import (
	"fmt"
	"net"

	"github.com/linuxsuren/api-testing/extensions/store-etcd/remote"
	"google.golang.org/grpc"
)

func main() {
	removeServer := NewRemoteServer(nil)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7071))
	if err != nil {
		fmt.Println(err)
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	gRPCServer.Serve(lis)
}
