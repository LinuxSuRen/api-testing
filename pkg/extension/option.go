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

package extension

import (
	"fmt"
	"net"
	"os"

	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

// Extension is the default command option of the extension
type Extension struct {
	Port   int
	Socket string

	name string
	port int
}

func NewExtension(name string, port int) *Extension {
	return &Extension{
		name: name,
		port: port,
	}
}

func (o *Extension) AddFlags(flags *pflag.FlagSet) {
	flags.IntVarP(&o.Port, "port", "p", o.port, "The port to listen on")
	flags.StringVarP(&o.Socket, "socket", "", "",
		fmt.Sprintf("The socket to listen on, for instance: /var/run/%s.sock", StoreName(o.name)))
}

func (o *Extension) GetListenAddress() (protocol, address string) {
	if o.Socket != "" {
		protocol = "unix"
		address = o.Socket
	} else {
		protocol = "tcp"
		address = fmt.Sprintf(":%d", o.Port)
	}
	return
}

func (o *Extension) GetFullName() string {
	return StoreName(o.name)
}

func CreateRunner(ext *Extension, c *cobra.Command, removeServer remote.LoaderServer) (err error) {
	protocol, address := ext.GetListenAddress()

	var lis net.Listener
	lis, err = net.Listen(protocol, address)
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, removeServer)
	c.Printf("%s@%s is running at %s\n", ext.GetFullName(), version.GetVersion(), address)

	RegisterStopSignal(c.Context(), func() {
		_ = os.Remove(ext.Socket)
	}, gRPCServer)

	err = gRPCServer.Serve(lis)
	return
}

func StoreName(name string) string {
	return fmt.Sprintf("atest-store-%s", name)
}
