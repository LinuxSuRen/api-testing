package testing

// Loader is an interface for test cases loader
type Loader interface {
	HasMore() bool
	Load() ([]byte, error)
	Put(string) (err error)
	GetContext() string
	GetCount() int
}
