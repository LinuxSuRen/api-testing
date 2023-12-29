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

package runner

import (
	"context"
	"net/http"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
)

type graphql struct {
	TestCaseRunner
}

func NewGraphQLRunner(parent TestCaseRunner) TestCaseRunner {
	return &graphql{
		TestCaseRunner: parent,
	}
}

func init() {
	RegisterRunner("graphql", func(*testing.TestSuite) TestCaseRunner {
		return NewGraphQLRunner(NewSimpleTestCaseRunner())
	})
}

func (r *graphql) RunTestCase(testcase *testing.TestCase, dataContext any, ctx context.Context) (
	output any, err error) {
	testcase.Request.Method = http.MethodPost

	if testcase.Request.Header == nil {
		testcase.Request.Header = make(map[string]string, 1)
	}
	testcase.Request.Header[util.ContentType] = util.JSON
	return r.TestCaseRunner.RunTestCase(testcase, dataContext, ctx)
}

func (s *graphql) WithSuite(suite *testing.TestSuite) {
	// not need this parameter
}
