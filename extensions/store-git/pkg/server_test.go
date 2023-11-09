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

func TestGetClient(t *testing.T) {
	ctx := remote.WithIncomingStoreContext(context.Background(), &atest.Store{})
	defaultGitClient := &gitClient{}
	t.Run("no context", func(t *testing.T) {
		opt, err := defaultGitClient.getClient(context.TODO())
		assert.Nil(t, opt)
		assert.Error(t, err)

		_, err = defaultGitClient.loadCache(context.TODO())
		assert.Error(t, err)

		err = defaultGitClient.pushCache(context.TODO())
		assert.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		opt, err := defaultGitClient.getClient(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, opt)
		assert.False(t, opt.cloneOptions.InsecureSkipTLS)
	})

	t.Run("verify", func(t *testing.T) {
		gitClient := NewRemoteServer()

		result, err := gitClient.Verify(ctx, &server.Empty{})
		assert.NoError(t, err)
		assert.False(t, result.Ready)
		assert.True(t, result.ReadOnly)
	})

	t.Run("ListTestSuite", func(t *testing.T) {
		_, err := defaultGitClient.ListTestSuite(ctx, &server.Empty{})
		assert.Error(t, err)
	})

	t.Run("CreateTestSuite", func(t *testing.T) {
		_, err := defaultGitClient.CreateTestSuite(ctx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("GetTestSuite", func(t *testing.T) {
		_, err := defaultGitClient.GetTestSuite(ctx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("UpdateTestSuite", func(t *testing.T) {
		_, err := defaultGitClient.UpdateTestSuite(ctx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("DeleteTestSuite", func(t *testing.T) {
		_, err := defaultGitClient.DeleteTestSuite(ctx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("ListTestCases", func(t *testing.T) {
		_, err := defaultGitClient.ListTestCases(ctx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("CreateTestCase", func(t *testing.T) {
		_, err := defaultGitClient.CreateTestCase(ctx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("GetTestCase", func(t *testing.T) {
		_, err := defaultGitClient.GetTestCase(ctx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("UpdateTestCase", func(t *testing.T) {
		_, err := defaultGitClient.UpdateTestCase(ctx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("DeleteTestCase", func(t *testing.T) {
		_, err := defaultGitClient.DeleteTestCase(ctx, &server.TestCase{})
		assert.Error(t, err)
	})
}
