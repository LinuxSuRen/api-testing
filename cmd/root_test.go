package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

func TestCreateRunCommand(t *testing.T) {
	cmd := createRunCommand()
	assert.Equal(t, "run", cmd.Use)

	init := createInitCommand(fakeruntime.FakeExecer{})
	assert.Equal(t, "init", init.Use)

	s := createServerCmd(&fakeGRPCServer{}, server.NewFakeHTTPServer())
	assert.NotNil(t, s)
	assert.Equal(t, "server", s.Use)

	root := NewRootCmd(fakeruntime.FakeExecer{}, NewFakeGRPCServer(), server.NewFakeHTTPServer())
	root.SetArgs([]string{"init", "-k=demo.yaml", "--wait-namespace", "demo", "--wait-resource", "demo"})
	err := root.Execute()
	assert.Nil(t, err)
}

func TestRootCmd(t *testing.T) {
	c := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, NewFakeGRPCServer(), server.NewFakeHTTPServer())
	assert.NotNil(t, c)
	assert.Equal(t, "atest", c.Use)
}
