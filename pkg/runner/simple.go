package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/andreyvit/diff"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/exec"
	"github.com/linuxsuren/api-testing/pkg/testing"
	unstructured "github.com/linuxsuren/unstructured/pkg"
)

type TestCaseRunner interface {
	RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error)
	WithOutputWriter(io.Writer) TestCaseRunner
	WithTestReporter(TestReporter) TestCaseRunner
}

// ReportRecord represents the raw data of a HTTP request
type ReportRecord struct {
	Method    string
	API       string
	BeginTime time.Time
	EndTime   time.Time
	Error     error
}

// Duration returns the duration between begin and end time
func (r *ReportRecord) Duration() time.Duration {
	return r.EndTime.Sub(r.BeginTime)
}

// ErrorCount returns the count number of errors
func (r *ReportRecord) ErrorCount() int {
	if r.Error == nil {
		return 0
	}
	return 1
}

// NewReportRecord creates a record, and set the begin time to be now
func NewReportRecord() *ReportRecord {
	return &ReportRecord{
		BeginTime: time.Now(),
	}
}

// ReportResult represents the report result of a set of the same API requests
type ReportResult struct {
	API     string
	Count   int
	Average time.Duration
	Max     time.Duration
	Min     time.Duration
	Error   int
}

// ReportResultSlice is the alias type of ReportResult slice
type ReportResultSlice []ReportResult

// Len returns the count of slice items
func (r ReportResultSlice) Len() int {
	return len(r)
}

// Less returns if i bigger than j
func (r ReportResultSlice) Less(i, j int) bool {
	return r[i].Average > r[j].Average
}

// Swap swaps the items
func (r ReportResultSlice) Swap(i, j int) {
	tmp := r[i]
	r[i] = r[j]
	r[j] = tmp
}

// ReportResultWriter is the interface of the report writer
type ReportResultWriter interface {
	Output([]ReportResult) error
}

// TestReporter is the interface of the report
type TestReporter interface {
	PutRecord(*ReportRecord)
	GetAllRecords() []*ReportRecord
	ExportAllReportResults() (ReportResultSlice, error)
}

type simpleTestCaseRunner struct {
	testReporter TestReporter
	writer       io.Writer
}

// NewSimpleTestCaseRunner creates the instance of the simple test case runner
func NewSimpleTestCaseRunner() TestCaseRunner {
	runner := &simpleTestCaseRunner{}
	return runner.WithOutputWriter(io.Discard).WithTestReporter(NewDiscardTestReporter())
}

// RunTestCase is the main entry point of a test case
func (r *simpleTestCaseRunner) RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error) {
	fmt.Fprintf(r.writer, "start to run: '%s'\n", testcase.Name)
	record := NewReportRecord()
	defer func(rr *ReportRecord) {
		rr.EndTime = time.Now()
		rr.Error = err
		rr.API = testcase.Request.API
		rr.Method = testcase.Request.Method
		r.testReporter.PutRecord(rr)
	}(record)

	if err = doPrepare(testcase); err != nil {
		err = fmt.Errorf("failed to prepare, error: %v", err)
		return
	}

	defer func() {
		if testcase.Clean.CleanPrepare {
			if err = doCleanPrepare(testcase); err != nil {
				return
			}
		}
	}()

	client := http.Client{}
	var requestBody io.Reader
	if testcase.Request.Body != "" {
		requestBody = bytes.NewBufferString(testcase.Request.Body)
	} else if testcase.Request.BodyFromFile != "" {
		var data []byte
		if data, err = os.ReadFile(testcase.Request.BodyFromFile); err != nil {
			return
		}
		requestBody = bytes.NewBufferString(string(data))
	}

	if err = testcase.Request.Render(dataContext); err != nil {
		return
	}

	if len(testcase.Request.Form) > 0 {
		if testcase.Request.Header["Content-Type"] == "multipart/form-data" {
			multiBody := &bytes.Buffer{}
			writer := multipart.NewWriter(multiBody)
			for key, val := range testcase.Request.Form {
				writer.WriteField(key, val)
			}

			_ = writer.Close()
			requestBody = multiBody
			testcase.Request.Header["Content-Type"] = writer.FormDataContentType()
		} else if testcase.Request.Header["Content-Type"] == "application/x-www-form-urlencoded" {
			data := url.Values{}
			for key, val := range testcase.Request.Form {
				data.Set(key, val)
			}
			requestBody = strings.NewReader(data.Encode())
		}
	}

	var request *http.Request
	if request, err = http.NewRequestWithContext(ctx, testcase.Request.Method, testcase.Request.API, requestBody); err != nil {
		return
	}

	// set headers
	for key, val := range testcase.Request.Header {
		request.Header.Add(key, val)
	}

	fmt.Fprintf(r.writer, "start to send request to %s\n", testcase.Request.API)

	// send the HTTP request
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	var responseBodyData []byte
	if responseBodyData, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	if err = testcase.Expect.Render(nil); err != nil {
		return
	}
	if err = expectInt(testcase.Name, testcase.Expect.StatusCode, resp.StatusCode); err != nil {
		err = fmt.Errorf("error is: %v\n%s", err, string(responseBodyData))
		return
	}

	for key, val := range testcase.Expect.Header {
		actualVal := resp.Header.Get(key)
		if err = expectString(testcase.Name, val, actualVal); err != nil {
			return
		}
	}

	if testcase.Expect.Body != "" {
		if string(responseBodyData) != strings.TrimSpace(testcase.Expect.Body) {
			err = fmt.Errorf("case: %s, got different response body, diff: \n%s", testcase.Name,
				diff.LineDiff(testcase.Expect.Body, string(responseBodyData)))
			return
		}
	}

	mapOutput := map[string]interface{}{}
	if err = json.Unmarshal(responseBodyData, &mapOutput); err != nil {
		switch b := err.(type) {
		case *json.UnmarshalTypeError:
			if b.Value != "array" {
				return
			}

			var arrayOutput []interface{}
			if err = json.Unmarshal(responseBodyData, &arrayOutput); err != nil {
				return
			}
			output = arrayOutput
		default:
			return
		}
	} else {
		output = mapOutput
	}

	for key, expectVal := range testcase.Expect.BodyFieldsExpect {
		var val interface{}
		var ok bool
		if val, ok, err = unstructured.NestedField(mapOutput, strings.Split(key, "/")...); err != nil {
			err = fmt.Errorf("failed to get field: %s, %v", key, err)
			return
		} else if !ok {
			err = fmt.Errorf("not found field: %s", key)
			return
		} else if !reflect.DeepEqual(expectVal, val) {
			if reflect.TypeOf(expectVal).Kind() == reflect.Int {
				if strings.Compare(fmt.Sprintf("%v", expectVal), fmt.Sprintf("%v", val)) == 0 {
					continue
				}
			}
			err = fmt.Errorf("field[%s] expect value: %v, actual: %v", key, expectVal, val)
			return
		}
	}

	for _, verify := range testcase.Expect.Verify {
		var program *vm.Program
		if program, err = expr.Compile(verify, expr.Env(output), expr.AsBool()); err != nil {
			return
		}

		var result interface{}
		if result, err = expr.Run(program, output); err != nil {
			return
		}

		if !result.(bool) {
			err = fmt.Errorf("failed to verify: %s", verify)
			break
		}
	}
	return
}

// WithOutputWriter sets the io.Writer
func (r *simpleTestCaseRunner) WithOutputWriter(writer io.Writer) TestCaseRunner {
	r.writer = writer
	return r
}

// WithTestReporter sets the TestReporter
func (r *simpleTestCaseRunner) WithTestReporter(reporter TestReporter) TestCaseRunner {
	r.testReporter = reporter
	return r
}

// Deprecated
// RunTestCase runs the test case.
func RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error) {
	return NewSimpleTestCaseRunner().WithOutputWriter(os.Stdout).RunTestCase(testcase, dataContext, ctx)
}

func doPrepare(testcase *testing.TestCase) (err error) {
	for i := range testcase.Prepare.Kubernetes {
		item := testcase.Prepare.Kubernetes[i]

		if err = exec.RunCommand("kubectl", "apply", "-f", item); err != nil {
			return
		}
	}
	return
}

func doCleanPrepare(testcase *testing.TestCase) (err error) {
	count := len(testcase.Prepare.Kubernetes)
	for i := count - 1; i >= 0; i-- {
		item := testcase.Prepare.Kubernetes[i]

		if err = exec.RunCommand("kubectl", "delete", "-f", item); err != nil {
			return
		}
	}
	return
}

func expectInt(name string, expect, actual int) (err error) {
	if expect != actual {
		err = fmt.Errorf("case: %s, expect %d, actual %d", name, expect, actual)
	}
	return
}

func expectString(name, expect, actual string) (err error) {
	if expect != actual {
		err = fmt.Errorf("case: %s, expect %s, actual %s", name, expect, actual)
	}
	return
}
