/*
Copyright 2025 API Testing Authors.

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

package home

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGetUserBinDir(t *testing.T) {
    assert.Contains(t, GetUserConfigDir(), "atest")
    assert.Contains(t, GetUserBinDir(), "bin")
    assert.Contains(t, GetUserDataDir(), "data")
    assert.Contains(t, GetExtensionSocketPath("fake"), "fake.sock")
}
