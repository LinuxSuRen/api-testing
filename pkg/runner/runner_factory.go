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
	"log"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"trpc.group/trpc-go/trpc-go/client"
)

// GetTestSuiteRunner returns a proper runner according to the given test suite.
func GetTestSuiteRunner(suite *testing.TestSuite) TestCaseRunner {
	// TODO: should be refactored to meet more types of runners
	kind := suite.Spec.Kind

	if suite.Spec.RPC != nil {
		switch strings.ToLower(kind) {
		case "", "grpc":
			return NewGRPCTestCaseRunner(suite.API, *suite.Spec.RPC)
		case "trpc":
			return NewTRPCTestCaseRunner(suite.API, *suite.Spec.RPC, client.New())
		default:
			log.Println("unknown test suite, try to use HTTP runner")
		}
	} else {
		switch strings.ToLower(kind) {
		case "graphql":
			return NewGraphQLRunner(NewSimpleTestCaseRunner())
		}
	}

	return NewSimpleTestCaseRunner()
}
