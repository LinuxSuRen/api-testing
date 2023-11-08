/**
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

package pkg

import (
	"context"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	atest "github.com/linuxsuren/api-testing/pkg/testing"

	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/stretchr/testify/assert"
)

func TestRemoteServer(t *testing.T) {
	ctx := remote.WithIncomingStoreContext(context.Background(), &atest.Store{
		Name: "test",
	})
	remoteServer := NewRemoteServer("endpoint", &fakeKV{
		data: map[string]string{},
	})

	t.Run("no context found", func(t *testing.T) {
		_, err := remoteServer.ListTestSuite(context.TODO(), nil)
		assert.Error(t, err)
	})

	t.Run("test suite list is empty", func(t *testing.T) {
		_, err := remoteServer.ListTestSuite(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("test suite", func(t *testing.T) {
		const name = "fake"
		createTestSuite(ctx, t, remoteServer, name)

		suites, err := remoteServer.ListTestSuite(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(suites.Data))

		_, err = remoteServer.UpdateTestSuite(ctx, &remote.TestSuite{
			Name: name,
			Api:  "http://fake.com",
		})
		assert.NoError(t, err)

		var suite *remote.TestSuite
		suite, err = remoteServer.GetTestSuite(ctx, &remote.TestSuite{Name: name})
		if assert.NoError(t, err) {
			assert.Equal(t, "http://fake.com", suite.Api)
		}

		_, err = remoteServer.DeleteTestSuite(ctx, &remote.TestSuite{Name: name})
		assert.NoError(t, err)

		// should not found
		_, err = remoteServer.GetTestSuite(ctx, &remote.TestSuite{Name: name})
		assert.Error(t, err)
	})

	t.Run("test case", func(t *testing.T) {
		const name = "foo"
		createTestSuite(ctx, t, remoteServer, name)

		_, err := remoteServer.CreateTestCase(ctx, &server.TestCase{
			SuiteName: name,
			Name:      "bar",
			Request: &server.Request{
				Api: "/foo",
			},
		})
		assert.NoError(t, err)

		var testcases *server.TestCases
		testcases, err = remoteServer.ListTestCases(ctx, &remote.TestSuite{
			Name: name,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(testcases.Data))
		assert.Equal(t, "bar", testcases.Data[0].Name)

		_, err = remoteServer.UpdateTestCase(ctx, &server.TestCase{
			SuiteName: name,
			Name:      "bar",
			Request: &server.Request{
				Api: "/bar",
			},
		})

		var testcase *server.TestCase
		testcase, err = remoteServer.GetTestCase(ctx, &server.TestCase{
			SuiteName: name,
			Name:      "bar",
		})
		if assert.NoError(t, err) {
			assert.Equal(t, "/bar", testcase.Request.Api)
		}

		_, err = remoteServer.DeleteTestCase(ctx, &server.TestCase{
			SuiteName: name,
			Name:      "bar",
		})
		assert.NoError(t, err)
	})

	verifyResult, err := remoteServer.Verify(ctx, &server.Empty{})
	if assert.NoError(t, err) {
		assert.True(t, verifyResult.Success)
	}
}

func createTestSuite(ctx context.Context, t *testing.T, server remote.LoaderServer, name string) {
	_, err := server.CreateTestSuite(ctx, &remote.TestSuite{
		Name: name,
	})
	assert.NoError(t, err)
}
