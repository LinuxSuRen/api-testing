package cmd

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

type runOption struct {
	pattern            string
	duration           time.Duration
	requestTimeout     time.Duration
	requestIgnoreError bool
	thread             int64
	context            context.Context
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
	flags.DurationVarP(&opt.duration, "duration", "", 0, "Running duration")
	flags.DurationVarP(&opt.requestTimeout, "request-timeout", "", time.Minute, "Timeout for per request")
	flags.BoolVarP(&opt.requestIgnoreError, "request-ignore-error", "", false, "Indicate if ignore the request error")
	flags.Int64VarP(&opt.thread, "thread", "", 1, "Threads of the execution")
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	var files []string
	o.context = cmd.Context()

	if files, err = filepath.Glob(o.pattern); err == nil {
		for i := range files {
			item := files[i]
			if err = o.runSuiteWithDuration(item); err != nil {
				return
			}
		}
	}
	return
}

func (o *runOption) runSuiteWithDuration(suite string) (err error) {
	sem := semaphore.NewWeighted(o.thread)
	stop := false
	var timeout *time.Ticker
	if o.duration > 0 {
		timeout = time.NewTicker(o.duration)
	} else {
		// make sure having a valid timer
		timeout = time.NewTicker(time.Second)
	}
	errChannel := make(chan error, 10*o.thread)
	var wait sync.WaitGroup

	for !stop {
		select {
		case <-timeout.C:
			stop = true
		case err = <-errChannel:
			if err != nil {
				stop = true
			}
		default:
			if err := sem.Acquire(o.context, 1); err != nil {
				continue
			}
			wait.Add(1)
			if o.duration <= 0 {
				stop = true
			}

			go func(ch chan error, sem *semaphore.Weighted) {
				now := time.Now()
				defer sem.Release(1)
				defer wait.Done()
				defer func() {
					fmt.Println("routing end with", time.Now().Sub(now))
				}()

				dataContext := getDefaultContext()
				ch <- o.runSuite(suite, dataContext, o.context)
			}(errChannel, sem)
		}
	}
	err = <-errChannel
	wait.Wait()
	return
}

func (o *runOption) runSuite(suite string, dataContext map[string]interface{}, ctx context.Context) (err error) {
	var testSuite *testing.TestSuite
	if testSuite, err = testing.Parse(suite); err != nil {
		return
	}

	var result string
	if result, err = render.Render("base api", testSuite.API, dataContext); err == nil {
		testSuite.API = result
		testSuite.API = strings.TrimSuffix(testSuite.API, "/")
	} else {
		return
	}

	for _, testCase := range testSuite.Items {
		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", testSuite.API, testCase.Request.API)
		}

		setRelativeDir(suite, &testCase)
		var output interface{}
		ctxWithTimeout, _ := context.WithTimeout(ctx, o.requestTimeout)
		if output, err = runner.RunTestCase(&testCase, dataContext, ctxWithTimeout); err != nil && !o.requestIgnoreError {
			return
		}
		dataContext[testCase.Name] = output
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
