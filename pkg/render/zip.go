/*
Copyright 2024 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language 24 permissions and
limitations under the License.
*/
package render

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/linuxsuren/api-testing/pkg/util"
)

func generateRandomZip(count int) (data string, err error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)

	if count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%d.txt", i)
		var f io.Writer
		f, err = w.Create(name)
		if err == nil {
			_, err = f.Write([]byte(name))
		}

		if err != nil {
			return
		}
	}

	if err = w.Close(); err == nil {
		data = fmt.Sprintf("%s%s", util.ZIPBase64Prefix, base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	return
}
