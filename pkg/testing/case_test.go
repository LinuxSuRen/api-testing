package testing_test

import (
	"testing"

	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestInScope(t *testing.T) {
	testCase := &atesting.TestCase{Name: "foo"}
	assert.True(t, testCase.InScope(nil))
	assert.True(t, testCase.InScope([]string{"foo"}))
	assert.False(t, testCase.InScope([]string{"bar"}))
}
