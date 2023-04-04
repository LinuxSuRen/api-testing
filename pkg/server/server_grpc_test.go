package server_test

import (
	"context"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	unimplemented := &server.UnimplementedRunnerServer{}
	_, err := unimplemented.Run(context.TODO(), nil)
	assert.NotNil(t, err)
}
