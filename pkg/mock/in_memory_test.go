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
package mock

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryServer(t *testing.T) {
	server := NewInMemoryServer(0)

	err := server.Start(NewLocalFileReader("testdata/api.yaml"), "/mock")
	assert.NoError(t, err)
	defer func() {
		server.Stop()
	}()

	api := "http://localhost:" + server.GetPort() + "/mock"

	_, err = http.Post(api+"/team", "", bytes.NewBufferString(`{
		"name": "test",
		"members": []
	}`))
	assert.NoError(t, err)

	var resp *http.Response
	resp, err = http.Get(api + "/team")
	if assert.NoError(t, err) {
		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, `[{"name":"someone"},{"members":[],"name":"test"}]`, string(data))
	}

	t.Run("check the /api.json", func(t *testing.T) {
		var resp *http.Response
		resp, err = http.Get(api + "/api.json")
		if assert.NoError(t, err) {
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.NotEmpty(t, string(data))
		}
	})

	t.Run("list with filter", func(t *testing.T) {
		var resp *http.Response
		resp, err = http.Get(api + "/team?name=someone")
		if assert.NoError(t, err) {
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, `[{"name":"someone"}]`, string(data))
		}
	})

	t.Run("update object", func(t *testing.T) {
		updateReq, err := http.NewRequest(http.MethodPut, api+"/team", bytes.NewBufferString(`{
			"name": "test",
			"members": [{
				"name": "rick"
			}]
		}`))
		assert.NoError(t, err)
		resp, err = http.DefaultClient.Do(updateReq)
		assert.NoError(t, err)
	})

	t.Run("get a single object", func(t *testing.T) {
		resp, err = http.Get(api + "/team/test")
		assert.NoError(t, err)

		var data []byte
		data, err = io.ReadAll(resp.Body)
		assert.NoError(t, err)

		assert.Equal(t, `{"members":[{"name":"rick"}],"name":"test"}`, string(data))
	})

	// delete object
	delReq, err := http.NewRequest(http.MethodDelete, api+"/team/test", nil)
	assert.NoError(t, err)
	resp, err = http.DefaultClient.Do(delReq)
	assert.NoError(t, err)

	t.Run("check if deleted", func(t *testing.T) {
		var resp *http.Response
		resp, err = http.Get(api + "/team")
		if assert.NoError(t, err) {
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, `[{"name":"someone"}]`, string(data))
		}

		resp, err = http.Get(api + "/team/test")
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		}
	})

	t.Run("invalid request method", func(t *testing.T) {
		delReq, err := http.NewRequest("fake", api+"/team", nil)
		assert.NoError(t, err)
		resp, err = http.DefaultClient.Do(delReq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("only accept GET method in getting a single object", func(t *testing.T) {
		wrongMethodReq, err := http.NewRequest(http.MethodPut, api+"/team/someone", nil)
		assert.NoError(t, err)
		resp, err = http.DefaultClient.Do(wrongMethodReq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("mock item", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, api+"/v1/repos/test/prs", nil)
		assert.NoError(t, err)
		req.Header.Set("name", "rick")

		resp, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "176", resp.Header.Get(util.ContentLength))
		assert.Equal(t, "mock", resp.Header.Get("Server"))
		assert.NotEmpty(t, resp.Header.Get(headerMockServer))

		data, _ := io.ReadAll(resp.Body)
		assert.True(t, strings.Contains(string(data), `"message": "mock"`), string(data))
	})

	t.Run("miss match header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, api+"/v1/repos/test/prs", nil)
		assert.NoError(t, err)

		resp, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("base64 encoder", func(t *testing.T) {
		resp, err = http.Get(api + "/v1/base64")
		assert.NoError(t, err)
		data, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "hello", string(data))
	})

	t.Run("not found config file", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewLocalFileReader("fake"), "/")
		assert.Error(t, err)
	})

	t.Run("invalid webhook", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewInMemoryReader(`webhooks:
  - timer: aa
    name: fake`), "/")
		assert.Error(t, err)
	})

	t.Run("missing name or timer in webhook", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewInMemoryReader(`webhooks:
  - timer: 1s`), "/")
		assert.Error(t, err)
	})

	t.Run("invalid webhook payload", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewInMemoryReader(`webhooks:
  - name: invalid
    timer: 1ms
    request:
      body: "{{.fake"`), "/")
		assert.Error(t, err)
	})

	t.Run("invalid webhook api template", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewInMemoryReader(`webhooks:
  - name: invalid
    timer: 1ms
    request:
      body: "{}"
      path: "{{.fake"`), "/")
		assert.NoError(t, err)
	})
}
