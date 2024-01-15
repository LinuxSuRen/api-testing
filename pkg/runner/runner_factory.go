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
	"fmt"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

// GetTestSuiteRunner returns a proper runner according to the given test suite.
func GetTestSuiteRunner(suite *testing.TestSuite) TestCaseRunner {
	// TODO: should be refactored to meet more types of runners
	kind := suite.Spec.Kind

	switch kind {
	case "swagger", "":
		kind = "http"
	}

	if suite.Spec.RPC != nil && kind == "" {
		kind = "grpc"
	}

	kind = strings.ToLower(kind)
	runner := runners[kind]
	if runner != nil {
		return runner(suite)
	}
	return nil
}

type RunnerCreator func(suite *testing.TestSuite) TestCaseRunner

var runners map[string]RunnerCreator = make(map[string]RunnerCreator, 4)

func RegisterRunner(kind string, runner RunnerCreator) error {
	if _, ok := runners[kind]; ok {
		return fmt.Errorf("duplicated kind %q", kind)
	}

	runners[kind] = runner
	return nil
}
