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
