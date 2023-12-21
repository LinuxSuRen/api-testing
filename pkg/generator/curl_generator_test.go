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
package generator

import (
	"net/http"
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
)

func TestCurlGenerator(t *testing.T) {
	tests := []struct {
		name     string
		testCase atest.TestCase
		expect   string
	}{{
		name: "basic HTTP GET",
		testCase: atest.TestCase{
			Request: atest.Request{
				API: fooForTest,
			},
		},
		expect: `curl -X GET 'http://foo'`,
	}, {
		name: "has query string",
		testCase: atest.TestCase{
			Request: atest.Request{
				API: fooForTest,
				Query: map[string]interface{}{
					"page": "1",
					"size": "10",
				},
			},
		},
		expect: `curl -X GET 'http://foo?page=1&size=10'`,
	}, {
		name: "basic HTTP POST",
		testCase: atest.TestCase{
			Request: atest.Request{
				API:    fooForTest,
				Method: http.MethodPost,
			},
		},
		expect: `curl -X POST 'http://foo'`,
	}, {
		name: "has header",
		testCase: atest.TestCase{
			Request: atest.Request{
				API: fooForTest,
				Header: map[string]string{
					"Content-Type": util.Plain,
					"Connection":   "keep-alive",
				},
			},
		},
		expect: `curl -X GET 'http://foo' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: text/plain'`,
	}, {
		name: "has body",
		testCase: atest.TestCase{
			Request: atest.Request{
				API:  fooForTest,
				Body: atest.NewRequestBody("hello"),
			},
		},
		expect: `curl -X GET 'http://foo' \
  --data-raw 'hello'`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &curlGenerator{}
			if got, err := g.Generate(nil, &tt.testCase); err != nil || got != tt.expect {
				t.Errorf("got %q, want %q, error: %v", got, tt.expect, err)
			}
		})
	}
}

const fooForTest = "http://foo"
