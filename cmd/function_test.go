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
			assert.Equal(t, "No such function\n", output)
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
				cmd.NewFakeGRPCServer(), server.NewFakeHTTPServer())

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
