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
