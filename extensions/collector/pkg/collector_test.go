package pkg_test

import (
	"net/http"
	"testing"

	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	defaultReq, _ := http.NewRequest(http.MethodGet, "http://foo.com", nil)

	tests := []struct {
		name    string
		Request *http.Request
	}{{
		name:    "normal",
		Request: defaultReq,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collects := pkg.NewCollects()
			collects.AddEvent(func(r *http.Request) {
				assert.Equal(t, tt.Request, r)
			})
			for i := 0; i < 10; i++ {
				collects.Add(tt.Request)
			}
			collects.Stop()
		})
	}
}
