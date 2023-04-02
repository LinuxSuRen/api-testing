package cmd

import (
	"github.com/linuxsuren/api-testing/sample"
	"github.com/spf13/cobra"
)

func createSampleCmd() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "sample",
		Short: "Generate a sample test case YAML file",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.Println(sample.TestSuiteGitLab)
			return
		},
	}
	return
}
