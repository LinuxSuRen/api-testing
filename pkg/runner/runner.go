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
	GetSuggestedAPIs(suite *testing.TestSuite, api string) ([]*testing.TestCase, error)
	WithSecure(secure *testing.Secure)
	WithOutputWriter(io.Writer)
	WithWriteLevel(level string)
	WithTestReporter(TestReporter)
	WithExecer(fakeruntime.Execer)
}

// HTTPResponseRecord represents a http response record
type ResponseRecord interface {
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
	Secure       *testing.Secure
}

func (r *UnimplementedRunner) RunTestCase(testcase *testing.TestCase, dataContext interface{}, ctx context.Context) (output interface{}, err error) {
	return nil, fmt.Errorf("unimplemented")
}

// WithOutputWriter sets the io.Writer
func (r *UnimplementedRunner) WithOutputWriter(writer io.Writer) {
	r.writer = writer
}

// WithWriteLevel sets the level writer
func (r *UnimplementedRunner) WithWriteLevel(level string) {
	if level != "" {
		r.log = NewDefaultLevelWriter(level, r.writer)
	}
}

// WithTestReporter sets the TestReporter
func (r *UnimplementedRunner) WithTestReporter(reporter TestReporter) {
	r.testReporter = reporter
}

// WithExecer sets the execer
func (r *UnimplementedRunner) WithExecer(execer fakeruntime.Execer) {
	r.execer = execer
}

// WithSecure sets the secure option.
func (r *UnimplementedRunner) WithSecure(secure *testing.Secure) {
	r.Secure = secure
}
