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

package compare

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNoEqualErr(t *testing.T) {
	err := newNoEqualErr("data", fmt.Errorf("this is msg"))
	err = newNoEqualErr("to", err)
	err = newNoEqualErr("path", err)
	assert.Equal(t, "compare: field path.to.data: this is msg", err.Error())
}
