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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ExtensionDownloader struct {
	OCIDownloader
	os, arch    string
	kind        string
	extFile     string
	imagePrefix string
}

type ExtensionDownloaderOption func(*ExtensionDownloader)

func NewStoreDownloader(opts ...ExtensionDownloaderOption) PlatformAwareOCIDownloader {
	ociDownloader := &ExtensionDownloader{
		OCIDownloader: NewDefaultOCIDownloader(),
	}
	ociDownloader.WithOS(runtime.GOOS)
	ociDownloader.WithArch(runtime.GOARCH)
	ociDownloader.WithImagePrefix("linuxsuren")
	ociDownloader.WithKind("store")
	for _, opt := range opts {
		opt(ociDownloader)
	}
	return ociDownloader
}

func (d *ExtensionDownloader) Download(name, tag, _ string) (reader io.Reader, err error) {
	name = strings.TrimPrefix(name, "atest-store-")
	d.extFile = d.GetTargetFile(name)
	image := fmt.Sprintf("%s/atest-ext-%s-%s", d.imagePrefix, d.kind, name)
	reader, err = d.OCIDownloader.Download(image, tag, d.extFile)
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

func (d *ExtensionDownloader) GetTargetFile(name string) string {
	name = strings.TrimPrefix(name, fmt.Sprintf("atest-%s-", d.kind))
	if d.os == "" {
		d.extFile = fmt.Sprintf("atest-%s-%s.tar.gz", d.kind, name)
	} else {
		d.extFile = fmt.Sprintf("atest-%s-%s_%s_%s/atest-%s-%s", d.kind, name, d.os, d.arch, d.kind, name)
		if d.os == "windows" {
			d.extFile = fmt.Sprintf("%s.exe", d.extFile)
		}
	}
	return d.extFile
}

func WithOS(os string) ExtensionDownloaderOption {
	return func(d *ExtensionDownloader) {
		d.os = os
	}
}

func WithArch(arch string) ExtensionDownloaderOption {
	return func(d *ExtensionDownloader) {
		d.arch = arch
		if d.arch == "amd64" {
			d.arch = "amd64_v1"
		}
	}
}

func WithImagePrefix(imagePrefix string) ExtensionDownloaderOption {
	return func(d *ExtensionDownloader) {
		d.imagePrefix = imagePrefix
	}
}

func WithKind(kind string) ExtensionDownloaderOption {
	return func(d *ExtensionDownloader) {
		d.kind = kind
	}
}

func (d *ExtensionDownloader) WithOS(os string) {
	d.os = os
}

func (d *ExtensionDownloader) WithImagePrefix(imagePrefix string) {
	d.imagePrefix = imagePrefix
}

func (d *ExtensionDownloader) WithArch(arch string) {
	d.arch = arch
	if d.arch == "amd64" {
		d.arch = "amd64_v1"
	}
}

func (d *ExtensionDownloader) WithKind(kind string) {
	d.kind = kind
}
