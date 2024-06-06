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

package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/downloader"
	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestExtensionCmd(t *testing.T) {
	t.Run("minimum one arg", func(t *testing.T) {
		command := createExtensionCommand(nil)
		err := command.Execute()
		assert.Error(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		d := downloader.NewStoreDownloader()
		server := mock.NewInMemoryServer(0)

		err := server.Start(mock.NewLocalFileReader("../pkg/downloader/testdata/registry.yaml"), "/v2")
		assert.NoError(t, err)
		defer func() {
			server.Stop()
		}()

		registry := fmt.Sprintf("127.0.0.1:%s", server.GetPort())
		d.WithRegistry(registry)
		d.WithInsecure(true)
		d.WithBasicAuth("", "")
		d.WithOS("linux")
		d.WithArch("amd64")
		d.WithRoundTripper(nil)

		var tmpDownloadDir string
		tmpDownloadDir, err = os.MkdirTemp(os.TempDir(), "download")
		defer os.RemoveAll(tmpDownloadDir)
		assert.NoError(t, err)

		command := createExtensionCommand(d)
		command.SetArgs([]string{"git", "--output", tmpDownloadDir, "--registry", registry})
		err = command.Execute()
		assert.NoError(t, err)

		// not found
		command.SetArgs([]string{"orm", "--output", tmpDownloadDir, "--registry", registry})
		err = command.Execute()
		assert.Error(t, err)
	})
}
