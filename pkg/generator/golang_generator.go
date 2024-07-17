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
package generator

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type golangGenerator struct {
}

func NewGolangGenerator() CodeGenerator {
	return &golangGenerator{}
}

func (g *golangGenerator) Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (result string, err error) {
	if testcase.Request.Method == "" {
		testcase.Request.Method = http.MethodGet
	}
	var tpl *template.Template
	if tpl, err = template.New("golang template").Funcs(template.FuncMap{"safeString": safeString}).Parse(golangTemplate); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, testcase); err == nil {
			result = buf.String()
			// Check if "bytes." is used in the generated code
			if strings.Contains(result, "bytes.") {
				// Only insert "bytes" import if it's not already included
				if !strings.Contains(result, "\"bytes\"") {
					result = strings.Replace(result, "import (", "import (\n\t\"bytes\"", 1)
				}
			}
		}
	}
	return
}

func init() {
	RegisterCodeGenerator("golang", NewGolangGenerator())
}

//go:embed data/main.go.tpl
var golangTemplate string
