package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRemoteServer(t *testing.T) {
	server, err := NewRemoteServer()
	assert.NotNil(t, server)
	assert.NoError(t, err)
}
