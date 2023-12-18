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
	"io"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	t.Run("not run", func(t *testing.T) {
		cmd := newRootCmdForTest()
		assert.NotNil(t, cmd)
		assert.Equal(t, "atest-store-etcd", cmd.Use)
		assert.Equal(t, "7073", cmd.Flags().Lookup("port").Value.String())
	})

	t.Run("invalid port", func(t *testing.T) {
		cmd := newRootCmdForTest()
		cmd.SetArgs([]string{"--port", "-1"})
		err := cmd.Execute()
		assert.Error(t, err)
	})
}

func newRootCmdForTest() *cobra.Command {
	cmd := NewRootCommand()
	cmd.SetOut(io.Discard)
	return cmd
}
