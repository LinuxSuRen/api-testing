package server

import (
	"context"
	"net/http"
	"testing"

	_ "embed"

	"github.com/h2non/gock"
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

//go:embed testdata/simple.yaml
var simpleSuite string

//go:embed testdata/simple_testcase.yaml
var simpleTestCase string
