package apispec

type fakeAPISpec struct {
	apis [][]string
}

// NewFakeAPISpec creates a new instance of fakeAPISpec
func NewFakeAPISpec(apis [][]string) APIConverage {
	return &fakeAPISpec{apis: apis}
}

// HaveAPI is fake method
func (f *fakeAPISpec) HaveAPI(path, method string) (exist bool) {
	for _, item := range f.apis {
		if len(item) >= 2 && item[0] == path && item[1] == method {
			exist = true
			break
		}
	}
	return
}

// APICount is fake method
func (f *fakeAPISpec) APICount() (count int) {
	count = len(f.apis)
	return
}
