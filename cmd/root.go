package main

import (
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type option struct {
	pattern string
}

func main() {
	opt := &option{}

	cmd := &cobra.Command{
		Use:  "atest",
		RunE: opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "testcase-*.yaml",
		"The file pattern which try to execute the test cases")

	// run command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var files []string
	if files, err = filepath.Glob(o.pattern); err == nil {
		for i := range files {
			item := files[i]

			var testcase *testing.TestCase
			if testcase, err = testing.Parse(item); err != nil {
				return
			}

			if err = runner.RunTestCase(testcase); err != nil {
				return
			}
		}
	}
	return
}
