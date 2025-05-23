/*
Copyright 2023-2025 API Testing Authors.

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
    "errors"
    "testing"
    "time"

    fakeruntime "github.com/linuxsuren/go-fake-runtime"
    "github.com/stretchr/testify/assert"
)

func TestStoreExtManager(t *testing.T) {
    t.Run("not found", func(t *testing.T) {
        mgr := NewStoreExtManager(&fakeruntime.FakeExecer{
            ExpectLookPathError: errors.New("not found"),
        })
        err := mgr.Start("fake", "")
        assert.Error(t, err)
    })

    t.Run("exist executable file", func(t *testing.T) {
        mgr := NewStoreExtManagerInstance(&fakeruntime.FakeExecer{
            ExpectLookPath: "/usr/local/bin/go",
        })
        err := mgr.Start("go", "")
        assert.NoError(t, err, err)

        time.Sleep(time.Microsecond * 100)
        err = mgr.Start("go", "")
        assert.NoError(t, err)

        err = mgr.StopAll()
        assert.NoError(t, err)
    })
}
