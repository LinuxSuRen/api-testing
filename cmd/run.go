package main

import (
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

type runOption struct {
	pattern string
}

func createRunCommand() (cmd *cobra.Command) {
	opt := &runOption{}
	cmd = &cobra.Command{
		Use:  "run",
		RunE: opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "testcase-*.yaml",
		"The file pattern which try to execute the test cases")
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	var files []string
	if files, err = filepath.Glob(o.pattern); err == nil {
		for i := range files {
			item := files[i]

			var testcase *testing.TestCase
			if testcase, err = testing.Parse(item); err != nil {
				return
			}

			setRelativeDir(item, testcase)

			if err = runner.RunTestCase(testcase); err != nil {
				return
			}
		}
	}
	return
}

func setRelativeDir(configFile string, testcase *testing.TestCase) {
	dir := filepath.Dir(configFile)

	for i := range testcase.Prepare.Kubernetes {
		testcase.Prepare.Kubernetes[i] = path.Join(dir, testcase.Prepare.Kubernetes[i])
	}
}
