package runner

import (
	"context"
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
