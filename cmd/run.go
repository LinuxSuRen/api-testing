package cmd

import (
	"path"
	"path/filepath"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
)

type runOption struct {
	pattern string
}

// CreateRunCommand returns the run command
func CreateRunCommand() (cmd *cobra.Command) {
	opt := &runOption{}
	cmd = &cobra.Command{
		Use:  "run",
		RunE: opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "test-suite-*.yaml",
		"The file pattern which try to execute the test cases")
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	var files []string

	ctx := map[string]interface{}{}

	if files, err = filepath.Glob(o.pattern); err == nil {
		for i := range files {
			item := files[i]

			var testSuite *testing.TestSuite
			if testSuite, err = testing.Parse(item); err != nil {
				return
			}

			for _, testCase := range testSuite.Items {
				setRelativeDir(item, &testCase)
				var output interface{}
				if output, err = runner.RunTestCase(&testCase, ctx); err != nil {
					return
				}
				ctx[testCase.Name] = output
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
