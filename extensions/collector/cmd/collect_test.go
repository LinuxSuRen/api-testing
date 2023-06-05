package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	c := NewRootCmd()
	assert.NotNil(t, c)
	assert.Equal(t, "atest-collector", c.Use)
}
