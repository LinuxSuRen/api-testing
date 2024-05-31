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
package testing_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
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

			readonly, err := loader.Verify()
			assert.NoError(t, err)
			assert.False(t, readonly)
		})
	}
}

func TestURLLoader(t *testing.T) {
	const json = `{"key": "value", "message": "sample"}`

	t.Run("normal HTTP GET request", func(t *testing.T) {
		loader := atest.NewFileLoader()
		defer gock.Off()

		err := loader.Put(urlFake)
		assert.NoError(t, err)

		gock.New(urlFake).Get("/").Reply(http.StatusOK).BodyString(json)

		assert.True(t, loader.HasMore())
		var data []byte
		data, err = loader.Load()
		assert.NoError(t, err)
		assert.Equal(t, json, string(data))
	})

	t.Run("HTTP POST request, lack of suite name", func(t *testing.T) {
		loader := atest.NewFileLoader()
		defer gock.Off()
		const api = "/server.Runner/ConvertTestSuite"
		const reqURL = urlFake + api

		err := loader.Put(reqURL)
		assert.NoError(t, err)

		gock.New(urlFake).Get(api).Reply(http.StatusOK).BodyString(json)

		assert.True(t, loader.HasMore())
		_, err = loader.Load()
		assert.Error(t, err)
	})

	t.Run("HTTP POST request", func(t *testing.T) {
		loader := atest.NewFileLoader()
		defer gock.Off()
		const api = "/server.Runner/ConvertTestSuite"
		const reqURL = urlFake + api + "?suite=sample"

		err := loader.Put(reqURL)
		assert.NoError(t, err)

		gock.New(urlFake).Post(api).BodyString(`{"TestSuite":"sample", "Generator":"raw"}`).
			Reply(http.StatusOK).BodyString(json)

		assert.True(t, loader.HasMore())

		var data []byte
		data, err = loader.Load()
		assert.NoError(t, err)
		assert.Equal(t, "sample", string(data))
	})
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

		fakeName := fmt.Sprintf("fake-%d", time.Now().Nanosecond())
		err = writer.CreateSuite(fakeName, urlFake)
		assert.NoError(t, err)

		err = writer.CreateSuite(fakeName, "")
		assert.Error(t, err)

		assert.Equal(t, 2, writer.GetCount())

		err = writer.DeleteSuite("test")
		assert.NoError(t, err)
		err = writer.DeleteSuite(fakeName)
		assert.NoError(t, err)

		assert.Equal(t, 0, writer.GetCount())

		err = writer.DeleteSuite(fakeName)
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

		var testcase atest.TestCase
		testcase, err = writer.GetTestCase("test", "login")
		if assert.NoError(t, err) {
			assert.Equal(t, urlTestLogin, testcase.Request.API)
		}

		var data []byte
		writer.Put("test")
		assert.True(t, writer.HasMore())
		data, err = writer.Load()
		assert.NoError(t, err)
		writer.Reset()

		var testSuiteYaml []byte
		testSuiteYaml, err = writer.GetTestSuiteYaml("test")
		if assert.NoError(t, err) {
			assert.Equal(t, data, testSuiteYaml)
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
