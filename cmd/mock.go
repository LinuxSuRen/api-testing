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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hpcloud/tail"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/spf13/cobra"
)

var (
	mockLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("mock")
	tailClient *tail.Tail
)

type mockOption struct {
	port       int
	prefix     string
	autoReload bool
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
	flags.BoolVarP(&opt.autoReload, "auto-reload", "", false, "Enable automatic refresh of mock config")
	return
}

func (o *mockOption) runE(c *cobra.Command, args []string) (err error) {

	reader := mock.NewLocalFileReader(args[0])
	server := mock.NewInMemoryServer(o.port)

	mockLogger.Info("Starting mock server", "port", o.port, "prefix", o.prefix, "autoReload", o.autoReload)
	if err = server.Start(reader, o.prefix); err != nil {
		return
	}

	// if mock auto refresh is enabled, start tailing the mock config file
	if o.autoReload {
		if err := o.enableAutoReload(args[0], server); err != nil {
			return err
		}
	}

	return o.shutdownServer(c, server)
}

func (o *mockOption) enableAutoReload(configFile string, server mock.DynamicServer) error {

	mockLogger.Info("Enable automatic refresh of mock config.")
	if tailClient == nil {
		initTail(configFile)
	}

	reader := processMockConfigFiles()
	if err := server.Stop(); err != nil {
		return fmt.Errorf("failed to stop server for reload: %w", err)
	}
	if err := server.Start(reader, o.prefix); err != nil {
		return fmt.Errorf("failed to restart server after reload: %w", err)
	}

	mockLogger.Info("Mock server restarted successfully by refresh configuration files.")
	return nil
}

func (o *mockOption) shutdownServer(c *cobra.Command, server mock.DynamicServer) error {

	clean := make(chan os.Signal, 1)
	signal.Notify(clean, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case <-c.Context().Done():
	case <-clean:
	}

	mockLogger.Info("Shutting down mock server!")
	if err := server.Stop(); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}

func initTail(file string) {

	// tail client configuration properties.
	tailConfig := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
	}

	// create tail client
	tc, _ := tail.TailFile(file, tailConfig)
	tailClient = tc
}

func processMockConfigFiles() mock.Reader {

	var reader mock.Reader
	if tailClient == nil {
		mockLogger.Error(nil, "Tail client is not initialized")
	}

	for range tailClient.Lines {
		reader = mock.NewLocalFileReader(tailClient.Filename)
	}

	return reader
}
