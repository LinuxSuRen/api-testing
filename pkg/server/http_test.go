/*
Copyright 2024 API Testing Authors.

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
package server_test

import (
	"net"
	"net/http"
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

func TestCombineHandler(t *testing.T) {
	defaultHandler := http.NewServeMux()
	fakeHandler := http.NewServeMux()

	t.Run("correct default handler", func(t *testing.T) {
		combineHandler := server.NewDefaultCombineHandler()

		combineHandler.PutHandler("", defaultHandler)
		combineHandler.PutHandler("/fake", fakeHandler)

		assert.NotEqual(t, defaultHandler, combineHandler.GetHandler())

		fakeServer := server.NewFakeHTTPServer()
		assert.Nil(t, fakeServer.Shutdown(nil))
	})

	t.Run("only one default handler", func(t *testing.T) {
		combineHandler := server.NewDefaultCombineHandler()

		combineHandler.PutHandler("", defaultHandler)

		assert.Equal(t, defaultHandler, combineHandler.GetHandler())
	})
}
