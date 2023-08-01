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
	"html/template"
	"net/http"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type golangGenerator struct {
}

func NewGolangGenerator() CodeGenerator {
	return &golangGenerator{}
}

func (g *golangGenerator) Generate(testcase *testing.TestCase) (result string, err error) {
	if testcase.Request.Method == "" {
		testcase.Request.Method = http.MethodGet
	}
	var tpl *template.Template
	if tpl, err = template.New("golang template").Parse(golangTemplate); err == nil {
		buf := new(bytes.Buffer)
		if err = tpl.Execute(buf, testcase); err == nil {
			result = buf.String()
		}
	}
	return
}

func init() {
	RegisterCodeGenerator("golang", NewGolangGenerator())
}

//go:embed data/main.go.tpl
var golangTemplate string
