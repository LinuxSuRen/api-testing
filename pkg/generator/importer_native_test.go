/*
Copyright 2024-2025 API Testing Authors.

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
	"bytes"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	_ "embed"
)

func TestNativeImporter(t *testing.T) {
	importer := NewNativeImporter()

	t.Run("simple native, from []byte", func(t *testing.T) {
		_, err := importer.Convert(simpleNativeData)
		assert.NoError(t, err)
	})

	t.Run("native inline", func(t *testing.T) {
		_, err := importer.Convert([]byte(`name: test
api: https://api.com
spec:
  kind: http
items:
  - name: name
    request:
      api: /octocat`))
		assert.NoError(t, err)
	})

	t.Run("read from file", func(t *testing.T) {
		_, err := importer.ConvertFromFile("testdata/native.json")
		assert.NoError(t, err)
	})

	t.Run("read from URL", func(t *testing.T) {
		defer gock.Off()
		gock.New(urlFoo).Get("/").Reply(http.StatusOK).Body(bytes.NewBuffer(simpleNativeData))

		_, err := importer.ConvertFromURL(urlFoo)
		assert.NoError(t, err)
	})
}

//go:embed testdata/native.json
var simpleNativeData []byte
