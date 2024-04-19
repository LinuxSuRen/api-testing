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
package generator_test

import (
	"testing"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/generator"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestCodeGeneratorManager(t *testing.T) {
	t.Run("GetCodeGenerators", func(t *testing.T) {
		// should returns a mutable map
		generators := generator.GetCodeGenerators()
		count := len(generators)

		generators["a-new-fake"] = nil
		assert.Equal(t, count, len(generator.GetCodeGenerators()))
	})

	t.Run("GetCodeGenerator", func(t *testing.T) {
		instance := generator.NewGolangGenerator()
		generator.RegisterCodeGenerator("fake", instance)
		assert.Equal(t, instance, generator.GetCodeGenerator("fake"))
	})
}

func TestGenerators(t *testing.T) {
	testcase := &atest.TestCase{
		Request: atest.Request{
			API: "https://www.baidu.com",
			Header: map[string]string{
				"User-Agent": "atest",
			},
		},
	}
	t.Run("golang", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("golang").Generate(nil, testcase)
		assert.NoError(t, err)
		assert.Equal(t, expectedGoCode, result)
	})

	t.Run("java", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("java").Generate(nil, testcase)
		assert.NoError(t, err)
		assert.Equal(t, expectedJavaCode, result)
	})

	formRequest := &atest.TestCase{Request: testcase.Request}
	formRequest.Request.Form = map[string]string{
		"key": "value",
	}
	t.Run("golang form HTTP request", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("golang").Generate(nil, formRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedFormRequestGoCode, result, result)
	})

	t.Run("java form HTTP request", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("java").Generate(nil, formRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedFormRequestJavaCode, result, result)
	})

	cookieRequest := &atest.TestCase{Request: formRequest.Request}
	cookieRequest.Request.Cookie = map[string]string{
		"name": "value",
	}
	t.Run("golang cookie HTTP request", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("golang").Generate(nil, cookieRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedCookieRequestGoCode, result, result)
	})

	t.Run("java cookie HTTP request", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("java").Generate(nil, cookieRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedCookieRequestJavaCode, result, result)
	})
}

//go:embed testdata/expected_go_code.txt
var expectedGoCode string

//go:embed testdata/expected_java_code.txt
var expectedJavaCode string

//go:embed testdata/expected_go_form_request_code.txt
var expectedFormRequestGoCode string

//go:embed testdata/expected_java_form_request_code.txt
var expectedFormRequestJavaCode string

//go:embed testdata/expected_go_cookie_request_code.txt
var expectedCookieRequestGoCode string

//go:embed testdata/expected_java_cookie_request_code.txt
var expectedCookieRequestJavaCode string
