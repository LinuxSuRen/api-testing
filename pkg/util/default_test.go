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

package util_test

import (
	"net/http"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMakeSureNotNil(t *testing.T) {
	var fun func()
	var mapStruct map[string]string

	assert.NotNil(t, util.MakeSureNotNil(fun))
	assert.NotNil(t, util.MakeSureNotNil(TestMakeSureNotNil))
	assert.NotNil(t, util.MakeSureNotNil(mapStruct))
	assert.NotNil(t, util.MakeSureNotNil(map[string]string{}))
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
			result := util.EmptyThenDefault(tt.val, tt.defVal)
			assert.Equal(t, tt.expect, result, result)
		})
	}

	assert.Equal(t, 1, util.ZeroThenDefault(0, 1))
	assert.Equal(t, 1, util.ZeroThenDefault(1, 2))
}

func TestGetFirstHeaderValue(t *testing.T) {
	tests := []struct {
		name   string
		header http.Header
		key    string
		expect string
	}{{
		name:   "empty header",
		header: http.Header{},
		key:    "abc",
		expect: "",
	}, {
		name: "not empty header",
		header: http.Header{
			"abc": []string{"def"},
		},
		key:    "abc",
		expect: "def",
	}, {
		name: "not empty header, has multiple values",
		header: http.Header{
			"abc": []string{"def", "ghi"},
		},
		key:    "abc",
		expect: "def",
	}, {
		name: "have ; in the value",
		header: http.Header{
			"abc": []string{"def;ghi"},
		},
		key:    "abc",
		expect: "def",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.GetFirstHeaderValue(tt.header, tt.key)
			assert.Equal(t, tt.expect, result, result)
		})
	}
}
