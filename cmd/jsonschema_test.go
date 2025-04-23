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

package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestJSONSchemaCmd(t *testing.T) {
	c := cmd.NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"},
		server.NewFakeHTTPServer())

	buf := new(bytes.Buffer)
	c.SetOut(buf)

	c.SetArgs([]string{"json"})
	err := c.Execute()
	assert.Nil(t, err)
	assert.True(t, strings.Contains(buf.String(), "schema"))
}
