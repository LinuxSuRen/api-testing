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

package generator

import (
	"net/http"
	"strings"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestPostmanImport(t *testing.T) {
	importer := NewPostmanImporter()

	converter := GetTestSuiteConverter("raw")
	if !assert.NotNil(t, converter) {
		return
	}

	t.Run("empty", func(t *testing.T) {
		suite, err := importer.Convert([]byte(emptyJSON))
		assert.NoError(t, err)

		var result string
		result, err = converter.Convert(suite)
		assert.NoError(t, err)
		assert.Equal(t, emptyJSON, strings.TrimSpace(result))
	})

	t.Run("simple postman, from []byte", func(t *testing.T) {
		suite, err := importer.Convert([]byte(simplePostman))
		assert.NoError(t, err)

		var result string
		result, err = converter.Convert(suite)
		assert.NoError(t, err)
		assert.Equal(t, expectedSuiteFromPostman, strings.TrimSpace(result), result)
	})

	t.Run("simple postman, from file", func(t *testing.T) {
		suite, err := importer.ConvertFromFile("testdata/postman.json")
		assert.NoError(t, err)

		var result string
		result, err = converter.Convert(suite)
		assert.NoError(t, err)
		assert.Equal(t, expectedSuiteFromPostman, strings.TrimSpace(result), result)
	})

	t.Run("simple postman, from URl", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Get("/").Reply(http.StatusOK).BodyString(simplePostman)

		suite, err := importer.ConvertFromURL(urlFoo)
		assert.NoError(t, err)

		var result string
		result, err = converter.Convert(suite)
		assert.NoError(t, err)
		assert.Equal(t, expectedSuiteFromPostman, strings.TrimSpace(result), result)
	})

	t.Run("nil data", func(t *testing.T) {
		_, err := importer.Convert(nil)
		assert.Error(t, err)
	})

	t.Run("pairs toMap", func(t *testing.T) {
		pairs := Paris{}
		assert.Equal(t, 0, len(pairs.ToMap()))
	})
}

const emptyJSON = "{}"
const urlFoo = "http://foo"

//go:embed testdata/postman.json
var simplePostman string

//go:embed testdata/expected_suite_from_postman.yaml
var expectedSuiteFromPostman string
