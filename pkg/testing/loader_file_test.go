package testing_test

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestFileLoader(t *testing.T) {
	tests := []struct {
		name   string
		items  []string
		verify func(t *testing.T, loader atest.Loader)
	}{{
		name:  "empty",
		items: []string{},
		verify: func(t *testing.T, loader atest.Loader) {
			assert.False(t, loader.HasMore())
			assert.Equal(t, 0, loader.GetCount())
		},
	}, {
		name:   "brace expansion path",
		items:  []string{"testdata/{invalid-,}testcase.yaml"},
		verify: defaultVerify,
	}, {
		name:   "glob path",
		items:  []string{"testdata/*testcase.yaml"},
		verify: defaultVerify,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loader := atest.NewFileLoader()
			for _, item := range tt.items {
				loader.Put(item)
			}
			tt.verify(t, loader)
		})
	}
}

func defaultVerify(t *testing.T, loader atest.Loader) {
	assert.True(t, loader.HasMore())
	data, err := loader.Load()
	assert.Nil(t, err)
	assert.Equal(t, invalidTestCaseContent, string(data))
	assert.Equal(t, "testdata", loader.GetContext())

	assert.True(t, loader.HasMore())
	data, err = loader.Load()
	assert.Nil(t, err)
	assert.Equal(t, testCaseContent, string(data))
	assert.Equal(t, "testdata", loader.GetContext())

	assert.False(t, loader.HasMore())
	loader.Reset()
	assert.True(t, loader.HasMore())
}
