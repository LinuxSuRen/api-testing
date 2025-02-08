/*
Copyright 2024-2025 API Testing Authors.

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

package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/spf13/cobra"
)

type mockOption struct {
	port    int
	prefix  string
	metrics bool
}

func createMockCmd() (c *cobra.Command) {
	opt := &mockOption{}

	c = &cobra.Command{
		Use:   "mock",
		Short: "Start a mock server",
		Args:  cobra.ExactArgs(1),
		RunE:  opt.runE,
	}

	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "", 6060, "The mock server port")
	flags.StringVarP(&opt.prefix, "prefix", "", "/mock", "The mock server API prefix")
	flags.BoolVarP(&opt.metrics, "metrics", "m", true, "Enable request metrics collection")
	return
}

func (o *mockOption) runE(c *cobra.Command, args []string) (err error) {
	reader := mock.NewLocalFileReader(args[0])
	server := mock.NewInMemoryServer(c.Context(), o.port)
	if o.metrics {
		server.EnableMetrics()
	}
	if err = server.Start(reader, o.prefix); err != nil {
		return
	}

	clean := make(chan os.Signal, 1)
	signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	printLocalIPs(c, o.port)
	if o.metrics {
		c.Printf("Metrics available at http://localhost:%d%s/metrics\n", o.port, o.prefix)
	}

	select {
	case <-c.Context().Done():
	case <-clean:
	}
	err = server.Stop()
	return
}

func printLocalIPs(c *cobra.Command, port int) {
	if ips, err := getLocalIPs(); err == nil {
		for _, ip := range ips {
			c.Printf("server is available at http://%s:%d\n", ip, port)
		}
	}
}

func getLocalIPs() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get interface addresses: %v", err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil && !ipNet.IP.IsLoopback() {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips, nil
}
