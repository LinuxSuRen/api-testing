// Package cmd provides all the commands
package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func createServerCmd() (c *cobra.Command) {
	opt := &serverOption{}
	c = &cobra.Command{
		Use:   "server",
		Short: "Run as a server mode",
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7070, "The RPC server port")
	flags.BoolVarP(&opt.printProto, "print-proto", "", false, "Print the proto content and exit")
	return
}

type serverOption struct {
	port       int
	printProto bool
}

func (o *serverOption) runE(cmd *cobra.Command, args []string) (err error) {
	if o.printProto {
		for _, val := range server.GetProtos() {
			cmd.Println(val)
		}
		return
	}

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}

	s := grpc.NewServer()
	server.RegisterRunnerServer(s, server.NewRemoteServer())
	log.Printf("server listening at %v", lis.Addr())
	s.Serve(lis)
	return
}
