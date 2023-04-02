package server_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestGetProtos(t *testing.T) {
	protos := server.GetProtos()
	assert.Equal(t, 1, len(protos))

	exists := []string{"server.proto"}
	for _, key := range exists {
		content, ok := protos[key]
		assert.True(t, ok)
		assert.NotEmpty(t, content)
	}
}
