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
package server

import (
	"testing"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
)

func TestStoreExtManager(t *testing.T) {
	mgr := NewStoreExtManager(fakeruntime.NewDefaultExecer())

	t.Run("not found", func(t *testing.T) {
		err := mgr.Start("fake", "")
		assert.Error(t, err)
	})

	t.Run("exist executable file", func(t *testing.T) {
		err := mgr.Start("go", "")
		assert.NoError(t, err)

		err = mgr.StopAll()
	})
}
