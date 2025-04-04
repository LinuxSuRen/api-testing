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
	"encoding/base64"
	"fmt"
	"os"

	"github.com/linuxsuren/api-testing/pkg/util"
)

func readFile(filename string) (data string, err error) {
	var rawData []byte
	if rawData, err = os.ReadFile(filename); err == nil {
		data = fmt.Sprintf("%s%s", util.BinaryBase64Prefix, base64.StdEncoding.EncodeToString(rawData))
	}
	return
}
