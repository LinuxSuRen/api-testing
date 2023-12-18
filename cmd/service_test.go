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
