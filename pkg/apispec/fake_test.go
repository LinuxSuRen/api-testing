package apispec_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/stretchr/testify/assert"
)

func TestFakeAPISpec(t *testing.T) {
	tests := []struct {
		name         string
		data         [][]string
		path, method string
		expectExist  bool
		expectCount  int
	}{{
		name: "normal",
		data: [][]string{{
			"/api", "get",
		}},
		path:        "/api",
		method:      "get",
		expectExist: true,
		expectCount: 1,
	}, {
		name:        "empty",
		data:        [][]string{},
		path:        "/api",
		method:      "post",
		expectExist: false,
		expectCount: 0,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coverage := apispec.NewFakeAPISpec(tt.data)
			exist := coverage.HaveAPI(tt.path, tt.method)
			count := coverage.APICount()
			assert.Equal(t, tt.expectExist, exist)
			assert.Equal(t, tt.expectCount, count)
		})
	}
}
