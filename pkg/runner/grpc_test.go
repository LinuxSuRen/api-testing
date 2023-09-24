/*
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

package runner

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/h2non/gock"
	testsrv "github.com/linuxsuren/api-testing/pkg/runner/grpc_test"
	atest "github.com/linuxsuren/api-testing/pkg/testing"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var cache sync.Map

type testUnit struct {
	name     string
	execer   fakeruntime.Execer
	testCase *atest.TestCase
	desc     *atest.GRPCDesc
	ctx      any
	prepare  func()
	verify   func(t *testing.T, output any, err error)
}

func TestGRPCTestCase(t *testing.T) {
	s := grpc.NewServer()
	testServer := &testsrv.TestServer{}
	testsrv.RegisterMainServer(s, testServer)

	l := runServer(t, s)

	doGRPCTest(t, l, &atest.GRPCDesc{
		ImportPath: []string{"grpc_test"},
		ProtoFile:  "test.proto",
	})
	s.Stop()
}

func TestGRPCProtoSetTestCase(t *testing.T) {
	s := grpc.NewServer()
	testServer := &testsrv.TestServer{}
	testsrv.RegisterMainServer(s, testServer)

	l := runServer(t, s)

	doGRPCTest(t, l, &atest.GRPCDesc{
		ProtoSet: "grpc_test/test.pb",
	})
	s.Stop()
}

func TestGRPCReflectTestCase(t *testing.T) {
	s := grpc.NewServer()
	testServer := &testsrv.TestServer{}
	testsrv.RegisterMainServer(s, testServer)
	reflection.RegisterV1(s)

	l := runServer(t, s)

	doGRPCTest(t, l, &atest.GRPCDesc{
		ServerReflection: true,
	})
	s.Stop()
}

func TestGRPCTestError(t *testing.T) {
	s := grpc.NewServer()
	testServer := &testsrv.TestServer{}
	testsrv.RegisterMainServer(s, testServer)

	l := runServer(t, s)

	tests := []testUnit{
		{
			name:   "test proto not found",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/Unary",
					Body: "{}",
				},
			},
			ctx: nil,
			desc: &atest.GRPCDesc{
				ProtoFile: "unknown",
			},
			prepare: func() {
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:   "test proto set found",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/Unary",
					Body: "{}",
				},
			},
			ctx: nil,
			desc: &atest.GRPCDesc{
				ProtoSet: "unknown",
			},
			prepare: func() {
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:   "test reflect on unsupported server",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/Unary",
					Body: "{}",
				},
			},
			ctx: nil,
			desc: &atest.GRPCDesc{
				ServerReflection: true,
			},
			prepare: func() {
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:   "test server is closed",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/Unary",
					Body: "{}",
				},
				Expect: atest.Response{
					Body: getJSONOrCache("unary", &testsrv.HelloReply{
						Message: "Hello!",
					}),
				},
			},
			ctx: nil,
			prepare: func() {
				s.Stop()
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	runUnits(tests, t, l, &atest.GRPCDesc{})
}

func runServer(t *testing.T, s *grpc.Server) net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err, "Listen port")

	log.Println("listening at", l.Addr().(*net.TCPAddr).Port)
	go s.Serve(l)
	return l
}

func doGRPCTest(t *testing.T, l net.Listener, desc *atest.GRPCDesc) {
	tests := []testUnit{
		{
			name:   "test unary rpc",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/Unary",
					Body: "{}",
				},
				Expect: atest.Response{
					Body: getJSONOrCache("unary", &testsrv.HelloReply{
						Message: "Hello!",
					}),
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "test client stream rpc",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API: "/grpctest.Main/ClientStream",
					Body: getJSONOrCache("stream", []*testsrv.StreamMessage{
						{MsgID: 1, ExpectLen: 3},
						{MsgID: 2, ExpectLen: 3},
						{MsgID: 3, ExpectLen: 3},
					}),
				},
				Expect: atest.Response{
					Body: getJSONOrCache("streamRepeted", &testsrv.StreamMessageRepeated{
						Data: []*testsrv.StreamMessage{
							{MsgID: 1, ExpectLen: 3},
							{MsgID: 2, ExpectLen: 3},
							{MsgID: 3, ExpectLen: 3},
						},
					}),
					Verify: []string{`len(data[0].data) == 3`},
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "test server stream rpc",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/ServerStream",
					Body: getJSONOrCache("streamRepeted", nil),
				},
				Expect: atest.Response{
					Body:   getJSONOrCache("stream", nil),
					Verify: []string{`len(data) == 3`},
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "test bid stream rpc",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/BidStream",
					Body: getJSONOrCache("stream", nil),
				},
				Expect: atest.Response{
					Body: getJSONOrCache("stream", nil),
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "test basic type",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API: "/grpctest.Main/TestBasicType",
					Body: getJSONOrCache("basic",
						&testsrv.BasicType{
							Int32:   rand.Int31(),
							Int64:   rand.Int63(),
							Uint32:  rand.Uint32(),
							Uint64:  rand.Uint64(),
							Float32: rand.Float32(),
							Float64: rand.Float64(),
							String_: time.Now().Format(time.RFC3339),
							Bool:    true,
						}),
				},
				Expect: atest.Response{
					Body: getJSONOrCache("basic", nil),
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "test advanced type",
			testCase: &atest.TestCase{
				Request: atest.Request{
					API: "/grpctest.Main/TestAdvancedType",
					Body: getJSONOrCache("advanced",
						&testsrv.AdvancedType{
							Int32Array:   []int32{rand.Int31(), rand.Int31()},
							Int64Array:   []int64{rand.Int63(), rand.Int63()},
							Uint32Array:  []uint32{rand.Uint32(), rand.Uint32()},
							Uint64Array:  []uint64{rand.Uint64(), rand.Uint64()},
							Float32Array: []float32{rand.Float32(), rand.Float32()},
							Float64Array: []float64{rand.NormFloat64(), rand.NormFloat64()},
							StringArray:  []string{time.Now().Format(time.RFC3339), time.Now().Format(time.RFC822)},
							BoolArray:    []bool{true, false},
							HelloReplyMap: map[string]*testsrv.HelloReply{"key": {
								Message: "Hello",
							}},
						}),
				},
				Expect: atest.Response{
					Body: getJSONOrCache("advanced", nil),
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name:   "test unknown rpc",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/UnknownName",
					Body: "{}",
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), "not found", "UnknownName")
			},
		},
		{
			name:   "test invalid input",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/TestBasicType",
					Body: "{",
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name:   "test wrong input type",
			execer: nil,
			testCase: &atest.TestCase{
				Request: atest.Request{
					API:  "/grpctest.Main/TestBasicType",
					Body: getJSONOrCache("unary", nil),
				},
			},
			verify: func(t *testing.T, output any, err error) {
				assert.NotNil(t, err)
			},
		},
	}

	runUnits(tests, t, l, desc)
}

func runUnits(tests []testUnit, t *testing.T, l net.Listener, desc *atest.GRPCDesc) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Clean()
			if tt.prepare != nil {
				tt.prepare()
			}
			if tt.verify == nil {
				tt.verify = hasError
			}
			if tt.desc != nil {
				desc = tt.desc
			}
			runner := NewGRPCTestCaseRunner(l.Addr().String(), *desc)
			runner.WithOutputWriter(os.Stdout)
			if tt.execer != nil {
				runner.WithExecer(tt.execer)
			}
			tt.testCase.Request.API = l.Addr().String() + tt.testCase.Request.API
			output, err := runner.RunTestCase(tt.testCase, tt.ctx, context.TODO())
			tt.verify(t, output, err)

			getter, ok := runner.(ResponseRecord)
			assert.True(t, ok)
			assert.NotNil(t, getter.GetResponseRecord())
		})
	}
}

func TestAPINameMatch(t *testing.T) {
	qn, err := splitFullQualifiedName("127.0.0.1:7070/server.Runner/GetVersion")
	assert.NoError(t, err)
	assert.Equal(t,
		protoreflect.FullName("server.Runner.GetVersion"),
		qn,
		"match full qualified name",
	)

	qn, err = splitFullQualifiedName("127.0.0.1:7070/server.v1.service/method")
	assert.NoError(t, err)
	assert.Equal(t,
		protoreflect.FullName("server.v1.service.method"),
		qn,
		"match full qualified name long",
	)

	qn, err = splitFullQualifiedName("127.0.0.1:7070//server.Runner/GetVersion")
	assert.NotNil(t,
		err,
		"unexpect leading character",
	)

	qn, err = splitFullQualifiedName("127.0.0.1:7070/server.Runner/GetVersion/")
	assert.NotNil(t,
		err,
		"unexpect trailing character",
	)
}

func getJSONOrCache(k string, s any) (msg string) {
	v, ok := cache.Load(k)
	if ok {
		return v.(string)
	}
	b, _ := json.Marshal(s)
	msg = string(b)
	cache.Store(k, msg)
	return
}
