// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"errors"
	"fmt"
	"os"
	reflect "reflect"
	"regexp"
	"strings"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/version"
	"github.com/linuxsuren/api-testing/sample"
	"gopkg.in/yaml.v3"
)

type server struct {
	UnimplementedRunnerServer
	loader testing.Writer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(loader testing.Writer) RunnerServer {
	return &server{loader: loader}
}

func withDefaultValue(old, defVal any) any {
	if old == "" || old == nil {
		old = defVal
	}
	return old
}

func parseSuiteWithItems(data []byte) (suite *testing.TestSuite, err error) {
	suite, err = testing.ParseFromData(data)
	if err == nil && (suite == nil || suite.Items == nil) {
		err = errNoTestSuiteFound
	}
	return
}

func (s *server) getSuiteFromTestTask(task *TestTask) (suite *testing.TestSuite, err error) {
	switch task.Kind {
	case "suite":
		suite, err = parseSuiteWithItems([]byte(task.Data))
	case "testcase":
		var testCase *testing.TestCase
		if testCase, err = testing.ParseTestCaseFromData([]byte(task.Data)); err != nil {
			return
		}
		suite = &testing.TestSuite{
			Items: []testing.TestCase{*testCase},
		}
	case "testcaseInSuite":
		suite, err = parseSuiteWithItems([]byte(task.Data))
		if err != nil {
			return
		}

		var targetTestcase *testing.TestCase
		for _, item := range suite.Items {
			if item.Name == task.CaseName {
				targetTestcase = &item
				break
			}
		}

		if targetTestcase != nil {
			parentCases := findParentTestCases(targetTestcase, suite)
			fmt.Printf("find %d parent cases\n", len(parentCases))
			suite.Items = append(parentCases, *targetTestcase)
		} else {
			err = fmt.Errorf("cannot found testcase %s", task.CaseName)
		}
	default:
		err = fmt.Errorf("not support '%s'", task.Kind)
	}
	return
}

func resetEnv(oldEnv map[string]string) {
	for key, val := range oldEnv {
		os.Setenv(key, val)
	}
}

// Run start to run the test task
func (s *server) Run(ctx context.Context, task *TestTask) (reply *TestResult, err error) {
	task.Level = withDefaultValue(task.Level, "info").(string)
	task.Env = withDefaultValue(task.Env, map[string]string{}).(map[string]string)

	var suite *testing.TestSuite

	// TODO may not safe in multiple threads
	oldEnv := map[string]string{}
	for key, val := range task.Env {
		oldEnv[key] = os.Getenv(key)
		os.Setenv(key, val)
	}

	defer func() {
		resetEnv(oldEnv)
	}()

	if suite, err = s.getSuiteFromTestTask(task); err != nil {
		return
	}

	fmt.Printf("prepare to run: %s, with level: %s\n", suite.Name, task.Level)
	fmt.Printf("task kind: %s, %d to run\n", task.Kind, len(suite.Items))
	dataContext := map[string]interface{}{}

	var result string
	if result, err = render.Render("base api", suite.API, dataContext); err == nil {
		suite.API = result
		suite.API = strings.TrimSuffix(suite.API, "/")
	} else {
		reply.Error = err.Error()
		err = nil
		return
	}

	buf := new(bytes.Buffer)
	reply = &TestResult{}

	for _, testCase := range suite.Items {
		simpleRunner := runner.NewSimpleTestCaseRunner()
		simpleRunner.WithOutputWriter(buf)
		simpleRunner.WithWriteLevel(task.Level)

		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", suite.API, testCase.Request.API)
		}

		output, testErr := simpleRunner.RunTestCase(&testCase, dataContext, ctx)
		if getter, ok := simpleRunner.(runner.HTTPResponseRecord); ok {
			resp := getter.GetResponseRecord()
			reply.TestCaseResult = append(reply.TestCaseResult, &TestCaseResult{
				StatusCode: int32(resp.StatusCode),
				Body:       resp.Body,
				Header:     mapToPair(resp.Header),
				Id:         testCase.ID,
				Output:     buf.String(),
			})
		}

		if testErr == nil {
			dataContext[testCase.Name] = output
		} else {
			reply.Error = testErr.Error()
			break
		}
	}

	if reply.Error != "" {
		fmt.Fprintln(buf, reply.Error)
	}
	reply.Message = buf.String()
	return
}

// GetVersion returns the version
func (s *server) GetVersion(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: version.GetVersion()}
	return
}

func (s *server) GetSuites(ctx context.Context, in *Empty) (reply *Suites, err error) {
	reply = &Suites{
		Data: make(map[string]*Items),
	}

	var suites []testing.TestSuite
	if suites, err = s.loader.ListTestSuite(); err == nil && suites != nil {
		for _, suite := range suites {
			items := &Items{}
			for _, item := range suite.Items {
				items.Data = append(items.Data, item.Name)
			}
			reply.Data[suite.Name] = items
		}
	}

	return
}

func (s *server) CreateTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	err = s.loader.CreateSuite(in.Name, in.Api)
	return
}

func (s *server) GetTestSuite(ctx context.Context, in *TestSuiteIdentity) (result *TestSuite, err error) {
	var suite *testing.TestSuite
	if suite, _, err = s.loader.GetSuite(in.Name); err == nil && suite != nil {
		result = &TestSuite{
			Name: suite.Name,
			Api:  suite.API,
			Spec: &APISpec{
				Kind: suite.Spec.Kind,
				Url:  suite.Spec.URL,
			},
		}
	}
	return
}

func convertToTestingTestSuite(in *TestSuite) (suite *testing.TestSuite) {
	suite = &testing.TestSuite{
		Name: in.Name,
		API:  in.Api,
	}
	if in.Spec != nil {
		suite.Spec = testing.APISpec{
			Kind: in.Spec.Kind,
			URL:  in.Spec.Url,
		}
	}
	return
}

func (s *server) UpdateTestSuite(ctx context.Context, in *TestSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = s.loader.UpdateSuite(*convertToTestingTestSuite(in))
	return
}

func (s *server) DeleteTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = s.loader.DeleteSuite(in.Name)
	return
}

func (s *server) ListTestCase(ctx context.Context, in *TestSuiteIdentity) (result *Suite, err error) {
	var items []testing.TestCase
	if items, err = s.loader.ListTestCase(in.Name); err == nil {
		result = &Suite{}
		for _, item := range items {
			result.Items = append(result.Items, convertToGRPCTestCase(item))
		}
	}
	return
}

func (s *server) GetTestCase(ctx context.Context, in *TestCaseIdentity) (reply *TestCase, err error) {
	var result testing.TestCase
	if result, err = s.loader.GetTestCase(in.Suite, in.Testcase); err == nil {
		reply = convertToGRPCTestCase(result)
	}
	return
}

func convertToGRPCTestCase(testCase testing.TestCase) (result *TestCase) {
	req := &Request{
		Api:    testCase.Request.API,
		Method: testCase.Request.Method,
		Query:  mapToPair(testCase.Request.Query),
		Header: mapToPair(testCase.Request.Header),
		Form:   mapToPair(testCase.Request.Form),
		Body:   testCase.Request.Body,
	}

	resp := &Response{
		StatusCode:       int32(testCase.Expect.StatusCode),
		Body:             testCase.Expect.Body,
		Header:           mapToPair(testCase.Expect.Header),
		BodyFieldsExpect: mapInterToPair(testCase.Expect.BodyFieldsExpect),
		Verify:           testCase.Expect.Verify,
		Schema:           testCase.Expect.Schema,
	}

	result = &TestCase{
		Name:     testCase.Name,
		Request:  req,
		Response: resp,
	}
	return
}

func (s *server) RunTestCase(ctx context.Context, in *TestCaseIdentity) (result *TestCaseResult, err error) {
	var targetTestSuite testing.TestSuite

	targetTestSuite, err = s.loader.GetTestSuite(in.Suite, true)
	if err != nil {
		err = nil
		result = &TestCaseResult{
			Error: fmt.Sprintf("not found suite: %s", in.Suite),
		}
		return
	}

	var data []byte
	if data, err = yaml.Marshal(targetTestSuite); err == nil {
		task := &TestTask{
			Kind:     "testcaseInSuite",
			Data:     string(data),
			CaseName: in.Testcase,
			Level:    "debug",
		}

		var reply *TestResult
		if reply, err = s.Run(ctx, task); err == nil && len(reply.TestCaseResult) > 0 {
			lastIndex := len(reply.TestCaseResult) - 1
			lastItem := reply.TestCaseResult[lastIndex]

			result = &TestCaseResult{
				Output:     reply.Message,
				Error:      reply.Error,
				Body:       lastItem.Body,
				Header:     lastItem.Header,
				StatusCode: lastItem.StatusCode,
			}
		} else if err != nil {
			result = &TestCaseResult{
				Error: err.Error(),
			}
		} else {
			result = &TestCaseResult{
				Output: reply.Message,
				Error:  reply.Error,
			}
		}
	}
	return
}

func mapInterToPair(data map[string]interface{}) (pairs []*Pair) {
	pairs = make([]*Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &Pair{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return
}

func mapToPair(data map[string]string) (pairs []*Pair) {
	pairs = make([]*Pair, 0)
	for k, v := range data {
		pairs = append(pairs, &Pair{
			Key:   k,
			Value: v,
		})
	}
	return
}

func pairToInterMap(pairs []*Pair) (data map[string]interface{}) {
	data = make(map[string]interface{})
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func pairToMap(pairs []*Pair) (data map[string]string) {
	data = make(map[string]string)
	for _, pair := range pairs {
		data[pair.Key] = pair.Value
	}
	return
}

func convertToTestingTestCase(in *TestCase) (result testing.TestCase) {
	result = testing.TestCase{
		Name: in.Name,
	}
	req := in.Request
	resp := in.Response

	if req != nil {
		result.Request.API = req.Api
		result.Request.Method = req.Method
		result.Request.Body = req.Body
		result.Request.Header = pairToMap(req.Header)
		result.Request.Form = pairToMap(req.Form)
		result.Request.Query = pairToMap(req.Query)
	}

	if resp != nil {
		result.Expect.Body = resp.Body
		result.Expect.Schema = resp.Schema
		result.Expect.StatusCode = int(resp.StatusCode)
		result.Expect.Verify = resp.Verify
		result.Expect.BodyFieldsExpect = pairToInterMap(resp.BodyFieldsExpect)
		result.Expect.Header = pairToMap(resp.Header)
	}
	return
}

func (s *server) CreateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = s.loader.CreateTestCase(in.SuiteName, convertToTestingTestCase(in.Data))
	return
}

func (s *server) UpdateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	if in.Data == nil {
		err = errors.New("data is required")
		return
	}
	err = s.loader.UpdateTestCase(in.SuiteName, convertToTestingTestCase(in.Data))
	return
}

func (s *server) DeleteTestCase(ctx context.Context, in *TestCaseIdentity) (reply *HelloReply, err error) {
	err = s.loader.DeleteTestCase(in.Suite, in.Testcase)
	return
}

// Sample returns a sample of the test task
func (s *server) Sample(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: sample.TestSuiteGitLab}
	return
}

// PopularHeaders returns a list of popular headers
func (s *server) PopularHeaders(ctx context.Context, in *Empty) (pairs *Pairs, err error) {
	pairs = &Pairs{
		Data: []*Pair{},
	}

	err = yaml.Unmarshal([]byte(popularHeaders), &pairs.Data)
	return
}

// GetSuggestedAPIs returns a list of suggested APIs
func (s *server) GetSuggestedAPIs(ctx context.Context, in *TestSuiteIdentity) (reply *TestCases, err error) {
	reply = &TestCases{}

	var suite *testing.TestSuite
	if suite, _, err = s.loader.GetSuite(in.Name); err != nil {
		return
	}

	if suite == nil || suite.Spec.URL == "" {
		return
	}

	reply.Data = []*TestCase{{
		Request: &Request{},
	}}

	var swaggerAPI *apispec.Swagger
	if swaggerAPI, err = apispec.ParseURLToSwagger(suite.Spec.URL); err == nil {
		for api, item := range swaggerAPI.Paths {
			for method, oper := range item {
				reply.Data = append(reply.Data, &TestCase{
					Name: oper.OperationId,
					Request: &Request{
						Api:    api,
						Method: strings.ToUpper(method),
					},
				})
			}
		}
	}
	return
}

// FunctionsQuery returns a list of functions
func (s *server) FunctionsQuery(ctx context.Context, in *SimpleQuery) (reply *Pairs, err error) {
	reply = &Pairs{}
	in.Name = strings.ToLower(in.Name)

	for name, fn := range render.FuncMap() {
		lowerCaseName := strings.ToLower(name)
		if in.Name == "" || strings.Contains(lowerCaseName, in.Name) {
			reply.Data = append(reply.Data, &Pair{
				Key:   name,
				Value: fmt.Sprintf("%v", reflect.TypeOf(fn)),
			})
		}
	}
	return
}

//go:embed data/headers.yaml
var popularHeaders string

func findParentTestCases(testcase *testing.TestCase, suite *testing.TestSuite) (testcases []testing.TestCase) {
	reg, matchErr := regexp.Compile(`(.*?\{\{.*\.\w*.*?\}\})`)
	targetReg, targetErr := regexp.Compile(`\.\w*`)

	expectNames := new(UniqueSlice[string])
	if matchErr == nil && targetErr == nil {
		var expectName string
		for _, val := range testcase.Request.Header {
			if matched := reg.MatchString(val); matched {
				expectName = targetReg.FindString(val)
				expectName = strings.TrimPrefix(expectName, ".")
				expectNames.Push(expectName)
			}
		}

		findExpectNames(testcase.Request.API, expectNames)
		findExpectNames(testcase.Request.Body, expectNames)

		fmt.Println("expect test case names", expectNames.GetAll())
		for _, item := range suite.Items {
			if expectNames.Exist(item.Name) {
				testcases = append(testcases, item)
			}
		}
	}
	return
}

func findExpectNames(target string, expectNames *UniqueSlice[string]) {
	reg, _ := regexp.Compile(`(.*?\{\{.*\.\w*.*?\}\})`)
	targetReg, _ := regexp.Compile(`\.\w*`)

	for _, sub := range reg.FindStringSubmatch(target) {
		// remove {{ and }}
		if left, leftErr := regexp.Compile(`.*\{\{`); leftErr == nil {
			body := left.ReplaceAllString(sub, "")

			expectName := targetReg.FindString(body)
			expectName = strings.TrimPrefix(expectName, ".")
			expectNames.Push(expectName)
		}
	}
}

// UniqueSlice represents an unique slice
type UniqueSlice[T comparable] struct {
	data []T
}

// Push pushes an item if it's not exist
func (s *UniqueSlice[T]) Push(item T) *UniqueSlice[T] {
	if s.data == nil {
		s.data = []T{item}
	} else {
		for _, it := range s.data {
			if it == item {
				return s
			}
		}
		s.data = append(s.data, item)
	}
	return s
}

// Exist checks if the item exist, return true it exists
func (s *UniqueSlice[T]) Exist(item T) bool {
	if s.data != nil {
		for _, it := range s.data {
			if it == item {
				return true
			}
		}
	}
	return false
}

// GetAll returns all the items
func (s *UniqueSlice[T]) GetAll() []T {
	return s.data
}

var errNoTestSuiteFound = errors.New("no test suite found")
