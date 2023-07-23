package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/limit"
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

	// for internal use
	loader testing.Loader
}

func newDefaultRunOption() *runOption {
	return &runOption{
		reporter:     runner.NewMemoryTestReporter(),
		reportWriter: runner.NewResultWriter(os.Stdout),
		loader:       testing.NewFileLoader(),
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
		Short:   "Run the test suite",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}

	// set flags
	flags := cmd.Flags()
	flags.StringVarP(&opt.pattern, "pattern", "p", "test-suite-*.yaml",
		"The file pattern which try to execute the test cases. Brace expansion is supported, such as: test-suite-{1,2}.yaml")
	flags.StringVarP(&opt.level, "level", "l", "info", "Set the output log level")
	flags.DurationVarP(&opt.duration, "duration", "", 0, "Running duration")
	flags.DurationVarP(&opt.requestTimeout, "request-timeout", "", time.Minute, "Timeout for per request")
	flags.BoolVarP(&opt.requestIgnoreError, "request-ignore-error", "", false, "Indicate if ignore the request error")
	flags.StringVarP(&opt.report, "report", "", "", "The type of target report. Supported: markdown, md, html, json, discard, std")
	flags.StringVarP(&opt.reportFile, "report-file", "", "", "The file path of the report")
	flags.BoolVarP(&opt.reportIgnore, "report-ignore", "", false, "Indicate if ignore the report output")
	flags.StringVarP(&opt.swaggerURL, "swagger-url", "", "", "The URL of swagger")
	flags.Int64VarP(&opt.thread, "thread", "", 1, "Threads of the execution")
	flags.Int32VarP(&opt.qps, "qps", "", 5, "QPS")
	flags.Int32VarP(&opt.burst, "burst", "", 5, "burst")
	return
}

func (o *runOption) preRunE(cmd *cobra.Command, args []string) (err error) {
	writer := cmd.OutOrStdout()

	if o.reportFile != "" {
		var reportFile *os.File
		if reportFile, err = os.Create(o.reportFile); err != nil {
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

	o.caseItems = args
	return
}

func (o *runOption) runE(cmd *cobra.Command, args []string) (err error) {
	o.startTime = time.Now()
	o.context = cmd.Context()
	o.limiter = limit.NewDefaultRateLimiter(o.qps, o.burst)
	defer func() {
		cmd.Printf("consume: %s\n", time.Since(o.startTime).String())
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
					fmt.Println("routing end with", time.Since(now))
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

	for _, testCase := range testSuite.Items {
		if !testCase.InScope(o.caseItems) {
			continue
		}

		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", testSuite.API, testCase.Request.API)
		}

		var output interface{}
		select {
		case <-stopSingal:
			return
		default:
			o.limiter.Accept()

			ctxWithTimeout, _ := context.WithTimeout(ctx, o.requestTimeout)
			ctxWithTimeout = context.WithValue(ctxWithTimeout, runner.ContextKey("").ParentDir(), loader.GetContext())

			simpleRunner := runner.NewSimpleTestCaseRunner()
			simpleRunner.WithTestReporter(o.reporter)
			if output, err = simpleRunner.RunTestCase(&testCase, dataContext, ctxWithTimeout); err != nil && !o.requestIgnoreError {
				err = fmt.Errorf("failed to run '%s', %v", testCase.Name, err)
				return
			} else {
				err = nil
			}
		}
		dataContext[testCase.Name] = output
	}
	return
}

func getDefaultContext() map[string]interface{} {
	return map[string]interface{}{}
}
