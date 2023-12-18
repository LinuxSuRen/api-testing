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

func FuzzZeroThenDefault(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b int) {
		val := util.ZeroThenDefault(a, b)
		if a == 0 {
			assert.Equal(t, b, val)
		} else {
			assert.Equal(t, a, val)
		}
	})
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
