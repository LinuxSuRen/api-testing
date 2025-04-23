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

	"github.com/linuxsuren/api-testing/pkg/downloader"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/version"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(execer fakeruntime.Execer, httpServer server.HTTPServer) (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "atest",
		Short: "API testing tool",
	}
	c.SetOut(os.Stdout)
	c.Version = "\n" + version.GetDetailedVersion()
	c.AddCommand(createInitCommand(execer),
		createRunCommand(), createSampleCmd(), createMockComposeCmd(),
		createServerCmd(execer, httpServer), createJSONSchemaCmd(),
		createServiceCommand(execer), createFunctionCmd(), createConvertCommand(),
		createMockCmd(), createExtensionCommand(downloader.NewStoreDownloader()))
	return
}

type printer interface {
	Println(i ...interface{})
}

func println(printer printer, err error, i ...interface{}) {
	if err != nil {
		printer.Println(i...)
	}
}
