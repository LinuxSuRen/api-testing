/*
Copyright 2023-2025 API Testing Authors.

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
	_ "embed"
	"io"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	c := cmd.NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"},
		server.NewFakeHTTPServer())
	c.SetOut(io.Discard)

	t.Run("normal", func(t *testing.T) {
		now := strconv.Itoa(int(time.Now().Unix()))
		tmpFile := path.Join(os.TempDir(), now)
		defer os.RemoveAll(tmpFile)

		c.SetArgs([]string{"convert", "-p=testdata/simple-suite.yaml", "--converter=jmeter", "--target", tmpFile})

		err := c.Execute()
		assert.NoError(t, err)

		var data []byte
		data, err = os.ReadFile(tmpFile)
		if assert.NoError(t, err) {
			assert.NotEmpty(t, string(data))
		}
	})

	t.Run("no testSuite", func(t *testing.T) {
		c.SetArgs([]string{"convert", "-p=testdata/fake.yaml", "--converter=jmeter"})

		err := c.Execute()
		assert.Error(t, err)
	})

	t.Run("no converter found", func(t *testing.T) {
		c.SetArgs([]string{"convert", "-p=testdata/simple-suite.yaml", "--converter=fake"})

		err := c.Execute()
		assert.Error(t, err)
	})

	t.Run("flag --pattern is required", func(t *testing.T) {
		c.SetArgs([]string{"convert", "--converter=fake"})

		err := c.Execute()
		assert.Error(t, err)
	})

	t.Run("flag --converter is required", func(t *testing.T) {
		c.SetArgs([]string{"convert", "-p=fake"})

		err := c.Execute()
		assert.Error(t, err)
	})

	t.Run("not supported source format", func(t *testing.T) {
		c.SetArgs([]string{"convert", "--source=fake"})
		err := c.Execute()
		assert.Error(t, err)
	})

	t.Run("convert from postmant", func(t *testing.T) {
		tmpFile := path.Join(os.TempDir(), time.Now().String())
		defer os.RemoveAll(tmpFile)

		c.SetArgs([]string{"convert", "--source=postman", "--target=", tmpFile, "-p=testdata/postman.json"})
		err := c.Execute()
		assert.NoError(t, err)
	})
}
