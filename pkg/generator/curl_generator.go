/*
Copyright 2023-2025 API Testing Authors.

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
package generator

import (
	"bytes"
	_ "embed"
	"net/http"
	"strings"
	"text/template"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type curlGenerator struct {
}

func NewCurlGenerator() CodeGenerator {
	return &curlGenerator{}
}

func (g *curlGenerator) Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (result string, err error) {
	if testcase.Request.Method == "" {
		testcase.Request.Method = http.MethodGet
	}

	if !strings.HasSuffix(testcase.Request.API, "?") {
		testcase.Request.API += "?"
	}

	queryKeys := testcase.Request.Query.Keys()
	for _, k := range queryKeys {
		testcase.Request.API += k + "=" + testcase.Request.Query.GetValue(k) + "&"
	}

	testcase.Request.API = strings.TrimSuffix(testcase.Request.API, "&")
	testcase.Request.API = strings.TrimSuffix(testcase.Request.API, "?")
	if err = testcase.Request.Render(nil, ""); err != nil {
		return
	}

	var tpl *template.Template
	if tpl, err = template.New("curl template").Parse(curlTemplate); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, testcase); err == nil {
			result = strings.TrimSpace(buf.String())

			result = strings.TrimSuffix(result, " \\")
		}
	}
	return
}

func init() {
	RegisterCodeGenerator("curl", NewCurlGenerator())
}

//go:embed data/curl.tpl
var curlTemplate string
