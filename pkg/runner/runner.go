package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/linuxsuren/api-testing/pkg/testing"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

// TestCaseRunner represents a test case runner
type TestCaseRunner interface {
	RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error)
	WithOutputWriter(io.Writer) TestCaseRunner
	WithWriteLevel(level string) TestCaseRunner
	WithTestReporter(TestReporter) TestCaseRunner
	WithExecer(fakeruntime.Execer) TestCaseRunner
}

// HTTPResponseRecord represents a http response record
type HTTPResponseRecord interface {
	GetResponseRecord() SimpleResponse
}

// SimpleResponse represents a simple response
type SimpleResponse struct {
	Header     map[string]string
	Body       string
	StatusCode int
}

// NewDefaultUnimplementedRunner initializes an unimplementedRunner using the default values.
func NewDefaultUnimplementedRunner() UnimplementedRunner {
	return UnimplementedRunner{
		testReporter: NewDiscardTestReporter(),
		writer:       io.Discard,
		log:          NewDefaultLevelWriter("info", io.Discard),
		execer:       fakeruntime.DefaultExecer{},
	}
}

// UnimplementedRunner implements interface TestCaseRunner except method RunTestCase.
//
// Generally, this struct can be inherited directly when implementing a new runner.
// It is recommended to use NewDefaultUnimplementedRunner to initalize rather than
// to fill it manully.
type UnimplementedRunner struct {
	testReporter TestReporter
	writer       io.Writer
	log          LevelWriter
	execer       fakeruntime.Execer
}

func (r *UnimplementedRunner) RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// WithOutputWriter sets the io.Writer
func (r *UnimplementedRunner) WithOutputWriter(writer io.Writer) TestCaseRunner {
	r.writer = writer
	return r
}

// WithWriteLevel sets the level writer
func (r *UnimplementedRunner) WithWriteLevel(level string) TestCaseRunner {
	if level != "" {
		r.log = NewDefaultLevelWriter(level, r.writer)
	}
	return r
}

// WithTestReporter sets the TestReporter
func (r *UnimplementedRunner) WithTestReporter(reporter TestReporter) TestCaseRunner {
	r.testReporter = reporter
	return r
}

// WithExecer sets the execer
func (r *UnimplementedRunner) WithExecer(execer fakeruntime.Execer) TestCaseRunner {
	r.execer = execer
	return r
}
