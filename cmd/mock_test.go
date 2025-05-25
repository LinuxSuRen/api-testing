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

package cmd

import (
	"context"
	"github.com/linuxsuren/api-testing/pkg/server"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestMockCommand(t *testing.T) {
	tt := []struct {
		name   string
		args   []string
		verify func(t *testing.T, err error)
	}{
		{
			name: "mock",
			args: []string{"mock"},
			verify: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "mock with file",
			args: []string{"mock", "testdata/stores.yaml", "--port=0"},
			verify: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			root := NewRootCmd(&fakeruntime.FakeExecer{ExpectOS: "linux"}, server.NewFakeHTTPServer())
			root.SetOut(io.Discard)
			root.SetArgs(tc.args)
			ctx, cancel := context.WithCancel(context.TODO())
			root.SetContext(ctx)
			go func() {
				time.Sleep(time.Second * 2)
				cancel()
			}()
			err := root.Execute()
			tc.verify(t, err)
		})
	}
}
