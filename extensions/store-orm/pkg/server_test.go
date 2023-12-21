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
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/stretchr/testify/assert"
)

func TestNewRemoteServer(t *testing.T) {
	remoteServer := NewRemoteServer()
	assert.NotNil(t, remoteServer)
	defaultCtx := context.Background()

	t.Run("ListTestSuite", func(t *testing.T) {
		_, err := remoteServer.ListTestSuite(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("CreateTestSuite", func(t *testing.T) {
		_, err := remoteServer.CreateTestSuite(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("GetTestSuite", func(t *testing.T) {
		_, err := remoteServer.GetTestSuite(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("UpdateTestSuite", func(t *testing.T) {
		_, err := remoteServer.UpdateTestSuite(defaultCtx, &remote.TestSuite{})
		assert.Error(t, err)
	})

	t.Run("DeleteTestSuite", func(t *testing.T) {
		_, err := remoteServer.DeleteTestSuite(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("ListTestCases", func(t *testing.T) {
		_, err := remoteServer.ListTestCases(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("CreateTestCase", func(t *testing.T) {
		_, err := remoteServer.CreateTestCase(defaultCtx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("GetTestCase", func(t *testing.T) {
		_, err := remoteServer.GetTestCase(defaultCtx, nil)
		assert.Error(t, err)
	})

	t.Run("UpdateTestCase", func(t *testing.T) {
		_, err := remoteServer.UpdateTestCase(defaultCtx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("DeleteTestCase", func(t *testing.T) {
		_, err := remoteServer.DeleteTestCase(defaultCtx, &server.TestCase{})
		assert.Error(t, err)
	})

	t.Run("Verify", func(t *testing.T) {
		reply, err := remoteServer.Verify(defaultCtx, nil)
		assert.NoError(t, err)
		assert.False(t, reply.Ready)
	})
}
