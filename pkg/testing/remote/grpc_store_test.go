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

		err = writer.Verify()
		assert.Error(t, err)
	})

	t.Run("NewGRPCloaderFromStore", func(t *testing.T) {
		assert.NotNil(t, NewGRPCloaderFromStore())
	})
}
