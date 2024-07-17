/*
Copyright 2023-2024 API Testing Authors.

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
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPrintProto(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		verify func(*testing.T, *bytes.Buffer, error)
	}{{
		name: "print ptoto only",
		args: []string{"server", "--print-proto"},
		verify: func(t *testing.T, buf *bytes.Buffer, err error) {
			assert.Nil(t, err)
			assert.True(t, strings.HasPrefix(buf.String(), `syntax = "proto3";`))
		},
	}, {
		name: "invalid port",
		args: []string{"server", "-p=-1"},
		verify: func(t *testing.T, buf *bytes.Buffer, err error) {
			assert.NotNil(t, err)
		},
	}, {
		name: "random port",
		args: []string{"server", "-p=0", "--http-port=0",
			"--local-storage=./*", "--secret-server=localhost:7073"},
		verify: func(t *testing.T, buf *bytes.Buffer, err error) {
			assert.Nil(t, err)
		},
	}, {
		name: "mock server, not found config",
		args: []string{"server", "--mock-config=fake", "-p=0", "--http-port=0"},
		verify: func(t *testing.T, buffer *bytes.Buffer, err error) {
			assert.Error(t, err)
		},
	}, {
		name: "mock server, normal",
		args: []string{"server", "--mock-config=testdata/invalid-api.yaml", "-p=0", "--http-port=0"},
		verify: func(t *testing.T, buffer *bytes.Buffer, err error) {
			assert.Error(t, err)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, server.NewFakeHTTPServer())
			root.SetOut(buf)
			root.SetArgs(append(tt.args, "--dry-run"))
			err := root.Execute()
			tt.verify(t, buf, err)
		})
	}
}

func TestFrontEndHandlerWithLocation(t *testing.T) {
	handler := frontEndHandlerWithLocation("testdata")
	const expect404 = "404 page not found\n"

	t.Run("404", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, expect404, resp.GetBody().String())
	})

	t.Run("get js", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/assets/index.js", nil)
		assert.NoError(t, err)
		defer func() {
			uiResourceJS = []byte("")
		}()

		resp := newFakeResponseWriter()

		uiResourceJS = []byte("js")
		handler(resp, req, map[string]string{})
		assert.Equal(t, uiResourceJS, resp.GetBody().Bytes())

		assert.Equal(t, "text/javascript; charset=utf-8", resp.Header().Get(util.ContentType))
	})

	t.Run("get css", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/assets/index.css", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, expect404, resp.GetBody().String())
	})

	t.Run("favicon", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/favicon.ico", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, "image/x-icon", resp.Header().Get(util.ContentType))
	})

	t.Run("healthz", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/healthz", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, "ok", resp.GetBody().String())
	})

	t.Run("pprof", func(t *testing.T) {
		apis := []string{"", "cmdline", "symbol",
			"trace", "profile",
			"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"}

		mu := runtime.NewServeMux()
		debugHandler(mu, nil)

		ready := make(chan struct{})
		var err error
		var listen net.Listener
		var port string
		go func() {
			listen, err = net.Listen("tcp", ":0")
			assert.NoError(t, err)

			addr := listen.Addr().String()
			items := strings.Split(addr, ":")
			port = items[len(items)-1]

			ready <- struct{}{}
			server := http.Server{Addr: addr, Handler: mu}
			server.Serve(listen)
		}()

		<-ready
		defer listen.Close()

		for _, name := range apis {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/debug/pprof/%s?seconds=1", port, name))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("download atest", func(t *testing.T) {
		opt := &serverOption{
			execer: fakeruntime.FakeExecer{
				ExpectOS:            "linux",
				ExpectLookPathError: errors.New("fake"),
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/get", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()

		opt.getAtestBinary(resp, req, map[string]string{})
		assert.Equal(t, `not found "atest"`, resp.GetBody().String())
	})

	t.Run("download atest, failed to read", func(t *testing.T) {
		opt := &serverOption{
			execer: fakeruntime.FakeExecer{
				ExpectOS: "linux",
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/get", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()

		opt.getAtestBinary(resp, req, map[string]string{})
		assert.Equal(t, `failed to read "atest": open : no such file or directory`, resp.GetBody().String())
	})
}

func TestProxy(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://localhost:8080").Post("/api/v1/echo").Reply(http.StatusOK)
		gock.New("http://localhost:9090").Post("/api/v1/echo").Reply(http.StatusOK)

		handle := postRequestProxy("http://localhost:9090/")
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/echo", strings.NewReader(`{"message": "hello"}`))
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handle(resp, req, map[string]string{})
	})

	t.Run("no proxy", func(t *testing.T) {
		defer gock.Off()

		gock.New("http://localhost:8080").Post("/api/v1/echo").Reply(http.StatusOK)

		handle := postRequestProxy("")
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/echo", strings.NewReader(`{"message": "hello"}`))
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handle(resp, req, map[string]string{})
	})
}

func TestOAuth(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		hasErr bool
	}{{

		name:   "invalid oauth provider",
		args:   []string{"server", "--auth=oauth", "--oauth-provider=fake"},
		hasErr: true,
	}, {
		name:   "client-id is missing",
		args:   []string{"server", "--auth=oauth", "--client-secret=fake"},
		hasErr: true,
	}, {
		name:   "client-secret is missing",
		args:   []string{"server", "--auth=oauth", "--client-id=fake"},
		hasErr: true,
	}, {
		name:   "oauth is ok",
		args:   []string{"server", "--auth=oauth", "--client-id=fake", "--client-secret=fake", "--http-port=0"},
		hasErr: false,
	}}
	for i, tt := range tests {
		buf := new(bytes.Buffer)
		root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, server.NewFakeHTTPServer())
		root.SetOut(buf)
		root.SetArgs(append(tt.args, "--dry-run"))
		err := root.Execute()
		if tt.hasErr {
			assert.Error(t, err, "should have error in case[%d] %q", i, tt.name)
		} else {
			assert.NoError(t, err, "should not have error in case[%d] %q", i, tt.name)
		}
	}
}

func TestStartPlugins(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "atest")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	err = os.WriteFile(filepath.Join(dir, "stores.yaml"), []byte(sampleStores), 0644)
	assert.NoError(t, err)

	t.Run("dry-run", func(t *testing.T) {
		rootCmd := &cobra.Command{
			Use: "atest",
		}
		rootCmd.SetOut(io.Discard)
		rootCmd.AddCommand(createServerCmd(
			fakeruntime.FakeExecer{ExpectOS: "linux", ExpectLookPathError: errors.New("not-found")},
			server.NewFakeHTTPServer(),
		))

		rootCmd.SetArgs([]string{"server", "--config-dir", dir, "--dry-run", "--port=0", "--http-port=0"})
		err = rootCmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		httpServer := server.NewDefaultHTTPServer()
		rootCmd := &cobra.Command{
			Use: "atest",
		}
		rootCmd.SetOut(io.Discard)
		rootCmd.AddCommand(createServerCmd(
			fakeruntime.FakeExecer{ExpectOS: "linux", ExpectLookPathError: errors.New("not-found")},
			httpServer,
		))

		rootCmd.SetArgs([]string{"server", "--config-dir", dir, "--port=0", "--http-port=0"})
		go func() {
			err = rootCmd.Execute()
			assert.NoError(t, err)
		}()

		for httpServer.GetPort() == "" {
			time.Sleep(time.Second)
		}

		defer func() {
			httpServer.Shutdown(context.Background())
		}()
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/suites", httpServer.GetPort()))
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}

		resp, err = http.Get(fmt.Sprintf("http://localhost:%s/metrics", httpServer.GetPort()))
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	})
}

type fakeResponseWriter struct {
	buf    *bytes.Buffer
	header http.Header
}

func newFakeResponseWriter() *fakeResponseWriter {
	return &fakeResponseWriter{
		buf:    new(bytes.Buffer),
		header: make(http.Header),
	}
}

func (w *fakeResponseWriter) Header() http.Header {
	return w.header
}
func (w *fakeResponseWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}
func (w *fakeResponseWriter) WriteHeader(int) {
	// do nothing due to this is a fake response writer
}
func (w *fakeResponseWriter) GetBody() *bytes.Buffer {
	return w.buf
}

//go:embed testdata/stores.yaml
var sampleStores string
