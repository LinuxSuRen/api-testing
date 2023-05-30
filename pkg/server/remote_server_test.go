package server

import (
	"context"
	"net/http"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/sample"
	"github.com/stretchr/testify/assert"
)

const (
	BearerToken = "Bearer {{.login.data.access_token}}"
)

func TestRemoteServer(t *testing.T) {
	server := NewRemoteServer()
	_, err := server.Run(context.TODO(), &TestTask{
		Kind: "fake",
	})
	assert.NotNil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind: "suite",
		Data: simpleSuite,
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind: "testcase",
		Data: simpleTestCase,
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind:     "testcaseInSuite",
		Data:     simpleSuite,
		CaseName: "get",
	})
	assert.Nil(t, err)

	gock.New(urlFoo).Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind:     "testcaseInSuite",
		Data:     simpleSuite,
		CaseName: "fake",
		Env: map[string]string{
			"SERVER": "http://localhost:9090",
		},
	})
	assert.NotNil(t, err)

	var ver *HelloReply
	ver, err = server.GetVersion(context.TODO(), &Empty{})
	assert.Empty(t, ver.Message)
	assert.Nil(t, err)

	ver, err = server.Sample(context.TODO(), &Empty{})
	assert.Nil(t, err)
	assert.Equal(t, sample.TestSuiteGitLab, ver.Message)
}

func TestFindParentTestCases(t *testing.T) {
	tests := []struct {
		name     string
		testcase *atesting.TestCase
		suite    *atesting.TestSuite
		expect   []atesting.TestCase
	}{{
		name: "normal",
		testcase: &atesting.TestCase{
			Request: atesting.Request{
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name: "body",
		testcase: &atesting.TestCase{
			Request: atesting.Request{
				Body: `{{.login.data}}`,
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atesting.Request{
					Body: `{{.login.data}}`,
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name:     "empty cases",
		testcase: &atesting.TestCase{},
		suite:    &atesting.TestSuite{},
	}, {
		name: "complex",
		testcase: &atesting.TestCase{
			Name: "user",
			Request: atesting.Request{
				API: "/users/{{(index .login 0).name}}",
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "user",
				Request: atesting.Request{
					API: "/users/{{(index .login 0).name}}",
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}},
	}, {
		name: "nest dep",
		testcase: &atesting.TestCase{
			Name: "user",
			Request: atesting.Request{
				API: "/users/{{(index .users 0).name}}{{randomKubernetesName}}",
				Header: map[string]string{
					"Authorization": BearerToken,
				},
			},
		},
		suite: &atesting.TestSuite{
			Items: []atesting.TestCase{{
				Name: "login",
			}, {
				Name: "users",
				Request: atesting.Request{
					API: "/users",
				},
			}, {
				Name: "user",
				Request: atesting.Request{
					API: "/users/{{(index .users 0).name}}",
					Header: map[string]string{
						"Authorization": BearerToken,
					},
				},
			}},
		},
		expect: []atesting.TestCase{{
			Name: "login",
		}, {
			Name: "users",
			Request: atesting.Request{
				API: "/users",
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findParentTestCases(tt.testcase, tt.suite)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUniqueSlice(t *testing.T) {
	uniqueSlice := new(UniqueSlice[string])
	uniqueSlice.Push("a").Push("a").Push("b")
	assert.Equal(t, []string{"a", "b"}, uniqueSlice.GetAll())
}

func TestWithDefaultValue(t *testing.T) {
	assert.Equal(t, withDefaultValue("a", "b"), "a")
	assert.Equal(t, withDefaultValue("", "b"), "b")
	assert.Equal(t, withDefaultValue(nil, map[string]string{"key": "val"}), map[string]string{"key": "val"})
	assert.Equal(t, withDefaultValue(map[string]string{"key": "val"}, map[string]string{"key": "value"}), map[string]string{"key": "val"})
}

//go:embed testdata/simple.yaml
var simpleSuite string

//go:embed testdata/simple_testcase.yaml
var simpleTestCase string

const urlFoo = "http://foo"
