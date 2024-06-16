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

package cmd

import (
	"os"
	"testing"

	"github.com/hpcloud/tail"
	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestInitTail(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "testing")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	initTail(tmpFile.Name())
	defer func() {
		if tailClient != nil {
			tailClient.Stop()
		}
	}()

	assert.NotNil(t, tailClient)
}

func TestProcessMockConfigFiles(t *testing.T) {

	var fileContent = "new file changes"

	tailClient = &tail.Tail{
		Lines: make(chan *tail.Line),
	}

	go func() {
		tailClient.Lines <- &tail.Line{Text: fileContent}
		close(tailClient.Lines)
	}()

	reader := processMockConfigFiles()
	data := reader.GetData()
	assert.Equal(t, fileContent, string(data))
	assert.NotNil(t, reader)

	_, ok := reader.(mock.Reader)
	assert.True(t, ok)
}
