package testing_test

import (
	"net/http"
	"os"
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestFileLoader(t *testing.T) {
	tests := []struct {
		name   string
		items  []string
		verify func(t *testing.T, loader atest.Loader)
	}{{
		name:  "empty",
		items: []string{},
		verify: func(t *testing.T, loader atest.Loader) {
			assert.False(t, loader.HasMore())
			assert.Equal(t, 0, loader.GetCount())
		},
	}, {
		name:   "brace expansion path",
		items:  []string{"testdata/{invalid-,}testcase.yaml"},
		verify: defaultVerify,
	}, {
		name:   "glob path",
		items:  []string{"testdata/*testcase.yaml"},
		verify: defaultVerify,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := atest.NewFileLoader()
			for _, item := range tt.items {
				loader.Put(item)
			}
			tt.verify(t, loader)

			assert.NoError(t, loader.Verify())
		})
	}
}

func defaultVerify(t *testing.T, loader atest.Loader) {
	assert.True(t, loader.HasMore())
	data, err := loader.Load()
	assert.Nil(t, err)
	assert.Equal(t, invalidTestCaseContent, string(data))
	assert.Equal(t, "testdata", loader.GetContext())

	assert.True(t, loader.HasMore())
	data, err = loader.Load()
	assert.Nil(t, err)
	assert.Equal(t, testCaseContent, string(data))
	assert.Equal(t, "testdata", loader.GetContext())

	assert.False(t, loader.HasMore())
	loader.Reset()
	assert.True(t, loader.HasMore())
}

func TestSuite(t *testing.T) {
	t.Run("create suite", func(t *testing.T) {
		writer := atest.NewFileWriter(os.TempDir())
		err := writer.CreateSuite("test", urlTest)
		assert.NoError(t, err)

		err = writer.UpdateSuite(atest.TestSuite{
			Name: "test",
			API:  urlFake,
		})
		assert.NoError(t, err)

		var suite *atest.TestSuite
		var absPath string
		suite, absPath, err = writer.GetSuite("test")
		assert.NoError(t, err)
		assert.NotEmpty(t, absPath)
		assert.Equal(t, urlFake, suite.API)

		err = writer.CreateSuite("fake", urlFake)
		assert.NoError(t, err)

		err = writer.CreateSuite("fake", "")
		assert.Error(t, err)

		assert.Equal(t, 2, writer.GetCount())

		err = writer.DeleteSuite("test")
		assert.NoError(t, err)
		err = writer.DeleteSuite("fake")
		assert.NoError(t, err)

		assert.Equal(t, 0, writer.GetCount())

		err = writer.DeleteSuite("fake")
		assert.Error(t, err)
	})

	t.Run("create case", func(t *testing.T) {
		writer := atest.NewFileWriter(os.TempDir())
		err := writer.CreateSuite("test", urlTest)
		assert.NoError(t, err)

		err = writer.CreateTestCase("test", atest.TestCase{
			Name: "login",
			Request: atest.Request{
				API: urlTestLogin,
			},
		})
		assert.NoError(t, err)

		var suite atest.TestSuite
		suite, err = writer.GetTestSuite("test", false)
		if assert.NoError(t, err) {
			assert.Equal(t, "test", suite.Name)
			assert.Equal(t, urlTest, suite.API)
		}

		err = writer.UpdateTestCase("test", atest.TestCase{
			Name: "login",
			Request: atest.Request{
				API:    urlTestLogin,
				Method: http.MethodPost,
			},
		})
		assert.NoError(t, err)

		var testcase atesting.TestCase
		testcase, err = writer.GetTestCase("test", "login")
		if assert.NoError(t, err) {
			assert.Equal(t, urlTestLogin, testcase.Request.API)
		}

		err = writer.DeleteTestCase("test", "login")
		assert.NoError(t, err)

		err = writer.DeleteTestCase("test", "login")
		assert.Error(t, err)

		err = writer.DeleteSuite("test")
		assert.NoError(t, err)
	})
}

const urlFake = "http://fake"
const urlTest = "http://test"
const urlTestLogin = "http://test/login"
