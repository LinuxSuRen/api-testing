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
