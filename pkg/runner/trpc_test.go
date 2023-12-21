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
	"testing"

	_ "embed"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestTRPC(t *testing.T) {
	fakeClient := &fakeClient{}
	tRPCRunner := NewTRPCTestCaseRunner("localhost:8080", atest.RPCDesc{
		Raw: sampleProto,
	}, fakeClient)
	assert.NotNil(t, tRPCRunner)

	t.Run("normal", func(t *testing.T) {
		testcase := &atest.TestCase{
			Name: "Unary",
			Request: atest.Request{
				API:  "/Main/Unary",
				Body: atest.NewRequestBody("{}"),
			},
		}

		_, err := tRPCRunner.RunTestCase(testcase, nil, context.Background())
		assert.NoError(t, err)
	})

	t.Run("no case found", func(t *testing.T) {
		_, err := tRPCRunner.RunTestCase(&atest.TestCase{
			Name: "Fake",
			Request: atest.Request{
				API:  "/Main/Fake",
				Body: atest.NewRequestBody("{}"),
			},
		}, nil, context.Background())
		assert.Error(t, err)
	})
}

//go:embed grpc_test/test.proto
var sampleProto string
