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
	"fmt"
	"io"
	"strings"

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
	WithSuite(*testing.TestSuite)
	WithAPISuggestLimit(int)
}

// HTTPResponseRecord represents a http response record
type ResponseRecord interface {
	GetResponseRecord() SimpleResponse
}

// SimpleResponse represents a simple response
type SimpleResponse struct {
	Header     map[string]string
	Body       string
	RawBody    []byte
	StatusCode int
}

func (s SimpleResponse) getFileName() string {
	for k, v := range s.Header {
		if k == "Content-Disposition" {
			return strings.TrimSuffix(strings.TrimPrefix(v, `attachment; filename="`), `"`)
		}
	}
	return ""
}

// NewDefaultUnimplementedRunner initializes an unimplementedRunner using the default values.
func NewDefaultUnimplementedRunner() UnimplementedRunner {
	return UnimplementedRunner{
		testReporter: NewDiscardTestReporter(),
		writer:       io.Discard,
		log:          NewDefaultLevelWriter("info", io.Discard),
		execer:       fakeruntime.NewDefaultExecer(),
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

func (r *UnimplementedRunner) GetSuggestedAPIs(suite *testing.TestSuite, api string) (result []*testing.TestCase, err error) {
	// empty implement
	return
}

func (r *UnimplementedRunner) WithAPISuggestLimit(int) {
	// empty implement
}

func (s *UnimplementedRunner) WithSuite(suite *testing.TestSuite) {
	// empty implement
}
