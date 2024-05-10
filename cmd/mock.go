/*
Copyright 2024 API Testing Authors.

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

type mockOption struct {
	port   int
	prefix string
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
	return
}

func (o *mockOption) runE(c *cobra.Command, args []string) (err error) {
	reader := mock.NewLocalFileReader(args[0])
	server := mock.NewInMemoryServer(o.port)

	c.Println("start listen", o.port)
	if err = server.Start(reader, o.prefix); err != nil {
		return
	}

	clean := make(chan os.Signal, 1)
	signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case <-c.Context().Done():
	case <-clean:
	}
	err = server.Stop()
	return
}
