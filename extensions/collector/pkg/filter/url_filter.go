package filter

import (
	"net/url"
	"strings"
)

// URLFilter is the interface of URL filter
type URLFilter interface {
	Filter(targetURL *url.URL) bool
}

// URLPathFilter filters the URL with path
type URLPathFilter struct {
	PathPrefix []string
}

// Filter implements the URLFilter
func (f *URLPathFilter) Filter(targetURL *url.URL) bool {
	for _, prefix := range f.PathPrefix {
		if strings.HasPrefix(targetURL.Path, prefix) {
			return true
		}
	}
	return false
}
