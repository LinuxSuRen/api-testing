/*
Copyright 2024-2025 API Testing Authors.

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
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistry(t *testing.T) {
	assert.Equal(t, DockerHubRegistry, getRegistry("alpine"))
	assert.Equal(t, DockerHubRegistry, getRegistry("library/alpine"))
	assert.Equal(t, DockerHubRegistry, getRegistry("docker.io/library/alpine"))
	assert.Equal(t, "ghcr.io", getRegistry("ghcr.io/library/alpine"))
}

func TestDetectAuthURL(t *testing.T) {
	t.Run("without registry", func(t *testing.T) {
		authURL, service, err := detectAuthURL("https", "linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://auth.docker.io/token", authURL)
		assert.Equal(t, "registry.docker.io", service)
	})

	t.Run("without docker.io", func(t *testing.T) {
		authURL, service, err := detectAuthURL("https", "docker.io/linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://auth.docker.io/token", authURL)
		assert.Equal(t, "registry.docker.io", service)
	})

	t.Run("without ghcr.io", func(t *testing.T) {
		authURL, service, err := detectAuthURL("https", "ghcr.io/linuxsuren/api-testing")
		assert.NoError(t, err)
		assert.Equal(t, "https://ghcr.io/token", authURL)
		assert.Equal(t, "ghcr.io", service)
	})
}

func TestDownload(t *testing.T) {
	server := mock.NewInMemoryServer(context.Background(), 0)
	err := server.Start(mock.NewLocalFileReader("testdata/registry.yaml"), "/v2")
	assert.NoError(t, err)
	defer func() {
		server.Stop()
	}()

	platforms := []string{
		"windows", "linux", "darwin",
	}
	for _, platform := range platforms {
		t.Run(fmt.Sprintf("on %s", platform), func(t *testing.T) {
			d := NewStoreDownloader()
			d.WithRegistry(fmt.Sprintf("127.0.0.1:%s", server.GetPort()))
			d.WithInsecure(true)
			d.WithOS(platform)
			d.WithArch("amd64")
			d.WithBasicAuth("", "")
			d.WithRoundTripper(nil)

			var reader io.Reader
			reader, err = d.Download("git", "", "")
			assert.NoError(t, err)
			assert.NotNil(t, reader)

			// download and verify it
			var tmpDownloadDir string
			tmpDownloadDir, err = os.MkdirTemp(os.TempDir(), "download")
			defer os.RemoveAll(tmpDownloadDir)
			assert.NoError(t, err)

			err = WriteTo(reader, tmpDownloadDir, "fake.txt")
			assert.NoError(t, err)

			var data []byte
			data, err = os.ReadFile(filepath.Join(tmpDownloadDir, "fake.txt"))
			assert.NoError(t, err)
			assert.Equal(t, "fake", string(data))

			assert.NotEmpty(t, d.GetTargetFile("git"))

			t.Run("not found", func(t *testing.T) {
				_, err = d.Download("orm", "", "")
				assert.Error(t, err)
			})
		})
	}
}
