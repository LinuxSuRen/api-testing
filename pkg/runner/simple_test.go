package runner

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestTestCase(t *testing.T) {
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
			Prepare: atest.Prepare{
				Kubernetes: []string{"demo.yaml"},
			},
		},
		execer: fakeruntime.FakeExecer{ExpectError: errors.New("fake")},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "normal, response is map",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
				Header: map[string]string{
					"key": "value",
				},
				Body: `{"foo":"bar"}`,
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
			Prepare: atest.Prepare{
				Kubernetes: []string{"demo.yaml"},
			},
			Clean: atest.Clean{
				CleanPrepare: true,
			},
		},
		execer: fakeruntime.FakeExecer{},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").
				MatchHeader("key", "value").
				Reply(http.StatusOK).
				SetHeader("type", "generic").
				File("testdata/generic_response.json")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.Nil(t, err)
			assert.Equal(t, map[string]interface{}{"name": "linuxsuren", "number": float64(1)}, output)
		},
	}, {
		name: "normal, response is slice",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
				Body:       `["foo", "bar"]`,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").
				Reply(http.StatusOK).
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
				API:          "http://localhost/foo",
				Method:       http.MethodPost,
				BodyFromFile: "testdata/generic_response.json",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Post("/foo").BodyString(genericBody).
				Reply(http.StatusOK).BodyString("123")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "response from a not found file",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API:          "http://localhost/foo",
				Method:       http.MethodPost,
				BodyFromFile: "testdata/fake.json",
			},
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "bad request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				StatusCode: http.StatusOK,
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusBadRequest)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "error with request",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").ReplyError(errors.New("error"))
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not match with body",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Body: "bar",
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString("foo")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not match with header",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				Header: map[string]string{
					"foo": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).SetHeader("foo", "value")
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "not found from fields",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(genericBody)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "body filed not match",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"name": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(genericBody)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "invalid filed finding",
		testCase: &atest.TestCase{
			Request: atest.Request{
				API: "http://localhost/foo",
			},
			Expect: atest.Response{
				BodyFieldsExpect: map[string]interface{}{
					"items[1]": "bar",
				},
			},
		},
		prepare: func() {
			gock.New("http://localhost").
				Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		},
		verify: func(t *testing.T, output interface{}, err error) {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "failed to get field")
		},
	},
		// {
		// 	name: "verify failed",
		// 	testCase: &atest.TestCase{
		// 		Request: atest.Request{
		// 			API: "http://localhost/foo",
		// 		},
		// 		Expect: atest.Response{
		// 			Verify: []string{
		// 				"len(data.items) > 0",
		// 			},
		// 		},
		// 	},
		// 	prepare: func() {
		// 		gock.New("http://localhost").
		// 			Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
		// 	},
		// 	verify: func(t *testing.T, output interface{}, err error) {
		// 		if assert.NotNil(t, err) {
		// 			assert.Contains(t, err.Error(), "failed to verify")
		// 		}
		// 	},
		// },
		{
			name: "failed to compile",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API: "http://localhost/foo",
				},
				Expect: atest.Response{
					Verify: []string{
						`println("12")`,
					},
				},
			},
			prepare: func() {
				gock.New("http://localhost").
					Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
			},
			verify: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "unknown name println")
			},
		}, {
			name: "failed to compile",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API: "http://localhost/foo",
				},
				Expect: atest.Response{
					Verify: []string{
						`1 + 1`,
					},
				},
			},
			prepare: func() {
				gock.New("http://localhost").
					Get("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
			},
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
					API:    "http://localhost/foo",
					Method: http.MethodPost,
					Header: map[string]string{
						"Content-Type": "multipart/form-data",
					},
					Form: map[string]string{
						"key": "value",
					},
				},
			},
			prepare: func() {
				gock.New("http://localhost").
					Post("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
			},
			verify: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, err)
			},
		}, {
			name: "normal form request",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:    "http://localhost/foo",
					Method: http.MethodPost,
					Header: map[string]string{
						"Content-Type": "application/x-www-form-urlencoded",
					},
					Form: map[string]string{
						"key": "value",
					},
				},
			},
			prepare: func() {
				gock.New("http://localhost").
					Post("/foo").Reply(http.StatusOK).BodyString(`{"items":[]}`)
			},
			verify: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, err)
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			if tt.prepare != nil {
				tt.prepare()
			}
			runner := NewSimpleTestCaseRunner().WithOutputWriter(os.Stdout)
			if tt.execer != nil {
				runner.WithExecer(tt.execer)
			}
			output, err := runner.RunTestCase(tt.testCase, tt.ctx, context.TODO())
			tt.verify(t, output, err)
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

const defaultSchemaForTest = `{"properties": {
	"name": {"type": "string"},
	"age": {"type": "integer"}
},
	"type":"object"
	}`

//go:embed testdata/generic_response.json
var genericBody string
