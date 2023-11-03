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

	formRequest := &atest.TestCase{Request: testcase.Request}
	formRequest.Request.Form = map[string]string{
		"key": "value",
	}
	t.Run("golang form HTTP request", func(t *testing.T) {
		result, err := generator.GetCodeGenerator("golang").Generate(nil, formRequest)
		assert.NoError(t, err)
		assert.Equal(t, expectedFormRequestGoCode, result, result)
	})
}

//go:embed testdata/expected_go_code.txt
var expectedGoCode string

//go:embed testdata/expected_go_form_request_code.txt
var expectedFormRequestGoCode string
