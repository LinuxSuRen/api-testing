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

package runner

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/h2non/gock"
)

func TestGithubPRCommentWriter(t *testing.T) {
	t.Run("lack of parameters", func(t *testing.T) {
		_, err := NewGithubPRCommentWriter(&GithubPRCommentOption{})
		assert.Error(t, err)

		_, err = NewGithubPRCommentWriter(&GithubPRCommentOption{
			Repo: "repo",
		})
		assert.Error(t, err)

		_, err = NewGithubPRCommentWriter(&GithubPRCommentOption{
			Identity: "id",
			Token:    "token",
		})
		assert.Error(t, err)

		_, err = NewGithubPRCommentWriter(&GithubPRCommentOption{
			Repo:     "repo",
			Identity: "id",
			Token:    "token",
		})
		assert.NoError(t, err)
	})

	t.Run("pr number is invalid", func(t *testing.T) {
		writer, err := NewGithubPRCommentWriter(&GithubPRCommentOption{
			Repo:     "linuxsuren/test",
			Identity: "id",
			Token:    "token",
		})
		assert.NoError(t, err)

		err = writer.Output(nil)
		assert.NoError(t, err)
		assert.Nil(t, writer.WithAPIConverage(nil))
		assert.NotNil(t, writer.WithResourceUsage(nil))
	})

	t.Run("error with getting comments", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.github.com").Get("/repos/linuxsuren/test/issues/1/comments").Reply(http.StatusBadRequest)
		writer := createWriter(t)

		err := writer.Output(nil)
		assert.Error(t, err)
	})

	tmpF, tErr := os.CreateTemp(os.TempDir(), "report")
	if !assert.NoError(t, tErr) {
		return
	}
	defer os.RemoveAll(tmpF.Name())

	t.Run("create new comment", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.github.com").Get("/repos/linuxsuren/test/issues/1/comments").Reply(http.StatusOK).JSON([]comment{})
		gock.New("https://api.github.com").Post("/repos/linuxsuren/test/issues/1/comments").Reply(http.StatusCreated)
		writer := createWriterWithReport(tmpF.Name(), t)

		err := writer.Output([]ReportResult{{
			API: "/api",
		}})
		assert.NoError(t, err)
	})

	t.Run("update comment", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://api.github.com").Get("/repos/linuxsuren/test/issues/1/comments").Reply(http.StatusOK).JSON([]comment{{
			ID:   1234,
			Body: "id",
		}})
		gock.New("https://api.github.com").Patch("/repos/linuxsuren/test/issues/comments/1234").Reply(http.StatusOK)
		writer := createWriterWithReport(tmpF.Name(), t)

		err := writer.Output([]ReportResult{{
			API: "/api",
		}})
		assert.NoError(t, err)
	})
}

func createWriter(t *testing.T) ReportResultWriter {
	return createWriterWithReport("", t)
}

func createWriterWithReport(report string, t *testing.T) ReportResultWriter {
	writer, err := NewGithubPRCommentWriter(&GithubPRCommentOption{
		Repo:       "linuxsuren/test",
		Identity:   "id",
		Token:      "token",
		PR:         1,
		ReportFile: report,
	})
	assert.NoError(t, err)
	return writer
}
