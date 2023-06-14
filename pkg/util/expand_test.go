package util_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestExpand(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{{
		name:   "without brace",
		input:  "/home",
		expect: []string{"/home"},
	}, {
		name:   "with brace",
		input:  "/home/{good,bad}",
		expect: []string{"/home/good", "/home/bad"},
	}, {
		name:   "with brace, have suffix",
		input:  "/home/{good,bad}.yaml",
		expect: []string{"/home/good.yaml", "/home/bad.yaml"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.Expand(tt.input)
			assert.Equal(t, tt.expect, got, got)
		})
	}
}
