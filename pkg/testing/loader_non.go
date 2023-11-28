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

func (l *nonLoader) Verify() (readOnly bool, err error) {
	// always be okay
	return
}
