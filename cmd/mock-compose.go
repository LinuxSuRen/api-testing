/*
Copyright 2025 API Testing Authors.

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
	"os"
	"os/signal"
	"syscall"

	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/spf13/cobra"
)

func createMockComposeCmd() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "mock-compose",
		Short: "Mock multiple servers",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reader := mock.NewLocalFileReader(args[0])

			var server *mock.Server
			if server, err = reader.Parse(); err != nil {
				return
			}

			var subServers []mock.DynamicServer
			for _, proxy := range server.Proxies {
				subProxy := &mock.Server{
					Proxies: []mock.Proxy{proxy},
				}

				subReader := mock.NewObjectReader(subProxy)
				subServer := mock.NewInMemoryServer(c.Context(), proxy.Port)
				if err = subServer.Start(subReader, proxy.Prefix); err != nil {
					return
				}
				subServers = append(subServers, subServer)
			}

			clean := make(chan os.Signal, 1)
			signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
			select {
			case <-c.Context().Done():
			case <-clean:
			}
			for _, server := range subServers {
				server.Stop()
			}
			return
		},
	}
	return
}
