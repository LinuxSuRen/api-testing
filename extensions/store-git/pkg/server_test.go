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
