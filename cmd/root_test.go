/*
Copyright 2025 API Testing Authors.

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

func TestCreateRunCommand(t *testing.T) {
	execer := fakeruntime.FakeExecer{}

	cmd := createRunCommand()
	assert.Equal(t, "run", cmd.Use)

	init := createInitCommand(execer)
	assert.Equal(t, "init", init.Use)

	s := createServerCmd(execer, server.NewFakeHTTPServer())
	assert.NotNil(t, s)
	assert.Equal(t, "server", s.Use)

	root := NewRootCmd(execer, server.NewFakeHTTPServer())
	root.SetArgs([]string{"init", "-k=demo.yaml", "--wait-namespace", "demo", "--wait-resource", "demo"})
	err := root.Execute()
	assert.Nil(t, err)
}

func TestRootCmd(t *testing.T) {
	c := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, server.NewFakeHTTPServer())
	assert.NotNil(t, c)
	assert.Equal(t, "atest", c.Use)
}
