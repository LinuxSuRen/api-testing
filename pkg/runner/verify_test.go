/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
