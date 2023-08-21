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

func TestEmptyThenDefault(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		defVal string
		expect string
	}{{
		name:   "empty string",
		val:    "",
		defVal: "abc",
		expect: "abc",
	}, {
		name:   "blank string",
		val:    " ",
		defVal: "abc",
		expect: "abc",
	}, {
		name:   "not empty or blank string",
		val:    "abc",
		defVal: "def",
		expect: "abc",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.EmptyThenDefault(tt.val, tt.defVal)
			assert.Equal(t, tt.expect, result, result)
		})
	}

	assert.Equal(t, 1, util.ZeroThenDefault(0, 1))
	assert.Equal(t, 1, util.ZeroThenDefault(1, 2))
}
