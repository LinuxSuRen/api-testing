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
	"io"
	"path/filepath"
	"runtime"

	"github.com/linuxsuren/api-testing/pkg/downloader"
	"github.com/spf13/cobra"
)

type extensionOption struct {
	ociDownloader downloader.PlatformAwareOCIDownloader
	output        string
	tag           string
	os            string
	arch          string
}

func createExtensionCommand(ociDownloader downloader.PlatformAwareOCIDownloader) (c *cobra.Command) {
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
	flags.StringVarP(&opt.output, "output", "", ".", "The target directory")
	flags.StringVarP(&opt.tag, "tag", "", "", "The extension image tag, try to find the latest one if this is empty")
	flags.StringVarP(&opt.os, "os", "", runtime.GOOS, "The OS")
	flags.StringVarP(&opt.arch, "arch", "", runtime.GOARCH, "The architecture")
	return
}

func (o *extensionOption) runE(cmd *cobra.Command, args []string) (err error) {
	o.ociDownloader.WithOS(o.os)
	o.ociDownloader.WithArch(o.arch)

	for _, arg := range args {
		var reader io.Reader
		if reader, err = o.ociDownloader.Download(arg, o.tag, ""); err != nil {
			return
		} else if reader == nil {
			err = fmt.Errorf("cannot find %s", arg)
			return
		}
		extFile := o.ociDownloader.GetTargetFile()
		cmd.Println("found target file", extFile)

		targetFile := filepath.Base(extFile)
		if err = downloader.WriteTo(reader, o.output, targetFile); err == nil {
			return
		}

		cmd.Println("downloaded", targetFile)
	}
	return
}
