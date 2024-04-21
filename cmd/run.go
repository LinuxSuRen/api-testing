/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/limit"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/runner/monitor"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/sync/semaphore"
	"google.golang.org/grpc"
)

type runOption struct {
	pattern            string
	duration           time.Duration
	requestTimeout     time.Duration
	requestIgnoreError bool
	thread             int64
	context            context.Context
	qps                int32
	burst              int32
	limiter            limit.RateLimiter
	startTime          time.Time
	reporter           runner.TestReporter
	reportFile         string
	reportWriter       runner.ReportResultWriter
	report             string
	reportIgnore       bool
	swaggerURL         string
	level              string
	caseItems          []string
	githubReportOption *runner.GithubPRCommentOption
	monitorDocker      string

	// for internal use
	loader testing.Loader
}

var (
	runLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("run")
)

func newDefaultRunOption() *runOption {
	return &runOption{
		reporter:           runner.NewMemoryTestReporter(nil, ""),
		reportWriter:       runner.NewResultWriter(os.Stdout),
		loader:             testing.NewFileLoader(),
		githubReportOption: &runner.GithubPRCommentOption{},
	}
}

func newDiscardRunOption() *runOption {
	return &runOption{
		reporter:     runner.NewDiscardTestReporter(),
		reportWriter: runner.NewDiscardResultWriter(),
	}
}

// createRunCommand returns the run command
func createRunCommand() (cmd *cobra.Command) {
	opt := newDefaultRunOption()
	cmd = &cobra.Command{
		Use:     "run",
		Aliases: []string{"r"},
		Example: `atest run -p sample.yaml
See also https://github.com/LinuxSuRen/api-testing/tree/master/sample`,
		Short:        "Run the test suite",
		PreRunE:      opt.preRunE,
		SilenceUsage: true,
		RunE:         opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "test-suite-*.yaml",
		"The file pattern which try to execute the test cases. Brace expansion is supported, such as: test-suite-{1,2}.yaml")
	flags.StringVarP(&opt.level, "level", "l", "info", "Set the output log level")
	flags.DurationVarP(&opt.duration, "duration", "", 0, "Running duration")
	flags.DurationVarP(&opt.requestTimeout, "request-timeout", "", time.Minute, "Timeout for per request")
	flags.BoolVarP(&opt.requestIgnoreError, "request-ignore-error", "", false, "Indicate if ignore the request error")
	flags.StringVarP(&opt.report, "report", "", "", "The type of target report. Supported: markdown, md, html, json, discard, std, prometheus")
	flags.StringVarP(&opt.reportFile, "report-file", "", "", "The file path of the report")
	flags.BoolVarP(&opt.reportIgnore, "report-ignore", "", false, "Indicate if ignore the report output")
	flags.StringVarP(&opt.swaggerURL, "swagger-url", "", "", "The URL of swagger")
	flags.Int64VarP(&opt.thread, "thread", "", 1, "Threads of the execution")
	flags.Int32VarP(&opt.qps, "qps", "", 5, "QPS")
	flags.Int32VarP(&opt.burst, "burst", "", 5, "burst")
	flags.StringVarP(&opt.monitorDocker, "monitor-docker", "", "", "The docker container name to monitor")
	addGitHubReportFlags(flags, opt.githubReportOption)
	return
}

func (o *runOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	o.context = cmd.Context()
	writer := cmd.OutOrStdout()

	if o.reportFile != "" && !strings.HasPrefix(o.reportFile, "http://") && !strings.HasPrefix(o.reportFile, "https://") {
		var reportFile *os.File
		if reportFile, err = os.OpenFile(o.reportFile, os.O_RDWR|os.O_CREATE, 0666); err != nil {
			return
		}

		writer = io.MultiWriter(writer, reportFile)
	}

	switch o.report {
	case "markdown", "md":
		o.reportWriter = runner.NewMarkdownResultWriter(writer)
	case "html":
		o.reportWriter = runner.NewHTMLResultWriter(writer)
	case "json":
		o.reportWriter = runner.NewJSONResultWriter(writer)
	case "discard":
		o.reportWriter = runner.NewDiscardResultWriter()
	case "", "std":
		o.reportWriter = runner.NewResultWriter(writer)
	case "pdf":
		o.reportWriter = runner.NewPDFResultWriter(writer)
	case "prometheus":
		if o.reportFile == "" {
			err = fmt.Errorf("report file is required for prometheus report")
			return
		}
		o.reporter = runner.NewPrometheusWriter(o.reportFile, false)
	case "github":
		o.githubReportOption.ReportFile = o.reportFile
		o.reportWriter, err = runner.NewGithubPRCommentWriter(o.githubReportOption)
	default:
		err = fmt.Errorf("not supported report type: '%s'", o.report)
	}

	if err == nil {
		var swaggerAPI apispec.APIConverage
		if o.swaggerURL != "" {
			if swaggerAPI, err = apispec.ParseURLToSwagger(o.swaggerURL); err == nil {
				o.reportWriter.WithAPIConverage(swaggerAPI)
			}
		}
	}

	if err == nil {
		err = o.startMonitor()
	}

	o.caseItems = args
	return
}

func (o *runOption) startMonitor() (err error) {
	if o.monitorDocker == "" {
		return
	}

	var monitorBin string
	if monitorBin, err = exec.LookPath("atest-monitor-docker"); err != nil {
		return
	}

	sockFile := os.ExpandEnv(fmt.Sprintf("$HOME/.config/atest/%s.sock", "atest-monitor-docker"))
	os.MkdirAll(filepath.Dir(sockFile), 0755)

	execer := fakeruntime.NewDefaultExecerWithContext(o.context)
	go func(socketURL, plugin string) {
		if err = execer.RunCommandWithIO(plugin, "", os.Stdout, os.Stderr, nil, "server", "--socket", socketURL); err != nil {
			runLogger.Info("failed to start %s, error: %v", socketURL, err)
		}
	}(sockFile, monitorBin)

	for i := 0; i < 6; i++ {
		_, fErr := os.Stat(sockFile)
		if fErr == nil {
			break
		}
		time.Sleep(time.Second)
	}

	var conn *grpc.ClientConn
	monitorServer := fmt.Sprintf("unix://%s", sockFile)
	if conn, err = grpc.Dial(monitorServer, grpc.WithInsecure()); err == nil {
		o.reporter = runner.NewMemoryTestReporter(monitor.NewMonitorClient(conn), o.monitorDocker)
	}
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	o.startTime = time.Now()
	o.limiter = limit.NewDefaultRateLimiter(o.qps, o.burst)
	defer func() {
		cmd.Printf("Consumed: %s\n", time.Since(o.startTime).String())
		o.limiter.Stop()
	}()

	if err = o.loader.Put(o.pattern); err != nil {
		return
	}

	cmd.Println("found suites:", o.loader.GetCount())
	for o.loader.HasMore() {
		if err = o.runSuiteWithDuration(o.loader); err != nil {
			break
		}
	}

	if o.reportIgnore {
		return
	}

	// print the report
	var reportErr error
	var results runner.ReportResultSlice
	if results, reportErr = o.reporter.ExportAllReportResults(); reportErr == nil {
		o.reportWriter.WithResourceUsage(o.reporter.GetResourceUsage())
		outputErr := o.reportWriter.Output(results)
		println(cmd, outputErr, "failed to Output all reports", outputErr)
	}
	println(cmd, reportErr, "failed to export all reports", reportErr)
	return
}

func (o *runOption) runSuiteWithDuration(loader testing.Loader) (err error) {
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
	stopSingal := make(chan struct{}, 1)
	var wait sync.WaitGroup

	for !stop {
		select {
		case <-timeout.C:
			stop = true
			stopSingal <- struct{}{}
		case err = <-errChannel:
			if err != nil {
				stop = true
			}
		default:
			if err := sem.Acquire(o.context, 1); err != nil {
				continue
			}
			wait.Add(1)

			go func(ch chan error, sem *semaphore.Weighted) {
				now := time.Now()
				defer sem.Release(1)
				defer wait.Done()
				defer func() {
					runLogger.Info("routing end with", time.Since(now))
				}()

				dataContext := getDefaultContext()
				ch <- o.runSuite(loader, dataContext, o.context, stopSingal)
			}(errChannel, sem)
			if o.duration <= 0 {
				stop = true
			}
		}
	}

	select {
	case err = <-errChannel:
	case <-stopSingal:
	}

	wait.Wait()
	return
}

func (o *runOption) runSuite(loader testing.Loader, dataContext map[string]interface{}, ctx context.Context, stopSingal chan struct{}) (err error) {
	var data []byte
	if data, err = loader.Load(); err != nil {
		return
	}

	var testSuite *testing.TestSuite
	if testSuite, err = testing.Parse(data); err != nil {
		return
	}

	if err = testSuite.Render(dataContext); err != nil {
		return
	}

	var errs []error
	suiteRunner := runner.GetTestSuiteRunner(testSuite)
	suiteRunner.WithTestReporter(o.reporter)
	suiteRunner.WithSecure(testSuite.Spec.Secure)
	suiteRunner.WithOutputWriter(os.Stdout)
	suiteRunner.WithWriteLevel(o.level)
	suiteRunner.WithSuite(testSuite)
	for _, testCase := range testSuite.Items {
		if !testCase.InScope(o.caseItems) {
			continue
		}

		testCase.Group = testSuite.Name
		testCase.Request.RenderAPI(testSuite.API)

		var output interface{}
		select {
		case <-stopSingal:
			return
		default:
			o.limiter.Accept()

			ctxWithTimeout, _ := context.WithTimeout(ctx, o.requestTimeout)
			ctxWithTimeout = context.WithValue(ctxWithTimeout, runner.ContextKey("").ParentDir(), loader.GetContext())

			output, err = suiteRunner.RunTestCase(&testCase, dataContext, ctxWithTimeout)
			if err = util.ErrorWrap(err, "failed to run '%s', %v", testCase.Name, err); err != nil {
				if o.requestIgnoreError {
					errs = append(errs, err)
				} else {
					return
				}
			}

			reverseRunner := runner.NewReverseHTTPRunner(suiteRunner)
			reverseRunner.WithTestReporter(runner.NewDiscardTestReporter())
			if _, err = reverseRunner.RunTestCase(
				&testCase, dataContext, ctxWithTimeout); err != nil {
				err = fmt.Errorf("got error in reverse test: %w", err)
				return
			}
			suiteRunner.WithTestReporter(o.reporter)
		}
		dataContext[testCase.Name] = output
	}

	if len(errs) > 0 {
		err = errors.Join(errs...)
	}
	return
}

func addGitHubReportFlags(flags *pflag.FlagSet, opt *runner.GithubPRCommentOption) {
	flags.StringVarP(&opt.Repo, "report-github-repo", "", "", "The GitHub repository for reporting, for instance: linuxsuren/api-testing")
	flags.IntVarP(&opt.PR, "report-github-pr", "", -1, "The GitHub pull-request number for reporting")
	flags.StringVarP(&opt.Identity, "report-github-identity", "", "Reported by api-testing.", "The identity for find the existing comment")
	flags.StringVarP(&opt.Token, "report-github-token", "", "", "GitHub token, take it from environment variable $GITHUB_TOKEN if this flag is empty")
}

func getDefaultContext() map[string]interface{} {
	return map[string]interface{}{}
}
