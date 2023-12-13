/*
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

package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/cmd/service"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, server.NewFakeHTTPServer())
	root.SetArgs([]string{"service", "fake"})
	root.SetOut(new(bytes.Buffer))
	err := root.Execute()
	assert.NotNil(t, err)

	notSupportedMode := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "fake"}, server.NewFakeHTTPServer())
	notSupportedMode.SetArgs([]string{"service", paramAction, "install", "--mode=fake"})
	notSupportedMode.SetOut(new(bytes.Buffer))
	assert.NotNil(t, notSupportedMode.Execute())

	tmpFile, err := os.CreateTemp(os.TempDir(), "service")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpFile.Name())
	}()

	tests := []struct {
		name         string
		action       string
		targetOS     string
		mode         string
		expectOutput string
	}{{
		name:         "action: start",
		action:       "start",
		targetOS:     "linux",
		expectOutput: "output1",
	}, {
		name:         "action: stop",
		action:       "stop",
		targetOS:     "linux",
		expectOutput: "output2",
	}, {
		name:         "action: restart",
		action:       "restart",
		targetOS:     "linux",
		expectOutput: "output3",
	}, {
		name:         "action: status",
		action:       "status",
		targetOS:     "linux",
		expectOutput: "output4",
	}, {
		name:         "action: install",
		action:       "install",
		targetOS:     "linux",
		expectOutput: "output4",
	}, {
		name:         "action: uninstall",
		action:       "uninstall",
		targetOS:     "linux",
		expectOutput: "output4",
	}, {
		name:         "action: start, macos",
		action:       "start",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}, {
		name:         "action: stop, macos",
		action:       "stop",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}, {
		name:         "action: restart, macos",
		action:       "restart",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}, {
		name:         "action: status, macos",
		action:       "status",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}, {
		name:         "action: install, macos",
		action:       "install",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}, {
		name:         "action: uninstall, macos",
		action:       "uninstall",
		targetOS:     fakeruntime.OSDarwin,
		expectOutput: "output4",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode == "" {
				tt.mode = string(service.ServiceModeOS)
			}

			buf := new(bytes.Buffer)
			normalRoot := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: tt.targetOS, ExpectOutput: tt.expectOutput},
				server.NewFakeHTTPServer())
			normalRoot.SetOut(buf)
			normalRoot.SetArgs([]string{"service",
				"--script-path", tmpFile.Name(), "--mode", tt.mode, "--image=",
				"--skywalking=http://localhost:8080",
				"--secret-server=http://localhost:9090",
				tt.action})
			fmt.Println([]string{"service",
				"--script-path", tmpFile.Name(), "--mode", tt.mode, "--image=",
				"--skywalking=http://localhost:8080",
				"--secret-server=http://localhost:9090",
				tt.action})
			err = normalRoot.Execute()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectOutput, strings.TrimSpace(buf.String()))
		})
	}
}

const paramAction = "--action"
