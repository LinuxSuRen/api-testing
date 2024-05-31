package cmd

import (
	"github.com/linuxsuren/api-testing/pkg/downloader"
	"os"

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
		createRunCommand(), createSampleCmd(),
		createServerCmd(execer, httpServer), createJSONSchemaCmd(),
		createServiceCommand(execer), createFunctionCmd(), createConvertCommand(),
		createMockCmd(), createExtensionCommand(downloader.NewDefaultOCIDownloader()))
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
