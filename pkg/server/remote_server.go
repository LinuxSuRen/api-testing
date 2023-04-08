// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/render"
	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/version"
)

type server struct {
	UnimplementedRunnerServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer() RunnerServer {
	return &server{}
}

// Run start to run the test task
func (s *server) Run(ctx context.Context, task *TestTask) (reply *HelloReply, err error) {
	if task.Level == "" {
		task.Level = "info"
	}

	var suite *testing.TestSuite
	if task.Env == nil {
		task.Env = map[string]string{}
	}

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
	reply.Message = buf.String()
	return
}

// GetVersion returns the version
func (s *server) GetVersion(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: version.GetVersion()}
	return
}

func findParentTestCases(testcase *testing.TestCase, suite *testing.TestSuite) (testcases []testing.TestCase) {
	reg, matchErr := regexp.Compile(`.*\{\{.*\.\w*.*}\}.*`)
	targetReg, targetErr := regexp.Compile(`\.\w*`)

	if matchErr == nil && targetErr == nil {
		expectName := ""
		for _, val := range testcase.Request.Header {
			if matched := reg.MatchString(val); matched {
				expectName = targetReg.FindString(val)
				expectName = strings.TrimPrefix(expectName, ".")
				break
			}
		}

		if expectName == "" {
			if mached := reg.MatchString(testcase.Request.API); mached {
				expectName = targetReg.FindString(testcase.Request.API)
				expectName = strings.TrimPrefix(expectName, ".")
			}
		}

		for _, item := range suite.Items {
			if item.Name == expectName {
				testcases = append(testcases, item)
			}
		}
	}
	return
}
