/*
Copyright 2024 API Testing Authors.

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
package server_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestGetProtos(t *testing.T) {
	protos := server.GetProtos()
	assert.Equal(t, 1, len(protos))

	exists := []string{"server.proto"}
	for _, key := range exists {
		content, ok := protos[key]
		assert.True(t, ok)
		assert.NotEmpty(t, content)
	}
}
