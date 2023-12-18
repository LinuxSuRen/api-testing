/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
				Body: "hello",
			},
		},
		expect: `curl -X GET 'http://foo' \
  --data-raw 'hello'`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &curlGenerator{}
			if got, err := g.Generate(nil, &tt.testCase); err != nil || got != tt.expect {
				t.Errorf("Generate() = %v, want %v", got, tt.expect)
			}
		})
	}
}

const fooForTest = "http://foo"
