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
	err := root.Execute()
	assert.NotNil(t, err)

	notLinux := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "fake"}, NewFakeGRPCServer())
	notLinux.SetArgs([]string{"service", paramAction, "install"})
	err = notLinux.Execute()
	assert.NotNil(t, err)

	tmpFile, err := os.CreateTemp(os.TempDir(), "service")
	assert.Nil(t, err)
	defer func() {
		os.RemoveAll(tmpFile.Name())
	}()

	targetScript := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, NewFakeGRPCServer())
	targetScript.SetArgs([]string{"service", paramAction, "install", "--script-path", tmpFile.Name()})
	err = targetScript.Execute()
	assert.Nil(t, err)
	data, err := os.ReadFile(tmpFile.Name())
	assert.Nil(t, err)
	assert.Equal(t, script, string(data))

	tests := []struct {
		name         string
		action       string
		expectOutput string
	}{{
		name:         "action: start",
		action:       "start",
		expectOutput: "output1",
	}, {
		name:         "action: stop",
		action:       "stop",
		expectOutput: "output2",
	}, {
		name:         "action: restart",
		action:       "restart",
		expectOutput: "output3",
	}, {
		name:         "action: status",
		action:       "status",
		expectOutput: "output4",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			normalRoot := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux", ExpectOutput: tt.expectOutput}, NewFakeGRPCServer())
			normalRoot.SetOut(buf)
			normalRoot.SetArgs([]string{"service", "--action", tt.action})
			err = normalRoot.Execute()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectOutput+"\n", buf.String())
		})
	}
}

const paramAction = "--action"
