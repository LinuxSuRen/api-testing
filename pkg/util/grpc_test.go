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
package util

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestLoadProtoFiles(t *testing.T) {
	t.Run("plain string proto file", func(t *testing.T) {
		targetProtoFile, importPath, _, _ := LoadProtoFiles("test.proto")
		assert.Equal(t, "test.proto", targetProtoFile)
		assert.Empty(t, importPath)
	})

	t.Run("URL with invalid status code", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://localhost").Get("/test.proto").Reply(http.StatusNotFound)

		_, _, _, err := LoadProtoFiles("http://localhost/test.proto")
		assert.Error(t, err)
	})

	t.Run("single file URL", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://localhost").Get("/test.proto").
			MatchParam("rand", "123").
			Reply(http.StatusOK).
			File("testdata/test.proto")

		targetProtoFile, importPath, _, err := LoadProtoFiles("http://localhost/test.proto?rand=123")
		defer os.Remove(targetProtoFile)

		assert.True(t, strings.HasPrefix(targetProtoFile, os.TempDir()), targetProtoFile)
		assert.Empty(t, importPath)
		assert.NoError(t, err)
	})

	t.Run("URL with zip file, the query is missing", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://localhost").Get("/test.proto").
			MatchParam("rand", "234").
			Reply(http.StatusOK).
			AddHeader(ContentType, ZIP)

		_, _, _, err := LoadProtoFiles("http://localhost/test.proto?rand=234")
		assert.Error(t, err)
	})

	t.Run("URL with zip file", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://localhost").Get("/test.proto").
			MatchParam("file", "testdata/report.html").
			Reply(http.StatusOK).
			AddHeader(ContentType, ZIP).
			AddHeader(ContentDisposition, "attachment; filename=test.zip").
			File("testdata/test.zip")

		_, _, targetProtoFileDir, err := LoadProtoFiles("http://localhost/test.proto?file=testdata/report.html")
		defer os.RemoveAll(targetProtoFileDir)

		assert.True(t, strings.HasPrefix(targetProtoFileDir, os.TempDir()))
		assert.NoError(t, err)
	})

	assert.Error(t, extractFiles("a", "", ""))
	assert.Error(t, extractFiles("", "b", ""))
}
