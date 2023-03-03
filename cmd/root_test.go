package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	atesting "github.com/linuxsuren/api-testing/pkg/testing"
)

func Test_setRelativeDir(t *testing.T) {
	type args struct {
		configFile string
		testcase   *atesting.TestCase
	}
	tests := []struct {
		name   string
		args   args
		verify func(*testing.T, *atesting.TestCase)
	}{{
		name: "normal",
		args: args{
			configFile: "a/b/c.yaml",
			testcase: &atesting.TestCase{
				Prepare: atesting.Prepare{
					Kubernetes: []string{"deploy.yaml"},
				},
			},
		},
		verify: func(t *testing.T, testCase *atesting.TestCase) {
			assert.Equal(t, "a/b/deploy.yaml", testCase.Prepare.Kubernetes[0])
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRelativeDir(tt.args.configFile, tt.args.testcase)
			tt.verify(t, tt.args.testcase)
		})
	}
}
