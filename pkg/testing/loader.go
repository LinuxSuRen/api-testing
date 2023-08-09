package testing

// Loader is an interface for test cases loader
type Loader interface {
	HasMore() bool
	Load() ([]byte, error)
	Put(string) (err error)
	GetContext() string
	GetCount() int
	Reset()

	Verify() (err error)
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
}
