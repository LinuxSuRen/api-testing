package filter_test

import (
	"net/url"
	"testing"

	"github.com/linuxsuren/api-testing/extensions/collector/pkg/filter"
	"github.com/stretchr/testify/assert"
)

func TestURLPathFilter(t *testing.T) {
	urlFilter := &filter.URLPathFilter{PathPrefix: "/api"}
	assert.True(t, urlFilter.Filter(&url.URL{Path: "/api/v1"}))
}
