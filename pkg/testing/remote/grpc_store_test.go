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

package remote

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestNewGRPCLoader(t *testing.T) {
	factory := NewGRPCloaderFromStore()

	t.Run("invalid address", func(t *testing.T) {
		writer, err := factory.NewInstance(atest.Store{})
		assert.Error(t, err)
		assert.Nil(t, writer)
	})

	t.Run("valid address", func(t *testing.T) {
		writer, err := factory.NewInstance(atest.Store{
			Kind: atest.StoreKind{
				URL: "localhost:7070",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, writer)

		assert.False(t, writer.HasMore())
		var data []byte
		data, err = writer.Load()
		assert.NoError(t, err)
		assert.Nil(t, data)

		assert.NoError(t, writer.Put("foo"))

		assert.Empty(t, writer.GetContext())

		assert.Equal(t, 0, writer.GetCount())
		writer.Reset()

		_, err = writer.ListTestCase("")
		assert.Error(t, err)

		_, err = writer.GetTestCase("", "")
		assert.Error(t, err)

		err = writer.CreateTestCase("", atest.TestCase{})
		assert.Error(t, err)

		err = writer.UpdateTestCase("", atest.TestCase{})
		assert.Error(t, err)

		err = writer.DeleteTestCase("", "")
		assert.Error(t, err)

		_, err = writer.ListTestSuite()
		assert.Error(t, err)

		_, err = writer.GetTestSuite("", false)
		assert.Error(t, err)

		err = writer.CreateSuite("", "")
		assert.Error(t, err)

		_, _, err = writer.GetSuite("")
		assert.Error(t, err)

		err = writer.UpdateSuite(atest.TestSuite{})
		assert.Error(t, err)

		err = writer.DeleteSuite("")
		assert.Error(t, err)

		_, err = writer.ListHistoryTestSuite()
		assert.Error(t, err)

		err = writer.CreateHistoryTestCase(atest.TestCaseResult{}, &atest.TestSuite{}, map[string]string{})
		assert.Error(t, err)

		_, err = writer.GetHistoryTestCase("")
		assert.Error(t, err)

		_, err = writer.GetHistoryTestCaseWithResult("")
		assert.Error(t, err)

		_, err = writer.GetTestCaseAllHistory("", "")
		assert.Error(t, err)

		err = writer.DeleteHistoryTestCase("")
		assert.Error(t, err)

		err = writer.DeleteAllHistoryTestCase("", "")
		assert.Error(t, err)

		var readonly bool
		readonly, err = writer.Verify()
		assert.Error(t, err)
		assert.False(t, readonly)
	})

	t.Run("NewGRPCloaderFromStore", func(t *testing.T) {
		assert.NotNil(t, NewGRPCloaderFromStore())
	})
}
