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
package cmd_test

import (
	"bytes"
	"testing"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestCreateFunctionCommand(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		verify func(t *testing.T, output string)
	}{{
		name: "normal",
		args: []string{"func"},
		verify: func(t *testing.T, output string) {
			assert.NotEmpty(t, output)
		},
	}, {
		name: "with function name",
		args: []string{"func", "clean"},
		verify: func(t *testing.T, output string) {
			assert.NotEmpty(t, output)
		},
	}, {
		name: "with not exit function",
		args: []string{"func", "fake"},
		verify: func(t *testing.T, output string) {
			assert.Equal(t, "No such function\n\nAll expr functions:\n", output)
		},
	}, {
		name: "query functions, not found",
		args: []string{"func", "--feature", `unknown`},
		verify: func(t *testing.T, output string) {
			assert.Equal(t, "\n", output)
		},
	}, {
		name: "query functions, not found",
		args: []string{"func", "--feature", `生成对象，字段包含 name`},
		verify: func(t *testing.T, output string) {
			assert.Equal(t, `{{generateJSONString "name"}}
`, output)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cmd.NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"},
				server.NewFakeHTTPServer())

			buf := new(bytes.Buffer)
			c.SetOut(buf)
			c.SetArgs(tt.args)

			err := c.Execute()
			assert.NoError(t, err)

			if tt.verify != nil {
				tt.verify(t, buf.String())
			}
		})
	}
}
