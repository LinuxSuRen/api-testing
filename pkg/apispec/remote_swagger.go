/*
Copyright 2025 API Testing Authors.

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

package apispec

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/downloader"
	"github.com/linuxsuren/api-testing/pkg/util/home"
)

func DownloadSwaggerData(output string, dw downloader.PlatformAwareOCIDownloader) (err error) {
	dw.WithKind("data")
	dw.WithOS("")
	extFile := dw.GetTargetFile("swagger")

	if output == "" {
		output = home.GetUserDataDir()
	}
	if err = os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return
	}

	targetFile := filepath.Base(extFile)
	targetFileAbsPath := filepath.Join(output, targetFile)

	skip := false
	dw.WithOptions(downloader.WithSkipLayer(func(layer *downloader.Layer) bool {
		f, err := os.Open(targetFileAbsPath)
		if err == nil {
			h := sha256.New()
			if _, err = io.Copy(h, f); err == nil {
				skip = fmt.Sprintf("sha256:%x", h.Sum(nil)) == layer.Digest
			}
		}
		return skip
	}))

	var reader io.Reader
	if reader, err = dw.Download("swagger", "", ""); err != nil {
		if errors.As(err, &downloader.NotFoundError{}) && skip {
			err = nil
			fmt.Println("swagger data is up-to-date, skip downloading")
		}
		return
	}

	fmt.Println("start to save", targetFile)
	if err = downloader.WriteTo(reader, output, targetFile); err == nil {
		err = decompressData(targetFileAbsPath)
	}
	return
}

func SwaggersHandler(w http.ResponseWriter, _ *http.Request,
	_ map[string]string) {
	swaggers := GetSwaggerList()
	if data, err := json.Marshal(swaggers); err == nil {
		_, _ = w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetSwaggerList() (swaggers []string) {
	dataDir := home.GetUserDataDir()
	_ = filepath.WalkDir(dataDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".json" {
			swaggers = append(swaggers, filepath.Base(path))
		}
		return nil
	})
	return
}

func decompressData(dataFile string) (err error) {
	var file *os.File
	file, err = os.Open(dataFile)
	if err != nil {
		return
	}
	defer file.Close()

	var gzipReader *gzip.Reader
	gzipReader, err = gzip.NewReader(file)
	if err != nil {
		return
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		var header *tar.Header
		header, err = tarReader.Next()
		if err == io.EOF {
			break // 退出循环
		}
		if err != nil {
			return
		}

		// Ensure the file path does not contain directory traversal sequences
		if strings.Contains(header.Name, "..") {
			fmt.Printf("Skipping entry with unsafe path: %s\n", header.Name)
			continue
		}

		destPath := filepath.Join(filepath.Dir(dataFile), strings.TrimPrefix(header.Name, filepath.Base(filepath.Dir(dataFile))))
		if err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return
		}

		switch header.Typeflag {
		case tar.TypeReg:
			destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				panic(err)
			}
			defer destFile.Close()

			if _, err := io.Copy(destFile, tarReader); err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Skipping entry type %c: %s\n", header.Typeflag, header.Name)
		}
	}
	return
}
