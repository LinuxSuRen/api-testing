package runner

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/andreyvit/diff"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

// LevelWriter represents a writer with level
type LevelWriter interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
}

// FormatPrinter represents a formart printer with level
type FormatPrinter interface {
	Fprintf(w io.Writer, level, format string, a ...any) (n int, err error)
}

type defaultLevelWriter struct {
	level int
	io.Writer
	FormatPrinter
}

// NewDefaultLevelWriter creates a default LevelWriter instance
func NewDefaultLevelWriter(level string, writer io.Writer) LevelWriter {
	result := &defaultLevelWriter{
		Writer: writer,
	}
	switch level {
	case "debug":
		result.level = 7
	case "info":
		result.level = 3
	}
	return result
}

// Fprintf implements interface FormatPrinter
func (w *defaultLevelWriter) Fprintf(writer io.Writer, level int, format string, a ...any) (n int, err error) {
	if level <= w.level {
		return fmt.Fprintf(writer, format, a...)
	}
	return
}

// Info writes the info level message
func (w *defaultLevelWriter) Info(format string, a ...any) {
	w.Fprintf(w.Writer, 3, format, a...)
}

// Debug writes the debug level message
func (w *defaultLevelWriter) Debug(format string, a ...any) {
	w.Fprintf(w.Writer, 7, format, a...)
}

// ReportResult represents the report result of a set of the same API requests
type ReportResult struct {
	API              string
	Count            int
	Average          time.Duration
	Max              time.Duration
	Min              time.Duration
	QPS              int
	Error            int
	LastErrorMessage string
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

type simpleTestCaseRunner struct {
	UnimplementedRunner
	simpleResponse SimpleResponse
}

// NewSimpleTestCaseRunner creates the instance of the simple test case runner
func NewSimpleTestCaseRunner() TestCaseRunner {
	runner := &simpleTestCaseRunner{
		UnimplementedRunner: NewDefaultUnimplementedRunner(),
		simpleResponse:      SimpleResponse{},
	}
	return runner
}

// ContextKey is the alias type of string for context key
type ContextKey string

// NewContextKeyBuilder returns an emtpy context key
func NewContextKeyBuilder() ContextKey {
	return ContextKey("")
}

// ParentDir returns the key of the parsent directory
func (c ContextKey) ParentDir() ContextKey {
	return ContextKey("parentDir")
}

// GetContextValueOrEmpty returns the value of the context key, if not exist, return empty string
func (c ContextKey) GetContextValueOrEmpty(ctx context.Context) string {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(string)
	}
	return ""
}

// RunTestCase is the main entry point of a test case
func (r *simpleTestCaseRunner) RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error) {
	r.log.Info("start to run: '%s'\n", testcase.Name)
	record := NewReportRecord()
	defer func(rr *ReportRecord) {
		rr.EndTime = time.Now()
		rr.Error = err
		rr.API = testcase.Request.API
		rr.Method = testcase.Request.Method
		r.testReporter.PutRecord(rr)
	}(record)

	defer func() {
		if err == nil {
			err = runJob(testcase.After)
		}
	}()

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	contextDir := NewContextKeyBuilder().ParentDir().GetContextValueOrEmpty(ctx)
	if err = testcase.Request.Render(dataContext, contextDir); err != nil {
		return
	}

	var requestBody io.Reader
	if requestBody, err = testcase.Request.GetBody(); err != nil {
		return
	}

	var request *http.Request
	if request, err = http.NewRequestWithContext(ctx, testcase.Request.Method, testcase.Request.API, requestBody); err != nil {
		return
	}

	// set headers
	for key, val := range testcase.Request.Header {
		request.Header.Add(key, val)
	}

	if err = runJob(testcase.Before); err != nil {
		return
	}

	r.log.Info("start to send request to %s\n", testcase.Request.API)

	// TODO only do this for unit testing, should remove it once we have a better way
	if strings.HasPrefix(testcase.Request.API, "http://") {
		client = *http.DefaultClient
	}

	// send the HTTP request
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	var responseBodyData []byte
	if responseBodyData, err = r.withResponseRecord(resp); err != nil {
		return
	}
	record.Body = string(responseBodyData)
	r.log.Debug("response body: %s\n", record.Body)

	if err = testcase.Expect.Render(nil); err != nil {
		return
	}
	if err = expectInt(testcase.Name, testcase.Expect.StatusCode, resp.StatusCode); err != nil {
		err = fmt.Errorf("error is: %v", err)
		return
	}

	for key, val := range testcase.Expect.Header {
		actualVal := resp.Header.Get(key)
		if err = expectString(testcase.Name, val, actualVal); err != nil {
			return
		}
	}

	if output, err = verifyResponseBodyData(testcase.Name, testcase.Expect, responseBodyData); err != nil {
		return
	}

	err = jsonSchemaValidation(testcase.Expect.Schema, responseBodyData)
	return
}

func (r *simpleTestCaseRunner) withResponseRecord(resp *http.Response) (responseBodyData []byte, err error) {
	responseBodyData, err = io.ReadAll(resp.Body)
	r.simpleResponse = SimpleResponse{
		StatusCode: resp.StatusCode,
		Header:     make(map[string]string),
		Body:       string(responseBodyData),
	}
	for key := range resp.Header {
		r.simpleResponse.Header[key] = resp.Header.Get(key)
	}
	return
}

// GetResponseRecord returns the response record
func (r *simpleTestCaseRunner) GetResponseRecord() SimpleResponse {
	return r.simpleResponse
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

func jsonSchemaValidation(schema string, body []byte) (err error) {
	if schema == "" {
		return
	}

	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonLoader := gojsonschema.NewBytesLoader(body)

	var result *gojsonschema.Result
	if result, err = gojsonschema.Validate(schemaLoader, jsonLoader); err == nil && !result.Valid() {
		err = fmt.Errorf("JSON schema validation failed: %v", result.Errors())
	}
	return
}

func verifyResponseBodyData(caseName string, expect testing.Response, responseBodyData []byte) (output interface{}, err error) {
	if expect.Body != "" {
		if string(responseBodyData) != strings.TrimSpace(expect.Body) {
			err = fmt.Errorf("case: %s, got different response body, diff: \n%s", caseName,
				diff.LineDiff(expect.Body, string(responseBodyData)))
			return
		}
	}

	var bodyMap map[string]interface{}
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
			mapOutput["data"] = arrayOutput
		default:
			return
		}
	} else {
		bodyMap = mapOutput
		output = mapOutput
		mapOutput = map[string]interface{}{
			"data": bodyMap,
		}
	}

	if err = bodyFieldsVerify(expect.BodyFieldsExpect, responseBodyData); err != nil {
		return
	}

	err = Verify(expect, mapOutput)
	return
}

func bodyFieldsVerify(bodyFieldsExpect map[string]interface{}, responseBodyData []byte) (err error) {
	for key, expectVal := range bodyFieldsExpect {
		result := gjson.Get(string(responseBodyData), key)
		if result.Exists() {
			err = valueCompare(expectVal, result, key)
		} else {
			err = fmt.Errorf("not found field: %s", key)
		}

		if err != nil {
			break
		}
	}
	return
}

func valueCompare(expect interface{}, acutalResult gjson.Result, key string) (err error) {
	var actual interface{}
	actual = acutalResult.Value()

	if !reflect.DeepEqual(expect, actual) {
		switch acutalResult.Type {
		case gjson.Number:
			expect = fmt.Sprintf("%v", expect)
			actual = fmt.Sprintf("%v", actual)

			if strings.Compare(expect.(string), actual.(string)) == 0 {
				return
			}
		}
		err = fmt.Errorf("field[%s] expect value: '%v', actual: '%v'", key, expect, actual)
	}
	return
}

func runJob(job *testing.Job) (err error) {
	if job == nil {
		return
	}
	var program *vm.Program
	env := struct{}{}

	for _, item := range job.Items {
		if program, err = expr.Compile(item, expr.Env(env),
			expr.Function("sleep", ExprFuncSleep)); err != nil {
			fmt.Printf("failed to compile: %s, %v\n", item, err)
			return
		}

		if _, err = expr.Run(program, env); err != nil {
			fmt.Printf("failed to Run: %s, %v\n", item, err)
			return
		}
	}
	return
}
