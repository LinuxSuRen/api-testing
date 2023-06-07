package filter

import (
	"fmt"
	"net/url"
	"strings"
)

// URLFilter is the interface of URL filter
type URLFilter interface {
	Filter(targetURL *url.URL) bool
}

// URLPathFilter filters the URL with path
type URLPathFilter struct {
	PathPrefix string
}

// Filter implements the URLFilter
func (f *URLPathFilter) Filter(targetURL *url.URL) bool {
	fmt.Println(targetURL.Path, f.PathPrefix)
	return strings.HasPrefix(targetURL.Path, f.PathPrefix)
}
