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

func TestRender(t *testing.T) {
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
