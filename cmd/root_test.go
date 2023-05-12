package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	exec "github.com/linuxsuren/go-fake-runtime"
)

func TestCreateRunCommand(t *testing.T) {
	cmd := createRunCommand()
	assert.Equal(t, "run", cmd.Use)

	init := createInitCommand(exec.FakeExecer{})
	assert.Equal(t, "init", init.Use)

	server := createServerCmd(&fakeGRPCServer{})
	assert.NotNil(t, server)
	assert.Equal(t, "server", server.Use)

	root := NewRootCmd(exec.FakeExecer{}, NewFakeGRPCServer())
	root.SetArgs([]string{"init", "-k=demo.yaml", "--wait-namespace", "demo", "--wait-resource", "demo"})
	err := root.Execute()
	assert.Nil(t, err)
}
