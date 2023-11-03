/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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

	"log"
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

	log.Printf("start to download proto file %q\n", protoFile)
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
