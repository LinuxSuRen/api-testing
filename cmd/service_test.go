package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, NewFakeGRPCServer(), server.NewFakeHTTPServer())
	root.SetArgs([]string{"service", "fake"})
	root.SetOut(new(bytes.Buffer))
	err := root.Execute()
	assert.NotNil(t, err)

	notLinux := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "fake"}, NewFakeGRPCServer(), server.NewFakeHTTPServer())
	notLinux.SetArgs([]string{"service", paramAction, "install"})
	notLinux.SetOut(new(bytes.Buffer))
	err = notLinux.Execute()
	assert.NotNil(t, err)

	notSupportedMode := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "fake"}, NewFakeGRPCServer(), server.NewFakeHTTPServer())
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
	}, {
		name:         "start in podman",
		action:       "start",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModePodman),
		expectOutput: "",
	}, {
		name:         "start in docker",
		action:       "start",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "stop in docker",
		action:       "stop",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "restart in docker",
		action:       "restart",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "status in docker",
		action:       "status",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "install in docker",
		action:       "install",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "uninstall in docker",
		action:       "uninstall",
		targetOS:     fakeruntime.OSLinux,
		mode:         string(ServiceModeDocker),
		expectOutput: "",
	}, {
		name:         "start in podman",
		action:       "start",
		targetOS:     fakeruntime.OSLinux,
		mode:         ServiceModePodman.String(),
		expectOutput: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode == "" {
				tt.mode = string(ServiceModeOS)
			}

			buf := new(bytes.Buffer)
			normalRoot := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: tt.targetOS, ExpectOutput: tt.expectOutput},
				NewFakeGRPCServer(), server.NewFakeHTTPServer())
			normalRoot.SetOut(buf)
			normalRoot.SetArgs([]string{"service", "--action", tt.action,
				"--script-path", tmpFile.Name(), "--mode", tt.mode, "--image="})
			err = normalRoot.Execute()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectOutput, strings.TrimSpace(buf.String()))
		})
	}
}

const paramAction = "--action"
