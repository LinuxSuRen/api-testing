package cmd

import (
	"fmt"
	"net"

	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewRootCmd(s3Creator pkg.S3Creator) (cmd *cobra.Command) {
	opt := &option{
		s3Creator: s3Creator,
	}
	cmd = &cobra.Command{
		Use:   "store-s3",
		Short: "S3 storage extension of api-testing",
		RunE:  opt.runE,
	}
	flags := cmd.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7072, "The port of gRPC server")
	return cmd
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	removeServer := pkg.NewRemoteServer(o.s3Creator)

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	cmd.Println("S3 storage extension is running at port", o.port)

	go func() {
		<-cmd.Context().Done()
		gRPCServer.Stop()
	}()

	err = gRPCServer.Serve(lis)
	return
}

type option struct {
	port int

	// inner fields
	s3Creator pkg.S3Creator
}
