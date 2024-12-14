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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	data, err := readFile("data/templateUsage.yaml")
	assert.Nil(t, err)
	assert.Equal(t, "data:application/octet-stream;base64,cmFuZEltYWdlOiB8CiAge3sgcmFuZEltYWdlIDEwMCAxMDAgfX0KcmFuZEFzY2lpOiB8CiAge3sgcmFuZEFzY2lpIDUgfX0KcmFuZFBkZjogfAogIHt7IHJhbmRQZGYgImNvbnRlbnQiIH19CnJhbmRaaXA6IHwKICB7eyByYW5kWmlwIDUgfX0K", data)
}
