package cmd

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/limit"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRunSuite(t *testing.T) {
	tests := []struct {
		name      string
		suiteFile string
		prepare   func()
		hasError  bool
	}{{
		name:      "simple",
		suiteFile: "testdata/simple-suite.yaml",
		prepare: func() {
			gock.New("http://foo").
				Get("/bar").
				Reply(http.StatusOK).
				JSON("{}")
		},
		hasError: false,
	}, {
		name:      "response is not JSON",
		suiteFile: "testdata/simple-suite.yaml",
		prepare: func() {
			gock.New("http://foo").
				Get("/bar").
				Reply(http.StatusOK)
		},
		hasError: true,
	}, {
		name:      "not found file",
		suiteFile: "testdata/fake.yaml",
		prepare:   func() {},
		hasError:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()

			tt.prepare()
			ctx := getDefaultContext()
			opt := newDiskCardRunOption()
			opt.requestTimeout = 30 * time.Second
			opt.limiter = limit.NewDefaultRateLimiter(0, 0)
			stopSingal := make(chan struct{}, 1)

			err := opt.runSuite(tt.suiteFile, ctx, context.TODO(), stopSingal)
			assert.Equal(t, tt.hasError, err != nil, err)
		})
	}
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		prepare func()
		hasErr  bool
	}{{
		name: "status code is not match",
		args: []string{"-p", "testdata/simple-suite.yaml"},
		prepare: func() {
			gock.New("http://foo").Get("/bar")
		},
		hasErr: true,
	}, {
		name:    "file not found",
		args:    []string{"--pattern", "fake"},
		prepare: func() {},
		hasErr:  false,
	}, {
		name: "normal case",
		args: []string{"-p", "testdata/simple-suite.yaml"},
		prepare: func() {
			gock.New("http://foo").Get("/bar").Reply(http.StatusOK).JSON("{}")
		},
		hasErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			tt.prepare()

			root := &cobra.Command{Use: "root"}
			root.AddCommand(CreateRunCommand())

			root.SetArgs(append([]string{"run"}, tt.args...))

			err := root.Execute()
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
