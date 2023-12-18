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
	"errors"
	"net/http"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestOkOrErrorMessage(t *testing.T) {
	assert.Equal(t, "OK", util.OKOrErrorMessage(nil))
	assert.Equal(t, "test", util.OKOrErrorMessage(errors.New("test")))
}

func TestIgnoreErrServerClosed(t *testing.T) {
	assert.Nil(t, util.IgnoreErrServerClosed(nil))
	assert.Nil(t, util.IgnoreErrServerClosed(http.ErrServerClosed))
	assert.Error(t, util.IgnoreErrServerClosed(errors.New("test")))
}
