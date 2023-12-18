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

package apispec

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed data/proto
var res embed.FS

func GetProtoFiles() (files map[string]string, err error) {
	efs := &res
	files = make(map[string]string)
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		var data []byte
		if data, err = fs.ReadFile(efs, path); err == nil {
			files[strings.TrimPrefix(path, "data/proto/")] = string(data)
		} else {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
