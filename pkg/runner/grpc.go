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
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/linuxsuren/api-testing/pkg/compare"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type gRPCTestCaseRunner struct {
	UnimplementedRunner
	host  string
	proto testing.GRPCDesc
	// fdCache sync.Map
}

var regexFullQualifiedName = regexp.MustCompile(`^([\w\.:]+)\/([\w\.]+)\/(\w+)$`)

func NewGRPCTestCaseRunner(host string, proto testing.GRPCDesc) TestCaseRunner {
	runner := &gRPCTestCaseRunner{
		UnimplementedRunner: NewDefaultUnimplementedRunner(),
		host:                host,
		proto:               proto,
	}
	return runner
}

func (r *gRPCTestCaseRunner) RunTestCase(testcase *testing.TestCase, dataContext any, ctx context.Context) (output any, err error) {
	r.log.Info("start to run: '%s'\n", testcase.Name)
	record := NewReportRecord()
	defer func(rr *ReportRecord) {
		rr.EndTime = time.Now()
		rr.Error = err
		rr.API = testcase.Request.API
		rr.Method = "gRPC"
		r.testReporter.PutRecord(rr)
	}(record)

	defer func() {
		if err == nil {
			err = runJob(testcase.After)
		}
	}()

	contextDir := NewContextKeyBuilder().ParentDir().GetContextValueOrEmpty(ctx)
	if err = testcase.Request.Render(dataContext, contextDir); err != nil {
		return
	}

	if err = runJob(testcase.Before); err != nil {
		return
	}

	r.log.Info("start to send request to %s\n", testcase.Request.API)
	conn, err := grpc.Dial(getHost(testcase.Request.API, r.host), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	md, err := getMethodDescriptor(ctx, r, testcase, conn)
	if err != nil {
		return nil, err
	}

	payload := testcase.Request.Body
	respsStr, err := invokeRequest(ctx, md, payload, conn)
	if err != nil {
		return nil, err
	}

	record.Body = strings.Join(respsStr, ",")
	r.log.Debug("response body: %s\n", record.Body)

	output, err = verifyResponsePayload(testcase.Name, testcase.Expect, respsStr)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func invokeRequest(ctx context.Context, md protoreflect.MethodDescriptor, payload string, conn *grpc.ClientConn) (respones []string, err error) {
	resps := make([]*dynamicpb.Message, 0)
	if md.IsStreamingClient() || md.IsStreamingServer() {
		gpayload := gjson.Parse(payload)
		if !gpayload.IsArray() {
			return nil, fmt.Errorf("payload is not a JSON array")
		}

		reqs := make([]*dynamicpb.Message, len(gpayload.Array()))
		for i, v := range gpayload.Array() {
			req, err := getReqMessagePb(md, v.Raw)
			if err != nil {
				return nil, err
			}
			reqs[i] = req
		}

		resps, err = invokeRPCStream(ctx, conn, md, reqs)
		if err != nil {
			return nil, err
		}
		return buildResponses(resps)
	}
	request, err := getReqMessagePb(md, payload)
	if err != nil {
		return nil, err
	}

	resp, err := invokeRPC(ctx, conn, md, request)
	if err != nil {
		return nil, err
	}
	resps = append(resps, resp)

	return buildResponses(resps)
}

func getReqMessagePb(md protoreflect.MethodDescriptor, message string) (messagepb *dynamicpb.Message, err error) {
	request := dynamicpb.NewMessage(md.Input())
	if message != "" {
		err := protojson.Unmarshal([]byte(message), request)
		if err != nil {
			return nil, err
		}
	}
	return request, nil
}

func buildResponses(resps []*dynamicpb.Message) ([]string, error) {
	respsStr := make([]string, 0)
	for i := range resps {
		respbR, err := protojson.Marshal(resps[i])
		if err != nil {
			return nil, err
		}
		respsStr = append(respsStr, string(respbR))
	}
	return respsStr, nil
}

func getMethodDescriptor(ctx context.Context, r *gRPCTestCaseRunner, testcase *testing.TestCase, conn *grpc.ClientConn) (protoreflect.MethodDescriptor, error) {
	fullname, err := splitFullQualifiedName(testcase.Request.API)
	if err != nil {
		return nil, err
	}

	var dp protoreflect.Descriptor
	// if fd, ok := r.fdCache.Load(fullname.Parent()); ok {
	// 	fmt.Println("hit cache",fullname)
	// 	return getMdFromFd(fd.(protoreflect.FileDescriptor), fullname)
	// }

	if r.proto.ServerReflection {
		dp, err = getByReflect(ctx, r, fullname, conn)
	} else {
		if r.proto.ProtoFile == "" {
			return nil, fmt.Errorf("missing descriptor source")
		}
		dp, err = getByProto(ctx, r, fullname)
	}

	if err != nil {
		return nil, err
	}

	if dp.IsPlaceholder() {
		return nil, fmt.Errorf("%s is not found", fullname)
	}

	if md, ok := dp.(protoreflect.MethodDescriptor); ok {
		return md, nil
	}
	return nil, fmt.Errorf("%s is not a gRPC method", fullname)
}

func getByProto(ctx context.Context, r *gRPCTestCaseRunner, fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(
			&protocompile.SourceResolver{
				ImportPaths: r.proto.ImportPath,
			},
		),
	}

	linker, err := compiler.Compile(ctx, r.proto.ProtoFile)
	if err != nil {
		return nil, err
	}

	dp, err := linker.AsResolver().FindDescriptorByName(fullName)
	if err != nil {
		return nil, err
	}

	// r.fdCache.Store(fullName.Parent(), dp.ParentFile())
	return dp, nil
}

func getByReflect(ctx context.Context, r *gRPCTestCaseRunner, fullName protoreflect.FullName, conn *grpc.ClientConn) (md protoreflect.Descriptor, err error) {
	reflectconn := grpc_reflection_v1.NewServerReflectionClient(conn)
	cli, err := reflectconn.ServerReflectionInfo(ctx)
	if err != nil {
		return nil, err
	}

	req := &grpc_reflection_v1.ServerReflectionRequest{
		Host: "",
		MessageRequest: &grpc_reflection_v1.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: string(fullName),
		},
	}

	err = cli.Send(req)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Recv()
	if err != nil {
		return nil, err
	}
	_ = cli.CloseSend()

	if resp := resp.GetErrorResponse(); resp != nil {
		return nil, fmt.Errorf(resp.GetErrorMessage())
	}

	fdresp := resp.GetFileDescriptorResponse()

	for _, fdb := range fdresp.FileDescriptorProto {
		fdp := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdb, fdp); err != nil {
			return nil, err
		}

		fd, err := protodesc.NewFile(fdp, nil)
		if err != nil {
			return nil, err
		}

		md, err = getMdFromFd(fd, fullName)
		if err == nil {
			// r.fdCache.Store(fullName.Parent(), fd)
			return md, nil
		}
	}

	return nil, fmt.Errorf("%s is not found", fullName)
}

func getMdFromFd(fd protoreflect.FileDescriptor, fullname protoreflect.FullName) (md protoreflect.MethodDescriptor, err error) {
	sd := fd.Services().ByName(fullname.Parent().Name())
	if sd == nil {
		return nil, fmt.Errorf("grpc service %s is not found in proto %s", fullname.Parent().Name(), fd.Name())
	}

	md = sd.Methods().ByName(fullname.Name())
	if md == nil {
		return nil, fmt.Errorf("method %s is not found in service %s", fullname.Name(), fullname.Parent().Name())
	}
	return md, nil
}

func splitFullQualifiedName(api string) (protoreflect.FullName, error) {
	qn := regexFullQualifiedName.FindStringSubmatch(api)
	if len(qn) == 0 {
		return "", fmt.Errorf("%s is not a valid gRPC api name", api)
	}
	fn := protoreflect.FullName(strings.Join(qn[2:], "."))
	if !fn.IsValid() {
		return "", fmt.Errorf("%s is not a valid gRPC api name", api)
	}
	return fn, nil
}

func getHost(api, fallback string) (host string) {
	qn := regexFullQualifiedName.FindStringSubmatch(api)
	if len(qn) == 0 {
		return fallback
	}
	return qn[1]
}

func getMethodName(md protoreflect.MethodDescriptor) string {
	return fmt.Sprintf("/%s/%s", md.Parent().FullName(), md.Name())
}

// invokeRPC sends an unary RPC to gRPC server.
func invokeRPC(ctx context.Context, conn grpc.ClientConnInterface, method protoreflect.MethodDescriptor, request *dynamicpb.Message) (resp *dynamicpb.Message, err error) {
	resp = dynamicpb.NewMessage(method.Output())
	if err := conn.Invoke(ctx, getMethodName(method), request, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// invokeRPCStream combine all three types of streaming rpc into a single function.
func invokeRPCStream(ctx context.Context, conn grpc.ClientConnInterface, method protoreflect.MethodDescriptor, requests []*dynamicpb.Message) (resps []*dynamicpb.Message, err error) {
	sd := &grpc.StreamDesc{
		StreamName:    string(method.Name()),
		ServerStreams: method.IsStreamingServer(),
		ClientStreams: method.IsStreamingClient(),
	}

	s, err := conn.NewStream(ctx, sd, getMethodName(method))
	if err != nil {
		return nil, err
	}

	i := 0

sendLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if len(requests) == i {
				break sendLoop
			}
			if err := s.SendMsg(requests[i]); err != nil {
				return nil, err
			}
			i++
		}
	}

	if err = s.CloseSend(); err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			resp := dynamicpb.NewMessage(method.Output())
			if err = s.RecvMsg(resp); err != nil {
				if err == io.EOF {
					return resps, nil
				}
				return nil, err
			}
			resps = append(resps, resp)
		}
	}
}

func verifyResponsePayload(caseName string, expect testing.Response, jsonPayload []string) (output any, err error) {
	mapOutput := map[string]any{
		"data": jsonPayload,
	}

	if err = payloadFieldsVerify(caseName, expect, jsonPayload); err != nil {
		return
	}

	err = Verify(expect, mapOutput)
	if err != nil {
		return nil, err
	}
	return
}

func payloadFieldsVerify(caseName string, expect testing.Response, jsonPayload []string) error {
	if expect.Body == "" {
		return nil
	}

	if !gjson.Valid(expect.Body) {
		fmt.Printf("expect.Body: %v\n", expect.Body)
		return fmt.Errorf("case %s: expect body is not a valid JSON", caseName)
	}

	exp := gjson.Parse(expect.Body)
	gjsonPayload := make([]gjson.Result, len(jsonPayload))
	for i := range jsonPayload {
		gjsonPayload[i] = gjson.Parse(jsonPayload[i])
	}

	if exp.IsArray() {
		return compare.Array(caseName, exp.Array(), gjsonPayload)
	}

	if exp.IsObject() {
		var msg string
		for i := range jsonPayload {
			err := compare.Object(fmt.Sprintf("%v[%v]", caseName, i),
				exp.Map(), gjsonPayload[i].Map())
			if err != nil {
				msg += err.Error()
			}
		}

		if msg != "" {
			return fmt.Errorf(msg)
		}
		return nil
	}
	return fmt.Errorf("case %s: unknown expect content", caseName)
}
