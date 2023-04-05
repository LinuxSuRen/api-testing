package cmd

import (
	"github.com/linuxsuren/api-testing/pkg/version"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command
func NewRootCmd() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "atest",
		Short: "API testing tool",
	}
	c.Version = version.GetVersion()
	c.AddCommand(createInitCommand(),
		createRunCommand(), createSampleCmd(),
		createServerCmd())
	return
}
