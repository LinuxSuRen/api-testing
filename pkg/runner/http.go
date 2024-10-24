/*
Copyright 2023-2024 API Testing Authors.

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

package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/spec"

	"github.com/andreyvit/diff"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/xeipuuv/gojsonschema"
)

// ReportResult represents the report result of a set of the same API requests
type ReportResult struct {
	Name             string
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
	simpleResponse  SimpleResponse
	cookies         []*http.Cookie
	apiSuggestLimit int
}

// NewSimpleTestCaseRunner creates the instance of the simple test case runner
func NewSimpleTestCaseRunner() TestCaseRunner {
	runner := &simpleTestCaseRunner{
		UnimplementedRunner: NewDefaultUnimplementedRunner(),
		simpleResponse:      SimpleResponse{},
		cookies:             []*http.Cookie{},
		apiSuggestLimit:     10,
	}
	return runner
}

func init() {
	RegisterRunner("http", func(*testing.TestSuite) TestCaseRunner {
		return NewSimpleTestCaseRunner()
	})
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

	q := request.URL.Query()
	for k := range testcase.Request.Query {
		q.Add(k, testcase.Request.Query.GetValue(k))
	}
	request.URL.RawQuery = q.Encode()

	// set headers
	for key, val := range testcase.Request.Header {
		request.Header.Add(key, val)
	}

	if err = runJob(testcase.Before, dataContext, nil); err != nil {
		return
	}

	for _, cookie := range r.cookies {
		request.AddCookie(cookie)
	}
	for k, v := range testcase.Request.Cookie {
		request.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	r.log.Info("request method: %s\n", request.Method)
	r.log.Info("request header %v\n", request.Header)
	r.log.Info("start to send request to %v\n", request.URL)

	// TODO only do this for unit testing, should remove it once we have a better way
	if strings.HasPrefix(testcase.Request.API, "http://") {
		client = http.DefaultClient
	}

	// send the HTTP request
	var resp *http.Response
	if resp, err = client.Do(request); err != nil {
		return
	}

	r.log.Debug("test case %q, status code: %d\n", testcase.Name, resp.StatusCode)

	if err = testcase.Expect.Render(dataContext); err != nil {
		return
	}
	if err = expectInt(testcase.Name, testcase.Expect.StatusCode, resp.StatusCode); err != nil {
		err = fmt.Errorf("error is: %v", err)
	}

	for key, val := range testcase.Expect.Header {
		actualVal := resp.Header.Get(key)
		err = errors.Join(err, expectString(testcase.Name, val, actualVal))
	}

	respType := util.GetFirstHeaderValue(resp.Header, util.ContentType)

	if isNonBinaryContent(respType) {
		var responseBodyData []byte
		var rErr error
		if responseBodyData, rErr = r.withResponseRecord(resp); rErr != nil {
			err = errors.Join(err, rErr)
			return
		}

		record.Body = string(responseBodyData)
		r.log.Trace("response body: %s\n", record.Body)

		if output, rErr = verifyResponseBodyData(testcase.Name, testcase.Expect, respType, responseBodyData); rErr != nil {
			err = errors.Join(err, rErr)
			return
		}

		err = errors.Join(err, jsonSchemaValidation(testcase.Expect.Schema, responseBodyData))
	} else {
		r.log.Debug("skip to read the body due to it is not struct content: %q\n", respType)
	}

	r.cookies = append(r.cookies, resp.Cookies()...)
	return
}

func ammendHeaders(headers http.Header, body []byte) {
	// add content-length if it's missing
	if val := headers.Get(util.ContentLength); val == "" {
		headers.Add(util.ContentLength, strconv.Itoa(len(body)))
		fmt.Printf("add content-length: %d\n", len(body))
	}
}

func (r *simpleTestCaseRunner) GetSuggestedAPIs(suite *testing.TestSuite, api string) (result []*testing.TestCase, err error) {
	if suite.Spec.URL == "" || suite.Spec.Kind != "swagger" {
		return
	}

	var swagger *spec.Swagger
	if swagger, err = apispec.ParseURLToSwagger(suite.Spec.URL); err == nil && swagger != nil {
		swaggerAPI := apispec.NewSwaggerAPI(swagger)
		for api, methods := range swaggerAPI.ApiMap {
			for _, method := range methods {
				testcase := &testing.TestCase{
					Name: swagger.ID,
					Request: testing.Request{
						API:    api,
						Method: strings.ToUpper(method),
						Query:  make(testing.SortedKeysStringMap),
					},
				}

				switch testcase.Request.Method {
				case http.MethodGet:
					for _, param := range swagger.Paths.Paths[api].Get.Parameters {
						switch param.In {
						case "query":
							// TODO should have a better way to provide the initial value
							(&(testcase.Request)).Query[param.Name] = generateRandomValue(param)
						}
					}
					testcase.Name = swagger.Paths.Paths[api].Get.ID
				case http.MethodPost:
					testcase.Name = swagger.Paths.Paths[api].Post.ID
				case http.MethodPut:
					testcase.Name = swagger.Paths.Paths[api].Put.ID
				case http.MethodDelete:
					testcase.Name = swagger.Paths.Paths[api].Delete.ID
				case http.MethodPatch:
					testcase.Name = swagger.Paths.Paths[api].Patch.ID
				}
				result = append(result, testcase)
				if len(result) >= r.apiSuggestLimit {
					return
				}
			}
		}
	}
	return
}

func generateRandomValue(param spec.Parameter) interface{} {
	switch param.Format {
	case "int32", "int64":
		return 101
	case "boolean":
		return true
	case "string":
		return "random"
	default:
		return "random"
	}
}

func (r *simpleTestCaseRunner) withResponseRecord(resp *http.Response) (responseBodyData []byte, err error) {
	responseBodyData, err = io.ReadAll(resp.Body)
	r.simpleResponse = SimpleResponse{
		StatusCode: resp.StatusCode,
		Header:     make(map[string]string),
		Body:       string(responseBodyData),
	}

	// add some headers for convienience
	ammendHeaders(resp.Header, responseBodyData)
	for key := range resp.Header {
		r.simpleResponse.Header[key] = resp.Header.Get(key)
	}
	return
}

// GetResponseRecord returns the response record
func (r *simpleTestCaseRunner) GetResponseRecord() SimpleResponse {
	return r.simpleResponse
}

func (r *simpleTestCaseRunner) WithAPISuggestLimit(limit int) {
	r.apiSuggestLimit = limit
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
		runnerLogger.Info("no body verify support with", "response type", responseType)
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
			fmt.Printf("failed to compile: %q, %v\n", exprText, err)
			return
		}

		if _, err = expr.Run(program, env); err != nil {
			fmt.Printf("failed to Run: %q, %v\n", exprText, err)
			return
		}
	}
	return
}

// isNonBinaryContent detect if the content belong to binary
func isNonBinaryContent(contentType string) bool {
	if IsJSONCompatileType(contentType) {
		return true
	}

	switch contentType {
	case util.JSON, util.YAML, util.Plain, util.OCIImageIndex:
		return true
	default:
		return false
	}
}

func IsJSONCompatileType(contentType string) bool {
	return strings.HasSuffix(contentType, "+json")
}
