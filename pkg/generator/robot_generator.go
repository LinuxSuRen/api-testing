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

type robotGenerator struct {
}

func NewRobotGenerator() CodeGenerator {
	return &robotGenerator{}
}

func (g *robotGenerator) Generate(testSuite *testing.TestSuite, testcase *testing.TestCase) (string, error) {
	tpl := robotTemplate
	if testcase == nil {
		tpl = robotSuiteTemplate
	}
	return generate(testSuite, testcase, "robot-framework", tpl)
}

func init() {
	RegisterCodeGenerator("robot-framework", NewRobotGenerator())
}

//go:embed data/robot.tpl
var robotTemplate string

//go:embed data/robot-suite.tpl
var robotSuiteTemplate string
