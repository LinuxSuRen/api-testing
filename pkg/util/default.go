/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Package util provides a set of common functions
package util

import (
	"net/http"
	"strings"
)

// MakeSureNotNil makes sure the parameter is not nil
func MakeSureNotNil[T any](inter T) T {
	switch val := any(inter).(type) {
	case func():
		if val == nil {
			val = func() {
				// only making sure this is not nil
			}
			return any(val).(T)
		}
	case map[string]string:
		if val == nil {
			val = map[string]string{}
			return any(val).(T)
		}
	}
	return inter
}

// ZeroThenDefault return the default value if the val is zero
func ZeroThenDefault(val, defVal int) int {
	if val == 0 {
		val = defVal
	}
	return val
}

// EmptyThenDefault return the default value if the val is empty
func EmptyThenDefault(val, defVal string) string {
	if strings.TrimSpace(val) == "" {
		val = defVal
	}
	return val
}

// GetFirstHeaderValue retursn the first value of the header
func GetFirstHeaderValue(header http.Header, key string) (val string) {
	values := header[key]
	if len(values) > 0 {
		val = values[0]
	}
	return
}

// ContentType is the HTTP header key
const (
	ContentType       = "Content-Type"
	MultiPartFormData = "multipart/form-data"
	Form              = "application/x-www-form-urlencoded"
	JSON              = "application/json"
)
