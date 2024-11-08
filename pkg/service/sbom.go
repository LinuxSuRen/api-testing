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

package service

import (
	"net/http"

	ui "github.com/linuxsuren/api-testing/console/atest-ui"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/linuxsuren/api-testing/pkg/version"
)

func SBomHandler(w http.ResponseWriter, r *http.Request,
	params map[string]string) {
	modMap, err := version.GetModVersions("")
	packageJSON := ui.GetPackageJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		data := make(map[string]interface{})
		data["js"] = packageJSON
		data["go"] = modMap
		util.WriteAsJSON(data, w)
	}
}
