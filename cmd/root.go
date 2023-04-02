package cmd

import "github.com/spf13/cobra"

// should be injected during the build process
var version string

// NewRootCmd creates the root command
func NewRootCmd() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "atest",
		Short: "API testing tool",
	}
	c.Version = version
	c.AddCommand(createInitCommand(),
		createRunCommand(), createSampleCmd(),
		createServerCmd())
	return
}
