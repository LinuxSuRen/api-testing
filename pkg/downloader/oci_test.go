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

package downloader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegistry(t *testing.T) {
	assert.Equal(t, "registry-1.docker.io", getRegistry("alpine"))
	assert.Equal(t, "registry-1.docker.io", getRegistry("library/alpine"))
	assert.Equal(t, "registry-1.docker.io", getRegistry("docker.io/library/alpine"))
	assert.Equal(t, "ghcr.io", getRegistry("ghcr.io/library/alpine"))
}

func TestDetectAuthURL(t *testing.T) {
	t.Run("without registry", func(t *testing.T) {
		authURL, service, err := detectAuthURL("linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://auth.docker.io/token", authURL)
		assert.Equal(t, "registry.docker.io", service)
	})

	t.Run("without docker.io", func(t *testing.T) {
		authURL, service, err := detectAuthURL("docker.io/linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://auth.docker.io/token", authURL)
		assert.Equal(t, "registry.docker.io", service)
	})

	t.Run("without ghcr.io", func(t *testing.T) {
		authURL, service, err := detectAuthURL("ghcr.io/linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://ghcr.io/token", authURL)
		assert.Equal(t, "ghcr.io", service)
	})
}
