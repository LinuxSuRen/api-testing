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
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
)

func TestRobotGenerator(t *testing.T) {
	tests := []struct {
		name      string
		testCase  *atest.TestCase
		testSuite *atest.TestSuite
		expect    string
	}{{
		name: "simple",
		testCase: &atest.TestCase{
			Name: "simple",
			Request: atest.Request{
				API: fooForTest,
			},
		},
		expect: simpleRobot,
	}, {
		name: "with header",
		testCase: &atest.TestCase{
			Name: "simple",
			Request: atest.Request{
				API: fooForTest,
				Header: map[string]string{
					"key": "value",
				},
			},
		},
		expect: headerRobot,
	}, {
		name: "test suite",
		testSuite: &atest.TestSuite{
			Items: []atest.TestCase{{
				Name: "one",
				Request: atest.Request{
					API: fooForTest,
					Header: map[string]string{
						"key1": "value1",
					},
				},
			}, {
				Name: "two",
				Request: atest.Request{
					API: fooForTest,
					Header: map[string]string{
						"key2": "value2",
					},
				},
			}},
		},
		expect: suiteRobot,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewRobotGenerator()
			if got, err := g.Generate(tt.testSuite, tt.testCase); err != nil || got != tt.expect {
				t.Errorf("got %q, want %q, error: %v", got, tt.expect, err)
			}
		})
	}
}

//go:embed testdata/simple.robot
var simpleRobot string

//go:embed testdata/with-headers.robot
var headerRobot string

//go:embed testdata/suite.robot
var suiteRobot string
