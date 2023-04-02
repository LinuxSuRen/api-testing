package cmd

import (
	"bytes"
	"strings"
	"testing"

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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			root := NewRootCmd()
			root.SetOut(buf)
			root.SetArgs(tt.args)
			err := root.Execute()
			tt.verify(t, buf, err)
		})
	}
}
