/*
Copyright 2023 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
		val = strings.Split(val, ";")[0]
	}
	return
}

// ContentType is the HTTP header key
const (
	ContentType        = "Content-Type"
	ContentDisposition = "Content-Disposition"
	MultiPartFormData  = "multipart/form-data"
	Form               = "application/x-www-form-urlencoded"
	JSON               = "application/json"
	OCIImageIndex      = "application/vnd.oci.image.index.v1+json"
	YAML               = "application/yaml"
	ZIP                = "application/zip"
	OctetStream        = "application/octet-stream"
	Plain              = "text/plain"
	Authorization      = "Authorization"
)
