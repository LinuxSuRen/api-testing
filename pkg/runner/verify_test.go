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

package runner_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestVerify(t *testing.T) {
	t.Run("conditionalVerify", func(t *testing.T) {
		err := runner.Verify(atest.Response{
			ConditionalVerify: []atest.ConditionalVerify{{
				Condition: []string{
					"1 == 1",
					"2 == 2",
				},
				Verify: []string{"1 == 2"},
			}},
		}, nil)
		assert.Error(t, err)

		err = runner.Verify(atest.Response{
			ConditionalVerify: []atest.ConditionalVerify{{
				Condition: []string{"1 != 1"},
				Verify:    []string{"1 == 2"},
			}},
		}, nil)
		assert.NoError(t, err)
	})

	t.Run("verify YAML contentType", func(t *testing.T) {
		assert.Nil(t, runner.NewBodyVerify("fake", nil))
		verfier := runner.NewBodyVerify(util.YAML, nil)
		assert.NotNil(t, verfier)

		obj, err := verfier.Parse([]byte(`name: linuxsuren`))
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{
			"name": "linuxsuren",
		}, obj)
		assert.NoError(t, verfier.Verify(nil))
	})
}
