package server_test

import (
	"net"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)

	fake := server.NewFakeHTTPServer()
	fake.WithHandler(nil)
	fake.Serve(lis)

	defaultHTTPServer := server.NewDefaultHTTPServer()
	defaultHTTPServer.WithHandler(nil)
}
