package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"service", "--action", "fake"})
	err := root.Execute()
	assert.NotNil(t, err)
}
