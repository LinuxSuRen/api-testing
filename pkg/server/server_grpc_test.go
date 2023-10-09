/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
	assert.Empty(t, reply.GetError())
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
	assert.Empty(t, reply.GetError())
	assert.Nil(t, err)

	var testResult *server.TestResult
	testResult, err = client.Run(context.Background(), &server.TestTask{})
	assert.NotNil(t, testResult)
	assert.Nil(t, err)

	reply, err = client.Sample(context.Background(), &server.Empty{})
	assert.Nil(t, err)
	assert.Empty(t, reply.GetMessage())

	clientWithErr, _ := server.NewFakeClient(context.Background(), "version", errors.New("fake"))
	reply, err = clientWithErr.GetVersion(context.Background(), &server.Empty{})
	assert.NotNil(t, err)
	assert.Nil(t, reply)

	testResult, err = clientWithErr.Run(context.Background(), &server.TestTask{})
	assert.NotNil(t, err)
	assert.Nil(t, testResult)
}
