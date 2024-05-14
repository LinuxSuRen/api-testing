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
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestTestCase(t *testing.T) {
	fooRequest := atest.Request{
		API: urlFoo,
	}
	defaultForm := map[string]string{
		"key": "value",
	}
	const defaultBody = `{"name":"hello"}`
	defaultPrepare := func() {
		gock.New(urlLocalhost).
			Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`).
			SetHeader(util.ContentType, util.JSON)
	}
	defaultPostPrepare := func() {
		gock.New(urlLocalhost).
			Post("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
	}

	tests := []struct {
		name     string
		execer   fakeruntime.Execer
		testCase *atest.TestCase
		ctx      interface{}
		prepare  func()
		verify   func(t *testing.T, output interface{}, err error)
	}{{
		name: "failed during the prepare stage",
		testCase: &atest.TestCase{
			Before: &atest.Job{
				Items: []string{"demo.yaml"},
			},
		},
		execer: fakeruntime.FakeExecer{ExpectError: errors.New("fake")},
	}, {
		name: "normal, response is map",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    urlFoo,
				Header: defaultForm,
				Query: map[string]interface{}{
					"name": "linuxsuren",
				},
				Cookie: map[string]string{
					"key": "value",
				},
				Body: atest.NewRequestBody(`{"foo":"bar"}`),
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				BodyFieldsExpect: map[string]interface{}{
					"name":   "linuxsuren",
					"number": 1,
				},
				Header: map[string]string{
					"type": "generic",
				},
				Verify: []string{
					`data.name == "linuxsuren"`,
				},
			},
			Before: &atest.Job{
				Items: []string{"sleep(1)"},
			},
		},
		execer: fakeruntime.FakeExecer{},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").
				MatchHeader("key", "value").
				Reply(http.StatusOK).
				SetHeader("type", "generic").
				SetHeader(util.ContentType, util.JSON).
				File("testdata/generic_response.json")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.Nil(t, err)
			assert.Equal(t, map[string]interface{}{"name": "linuxsuren", "number": float64(1)}, output)
		},
	}, {
		name: "normal, response is slice",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				Body:       `["foo", "bar"]`,
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").
				Reply(http.StatusOK).
				SetHeader(util.ContentType, util.JSON).
				BodyString(`["foo", "bar"]`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.Nil(t, err)
			assert.Equal(t, []interface{}{"foo", "bar"}, output)
		},
	}, {
		name: "normal, response from file",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:          urlFoo,
				Method:       http.MethodPost,
				BodyFromFile: "testdata/generic_response.json",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Post("/foo").BodyString(genericBody).
				Reply(http.StatusOK).BodyString("123").
				SetHeader(util.ContentType, util.JSON)
		},
	}, {
		name: "response from a not found file",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:          urlFoo,
				Method:       http.MethodPost,
				BodyFromFile: "testdata/fake.json",
			},
		},
	}, {
		name: "bad request",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").Reply(http.StatusBadRequest)
		},
	}, {
		name: "error with request",
		testCase: &atest.TestCase{
			Request: fooRequest,
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").ReplyError(errors.New("error"))
		},
	}, {
		name: "not match with body",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				Body: "bar",
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").Reply(http.StatusOK).
				SetHeader(util.ContentType, util.Plain).BodyString("foo")
		},
	}, {
		name: "not match with header",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				Header: map[string]string{
					"foo": "bar",
				},
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Get("/foo").Reply(http.StatusOK).SetHeader("foo", "value")
		},
	}, {
		name: "not found from fields",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		prepare: prepareForFoo,
	}, {
		name: "body filed not match",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"name": "bar",
				},
			},
		},
		prepare: prepareForFoo,
	}, {
		name: "invalid filed finding",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"0.items": "bar",
				},
			},
		},
		prepare: defaultPrepare,
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "not found field")
		},
	}, {
		name: "verify failed",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: urlFoo,
			},
			Expect: atest.Response{
				Verify: []string{
					"len(data.items) > 0",
				},
			},
		},
		prepare: defaultPrepare,
		verify: func(t *testing.T, output interface{}, err error) {
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), "failed to verify")
			}
		},
	}, {
		name: "failed to compile",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				Verify: []string{
					`println("12")`,
				},
			},
		},
		prepare: defaultPrepare,
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "unknown name println")
		},
	}, {
		name: "failed to compile",
		testCase: &atest.TestCase{
			Request: fooRequest,
			Expect: atest.Response{
				Verify: []string{
					`1 + 1`,
				},
			},
		},
		prepare: defaultPrepare,
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "expected bool, but got int")
		},
	}, {
		name: "wrong API format",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    "ssh://localhost/foo",
				Method: "fake,fake",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "invalid method")
		},
	}, {
		name: "failed to render API",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo/{{.abc}",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "template: api:1:")
		},
	}, {
		name: "multipart form request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    urlFoo,
				Method: http.MethodPost,
				Header: map[string]string{
					util.ContentType: "multipart/form-data",
				},
				Form: defaultForm,
			},
		},
		prepare: defaultPostPrepare,
		verify:  noError,
	}, {
		name: "normal form request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    urlFoo,
				Method: http.MethodPost,
				Header: map[string]string{
					util.ContentType: "application/x-www-form-urlencoded",
				},
				Form: defaultForm,
			},
		},
		prepare: defaultPostPrepare,
		verify:  noError,
	}, {
		name: "body is a template",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    urlFoo,
				Method: http.MethodPost,
				Body:   atest.NewRequestBody(`{"name":"{{lower "HELLO"}}"}`),
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Post("/foo").BodyString(`{"name":"hello"}`).
				Reply(http.StatusOK).BodyString(`{}`)
		},
		verify: noError,
	}, {
		name: "status code not match",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:    urlFoo,
				Method: http.MethodPost,
			},
			Expect: atest.Response{
				StatusCode: http.StatusBadRequest,
			},
		},
		prepare: func() {
			gock.New(urlLocalhost).
				Post("/foo").
				Reply(http.StatusOK).
				SetHeader(util.ContentType, util.JSON).
				BodyString(defaultBody)
		},
		verify: func(t *testing.T, _ interface{}, err error) {
			assert.Error(t, err)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.verify == nil {
				tt.verify = hasError
			}
			runner := NewSimpleTestCaseRunner()
			runner.WithOutputWriter(io.Discard)
			if tt.execer != nil {
				runner.WithExecer(tt.execer)
			}
			output, err := runner.RunTestCase(tt.testCase, tt.ctx, context.TODO())
			tt.verify(t, output, err)

			getter, ok := runner.(ResponseRecord)
			assert.True(t, ok)
			assert.NotNil(t, getter.GetResponseRecord())
		})
	}
}

func TestLevelWriter(t *testing.T) {
	tests := []struct {
		name   string
		buf    *bytes.Buffer
		level  string
		expect string
	}{{
		name:   "debug",
		buf:    new(bytes.Buffer),
		level:  "debug",
		expect: "debuginfo",
	}, {
		name:   "info",
		buf:    new(bytes.Buffer),
		level:  "info",
		expect: "info",
	}}
	for _, tt := range tests {
		writer := NewDefaultLevelWriter(tt.level, tt.buf)
		if assert.NotNil(t, writer) {
			writer.Debug("debug")
			writer.Info("info")

			assert.Equal(t, tt.expect, tt.buf.String())
		}
	}
}

func TestJSONSchemaValidation(t *testing.T) {
	tests := []struct {
		name   string
		schema string
		body   string
		hasErr bool
	}{{
		name:   "normal",
		schema: defaultSchemaForTest,
		body:   `{"name": "linuxsuren", "age": 100}`,
		hasErr: false,
	}, {
		name:   "schema is empty",
		schema: "",
		hasErr: false,
	}, {
		name:   "failed to validate",
		schema: defaultSchemaForTest,
		body:   `{"name": "linuxsuren", "age": "100"}`,
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := jsonSchemaValidation(tt.schema, []byte(tt.body))
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestRunJob(t *testing.T) {
	tests := []struct {
		name   string
		job    atest.Job
		hasErr bool
	}{{
		name: "sleep 1s",
		job: atest.Job{
			Items: []string{"sleep(1)"},
		},
		hasErr: false,
	}, {
		name: "no params",
		job: atest.Job{
			Items: []string{"sleep()"},
		},
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runJob(&tt.job, nil, nil)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestContextKey(t *testing.T) {
	assert.Equal(t, ContextKey("parentDir"), NewContextKeyBuilder().ParentDir())

	ctx := context.WithValue(context.Background(), NewContextKeyBuilder().ParentDir(), "/tmp")
	assert.Equal(t, "/tmp", NewContextKeyBuilder().ParentDir().GetContextValueOrEmpty(ctx))
	assert.Empty(t, ContextKey("fake").GetContextValueOrEmpty(ctx))
}

func TestBodyFiledsVerify(t *testing.T) {
	tests := []struct {
		name       string
		bodyFields map[string]interface{}
		body       string
		hasErr     bool
	}{{
		name: "normal",
		bodyFields: map[string]interface{}{
			"name":   "linuxsuren",
			"number": 1,
		},
		body:   genericBody,
		hasErr: false,
	}, {
		name: "field not found",
		bodyFields: map[string]interface{}{
			"project": "",
		},
		body:   genericBody,
		hasErr: true,
	}, {
		name: "number is not equal",
		bodyFields: map[string]interface{}{
			"number": 2,
		},
		body:   genericBody,
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := &jsonBodyVerifier{
				body: &atest.Response{
					BodyFieldsExpect: tt.bodyFields,
				},
			}

			err := verifier.Verify([]byte(tt.body))
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestGetSuggestedAPIs(t *testing.T) {
	runner := NewSimpleTestCaseRunner()
	// not a swagger
	result, err := runner.GetSuggestedAPIs(&atest.TestSuite{}, "")
	assert.NoError(t, err, err)
	assert.Empty(t, result)

	// swagger
	gock.Off()
	gock.New(urlFoo).Get("swagger.json").Reply(http.StatusOK).File("testdata/swagger.json")
	result, err = runner.GetSuggestedAPIs(&atest.TestSuite{
		Spec: atest.APISpec{
			Kind: "swagger",
			URL:  urlFoo + "/swagger.json",
		},
	}, "")
	assert.NoError(t, err, err)
	assert.NotEmpty(t, result)
	method := result[0].Request.Method
	assert.Equal(t, strings.ToUpper(method), method)
}

func TestIsStructContent(t *testing.T) {
	tests := []struct {
		contentType string
		expectOk    bool
	}{{
		contentType: util.JSON,
		expectOk:    true,
	}, {
		contentType: util.YAML,
		expectOk:    true,
	}, {
		contentType: util.OctetStream,
		expectOk:    false,
	}, {
		contentType: util.OCIImageIndex,
		expectOk:    true,
	}}
	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			ok := isNonBinaryContent(tt.contentType)
			assert.Equal(t, tt.expectOk, ok)
		})
	}
}

const defaultSchemaForTest = `{"properties": {
	"name": {"type": "string"},
	"age": {"type": "integer"}
},
	"type":"object"
	}`

func hasError(t *testing.T, output interface{}, err error) {
	assert.NotNil(t, err)
}

func noError(t *testing.T, output interface{}, err error) {
	assert.Nil(t, err)
}

func prepareForFoo() {
	gock.New(urlLocalhost).
		Get("/foo").Reply(http.StatusOK).BodyString(genericBody).
		SetHeader(util.ContentType, util.JSON)
}

//go:embed testdata/generic_response.json
var genericBody string

const urlFoo = "http://localhost/foo"
const urlLocalhost = "http://localhost"
