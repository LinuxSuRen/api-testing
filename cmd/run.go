package cmd

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

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
		Use:     "run",
		Aliases: []string{"r"},
		Example: `atest run -p sample.yaml
See also https://github.com/LinuxSuRen/api-testing/tree/master/sample`,
		Short: "Run the test suite",
		RunE:  opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "test-suite-*.yaml",
		"The file pattern which try to execute the test cases")
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	var files []string
	ctx := getDefaultContext()

	if files, err = filepath.Glob(o.pattern); err == nil {
		for i := range files {
			item := files[i]
			if err = runSuite(item, ctx); err != nil {
				return
			}
		}
	}
	return
}

func runSuite(suite string, ctx map[string]interface{}) (err error) {
	var testSuite *testing.TestSuite
	if testSuite, err = testing.Parse(suite); err != nil {
		return
	}

	testSuite.API = strings.TrimSuffix(testSuite.API, "/")
	for _, testCase := range testSuite.Items {
		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", testSuite.API, testCase.Request.API)
		}

		setRelativeDir(suite, &testCase)
		var output interface{}
		if output, err = runner.RunTestCase(&testCase, ctx); err != nil {
			return
		}
		ctx[testCase.Name] = output
	}
	return
}

func getDefaultContext() map[string]interface{} {
	return map[string]interface{}{}
}

func setRelativeDir(configFile string, testcase *testing.TestCase) {
	dir := filepath.Dir(configFile)

	for i := range testcase.Prepare.Kubernetes {
		testcase.Prepare.Kubernetes[i] = path.Join(dir, testcase.Prepare.Kubernetes[i])
	}
}
