/**
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

package cmd

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	t.Run("not run", func(t *testing.T) {
		cmd := newRootCmdForTest()
		assert.NotNil(t, cmd)
		assert.Equal(t, "atest-store-git", cmd.Use)
		assert.Equal(t, "7074", cmd.Flags().Lookup("port").Value.String())
	})

	t.Run("invalid port", func(t *testing.T) {
		cmd := newRootCmdForTest()
		cmd.SetArgs([]string{"--port", "-1"})
		err := cmd.Execute()
		assert.Error(t, err)
	})

	t.Run("stop the command", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()

		cmd := newRootCmdForTest()
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"--port", "0"})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func newRootCmdForTest() *cobra.Command {
	cmd := NewRootCommand()
	cmd.SetOut(io.Discard)
	return cmd
}
