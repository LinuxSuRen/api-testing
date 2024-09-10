/*
Copyright 2023 API Testing Authors.

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
package server

import (
	context "context"
	"io"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/linuxsuren/api-testing/pkg/testing"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	grpc "google.golang.org/grpc"
)

type remoteRunnerAdapter struct {
	suite   *testing.TestSuite
	address string
}

func (s *remoteRunnerAdapter) RunTestCase(testcase *testing.TestCase,
	ctx context.Context) (output interface{}, err error) {

	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(s.address, grpc.WithInsecure()); err != nil {
		return
	}

	suiteWithCase := &TestSuiteWithCase{
		Suite: ToGRPCSuite(s.suite),
		Case:  ToGRPCTestCase(*testcase),
	}

	client := NewRunnerExtensionClient(conn)
	_, err = client.Run(ctx, suiteWithCase)
	return
}

func (s *remoteRunnerAdapter) GetSuggestedAPIs(suite *testing.TestSuite, api string) (
	cases []*testing.TestCase, err error) {
	return
}
func (s *remoteRunnerAdapter) WithSecure(secure *testing.Secure) {
}
func (s *remoteRunnerAdapter) WithOutputWriter(io.Writer) {
}
func (s *remoteRunnerAdapter) WithWriteLevel(level string) {
}
func (s *remoteRunnerAdapter) WithTestReporter(runner.TestReporter) {
}
func (s *remoteRunnerAdapter) WithExecer(fakeruntime.Execer) {
}
func (s *remoteRunnerAdapter) WithSuite(suite *testing.TestSuite) {
	s.suite = suite
}
func (s *remoteRunnerAdapter) WithAPISuggestLimit(limit int) {
}

func init() {
	env := os.Environ()
	runners := make(map[string]string, 0)

	for _, e := range env {
		if strings.HasPrefix(e, "REMOTE_RUNNER_") {
			pair := strings.Split(strings.TrimPrefix(e, "REMOTE_RUNNER_"), "=")
			runners[strings.ToLower(pair[0])] = pair[1]
		}
	}

	for k, v := range runners {
		runner.RegisterRunner(k, func(suite *testing.TestSuite) runner.TestCaseRunner {
			return &remoteRunnerAdapter{address: v}
		})
	}
}
