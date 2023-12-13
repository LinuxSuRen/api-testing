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
}

type Writer interface {
	Loader

	ListTestCase(suite string) (testcases []TestCase, err error)
	GetTestCase(suite, name string) (testcase TestCase, err error)
	CreateTestCase(suite string, testcase TestCase) (err error)
	UpdateTestCase(suite string, testcase TestCase) (err error)
	DeleteTestCase(suite, testcase string) (err error)

	ListTestSuite() (suites []TestSuite, err error)
	GetTestSuite(name string, full bool) (suite TestSuite, err error)
	CreateSuite(name, api string) (err error)
	GetSuite(name string) (*TestSuite, string, error)
	UpdateSuite(TestSuite) (err error)
	DeleteSuite(name string) (err error)
	Close()
}
