package cmd

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/linuxsuren/api-testing/extensions/collector/pkg"
	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	c := NewRootCmd()
	assert.NotNil(t, c)
	assert.Equal(t, "atest-collector", c.Use)
}

func TestResponseFilter(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		},
		Request: &http.Request{
			URL: &url.URL{},
		},
	}
	emptyResp := &http.Response{}

	filter := &responseFilter{
		urlFilter: &filter.URLPathFilter{},
		collects:  pkg.NewCollects(),
	}
	filter.filter(emptyResp, nil)
	filter.filter(resp, nil)
}
