package main

import (
	"fmt"
	"net"
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func main() {
	opt := &option{}
	cmd := &cobra.Command{
		Use:   "store-s3",
		Short: "S3 storage extension of api-testing",
		RunE:  opt.runE,
	}
	flags := cmd.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7072, "The port of gRPC server")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var removeServer remote.LoaderServer
	if removeServer, err = pkg.NewRemoteServer(); err != nil {
		return
	}

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	cmd.Println("S3 storage extension is running at port", o.port)
	err = gRPCServer.Serve(lis)
	return
}

type option struct {
	port int
}
