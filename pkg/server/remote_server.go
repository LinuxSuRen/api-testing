// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

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

// Run start to run the test task
func (s *server) Run(ctx context.Context, task *TestTask) (reply *HelloReply, err error) {
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
		for key, val := range oldEnv {
			os.Setenv(key, val)
		}
	}()

	switch task.Kind {
	case "suite":
		if suite, err = testing.ParseFromData([]byte(task.Data)); err != nil {
			return
		} else if suite == nil || suite.Items == nil {
			err = fmt.Errorf("no test suite found")
			return
		}
	case "testcase":
		var testCase *testing.TestCase
		if testCase, err = testing.ParseTestCaseFromData([]byte(task.Data)); err != nil {
			return
		}
		suite = &testing.TestSuite{
			Items: []testing.TestCase{*testCase},
		}
	case "testcaseInSuite":
		if suite, err = testing.ParseFromData([]byte(task.Data)); err != nil {
			return
		} else if suite == nil || suite.Items == nil {
			err = fmt.Errorf("no test suite found")
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
			return
		}
	default:
		err = fmt.Errorf("not support '%s'", task.Kind)
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
	reply = &HelloReply{}

	for _, testCase := range suite.Items {
		simpleRunner := runner.NewSimpleTestCaseRunner()
		simpleRunner.WithOutputWriter(buf)
		simpleRunner.WithWriteLevel(task.Level)

		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", suite.API, testCase.Request.API)
		}

		if output, testErr := simpleRunner.RunTestCase(&testCase, dataContext, ctx); testErr == nil {
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
	defer func() {
		s.loader.Reset()
	}()

	reply = &Suites{
		Data: make(map[string]*Items),
	}
	for s.loader.HasMore() {
		var data []byte
		if data, err = s.loader.Load(); err != nil {
			continue
		}

		var testSuite *testing.TestSuite
		if testSuite, err = testing.Parse(data); err != nil {
			return
		}

		items := &Items{}
		for _, item := range testSuite.Items {
			items.Data = append(items.Data, item.Name)
		}
		reply.Data[testSuite.Name] = items
	}
	return
}

func (s *server) GetTestCase(ctx context.Context, in *TestCaseIdentity) (reply *TestCase, err error) {
	defer func() {
		s.loader.Reset()
	}()

	for s.loader.HasMore() {
		var data []byte
		if data, err = s.loader.Load(); err != nil {
			continue
		}

		var testSuite *testing.TestSuite
		if testSuite, err = testing.Parse(data); err != nil {
			return
		}

		if testSuite.Name != in.Suite {
			continue
		}

		for _, testCase := range testSuite.Items {
			if testCase.Name != in.Testcase {
				continue
			}

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

			reply = &TestCase{
				Name:     testCase.Name,
				Request:  req,
				Response: resp,
			}
			break
		}
	}
	return
}

func (s *server) RunTestCase(ctx context.Context, in *TestCaseIdentity) (result *TestCaseResult, err error) {
	defer func() {
		s.loader.Reset()
	}()

	var targetTestSuite *testing.TestSuite
	for s.loader.HasMore() {
		var data []byte
		if data, err = s.loader.Load(); err != nil {
			continue
		}

		var testSuite *testing.TestSuite
		if testSuite, err = testing.Parse(data); err != nil {
			continue
		}

		if testSuite.Name == in.Suite {
			targetTestSuite = testSuite
			break
		}
	}

	if targetTestSuite != nil {
		var data []byte
		if data, err = yaml.Marshal(targetTestSuite); err == nil {
			task := &TestTask{
				Kind:     "testcaseInSuite",
				Data:     string(data),
				CaseName: in.Testcase,
				Level:    "debug",
			}

			var reply *HelloReply
			if reply, err = s.Run(ctx, task); err == nil {
				result = &TestCaseResult{
					Body:  reply.Message,
					Error: reply.Error,
				}
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

func (s *server) UpdateTestCase(ctx context.Context, in *TestCaseWithSuite) (reply *HelloReply, err error) {
	defer func() {
		s.loader.Reset()
	}()

	if in.Data == nil {
		err = errors.New("data is required")
		return
	}

	var targetTestSuite *testing.TestSuite
	for s.loader.HasMore() {
		var data []byte
		if data, err = s.loader.Load(); err != nil {
			continue
		}

		var testSuite *testing.TestSuite
		if testSuite, err = testing.Parse(data); err != nil {
			continue
		}

		if testSuite.Name == in.SuiteName {
			targetTestSuite = testSuite
			break
		}
	}

	if targetTestSuite == nil {
		err = errors.New("no test suite found")
		return
	}

	found := false
	for i := range targetTestSuite.Items {
		item := targetTestSuite.Items[i]
		if item.Name == in.Data.Name {
			item.Request = grpcRequestToRaw(in.Data.Request)
			item.Expect = grpcResponseToRaw(in.Data.Response)

			err = s.loader.UpdateTestCase(in.SuiteName, item)
			found = true
			break
		}
	}

	if !found {
		item := testing.TestCase{
			Name:    in.Data.Name,
			Request: grpcRequestToRaw(in.Data.Request),
			Expect:  grpcResponseToRaw(in.Data.Response),
		}
		err = s.loader.UpdateTestCase(in.SuiteName, item)
	}
	return
}

func grpcRequestToRaw(request *Request) (req testing.Request) {
	if request == nil {
		return
	}
	req.API = request.Api
	req.Method = request.Method
	req.Header = pairToMap(request.Header)
	req.Query = pairToMap(request.Query)
	req.Form = pairToMap(request.Form)
	req.Body = request.Body
	return
}

func grpcResponseToRaw(response *Response) (req testing.Response) {
	if response == nil {
		return
	}
	req.StatusCode = int(response.StatusCode)
	req.Body = response.Body
	req.Header = pairToMap(response.Header)
	req.BodyFieldsExpect = pairToInterMap(response.BodyFieldsExpect)
	req.Verify = response.Verify
	req.Schema = response.Schema
	return
}

func (s *server) CreateTestSuite(ctx context.Context, in *TestSuiteIdentity) (reply *HelloReply, err error) {
	err = s.loader.CreateSuite(in.Name, in.Api)
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
