package cmd_test

import (
	"bytes"
	"testing"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/sample"
	"github.com/stretchr/testify/assert"
)

func TestSampleCmd(t *testing.T) {
	c := cmd.NewRootCmd()

	buf := new(bytes.Buffer)
	c.SetOut(buf)

	c.SetArgs([]string{"sample"})
	err := c.Execute()
	assert.Nil(t, err)
	assert.Equal(t, sample.TestSuiteGitLab+"\n", buf.String())
}
