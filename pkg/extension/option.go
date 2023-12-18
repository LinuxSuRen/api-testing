/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package extension

import (
	"fmt"
	"net"
	"os"

	"github.com/linuxsuren/api-testing/pkg/runner/monitor"
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

	kind string
	name string
	port int
}

func NewExtension(name, kind string, port int) *Extension {
	return &Extension{
		name: name,
		kind: kind,
		port: port,
	}
}

func (o *Extension) AddFlags(flags *pflag.FlagSet) {
	flags.IntVarP(&o.Port, "port", "p", o.port, "The port to listen on")
	flags.StringVarP(&o.Socket, "socket", "", "",
		fmt.Sprintf("The socket to listen on, for instance: /var/run/%s.sock", o.GetFullName()))
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
	return fmt.Sprintf("atest-%s-%s", o.kind, o.name)
}

func CreateRunner(ext *Extension, c *cobra.Command, remoteServer remote.LoaderServer) (err error) {
	protocol, address := ext.GetListenAddress()

	var lis net.Listener
	lis, err = net.Listen(protocol, address)
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	remote.RegisterLoaderServer(gRPCServer, remoteServer)
	c.Printf("%s@%s is running at %s\n", ext.GetFullName(), version.GetVersion(), address)

	RegisterStopSignal(c.Context(), func() {
		_ = os.Remove(ext.Socket)
	}, gRPCServer)

	err = gRPCServer.Serve(lis)
	return
}

func CreateMonitor(ext *Extension, c *cobra.Command, remoteServer monitor.MonitorServer) (err error) {
	protocol, address := ext.GetListenAddress()

	var lis net.Listener
	lis, err = net.Listen(protocol, address)
	if err != nil {
		return
	}

	gRPCServer := grpc.NewServer()
	monitor.RegisterMonitorServer(gRPCServer, remoteServer)
	c.Printf("%s@%s is running at %s\n", ext.GetFullName(), version.GetVersion(), address)

	RegisterStopSignal(c.Context(), func() {
		_ = os.Remove(ext.Socket)
	}, gRPCServer)

	err = gRPCServer.Serve(lis)
	return
}
