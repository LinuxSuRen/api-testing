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
	"os"
	"os/signal"

	"github.com/linuxsuren/api-testing/extensions/store-git/pkg"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewRootCommand returns the root Command
func NewRootCommand() (c *cobra.Command) {
	opt := &options{}
	c = &cobra.Command{
		Use:   "atest-store-git",
		Short: "A store extension for git",
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "p", 7074, "The port to listen on")
	flags.StringVarP(&opt.socket, "socket", "", "", "The socket to listen on, for instance: /var/run/atest-ext-store-git.sock")
	return
}

type options struct {
	port   int
	socket string
}

func (o *options) getListenAddress() (protocol, address string) {
	if o.socket != "" {
		protocol = "unix"
		address = o.socket
	} else {
		protocol = "tcp"
		address = fmt.Sprintf(":%d", o.port)
	}
	return
}

func (o *options) runE(c *cobra.Command, args []string) (err error) {
	removeServer := pkg.NewRemoteServer()
	protocol, address := o.getListenAddress()

	var lis net.Listener
	lis, err = net.Listen(protocol, address)
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	c.Printf("Git storage extension is running at %s\n", address)

	endChan := make(chan os.Signal, 1)
	signal.Notify(endChan, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-endChan:
		case <-c.Context().Done():
		}
		fmt.Println("Stopping the server...")
		_ = os.Remove(o.socket)
		gRPCServer.Stop()
	}()

	err = gRPCServer.Serve(lis)
	return
}
