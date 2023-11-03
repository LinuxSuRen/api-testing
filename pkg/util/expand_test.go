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
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{{
		name:   "without brace",
		input:  "/home",
		expect: []string{"/home"},
	}, {
		name:   "with brace",
		input:  "/home/{good,bad}",
		expect: []string{"/home/good", "/home/bad"},
	}, {
		name:   "with brace, have suffix",
		input:  "/home/{good,bad}.yaml",
		expect: []string{"/home/good.yaml", "/home/bad.yaml"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.Expand(tt.input)
			assert.Equal(t, tt.expect, got, got)
		})
	}
}
