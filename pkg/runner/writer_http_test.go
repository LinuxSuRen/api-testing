/*
Copyright 2024 API Testing Authors.

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
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestHTTPResultWriter(t *testing.T) {
	t.Run("test get request", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://test.com").Get("/result/get").Reply(http.StatusOK).JSON([]comment{})

		writer := NewHTTPResultWriter("GET", "https://test.com/result/get", nil, nil)

		err := writer.Output([]ReportResult{{
			API: "/api",
		}})
		assert.NoError(t, err)
	})

	t.Run("test post request", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://test.com").Post("/result/post").Reply(http.StatusOK).JSON([]comment{})

		writer := NewHTTPResultWriter("POST", "https://test.com/result/post", nil, nil)

		err := writer.Output([]ReportResult{{
			API: "/api",
		}})
		assert.NoError(t, err)
	})

	t.Run("test parameters", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://test.com/result/post?username=1&pwd=2").Post("").Reply(http.StatusOK).JSON([]comment{})

		parameters := map[string]string{
			"username": "1",
			"pwd":      "2",
		}

		writer := NewHTTPResultWriter("POST", "https://test.com/result/post", parameters, nil)

		err := writer.Output([]ReportResult{{
			API: "/api",
		}})
		assert.NoError(t, err)
	})

	t.Run("test user does not send template file", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://test.com/result/post?username=1&pwd=2").Post("").Reply(http.StatusOK).JSON([]comment{})

		parameters := map[string]string{
			"username": "1",
			"pwd":      "2",
		}

		writer := NewHTTPResultWriter("POST", "https://test.com/result/post", parameters, nil)

		err := writer.Output([]ReportResult{{
			Name:  "test",
			API:   "/api",
			Count: 1,
		}})
		assert.NoError(t, err)
	})

	t.Run("test user send template file", func(t *testing.T) {
		defer gock.Off()
		gock.New("https://test.com/result/post?username=1&pwd=2").Post("").Reply(http.StatusOK).JSON([]comment{})

		parameters := map[string]string{
			"username": "1",
			"pwd":      "2",
		}
		templateOption := NewTemplateOption("./writer_templates/example.tpl", "json")
		writer := NewHTTPResultWriter("POST", "https://test.com/result/post", parameters, templateOption)

		err := writer.Output([]ReportResult{{
			Name:  "test",
			API:   "/api",
			Count: 1,
		}})
		assert.NoError(t, err)
	})
}
