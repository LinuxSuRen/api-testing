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

package cmd_test

import (
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	c := cmd.NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"},
		cmd.NewFakeGRPCServer(), server.NewFakeHTTPServer())
	c.SetOut(io.Discard)

	t.Run("normal", func(t *testing.T) {
		tmpFile := path.Join(os.TempDir(), time.Now().String())
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
}
