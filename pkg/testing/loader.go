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
package testing

// Loader is an interface for test cases loader
type Loader interface {
	HasMore() bool
	Load() ([]byte, error)
	Put(string) (err error)
	GetContext() string
	GetCount() int
	Reset()

	Verify() (readOnly bool, err error)
	PProf(name string) []byte

	Query(query map[string]string) (result DataResult, err error)
}

type DataResult struct {
	Pairs           map[string]string
	Rows            []map[string]string
	Databases       []string
	Tables          []string
	CurrentDatabase string
}

type Writer interface {
	Loader

	ListTestCase(suite string) (testcases []TestCase, err error)
	GetTestCase(suite, name string) (testcase TestCase, err error)
	CreateTestCase(suite string, testcase TestCase) (err error)
	UpdateTestCase(suite string, testcase TestCase) (err error)
	DeleteTestCase(suite, testcase string) (err error)

	ListHistoryTestSuite() (suites []HistoryTestSuite, err error)
	CreateHistoryTestCase(testcaseResult TestCaseResult, suite *TestSuite, historyHeader map[string]string) (err error)
	GetHistoryTestCaseWithResult(id string) (historyTestCase HistoryTestResult, err error)
	GetHistoryTestCase(id string) (historyTestCase HistoryTestCase, err error)
	DeleteHistoryTestCase(id string) (err error)
	DeleteAllHistoryTestCase(suite, name string) (err error)
	RenameTestCase(suite, oldName, newName string) (err error)
	GetTestCaseAllHistory(suite, name string) (historyTestCase []HistoryTestCase, err error)

	ListTestSuite() (suites []TestSuite, err error)
	GetTestSuite(name string, full bool) (suite TestSuite, err error)
	GetTestSuiteYaml(name string) (testSuiteYaml []byte, err error)
	CreateSuite(name, api string) (err error)
	GetSuite(name string) (*TestSuite, string, error)
	UpdateSuite(TestSuite) (err error)
	DeleteSuite(name string) (err error)
	RenameTestSuite(oldName, newName string) error
	Close()
}
