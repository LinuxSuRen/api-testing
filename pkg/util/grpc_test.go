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
