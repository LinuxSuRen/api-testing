/*
Copyright 2023 API Testing Authors.

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
	"context"
	"fmt"
	"io"
	"os"
	"path"
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

	t.Run("stop the command with socket", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()

		tempSocket := path.Join(os.TempDir(), fmt.Sprintf("atest-store-git-%d.sock", time.Now().UnixMicro()))
		cmd := newRootCmdForTest()
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"--socket", tempSocket})
		err := cmd.Execute()
		assert.NoError(t, err)
	})
}

func newRootCmdForTest() *cobra.Command {
	cmd := NewRootCommand()
	cmd.SetOut(io.Discard)
	return cmd
}
