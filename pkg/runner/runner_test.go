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
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestRunnerFactory(t *testing.T) {
	runner := GetTestSuiteRunner(&atest.TestSuite{})
	assert.IsType(t, NewSimpleTestCaseRunner(), runner)

	runner = GetTestSuiteRunner(&atest.TestSuite{Spec: atest.APISpec{RPC: &atest.RPCDesc{}}})
	assert.IsType(t, NewGRPCTestCaseRunner("", atest.RPCDesc{}), runner)
}
