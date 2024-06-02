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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type storeDownloader struct {
	*defaultOCIDownloader
	os, arch string
	extFile  string
}

func NewStoreDownloader() PlatformAwareOCIDownloader {
	ociDownloader := &storeDownloader{}
	ociDownloader.WithOS(runtime.GOOS)
	ociDownloader.WithArch(runtime.GOARCH)
	return ociDownloader
}

func (d *storeDownloader) Download(name, tag, _ string) (reader io.Reader, err error) {
	name = strings.TrimPrefix(name, "atest-store-")
	d.extFile = fmt.Sprintf("atest-store-%s_%s_%s/atest-store-%s", name, d.os, d.arch, name)
	image := fmt.Sprintf("linuxsuren/atest-ext-store-%s", name)
	reader, err = d.defaultOCIDownloader.Download(image, tag, d.extFile)
	return
}

func WriteTo(reader io.Reader, dir, file string) (err error) {
	var data []byte
	if data, err = io.ReadAll(reader); err == nil {
		if err = os.MkdirAll(dir, 0755); err == nil {
			targetFile := filepath.Join(dir, file)
			err = os.WriteFile(targetFile, data, 0755)
		}
	}
	return
}

func (d *storeDownloader) GetTargetFile() string {
	return d.extFile
}

func (d *storeDownloader) WithOS(os string) {
	d.os = os
}

func (d *storeDownloader) WithArch(arch string) {
	d.arch = arch
	if d.arch == "amd64" {
		d.arch = "amd64_v1"
	}
}
