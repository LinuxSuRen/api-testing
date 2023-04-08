package server_test

import (
	"context"
	"errors"
	"testing"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestUnimplement(t *testing.T) {
	unimplemented := &server.UnimplementedRunnerServer{}
	_, err := unimplemented.Run(context.TODO(), nil)
	assert.NotNil(t, err)

	_, err = unimplemented.GetVersion(context.Background(), nil)
	assert.NotNil(t, err)

	var reply *server.HelloReply
	assert.Empty(t, reply.GetMessage())
	assert.Empty(t, &server.Empty{})

	var task *server.TestTask
	assert.Empty(t, task.GetData())
	assert.Empty(t, task.GetKind())
	assert.Empty(t, task.GetCaseName())
	assert.Empty(t, task.GetLevel())
	assert.Nil(t, task.GetEnv())

	task = &server.TestTask{
		Data:     "data",
		Kind:     "kind",
		CaseName: "casename",
		Level:    "level",
		Env:      map[string]string{},
	}
	assert.Equal(t, "data", task.GetData())
	assert.Equal(t, "kind", task.GetKind())
	assert.Equal(t, "casename", task.GetCaseName())
	assert.Equal(t, "level", task.GetLevel())
	assert.Equal(t, map[string]string{}, task.GetEnv())
}

func TestServer(t *testing.T) {
	client, _ := server.NewFakeClient(context.Background(), "version", nil)
	reply, err := client.GetVersion(context.Background(), &server.Empty{})
	assert.NotNil(t, reply)
	assert.Equal(t, "version", reply.GetMessage())
	assert.Nil(t, err)

	reply, err = client.Run(context.Background(), &server.TestTask{})
	assert.NotNil(t, reply)
	assert.Nil(t, err)

	clientWithErr, _ := server.NewFakeClient(context.Background(), "version", errors.New("fake"))
	reply, err = clientWithErr.GetVersion(context.Background(), &server.Empty{})
	assert.NotNil(t, err)
	assert.Nil(t, reply)

	reply, err = clientWithErr.Run(context.Background(), &server.TestTask{})
	assert.NotNil(t, err)
	assert.Nil(t, reply)
}
