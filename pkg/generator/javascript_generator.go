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
package generator

import (
	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/testing"
)

type javascriptGenerator struct {
}

func NewJavaScriptGenerator() CodeGenerator {
	return &javascriptGenerator{}
}

func (g *javascriptGenerator) Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (string, error) {
	return generate(testSuite, testcase, "javascript template", javascriptTemplate)
}

func init() {
	RegisterCodeGenerator("JavaScript", NewJavaScriptGenerator())
}

//go:embed data/main.javascript.tpl
var javascriptTemplate string
