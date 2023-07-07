package testing

// Loader is an interface for test cases loader
type Loader interface {
	HasMore() bool
	Load() ([]byte, error)
	Put(string) (err error)
	GetContext() string
	GetCount() int
	Reset()
}

type Writer interface {
	Loader

	CreateTestCase(suite string, testcase TestCase) (err error)
	UpdateTestCase(suite string, testcase TestCase) (err error)
	DeleteTestCase(suite, testcase string) (err error)

	CreateSuite(name, api string) (err error)
	GetSuite(name string) (*TestSuite, string, error)
	UpdateSuite(name, api string) (err error)
	DeleteSuite(name string) (err error)
}
