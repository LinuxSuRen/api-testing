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
