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
	"errors"

	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/spf13/cobra"
)

type mockOption struct {
	port   int
	prefix string
	files  []string
}

func createMockCmd() (c *cobra.Command) {
	opt := &mockOption{}

	c = &cobra.Command{
		Use:     "mock",
		Short:   "Start a mock server",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}

	flags := c.Flags()
	flags.IntVarP(&opt.port, "port", "", 6060, "The mock server port")
	flags.StringVarP(&opt.prefix, "prefix", "", "/mock", "The mock server API prefix")
	flags.StringSliceVarP(&opt.files, "files", "", nil, "The mock config files")
	return
}

func (o *mockOption) preRunE(c *cobra.Command, args []string) (err error) {
	if len(o.files) == 0 {
		err = errors.New("at least one file is required")
	}
	return
}

func (o *mockOption) runE(c *cobra.Command, args []string) (err error) {
	reader := mock.NewLocalFileReader(o.files[0])

	server := mock.NewInMemoryServer(o.port)

	if err = server.Start(reader, o.prefix); err != nil {
		return
	}

	<-c.Context().Done()
	server.Stop()
	return
}
