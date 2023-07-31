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
		assert.Equal(t, "store-s3", cmd.Use)
		assert.Equal(t, "7072", cmd.Flags().Lookup("port").Value.String())
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
	cmd := NewRootCmd(nil)
	cmd.SetOut(io.Discard)
	return cmd
}
