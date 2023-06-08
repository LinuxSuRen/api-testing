package filter_test

import (
	"net/url"
	"testing"

	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	"github.com/stretchr/testify/assert"
)

func TestURLPathFilter(t *testing.T) {
	urlFilter := &filter.URLPathFilter{PathPrefix: []string{"/api/v1", "/api/v2"}}
	assert.True(t, urlFilter.Filter(&url.URL{Path: "/api/v1"}))
	assert.True(t, urlFilter.Filter(&url.URL{Path: "/api/v2"}))
	assert.False(t, urlFilter.Filter(&url.URL{Path: "/api/v3"}))
}
