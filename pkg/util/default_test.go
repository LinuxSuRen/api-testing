package util_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestMakeSureNotNil(t *testing.T) {
	var fun func()
	var mapStruct map[string]string

	assert.NotNil(t, util.MakeSureNotNil(fun))
	assert.NotNil(t, util.MakeSureNotNil(TestMakeSureNotNil))
	assert.NotNil(t, util.MakeSureNotNil(mapStruct))
	assert.NotNil(t, util.MakeSureNotNil(map[string]string{}))
}
