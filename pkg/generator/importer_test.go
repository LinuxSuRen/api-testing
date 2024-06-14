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

	t.Run("sub postman, from file", func(t *testing.T) {
		suite, err := importer.ConvertFromFile("testdata/postman-sub.json")
		assert.NoError(t, err)

		var result string
		result, err = converter.Convert(suite)
		assert.NoError(t, err)
		assert.Equal(t, expectedSuiteFromSubPostman, strings.TrimSpace(result), result)
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

//go:embed testdata/expected_suite_from_sub_postman.yaml
var expectedSuiteFromSubPostman string
