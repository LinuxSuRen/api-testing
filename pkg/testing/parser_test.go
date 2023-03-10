package testing

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	suite, err := Parse("../../sample/testsuite-gitlab.yaml")
	if assert.Nil(t, err) && assert.NotNil(t, suite) {
		assert.Equal(t, "Gitlab", suite.Name)
		assert.Equal(t, 2, len(suite.Items))
		assert.Equal(t, TestCase{
			Name: "projects",
			Request: Request{
				API: "https://gitlab.com/api/v4/projects",
			},
			Expect: Response{
				StatusCode: http.StatusOK,
			},
		}, suite.Items[0])
	}
}

func TestRequestRender(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
		verify  func(t *testing.T, req *Request)
		ctx     interface{}
		hasErr  bool
	}{{
		name: "slice as context",
		request: &Request{
			API:  "http://localhost/{{index . 0}}",
			Body: "{{index . 1}}",
		},
		ctx:    []string{"foo", "bar"},
		hasErr: false,
		verify: func(t *testing.T, req *Request) {
			assert.Equal(t, "http://localhost/foo", req.API)
			assert.Equal(t, "bar", req.Body)
		},
	}, {
		name:    "default values",
		request: &Request{},
		verify: func(t *testing.T, req *Request) {
			assert.Equal(t, http.MethodGet, req.Method)
		},
		hasErr: false,
	}, {
		name:    "context is nil",
		request: &Request{},
		ctx:     nil,
		hasErr:  false,
	}, {
		name: "body from file",
		request: &Request{
			BodyFromFile: "testdata/generic_body.json",
		},
		ctx: TestCase{
			Name: "linuxsuren",
		},
		hasErr: false,
		verify: func(t *testing.T, req *Request) {
			assert.Equal(t, `{"name": "linuxsuren"}`, req.Body)
		},
	}, {
		name: "body file not found",
		request: &Request{
			BodyFromFile: "testdata/fake",
		},
		hasErr: true,
	}, {
		name: "invalid API as template",
		request: &Request{
			API: "{{.name}",
		},
		hasErr: true,
	}, {
		name: "failed with API render",
		request: &Request{
			API: "{{.name}}",
		},
		ctx:    TestCase{},
		hasErr: true,
	}, {
		name: "invalid body as template",
		request: &Request{
			Body: "{{.name}",
		},
		hasErr: true,
	}, {
		name: "failed with body render",
		request: &Request{
			Body: "{{.name}}",
		},
		ctx:    TestCase{},
		hasErr: true,
	}, {
		name: "form render",
		request: &Request{
			Form: map[string]string{
				"key": "{{.Name}}",
			},
		},
		ctx: TestCase{Name: "linuxsuren"},
		verify: func(t *testing.T, req *Request) {
			assert.Equal(t, "linuxsuren", req.Form["key"])
		},
		hasErr: false,
	}, {
		name: "header render",
		request: &Request{
			Header: map[string]string{
				"key": "{{.Name}}",
			},
		},
		ctx: TestCase{Name: "linuxsuren"},
		verify: func(t *testing.T, req *Request) {
			assert.Equal(t, "linuxsuren", req.Header["key"])
		},
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Render(tt.ctx)
			if assert.Equal(t, tt.hasErr, err != nil, err) && tt.verify != nil {
				tt.verify(t, tt.request)
			}
		})
	}
}

func TestResponseRender(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
		verify   func(t *testing.T, req *Response)
		ctx      interface{}
		hasErr   bool
	}{{
		name:     "blank response",
		response: &Response{},
		verify: func(t *testing.T, req *Response) {
			assert.Equal(t, http.StatusOK, req.StatusCode)
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
}

func TestEmptyThenDefault(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		defVal string
		expect string
	}{{
		name:   "empty string",
		val:    "",
		defVal: "abc",
		expect: "abc",
	}, {
		name:   "blank string",
		val:    " ",
		defVal: "abc",
		expect: "abc",
	}, {
		name:   "not empty or blank string",
		val:    "abc",
		defVal: "def",
		expect: "abc",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := emptyThenDefault(tt.val, tt.defVal)
			assert.Equal(t, tt.expect, result, result)
		})
	}

	assert.Equal(t, 1, zeroThenDefault(0, 1))
	assert.Equal(t, 1, zeroThenDefault(1, 2))
}
