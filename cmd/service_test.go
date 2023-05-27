package cmd

import (
	"bytes"
	"os"
	"testing"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, NewFakeGRPCServer())
	root.SetArgs([]string{"service", "fake"})
	root.SetOut(new(bytes.Buffer))
	err := root.Execute()
	assert.NotNil(t, err)

	notLinux := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "fake"}, NewFakeGRPCServer())
	notLinux.SetArgs([]string{"service", paramAction, "install"})
	notLinux.SetOut(new(bytes.Buffer))
	err = notLinux.Execute()
	assert.NotNil(t, err)

	tmpFile, err := os.CreateTemp(os.TempDir(), "service")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpFile.Name())
	}()

	tests := []struct {
		name         string
		action       string
		targetOS     string
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			normalRoot := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: tt.targetOS, ExpectOutput: tt.expectOutput}, NewFakeGRPCServer())
			normalRoot.SetOut(buf)
			normalRoot.SetArgs([]string{"service", "--action", tt.action, "--script-path", tmpFile.Name()})
			err = normalRoot.Execute()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectOutput+"\n", buf.String())
		})
	}
}

const paramAction = "--action"
