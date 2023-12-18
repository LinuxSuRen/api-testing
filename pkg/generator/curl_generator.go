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

	var tpl *template.Template
	if tpl, err = template.New("curl template").Parse(curlTemplate); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, testcase); err == nil {
			result = strings.TrimSpace(buf.String())

			result = strings.ReplaceAll(result, "\n", " \\\n")
		}
	}
	return
}

func init() {
	RegisterCodeGenerator("curl", NewCurlGenerator())
}

//go:embed data/curl.tpl
var curlTemplate string
