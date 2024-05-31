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
	"github.com/linuxsuren/api-testing/pkg/downloader"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type extensionOption struct {
	ociDownloader downloader.OCIDownloader
	image         string
	tag           string
	os            string
	arch          string
}

func createExtensionCommand(ociDownloader downloader.OCIDownloader) (c *cobra.Command) {
	opt := &extensionOption{
		ociDownloader: ociDownloader,
	}
	c = &cobra.Command{
		Use:   "extension",
		Short: "Manage extension",
		Long:  "Download the store extension file",
		Args:  cobra.MinimumNArgs(1),
		RunE:  opt.runE,
	}
	flags := c.Flags()
	flags.StringVarP(&opt.image, "image", "", "linuxsuren/atest-ext-store", "The image name")
	flags.StringVarP(&opt.tag, "tag", "", "0.0.2", "The image tag")
	flags.StringVarP(&opt.os, "os", "", runtime.GOOS, "The OS")
	flags.StringVarP(&opt.arch, "arch", "", runtime.GOARCH, "The architecture")
	return
}

func (o *extensionOption) runE(cmd *cobra.Command, args []string) (err error) {
	if o.arch == "amd64" {
		o.arch = "amd64_v1"
	}

	for _, arg := range args {
		extFile := fmt.Sprintf("atest-store-%s_%s_%s/atest-store-%s", arg, o.os, o.arch, arg)
		image := fmt.Sprintf("linuxsuren/atest-ext-store-%s", arg)

		var reader io.Reader
		if reader, err = o.ociDownloader.Download(image, o.tag, extFile); err != nil {
			return
		}
		cmd.Println("found target file")

		if reader == nil {
			err = fmt.Errorf("cannot find %s", arg)
			return
		}

		var data []byte
		if data, err = io.ReadAll(reader); err != nil {
			return
		}

		targetFile := filepath.Base(extFile)
		if err = os.WriteFile(targetFile, data, 0755); err != nil {
			return
		}

		cmd.Println("downloaded", targetFile)
	}
	return
}
