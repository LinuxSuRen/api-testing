// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"fmt"
	"os"
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
			suite.Items = []testing.TestCase{*targetTestcase}
		} else {
			err = fmt.Errorf("cannot found testcase %s", task.CaseName)
			return
		}
	default:
		err = fmt.Errorf("not support '%s'", task.Kind)
		return
	}

	fmt.Println("prepare to run:", suite.Name)
	dataContext := map[string]interface{}{}

	var result string
	if result, err = render.Render("base api", suite.API, dataContext); err == nil {
		suite.API = result
		suite.API = strings.TrimSuffix(suite.API, "/")
	} else {
		return
	}

	buf := new(bytes.Buffer)

	for _, testCase := range suite.Items {
		simpleRunner := runner.NewSimpleTestCaseRunner()
		simpleRunner.WithOutputWriter(buf)

		// reuse the API prefix
		if strings.HasPrefix(testCase.Request.API, "/") {
			testCase.Request.API = fmt.Sprintf("%s%s", suite.API, testCase.Request.API)
		}

		var output interface{}
		if output, err = simpleRunner.RunTestCase(&testCase, dataContext, ctx); err == nil {
			dataContext[testCase.Name] = output
		} else {
			break
		}
	}
	reply = &HelloReply{Message: buf.String()}
	return
}

// GetVersion returns the version
func (s *server) GetVersion(ctx context.Context, in *Empty) (reply *HelloReply, err error) {
	reply = &HelloReply{Message: version.GetVersion()}
	return
}
