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

package testing_test

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestNonLoader(t *testing.T) {
	loader := atest.NewNonWriter()
	defer loader.Close()

	assert.False(t, loader.HasMore())

	data, err := loader.Load()
	assert.NoError(t, err)
	assert.Nil(t, data)

	assert.NoError(t, loader.Put(""))
	assert.Empty(t, loader.GetContext())
	assert.Equal(t, 0, loader.GetCount())

	loader.Reset()

	var suites []atest.TestSuite
	suites, err = loader.ListTestSuite()
	assert.NoError(t, err)
	assert.Empty(t, suites)

	_, err = loader.GetTestSuite("test", false)
	assert.NoError(t, err)

	assert.NoError(t, loader.CreateSuite("", ""))

	var suite *atest.TestSuite
	var absPath string
	suite, absPath, err = loader.GetSuite("test")
	assert.NoError(t, err)
	assert.Nil(t, suite)
	assert.Empty(t, absPath)

	assert.NoError(t, loader.UpdateSuite(atest.TestSuite{}))
	assert.NoError(t, loader.DeleteSuite(""))

	var testCases []atest.TestCase
	testCases, err = loader.ListTestCase("")
	assert.NoError(t, err)
	assert.Empty(t, testCases)

	data, err = loader.GetTestSuiteYaml("")
	assert.NoError(t, err)
	assert.Nil(t, data)

	_, err = loader.GetTestCase("", "")
	assert.NoError(t, err)
	assert.NoError(t, loader.CreateTestCase("", atest.TestCase{}))
	assert.NoError(t, loader.UpdateTestCase("", atest.TestCase{}))
	assert.NoError(t, loader.DeleteTestCase("", ""))
	assert.NoError(t, loader.DeleteAllHistoryTestCase("", ""))
	assert.NoError(t, loader.DeleteHistoryTestCase(""))
	assert.NoError(t, loader.RenameTestSuite("", ""))
	assert.NoError(t, loader.RenameTestCase("", "", ""))

	var readonly bool
	readonly, _, err = loader.Verify()
	assert.NoError(t, err)
	assert.False(t, readonly)
	assert.Empty(t, loader.PProf(""))
}
