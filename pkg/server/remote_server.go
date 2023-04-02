// Package server provides a GRPC based server
package server

import (
	"bytes"
	context "context"
	"fmt"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
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
	default:
		err = fmt.Errorf("not support '%s'", task.Kind)
		return
	}

	dataContext := map[string]interface{}{}
	buf := new(bytes.Buffer)

	for _, testCase := range suite.Items {
		simpleRunner := runner.NewSimpleTestCaseRunner()
		simpleRunner.WithOutputWriter(buf)

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
