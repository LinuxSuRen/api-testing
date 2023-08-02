/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"
	"net"

	"github.com/linuxsuren/api-testing/extensions/store-orm/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func NewRootCommand() (c *cobra.Command) {
	opt := &option{}
	c = &cobra.Command{
		Use:   "store-orm",
		Short: "Storage extension of api-testing",
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7071, "The port of gRPC server")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	removeServer := pkg.NewRemoteServer()

	var lis net.Listener
	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", o.port))
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)

	go func() {
		<-cmd.Context().Done()
		gRPCServer.Stop()
	}()
	cmd.Println("ORM storage extension is running at port", o.port)
	err = gRPCServer.Serve(lis)
	return
}

type option struct {
	port int
}
