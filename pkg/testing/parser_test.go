/*
Copyright 2023-2024 API Testing Authors.

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
package testing_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	_ "embed"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	data, err := os.ReadFile("../../sample/testsuite-gitlab.yaml")
	if !assert.NoError(t, err) {
		return
	}

	suite, err := atest.Parse(data)
	if assert.Nil(t, err) && assert.NotNil(t, suite) {
		assert.Equal(t, "Gitlab", suite.Name)
		assert.Equal(t, 2, len(suite.Items))
		assert.Equal(t, atest.TestCase{
			Name: "projects",
			Request: atest.Request{
				API: "/projects",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				Schema: `{
  "type": "array"
}
`,
			},
			Before: &atest.Job{
				Items: []string{"sleep(1)"},
			},
			After: &atest.Job{
				Items: []string{"sleep(1)"},
			},
		}, suite.Items[0])

		// render suite
		data := map[string]interface{}{}
		err = suite.Render(data)
		assert.NoError(t, err)
		assert.Equal(t, "https://gitlab.com/api/v4", suite.API)
	}

	var anotherSuite *atest.TestSuite
	anotherSuite, err = atest.ParseFromStream(bytes.NewBuffer(data))
	assert.Equal(t, suite, anotherSuite)

	_, err = atest.Parse([]byte(invalidTestCaseContent))
	assert.NotNil(t, err)
}

func TestDuplicatedNames(t *testing.T) {
	data, err := os.ReadFile("testdata/duplicated-names.yaml")
	if !assert.NoError(t, err) {
		return
	}

	_, err = atest.Parse(data)
	assert.NotNil(t, err)

	_, err = atest.ParseFromData([]byte("fake"))
	assert.NotNil(t, err)
}

func TestRequestRender(t *testing.T) {
	validMap := map[string]string{
		"key": "{{.Name}}",
	}
	invalidMap := map[string]string{
		"key": "{{.name}}",
	}

	tests := []struct {
		name    string
		request *atest.Request
		verify  func(t *testing.T, req *atest.Request)
		ctx     interface{}
		hasErr  bool
	}{{
		name: "slice as context",
		request: &atest.Request{
			API:  "http://localhost/{{index . 0}}",
			Body: atest.NewRequestBody("{{index . 1}}"),
		},
		ctx:    []string{"foo", "bar"},
		hasErr: false,
		verify: func(t *testing.T, req *atest.Request) {
			assert.Equal(t, "http://localhost/foo", req.API)
			assert.Equal(t, "bar", req.Body.String())
		},
	}, {
		name:    "default values",
		request: &atest.Request{},
		verify: func(t *testing.T, req *atest.Request) {
			assert.Equal(t, http.MethodGet, req.Method)
		},
		hasErr: false,
	}, {
		name:    "context is nil",
		request: &atest.Request{},
		ctx:     nil,
		hasErr:  false,
	}, {
		name: "body from file",
		request: &atest.Request{
			BodyFromFile: "testdata/generic_body.json",
		},
		ctx: atest.TestCase{
			Name: "linuxsuren",
		},
		hasErr: false,
		verify: func(t *testing.T, req *atest.Request) {
			assert.Equal(t, `{"name": "linuxsuren"}`, req.Body.String())
		},
	}, {
		name: "body file not found",
		request: &atest.Request{
			BodyFromFile: "testdata/fake",
		},
		hasErr: true,
	}, {
		name: "invalid API as template",
		request: &atest.Request{
			API: "{{.name}",
		},
		hasErr: true,
	}, {
		name: "failed with API render",
		request: &atest.Request{
			API: "{{.name}}",
		},
		ctx:    atest.TestCase{},
		hasErr: true,
	}, {
		name: "invalid body as template",
		request: &atest.Request{
			Body: atest.NewRequestBody("{{.name}"),
		},
		hasErr: true,
	}, {
		name: "failed with body render",
		request: &atest.Request{
			Body: atest.NewRequestBody("{{.name}"),
		},
		ctx:    atest.TestCase{},
		hasErr: true,
	}, {
		name: "failed with header render",
		request: &atest.Request{
			Header: map[string]string{
				"key": "{{.name}}",
			},
		},
		ctx:    atest.TestCase{},
		hasErr: true,
	}, {
		name: "failed with form render",
		request: &atest.Request{
			Form: invalidMap,
		},
		ctx:    atest.TestCase{},
		hasErr: true,
	}, {
		name: "form render",
		request: &atest.Request{
			Form: validMap,
		},
		ctx: atest.TestCase{Name: "linuxsuren"},
		verify: func(t *testing.T, req *atest.Request) {
			assert.Equal(t, "linuxsuren", req.Form["key"])
		},
		hasErr: false,
	}, {
		name: "header render",
		request: &atest.Request{
			Header: validMap,
		},
		ctx: atest.TestCase{Name: "linuxsuren"},
		verify: func(t *testing.T, req *atest.Request) {
			assert.Equal(t, "linuxsuren", req.Header["key"])
		},
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Render(tt.ctx, "")
			if assert.Equal(t, tt.hasErr, err != nil, err) && tt.verify != nil {
				tt.verify(t, tt.request)
			}
		})
	}
}

func TestResponseRender(t *testing.T) {
	tests := []struct {
		name     string
		response *atest.Response
		verify   func(t *testing.T, req *atest.Response)
		ctx      interface{}
		hasErr   bool
	}{{
		name:     "blank response",
		response: &atest.Response{},
		verify: func(t *testing.T, req *atest.Response) {
			assert.Equal(t, http.StatusOK, req.StatusCode)
		},
		hasErr: false,
	}, {
		name: "have bodyFieldsExpect",
		ctx: map[string]string{
			"name": "linuxsuren",
		},
		response: &atest.Response{
			BodyFieldsExpect: map[string]interface{}{
				"{{.name}}": "{{.name}}",
			},
		},
		verify: func(t *testing.T, req *atest.Response) {
			assert.Equal(t, "linuxsuren", req.BodyFieldsExpect["linuxsuren"])
		},
		hasErr: false,
	}, {
		name: "bodyFieldsExpect, template syntax error with key",
		response: &atest.Response{
			BodyFieldsExpect: map[string]interface{}{
				"{{.name": "{{.name}",
			},
		},
		hasErr: true,
	}, {
		name: "bodyFieldsExpect, template syntax error with value",
		response: &atest.Response{
			BodyFieldsExpect: map[string]interface{}{
				"name": "{{.name}",
			},
		},
		hasErr: true,
	}, {
		name: "bodyFieldsExpect value is not string",
		response: &atest.Response{
			BodyFieldsExpect: map[string]interface{}{
				"name": []string{},
			},
		},
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.response.Render(tt.ctx)
			if assert.Equal(t, tt.hasErr, err != nil, err) && tt.verify != nil {
				tt.verify(t, tt.response)
			}
		})
	}

	req := &atest.Request{
		API: "/good",
	}
	req.RenderAPI("/api/v1")
	assert.Equal(t, "/api/v1/good", req.API)
}

func TestTestCase(t *testing.T) {
	testCase, err := atest.ParseTestCaseFromData([]byte(testCaseContent))
	assert.Nil(t, err)
	assert.Equal(t, &atest.TestCase{
		Name: "projects",
		Request: atest.Request{
			API: "https://foo",
		},
		Expect: atest.Response{
			StatusCode: http.StatusOK,
		},
	}, testCase)
}

func TestGetBody(t *testing.T) {
	defaultBody := "fake body"

	tests := []struct {
		name        string
		req         *atest.Request
		expectBody  string
		containBody string
		expectErr   bool
	}{{
		name:       "normal body",
		req:        &atest.Request{Body: atest.NewRequestBody(defaultBody)},
		expectBody: defaultBody,
	}, {
		name:       "body from file",
		req:        &atest.Request{BodyFromFile: "testdata/testcase.yaml"},
		expectBody: testCaseContent,
	}, {
		name: "multipart form data",
		req: &atest.Request{
			Header: map[string]string{
				util.ContentType: util.MultiPartFormData,
			},
			Form: map[string]string{
				"key": "value",
			},
		},
		containBody: "name=\"key\"\r\n\r\nvalue\r\n",
	}, {
		name: "normal form",
		req: &atest.Request{
			Header: map[string]string{
				util.ContentType: util.Form,
			},
			Form: map[string]string{
				"name": "linuxsuren",
			},
		},
		expectBody: "name=linuxsuren",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := tt.req.GetBody()
			if tt.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, reader)
				data, err := io.ReadAll(reader)
				assert.Nil(t, err)
				if tt.expectBody != "" {
					assert.Equal(t, tt.expectBody, string(data))
				} else {
					assert.Contains(t, string(data), tt.containBody)
				}
			}
		})
	}
}

//go:embed testdata/testcase.yaml
var testCaseContent string

//go:embed testdata/invalid-testcase.yaml
var invalidTestCaseContent string
