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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/linuxsuren/api-testing/pkg/testing"
	"google.golang.org/protobuf/reflect/protoregistry"
	"trpc.group/trpc-go/trpc-cmdline/descriptor"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"

	"trpc.group/trpc-go/trpc-cmdline/parser"
)

type tRPCTestCaseRunner struct {
	UnimplementedRunner
	host     string
	proto    testing.RPCDesc
	response SimpleResponse
	cc       client.Client
}

func NewTRPCTestCaseRunner(host string, proto testing.RPCDesc, cc client.Client) TestCaseRunner {
	runner := &tRPCTestCaseRunner{
		UnimplementedRunner: NewDefaultUnimplementedRunner(),
		host:                host,
		proto:               proto,
		cc:                  cc,
	}
	return runner
}

func (r *tRPCTestCaseRunner) RunTestCase(testcase *testing.TestCase, dataContext any, ctx context.Context) (output any, err error) {
	r.log.Info("start to run: '%s'\n", testcase.Name)
	record := NewReportRecord()
	defer func(rr *ReportRecord) {
		rr.EndTime = time.Now()
		rr.Error = err
		rr.API = testcase.Request.API
		rr.Method = "tRPC"
		r.testReporter.PutRecord(rr)
	}(record)

	defer func() {
		if err == nil {
			err = runJob(testcase.After, dataContext, output)
		}
	}()

	contextDir := NewContextKeyBuilder().ParentDir().GetContextValueOrEmpty(ctx)
	if err = testcase.Request.Render(dataContext, contextDir); err != nil {
		return
	}

	if err = runJob(testcase.Before, dataContext, nil); err != nil {
		return
	}

	r.log.Info("start to send request to %s\n", testcase.Request.API)

	var service string
	service, md, err := getTRPCMethodDescriptor(r.proto, testcase)
	if err != nil {
		if err == protoregistry.NotFound {
			return nil, fmt.Errorf("API %q is not found", testcase.Request.API)
		}
		return nil, err
	}
	if md == nil {
		return nil, fmt.Errorf("API %q is not found", testcase.Request.API)
	}

	payload := testcase.Request.Body
	resp, err := invokeTRPCRequest(ctx, r.cc, service, md, payload, r.host)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(resp); err == nil {
		record.Body = string(data)
		r.response = SimpleResponse{
			Body: record.Body,
		}
	}

	r.log.Debug("response body: %s\n", record.Body)

	err = Verify(testcase.Expect, map[string]any{
		"data": resp,
	})
	return
}

func (r *tRPCTestCaseRunner) GetSuggestedAPIs(suite *testing.TestSuite, api string) (result []*testing.TestCase, err error) {
	// TODO need to implement
	return
}

func (r *tRPCTestCaseRunner) GetResponseRecord() SimpleResponse {
	return r.response
}

func getTRPCMethodDescriptor(proto testing.RPCDesc, testcase *testing.TestCase) (service string, d *descriptor.RPCDescriptor, err error) {
	opts := []parser.Option{
		parser.WithAliasOn(false),
		parser.WithAPPName(""),
		parser.WithServerName(""),
		parser.WithAliasAsClientRPCName(false),
		parser.WithLanguage("Go"),
		parser.WithRPCOnly(true),
		parser.WithMultiVersion(false),
	}

	if proto.Raw != "" {
		var tempF *os.File
		if tempF, err = os.CreateTemp(os.TempDir(), "proto"); err != nil {
			return
		}
		defer func() {
			_ = os.Remove(tempF.Name())
		}()

		if err = os.WriteFile(tempF.Name(), []byte(proto.Raw), 0644); err != nil {
			return
		}
		proto.ProtoFile = tempF.Name()
	}

	var fd *descriptor.FileDescriptor
	var method string
	service, method = splitServiceAndMethod(testcase.Request.API)
	if fd, err = parser.Parse(
		proto.ProtoFile,
		[]string{},
		0,
		opts...,
	); err == nil {
		for _, svc := range fd.Services {
			if svc.Name == service {
				d = svc.MethodRPC[method]
				break
			}
		}
	}
	return
}

func splitServiceAndMethod(api string) (service, method string) {
	parts := strings.Split(api, "/")
	if len(parts) >= 2 {
		service = parts[len(parts)-2]
		method = parts[len(parts)-1]
	}
	return
}

func invokeTRPCRequest(ctx context.Context, cc client.Client, serviceName string, md *descriptor.RPCDescriptor, payload string, host string) (
	resp map[string]string, err error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	defer codec.PutBackMessage(msg)

	msg.WithClientRPCName(fmt.Sprintf("/%s.%s/%s", md.RequestTypePkgDirective, serviceName, md.Name))
	msg.WithCalleeServiceName(md.RequestTypePkgDirective + "." + serviceName)
	msg.WithCalleeApp("")
	msg.WithCalleeServer("")
	msg.WithCalleeService("")
	msg.WithCalleeMethod("")
	msg.WithSerializationType(codec.SerializationTypeJSON)
	callopts := []client.Option{}
	callopts = append(callopts, client.WithTarget(host))

	ccc := codec.GetClient(trpc.ProtocolName)

	_, err = ccc.Encode(msg, []byte(payload))

	req := map[string]string{}
	if err = json.Unmarshal([]byte(payload), &req); err != nil {
		err = fmt.Errorf("failed to unmarshal payload, error: %v", err)
		return
	}

	resp = make(map[string]string)

	c := cc //client.New()
	err = c.Invoke(ctx, req, &resp, callopts...)
	return
}
