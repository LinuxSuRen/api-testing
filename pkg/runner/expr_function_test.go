package runner_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/runner"
	"github.com/stretchr/testify/assert"
)

func TestExprFuncSleep(t *testing.T) {
	tests := []struct {
		name   string
		params []interface{}
		hasErr bool
	}{{
		name:   "string format duration",
		params: []interface{}{"0.01s"},
		hasErr: false,
	}, {
		name:   "without params",
		hasErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := runner.ExprFuncSleep(tt.params...)
			if tt.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
