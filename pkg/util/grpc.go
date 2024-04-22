/*
Copyright 2023 API Testing Authors.

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
package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/logging"
)

var (
	utilLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("memory")
)

func LoadProtoFiles(protoFile string) (targetProtoFile string, importPath []string, protoParentDir string, err error) {
	if !strings.HasPrefix(protoFile, "http://") && !strings.HasPrefix(protoFile, "https://") {
		targetProtoFile = protoFile
		return
	}

	var protoURL *url.URL
	if protoURL, err = url.Parse(protoFile); err != nil {
		return
	}

	utilLogger.Info("start to download proto file", "file", protoFile)
	resp, err := GetDefaultCachedHTTPClient().Get(protoFile)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code %d with %q", resp.StatusCode, protoFile)
		return
	}

	var f *os.File
	contentType := resp.Header.Get(ContentType)
	if contentType != ZIP {
		var data []byte
		if data, err = io.ReadAll(resp.Body); err == nil {
			if f, err = os.CreateTemp(os.TempDir(), "proto"); err == nil {
				_, err = f.Write(data)
				targetProtoFile = f.Name()
			}
		}
	} else {
		targetProtoFile = protoURL.Query().Get("file")
		if targetProtoFile == "" {
			err = errors.New("query parameter file is empty")
			return
		}

		attachment := resp.Header.Get(ContentDisposition)
		filename := strings.TrimPrefix(attachment, "attachment; filename=")
		name := strings.TrimSuffix(filename, filepath.Ext(filename))

		parentDir := os.TempDir()
		if f, err = os.CreateTemp(parentDir, filename); err == nil {
			_, err = io.Copy(f, resp.Body)

			protoParentDir = filepath.Join(parentDir, name)
			err = extractFiles(f.Name(), protoParentDir, targetProtoFile)
			if err != nil {
				return
			}
		}
	}
	return
}

func extractFiles(sourceFile, targetDir, filter string) (err error) {
	if sourceFile == "" || targetDir == "" {
		err = errors.New("source or target filename is empty")
		return
	}

	var archive *zip.ReadCloser
	if archive, err = zip.OpenReader(sourceFile); err != nil {
		return
	}
	defer func() {
		_ = archive.Close()
	}()

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}

		targetFilePath := filepath.Join(targetDir, f.Name)
		if err = os.MkdirAll(filepath.Dir(targetFilePath), os.ModePerm); err != nil {
			return
		}

		var targetFile *os.File
		if targetFile, err = os.OpenFile(targetFilePath,
			os.O_CREATE|os.O_RDWR, f.Mode()); err != nil {
			continue
		}

		var fileInArchive io.ReadCloser
		fileInArchive, err = f.Open()
		if err != nil {
			continue
		}

		_, err = io.Copy(targetFile, fileInArchive)
		_ = targetFile.Close()
	}
	return
}
