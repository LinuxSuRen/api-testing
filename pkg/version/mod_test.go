/*
Copyright 2024 API Testing Authors.

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
package version

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetModVersion(t *testing.T) {
	t.Run("empty mod", func(t *testing.T) {
		SetMod("")
		_, err := GetModVersion("", "")
		assert.NoError(t, err)
	})

	t.Run("a simple mod", func(t *testing.T) {
		ver, err := GetModVersion("github.com/a/b", simpleMod)
		assert.NoError(t, err)
		assert.Equal(t, "v0.0.1", ver)
	})

	t.Run("not found in mod", func(t *testing.T) {
		ver, err := GetModVersion("github.com/a/b/c", simpleMod)
		assert.NoError(t, err)
		assert.Equal(t, "", ver)
	})

	t.Run("invalid mod", func(t *testing.T) {
		_, err := GetModVersion("github.com/a/b", `invalid`)
		assert.Error(t, err)
	})
}

//go:embed testdata/go.mod.txt
var simpleMod string
