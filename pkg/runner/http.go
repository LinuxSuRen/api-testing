/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package runner

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andreyvit/diff"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/xeipuuv/gojsonschema"
)

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
		rr.Group = testcase.Group
		rr.Name = testcase.Name
		rr.EndTime = time.Now()
		rr.Error = err
		rr.API = testcase.Request.API
		rr.Method = testcase.Request.Method
		r.testReporter.PutRecord(rr)
	}(record)

	defer func() {
		if err == nil {
			err = runJob(testcase.After, dataContext, output)
		}
	}()

	client := util.TlsAwareHTTPClient(true) // TODO should have a way to change it
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

	if err = runJob(testcase.Before, dataContext, nil); err != nil {
		return
	}

	r.log.Info("start to send request to %s\n", testcase.Request.API)

	// TODO only do this for unit testing, should remove it once we have a better way
	if strings.HasPrefix(testcase.Request.API, "http://") {
		client = http.DefaultClient
	}

	// send the HTTP request
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	r.log.Debug("status code: %d\n", resp.StatusCode)

	var responseBodyData []byte
	if responseBodyData, err = r.withResponseRecord(resp); err != nil {
		return
	}
	record.Body = string(responseBodyData)
	r.log.Trace("response body: %s\n", record.Body)

	if err = testcase.Expect.Render(dataContext); err != nil {
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

	respType := util.GetFirstHeaderValue(resp.Header, util.ContentType)
	if output, err = verifyResponseBodyData(testcase.Name, testcase.Expect, respType, responseBodyData); err != nil {
		return
	}

	err = jsonSchemaValidation(testcase.Expect.Schema, responseBodyData)
	return
}

func (r *simpleTestCaseRunner) GetSuggestedAPIs(suite *testing.TestSuite, api string) (result []*testing.TestCase, err error) {
	if suite.Spec.URL == "" || suite.Spec.Kind != "swagger" {
		return
	}

	var swaggerAPI *apispec.Swagger
	if swaggerAPI, err = apispec.ParseURLToSwagger(suite.Spec.URL); err == nil && swaggerAPI != nil {
		result = []*testing.TestCase{}
		for api, item := range swaggerAPI.Paths {
			for method, oper := range item {
				result = append(result, &testing.TestCase{
					Name: oper.OperationId,
					Request: testing.Request{
						API:    api,
						Method: strings.ToUpper(method),
					},
				})
			}
		}
	}
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

func verifyResponseBodyData(caseName string, expect testing.Response, responseType string, responseBodyData []byte) (output interface{}, err error) {
	if expect.Body != "" {
		if string(responseBodyData) != strings.TrimSpace(expect.Body) {
			err = fmt.Errorf("case: %s, got different response body, diff: \n%s", caseName,
				diff.LineDiff(expect.Body, string(responseBodyData)))
			return
		}
	}

	verifier := NewBodyVerify(responseType, expect)
	if verifier == nil {
		log.Printf("no body verify support with %q\n", responseType)
		return
	}

	if output, err = verifier.Parse(responseBodyData); err != nil {
		return
	}

	mapOutput := map[string]interface{}{
		"data": output,
	}
	if err = verifier.Verify(responseBodyData); err == nil {
		err = Verify(expect, mapOutput)
	}
	return
}

func runJob(job *testing.Job, ctx interface{}, current interface{}) (err error) {
	if job == nil {
		return
	}
	var program *vm.Program
	env := map[string]interface{}{
		"ctx":     ctx,
		"current": current,
	}

	for _, item := range job.Items {
		var exprText string
		if exprText, err = render.Render("job", item, ctx); err != nil {
			err = fmt.Errorf("failed to render: %q, error is: %v", item, err)
			break
		}

		if program, err = expr.Compile(exprText, expr.Env(env)); err != nil {
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
