package cmd

import (
	"os"

	"github.com/linuxsuren/api-testing/pkg/version"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd(execer fakeruntime.Execer) (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "atest",
		Short: "API testing tool",
	}
	c.SetOut(os.Stdout)
	c.Version = version.GetVersion()
	c.AddCommand(createInitCommand(),
		createRunCommand(), createSampleCmd(),
		createServerCmd(), createJSONSchemaCmd(),
		createServiceCommand(execer))
	return
}
