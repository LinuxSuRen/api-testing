package filter

import (
	"fmt"
	"net/url"
	"strings"
)

type URLFilter interface {
	Filter(targetURL *url.URL) bool
}

type URLPathFilter struct {
	PathPrefix string
}

func (f *URLPathFilter) Filter(targetURL *url.URL) bool {
	fmt.Println(targetURL.Path, f.PathPrefix)
	return strings.HasPrefix(targetURL.Path, f.PathPrefix)
}
