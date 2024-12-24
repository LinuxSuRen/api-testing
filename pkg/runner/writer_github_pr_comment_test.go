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
		assert.Nil(t, writer.WithAPICoverage(nil))
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
