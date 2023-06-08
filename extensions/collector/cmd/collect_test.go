package cmd

import (
	"bytes"
	"context"
	"io"
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
	targetURL, err := url.Parse("http://foo.com/api/v1")
	assert.NoError(t, err)

	resp := &http.Response{
		Header: http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		},
		Request: &http.Request{
			URL: targetURL,
		},
		Body: io.NopCloser(bytes.NewBuffer([]byte("hello"))),
	}
	emptyResp := &http.Response{}

	filter := &responseFilter{
		urlFilter: &filter.URLPathFilter{
			PathPrefix: []string{"/api/v1"},
		},
		collects: pkg.NewCollects(),
		ctx:      context.Background(),
	}
	filter.filter(emptyResp, nil)
	filter.filter(resp, nil)
}
