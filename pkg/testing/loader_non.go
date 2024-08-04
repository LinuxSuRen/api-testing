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

// nonLoader is an empty implement for avoid too many nil check
type nonLoader struct{}

func NewNonWriter() Writer {
	return &nonLoader{}
}

// HasMore returns if there are more test cases
func (l *nonLoader) HasMore() bool {
	return false
}

// Load returns the test case content
func (l *nonLoader) Load() (data []byte, err error) {
	return
}

// Put adds the test case path
func (l *nonLoader) Put(item string) (err error) {
	return
}

// GetContext returns the context of current test case
func (l *nonLoader) GetContext() string {
	return ""
}

// GetCount returns the count of test cases
func (l *nonLoader) GetCount() int {
	return 0
}

// Reset resets the index
func (l *nonLoader) Reset() {
	// non-implement
}

func (l *nonLoader) ListTestSuite() (suites []TestSuite, err error) {
	return
}
func (l *nonLoader) GetTestSuite(name string, full bool) (suite TestSuite, err error) {
	return
}

func (l *nonLoader) CreateSuite(name, api string) (err error) {
	return
}

func (l *nonLoader) GetSuite(name string) (suite *TestSuite, absPath string, err error) {
	return
}

// UpdateSuite updates the suite
func (l *nonLoader) UpdateSuite(suite TestSuite) (err error) {
	return
}

func (l *nonLoader) DeleteSuite(name string) (err error) {
	return
}

func (l *nonLoader) ListTestCase(suite string) (testcases []TestCase, err error) {
	return
}

func (l *nonLoader) GetTestSuiteYaml(suite string) (testSuiteYaml []byte, err error) {
	return
}

func (l *nonLoader) GetTestCase(suite, name string) (testcase TestCase, err error) {
	return
}

func (l *nonLoader) CreateTestCase(suiteName string, testcase TestCase) (err error) {
	return
}

func (l *nonLoader) UpdateTestCase(suite string, testcase TestCase) (err error) {
	return
}

func (l *nonLoader) DeleteTestCase(suiteName, testcase string) (err error) {
	return
}

func (l *nonLoader) CreateHistoryTestCase(testcaseResult TestCaseResult, suiteName *TestSuite) (err error) {
	return
}

func (l *nonLoader) ListHistoryTestSuite()(suites []HistoryTestSuite, err error) {
	return
}

func (l *nonLoader) GetHistoryTestCaseWithResult(id string) (testcase HistoryTestResult, err error) {
	return
}

func (l *nonLoader) GetHistoryTestCase(id string) (testcase HistoryTestCase, err error) {
	return
}
func (l *nonLoader) DeleteHistoryTestCase(id string) (err error) {
	return
}

func (l *nonLoader) Verify() (readOnly bool, err error) {
	// always be okay
	return
}

func (l *nonLoader) PProf(string) []byte {
	// not support
	return nil
}

func (l *nonLoader) Close() {
	// not support
}
