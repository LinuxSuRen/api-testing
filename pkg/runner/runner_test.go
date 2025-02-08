/*
Copyright 2023-2025 API Testing Authors.

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

    runner = GetTestSuiteRunner(&atest.TestSuite{Spec: atest.APISpec{Kind: "grpc", RPC: &atest.RPCDesc{}}})
    assert.IsType(t, NewGRPCTestCaseRunner("", atest.RPCDesc{}), runner)
}

func TestUnimplementedRunner(t *testing.T) {
    runner := NewDefaultUnimplementedRunner()
    output, err := runner.RunTestCase(&atest.TestCase{}, nil, nil)
    assert.Nil(t, output)
    assert.Error(t, err)

    runner.WithWriteLevel("debug")
    runner.WithTestReporter(nil)

    var results []*atest.TestCase
    results, err = runner.GetSuggestedAPIs(nil, "")
    assert.Nil(t, results)
    assert.NoError(t, err)

    runner.WithAPISuggestLimit(0)
}

func TestSimpleResponse(t *testing.T) {
    t.Run("get fileName", func(t *testing.T) {
        // without filename
        assert.Empty(t, SimpleResponse{}.getFileName())

        // normal case
        assert.Equal(t, "a.txt", SimpleResponse{
            Header: map[string]string{
                "Content-Disposition": `attachment; filename="a.txt"`,
            },
        }.getFileName())

        // without space
        assert.Equal(t, "a.txt", SimpleResponse{
            Header: map[string]string{
                "Content-Disposition": `attachment;filename="a.txt"`,
            },
        }.getFileName())

        // without quote
        assert.Equal(t, "a.txt", SimpleResponse{
            Header: map[string]string{
                "Content-Disposition": `attachment; filename=a.txt`,
            },
        }.getFileName())

        // without quote and space
        assert.Equal(t, "a.txt", SimpleResponse{
            Header: map[string]string{
                "Content-Disposition": `attachment;filename=a.txt`,
            },
        }.getFileName())
    })
}
