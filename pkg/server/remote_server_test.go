package server

import (
	"context"
	"net/http"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
	atesting "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestRemoteServer(t *testing.T) {
	server := NewRemoteServer()
	_, err := server.Run(context.TODO(), &TestTask{
		Kind: "fake",
	})
	assert.NotNil(t, err)

	gock.New("http://foo").Get("/").Reply(http.StatusOK).JSON(&server)
	gock.New("http://foo").Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind: "suite",
		Data: simpleSuite,
	})
	assert.Nil(t, err)

	gock.New("http://bar").Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind: "testcase",
		Data: simpleTestCase,
	})
	assert.Nil(t, err)

	gock.New("http://foo").Get("/").Reply(http.StatusOK).JSON(&server)
	_, err = server.Run(context.TODO(), &TestTask{
		Kind:     "testcaseInSuite",
		Data:     simpleSuite,
		CaseName: "get",
	})
	assert.Nil(t, err)

	gock.New("http://foo").Get("/").Reply(http.StatusOK).JSON(&server)
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
					"Authorization": "Bearer {{.login.data.access_token}}",
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
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findParentTestCases(tt.testcase, tt.suite)
			assert.Equal(t, tt.expect, result)
		})
	}
}

//go:embed testdata/simple.yaml
var simpleSuite string

//go:embed testdata/simple_testcase.yaml
var simpleTestCase string
