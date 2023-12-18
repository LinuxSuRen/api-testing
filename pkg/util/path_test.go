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
package util_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		ok, err := util.PathExists(os.TempDir())
		assert.True(t, ok)
		assert.NoError(t, err)
	})

	t.Run("not exist", func(t *testing.T) {
		ok, err := util.PathExists(path.Join(os.TempDir(), time.Now().String()))
		assert.False(t, ok)
		assert.NoError(t, err)
	})
}
