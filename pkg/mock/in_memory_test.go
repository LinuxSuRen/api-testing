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
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestInMemoryServer(t *testing.T) {
	server := NewInMemoryServer(0)

	err := server.Start(NewLocalFileReader("data/api.yaml"))
	assert.NoError(t, err)
	defer func() {
		server.Stop()
	}()

	api := "http://localhost:" + server.GetPort()

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
	delReq, err := http.NewRequest(http.MethodDelete, api+"/team", bytes.NewBufferString(`{
		"name": "test",
		"members": []
	}`))
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
	})

	t.Run("invalid request method", func(t *testing.T) {
		delReq, err := http.NewRequest("fake", api+"/team", nil)
		assert.NoError(t, err)
		resp, err = http.DefaultClient.Do(delReq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("only accept GET method in getting a single object", func(t *testing.T) {
		wrongMethodReq, err := http.NewRequest(http.MethodPut, api+"/team/test", nil)
		assert.NoError(t, err)
		resp, err = http.DefaultClient.Do(wrongMethodReq)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("not found config file", func(t *testing.T) {
		server := NewInMemoryServer(0)
		err := server.Start(NewLocalFileReader("fake"))
		assert.Error(t, err)
	})
}
