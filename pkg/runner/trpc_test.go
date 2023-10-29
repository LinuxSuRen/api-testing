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
				Body: "{}",
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
				Body: "{}",
			},
		}, nil, context.Background())
		assert.Error(t, err)
	})
}

//go:embed grpc_test/test.proto
var sampleProto string
