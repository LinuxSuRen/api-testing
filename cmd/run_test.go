package cmd

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/linuxsuren/api-testing/pkg/limit"
	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
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
		suiteFile: simpleSuite,
		prepare: func() {
			gock.New(urlFoo).
				Get("/bar").
				Reply(http.StatusOK).
				JSON("{}")
		},
		hasError: false,
	}, {
		name:      "response is not JSON",
		suiteFile: simpleSuite,
		prepare: func() {
			gock.New(urlFoo).
				Get("/bar").
				Reply(http.StatusOK)
		},
		hasError: true,
	}, {
		name:      "not found file",
		suiteFile: "testdata/fake.yaml",
		hasError:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			util.MakeSureNotNil(tt.prepare)()
			ctx := getDefaultContext()
			opt := newDiscardRunOption()
			opt.requestTimeout = 30 * time.Second
			opt.limiter = limit.NewDefaultRateLimiter(0, 0)
			stopSingal := make(chan struct{}, 1)

			err := opt.runSuite(tt.suiteFile, ctx, context.TODO(), stopSingal)
			assert.Equal(t, tt.hasError, err != nil, err)
		})
	}
}

func TestRunCommand(t *testing.T) {
	fooPrepare := func() {
		gock.New(urlFoo).Get("/bar").Reply(http.StatusOK).JSON("{}")
	}
	tmpFile, err := os.CreateTemp(os.TempDir(), "api-testing")
	if !assert.Nil(t, err) {
		return
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	tests := []struct {
		name    string
		args    []string
		prepare func()
		hasErr  bool
	}{{
		name: "status code is not match",
		args: []string{"-p", simpleSuite},
		prepare: func() {
			gock.New(urlFoo).Get("/bar")
		},
		hasErr: true,
	}, {
		name: "file not found",
		args: []string{"--pattern", "fake"},
	}, {
		name:    "normal case",
		args:    []string{"-p", simpleSuite},
		prepare: fooPrepare,
	}, {
		name:    "report ignore",
		args:    []string{"-p", simpleSuite, "--report-ignore"},
		prepare: fooPrepare,
	}, {
		name: "specify a test case",
		args: []string{"-p", simpleSuite, "fake"},
	}, {
		name:   "invalid api",
		args:   []string{"-p", "testdata/invalid-api.yaml"},
		hasErr: true,
	}, {
		name:   "invalid schema",
		args:   []string{"-p", "testdata/invalid-schema.yaml"},
		hasErr: true,
	}, {
		name:    "report file",
		prepare: fooPrepare,
		args:    []string{"-p", simpleSuite, "--report", "md", "--report-file", tmpFile.Name()},
		hasErr:  false,
	}, {
		name: "report with swagger URL",
		prepare: func() {
			fooPrepare()
			fooPrepare()
		},
		args:   []string{"-p", simpleSuite, "--swagger-url", urlFoo + "/bar"},
		hasErr: false,
	}, {
		name:    "report file with error",
		prepare: fooPrepare,
		args:    []string{"-p", simpleSuite, "--report", "md", "--report-file", path.Join(tmpFile.Name(), "fake")},
		hasErr:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			buf := new(bytes.Buffer)
			util.MakeSureNotNil(tt.prepare)()
			root := &cobra.Command{Use: "root"}
			root.SetOut(buf)
			root.AddCommand(createRunCommand())

			root.SetArgs(append([]string{"run"}, tt.args...))

			err := root.Execute()
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestRootCmd(t *testing.T) {
	c := NewRootCmd(fakeruntime.FakeExecer{ExpectOS: "linux"}, NewFakeGRPCServer())
	assert.NotNil(t, c)
	assert.Equal(t, "atest", c.Use)
}

func TestPreRunE(t *testing.T) {
	tests := []struct {
		name   string
		opt    *runOption
		verify func(*testing.T, *runOption, error)
	}{{
		name: "markdown report",
		opt: &runOption{
			report: "markdown",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "md report",
		opt: &runOption{
			report: "md",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "discard report",
		opt: &runOption{
			report: "discard",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "std report",
		opt: &runOption{
			report: "std",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "html report",
		opt: &runOption{
			report: "html",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "empty report",
		opt: &runOption{
			report: "",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.Nil(t, err)
			assert.NotNil(t, ro.reportWriter)
		},
	}, {
		name: "invalid report",
		opt: &runOption{
			report: "fake",
		},
		verify: func(t *testing.T, ro *runOption, err error) {
			assert.NotNil(t, err)
			assert.Nil(t, ro.reportWriter)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cobra.Command{}
			err := tt.opt.preRunE(c, nil)
			tt.verify(t, tt.opt, err)
		})
	}
}

func TestPrinter(t *testing.T) {
	buf := new(bytes.Buffer)
	c := &cobra.Command{}
	c.SetOutput(buf)

	println(c, nil, "foo")
	assert.Empty(t, buf.String())

	println(c, errors.New("bar"), "foo")
	assert.Equal(t, "foo\n", buf.String())
}

const urlFoo = "http://foo"
const simpleSuite = "testdata/simple-suite.yaml"
