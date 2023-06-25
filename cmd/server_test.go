package cmd

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
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
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	handler(&fakeResponseWriter{buf: buf}, req, map[string]string{})
	assert.Equal(t, "404 page not found\n", buf.String())
}

type fakeResponseWriter struct {
	buf *bytes.Buffer
}

func (w *fakeResponseWriter) Header() http.Header {
	return make(http.Header)
}
func (w *fakeResponseWriter) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}
func (w *fakeResponseWriter) WriteHeader(int) {
}
