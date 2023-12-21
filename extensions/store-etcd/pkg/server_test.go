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
	remoteServer := NewRemoteServer(&fakeKV{
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
		assert.True(t, verifyResult.Ready)
	}
}

func createTestSuite(ctx context.Context, t *testing.T, server remote.LoaderServer, name string) {
	_, err := server.CreateTestSuite(ctx, &remote.TestSuite{
		Name: name,
	})
	assert.NoError(t, err)
}
