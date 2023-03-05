package runner

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestTestCase(t *testing.T) {
	tests := []struct {
		name     string
		testCase *atest.TestCase
		ctx      interface{}
		prepare  func()
		verify   func(t *testing.T, output interface{}, err error)
	}{{
		name: "normal",
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
				BodyFieldsExpect: map[string]string{
					"name": "linuxsuren",
				},
				Header: map[string]string{
					"type": "generic",
				},
				Verify: []string{
					`name == "linuxsuren"`,
				},
			},
		},
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
			assert.Equal(t, map[string]interface{}{"name": "linuxsuren"}, output)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			tt.prepare()
			output, err := RunTestCase(tt.testCase, tt.ctx)
			tt.verify(t, output, err)
		})
	}
}
