package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
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
		args: []string{"server", "-p=0", "--http-port=0", "--local-storage=./*"},
		verify: func(t *testing.T, buf *bytes.Buffer, err error) {
			assert.Nil(t, err)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			root := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"},
				&fakeGRPCServer{}, server.NewFakeHTTPServer())
			root.SetOut(buf)
			root.SetArgs(tt.args)
			err := root.Execute()
			tt.verify(t, buf, err)
		})
	}
}

func TestFrontEndHandlerWithLocation(t *testing.T) {
	handler := frontEndHandlerWithLocation("testdata")
	const expect404 = "404 page not found\n"

	t.Run("404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, expect404, resp.GetBody().String())
	})

	t.Run("get js", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/assets/index.js", nil)
		assert.NoError(t, err)
		defer func() {
			uiResourceJS = ""
		}()

		resp := newFakeResponseWriter()

		uiResourceJS = "js"
		handler(resp, req, map[string]string{})
		assert.Equal(t, uiResourceJS, resp.GetBody().String())

		fmt.Println(resp.Header())
		assert.Equal(t, "text/javascript; charset=utf-8", resp.Header().Get(util.ContentType))
	})

	t.Run("get css", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/assets/index.css", nil)
		assert.NoError(t, err)

		resp := newFakeResponseWriter()
		handler(resp, req, map[string]string{})
		assert.Equal(t, expect404, resp.GetBody().String())
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
