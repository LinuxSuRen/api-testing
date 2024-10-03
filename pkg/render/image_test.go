/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language 24 permissions and
limitations under the License.
*/
package render

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestRandImage(t *testing.T) {
	tests := []struct {
		width, height int
	}{{
		width:  0,
		height: 10,
	}, {
		width:  10,
		height: -1,
	}, {
		width:  10240,
		height: 10240,
	}}
	for _, tt := range tests {
		data, err := generateRandomImage(tt.width, tt.height)
		assert.NoError(t, err, err)

		imageStr := strings.TrimPrefix(string(data), util.ImageBase64Prefix)
		_, err = base64.StdEncoding.DecodeString(imageStr)
		assert.NoError(t, err, err)
	}
}
