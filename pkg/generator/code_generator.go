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
	"fmt"
	"net/http"
	"text/template"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

// CodeGenerator is the interface of code generator
type CodeGenerator interface {
	Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (result string, err error)
}

var codeGenerators = map[string]CodeGenerator{}

func GetCodeGenerator(name string) CodeGenerator {
	return codeGenerators[name]
}

func RegisterCodeGenerator(name string, generator CodeGenerator) {
	codeGenerators[name] = generator
}

func GetCodeGenerators() (result map[string]CodeGenerator) {
	// returns an immutable map
	result = make(map[string]CodeGenerator, len(codeGenerators))
	for k, v := range codeGenerators {
		result[k] = v
	}
	return
}

// TestSuiteConverter is the interface of test suite converter
type TestSuiteConverter interface {
	Convert(*testing.TestSuite) (result string, err error)
}

var converters = map[string]TestSuiteConverter{}

func GetTestSuiteConverter(name string) TestSuiteConverter {
	return converters[name]
}

func RegisterTestSuiteConverter(name string, converter TestSuiteConverter) {
	converters[name] = converter
}

func GetTestSuiteConverters() (result map[string]TestSuiteConverter) {
	// returns an immutable map
	result = make(map[string]TestSuiteConverter, len(converters))
	for k, v := range converters {
		result[k] = v
	}
	return
}

func generate(testsuite *testing.TestSuite, testcase *testing.TestCase, templateName, templateText string) (result string, err error) {
	if testcase != nil && testcase.Request.Method == "" {
		testcase.Request.Method = http.MethodGet
	}
	if testsuite != nil && testsuite.Items != nil {
		for i, _ := range testsuite.Items {
			if testsuite.Items[i].Request.Method == "" {
				testsuite.Items[i].Request.Method = http.MethodGet
			}
		}
	}
	var tpl *template.Template
	if tpl, err = template.New(templateName).
		Funcs(template.FuncMap{"safeString": safeString}).
		Parse(templateText); err == nil {
		buf := new(bytes.Buffer)
		var ctx interface{}
		if testcase == nil {
			ctx = testsuite
		} else {
			ctx = testcase
		}

		if err = tpl.Execute(buf, ctx); err == nil {
			result = buf.String()
		}
	} else {
		fmt.Println(err)
	}
	return
}
