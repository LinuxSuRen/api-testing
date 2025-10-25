/*
Copyright 2023-2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/linuxsuren/api-testing/pkg/apispec"
	"github.com/linuxsuren/api-testing/pkg/compare"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/util"

	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	grpcRunnerLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("memory")
)

type gRPCTestCaseRunner struct {
	UnimplementedRunner
	host     string
	proto    testing.RPCDesc
	response SimpleResponse
	// fdCache sync.Map
}

var regexFullQualifiedName = regexp.MustCompile(`^([\w\.:]+)\/([\w\.]+)\/(\w+)$`)
var regexURLPrefix = regexp.MustCompile(`^https?://`)

func NewGRPCTestCaseRunner(host string, proto testing.RPCDesc) TestCaseRunner {
	runner := &gRPCTestCaseRunner{
		UnimplementedRunner: NewDefaultUnimplementedRunner(),
		host:                host,
		proto:               proto,
		response:            SimpleResponse{},
	}
	return runner
}

func init() {
	RegisterRunner("grpc", func(suite *testing.TestSuite) TestCaseRunner {
		if suite.Spec.RPC == nil {
			suite.Spec.RPC = &testing.RPCDesc{}
		}
		return NewGRPCTestCaseRunner(suite.API, *suite.Spec.RPC)
	})
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

	server := getHost(testcase.Request.API, r.host)
	r.log.Info("start to send request to %s\n", server)

	var conn *grpc.ClientConn
	if conn, err = r.getConnection(server); err != nil {
		return
	}
	defer conn.Close()

	md, err := getMethodDescriptor(ctx, r, testcase, conn)
	if err != nil {
		if err == protoregistry.NotFound {
			return nil, fmt.Errorf("api %q is not found", testcase.Request.API)
		}
		return nil, fmt.Errorf("failed to get method descriptor: %v", err)
	}

	// pass the headers into gRPC request metadata
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(testcase.Request.Header))

	payload := testcase.Request.Body
	respsStr, err := invokeRequest(ctx, md, payload.String(), conn)
	if err != nil {
		return nil, err
	}

	if len(respsStr) == 0 {
		record.Body = strings.Join(respsStr, ",")
	} else {
		record.Body = respsStr[0]
	}
	r.response.Body = record.Body
	r.log.Debug("response body: %s\n", record.Body)

	output, err = verifyResponsePayload(md, testcase.Name, testcase.Expect, respsStr)
	if err != nil {
		return nil, err
	}

	if output == nil {
		output, err = NewBodyVerify(util.JSON, nil).Parse([]byte(record.Body))
	}
	return
}

func (r *gRPCTestCaseRunner) getConnection(host string) (conn *grpc.ClientConn, err error) {
	if r.Secure == nil || r.Secure.Insecure {
		conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		var cred credentials.TransportCredentials
		cred, err = credentials.NewClientTLSFromFile(r.Secure.CertFile, r.Secure.ServerName)
		if err == nil {
			conn, err = grpc.Dial(host, grpc.WithTransportCredentials(cred))
		}
	}
	return
}

func (r *gRPCTestCaseRunner) GetSuggestedAPIs(suite *testing.TestSuite, api string) (result []*testing.TestCase, err error) {
	if suite.Spec.RPC.ServerReflection {
		var conn *grpc.ClientConn
		if conn, err = r.getConnection(suite.API); err != nil {
			return
		}

		stub := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
		client := grpcreflect.NewClientV1Alpha(context.Background(), stub)

		var allSvcs []string
		allSvcs, err = client.ListServices()
		if err != nil {
			return
		}
		for _, svc := range allSvcs {
			var fd *desc.FileDescriptor
			if fd, err = client.FileContainingSymbol(svc); err != nil {
				return
			}
			svcs := fd.GetServices()

			for _, item := range svcs {
				for _, fdb := range item.GetMethods() {
					result = append(result, &testing.TestCase{
						Name: fdb.GetName(),
						Request: testing.Request{
							API: fmt.Sprintf("/%s/%s", svc, fdb.GetName()),
						},
					})
				}
			}
		}
		return
	}

	var linkerFiles linker.Files
	linkerFiles, err = compileProto(context.Background(), r)
	if err != nil {
		return
	}

	for _, f := range linkerFiles {
		for j := 0; j < f.Services().Len(); j++ {
			svc := f.Services().Get(j)

			methodCount := svc.Methods().Len()
			for m := 0; m < methodCount; m++ {
				method := svc.Methods().Get(m)
				methodName := string(method.Name())
				api := "/" + string(method.FullName())
				api = strings.ReplaceAll(api, "."+methodName, "/"+methodName)

				result = append(result, &testing.TestCase{
					Name: string(methodName),
					Request: testing.Request{
						API: api,
					},
				})
			}
		}
	}
	return
}

func (r *gRPCTestCaseRunner) GetResponseRecord() SimpleResponse {
	return r.response
}
func (s *gRPCTestCaseRunner) WithSuite(suite *testing.TestSuite) {
	// not need this parameter
}

func invokeRequest(ctx context.Context, md protoreflect.MethodDescriptor, payload string, conn *grpc.ClientConn) (response []string, err error) {
	resps := make([]*dynamicpb.Message, 0)
	if md.IsStreamingClient() || md.IsStreamingServer() {
		reqs, err := getStreamMessagepb(md.Input(), payload)
		if err != nil {
			return nil, err
		}

		resps, err = invokeRPCStream(ctx, conn, md, reqs)
		if err != nil {
			return nil, err
		}
		return buildResponses(resps)
	}
	request, err := getMessagePb(md.Input(), payload)
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

func getStreamMessagepb(md protoreflect.MessageDescriptor, messages string) ([]*dynamicpb.Message, error) {
	gpayload := gjson.Parse(messages)
	var garray []gjson.Result

	if !gpayload.IsArray() {
		garray = []gjson.Result{gpayload}
	} else {
		garray = gpayload.Array()
	}
	reqs := make([]*dynamicpb.Message, len(garray))

	for i, v := range garray {
		req, err := getMessagePb(md, v.Raw)
		if err != nil {
			return nil, err
		}
		reqs[i] = req
	}

	return reqs, nil
}

func getMessagePb(md protoreflect.MessageDescriptor, message string) (messagepb *dynamicpb.Message, err error) {
	messagepb = dynamicpb.NewMessage(md)
	if message != "" {
		err = protojson.Unmarshal([]byte(message), messagepb)
		if err != nil {

			err = fmt.Errorf("failed to unmarshal %q message: %v", messagepb.Descriptor().Name(), err)
		}
	}
	return
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
	// 	grpcRunnerLogger.Info("hit cache",fullname)
	// 	return getMdFromFd(fd.(protoreflect.FileDescriptor), fullname)
	// }

	if r.proto.ServerReflection {
		dp, err = getByReflect(ctx, r, fullname, conn)
	} else {
		if r.proto.ProtoFile == "" && r.proto.ProtoSet == "" && r.proto.Raw == "" {
			return nil, fmt.Errorf("missing descriptor source")
		}
		dp, err = getByProto(ctx, r, fullname)
	}

	if err != nil {
		return nil, err
	}

	if dp.IsPlaceholder() {
		return nil, protoregistry.NotFound
	}

	if md, ok := dp.(protoreflect.MethodDescriptor); ok {
		return md, nil
	}
	return nil, protoregistry.NotFound
}

func compileProto(ctx context.Context, r *gRPCTestCaseRunner) (fileLinker linker.Files, err error) {
	var protoFile string
	var importPath []string
	var parentProtoDir string
	protoFile, importPath, parentProtoDir, err = util.LoadProtoFiles(r.proto.ProtoFile)
	if err != nil {
		return
	}

	if len(importPath) == 0 {
		importPath = r.proto.ImportPath
	}

	if parentProtoDir != "" {
		for i, p := range importPath {
			importPath[i] = filepath.Join(parentProtoDir, p)
		}
		if len(importPath) == 0 {
			importPath = append(importPath, parentProtoDir)
		}
	}

	var protoLibrary map[string]string
	if protoLibrary, err = apispec.GetProtoFiles(); err != nil {
		return
	}

	grpcRunnerLogger.Info("proto import files", "files", importPath)
	compiler := protocompile.Compiler{
		Resolver: protocompile.WithStandardImports(
			&protocompile.SourceResolver{
				ImportPaths: importPath,
				Accessor: func(path string) (io.ReadCloser, error) {
					if content, ok := protoLibrary[strings.TrimPrefix(path, parentProtoDir+"/")]; ok {
						return io.NopCloser(strings.NewReader(content)), nil
					}
					return os.Open(path)
				},
			},
		),
	}

	// save the proto to a temp file if the raw content given
	if r.proto.Raw != "" {
		var f *os.File
		f, err = os.CreateTemp(os.TempDir(), "proto")
		if err != nil {
			err = fmt.Errorf("failed to create temp file when saving proto content: %v", err)
			return
		}
		defer os.Remove(f.Name())

		_, err = f.WriteString(r.proto.Raw)
		if err != nil {
			err = fmt.Errorf("failed to write proto content to file %q: %v", f.Name(), err)
			return
		}
		protoFile = f.Name()
	}

	return compiler.Compile(ctx, protoFile)
}

func getByProto(ctx context.Context, r *gRPCTestCaseRunner, fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	if r.proto.ProtoSet != "" {
		return getByProtoSet(ctx, r, fullName)
	}

	linker, err := compileProto(ctx, r)
	if err != nil {
		err = fmt.Errorf("failed to compile proto: %v", err)
		return nil, err
	}

	dp, err := linker.AsResolver().FindDescriptorByName(fullName)
	if err != nil {
		return nil, err
	}

	// r.fdCache.Store(fullName.Parent(), dp.ParentFile())
	return dp, nil
}

func getByProtoSet(ctx context.Context, r *gRPCTestCaseRunner, fullName protoreflect.FullName) (protoreflect.Descriptor, error) {
	var decs []byte
	var err error
	if regexURLPrefix.FindString(r.proto.ProtoSet) != "" {
		resp, err := http.Get(r.proto.ProtoSet)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		decs, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		decs, err = os.ReadFile(r.proto.ProtoSet)
		if err != nil {
			return nil, err
		}
	}

	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(decs, fds)
	if err != nil {
		return nil, err
	}

	prfs, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, err
	}

	dp, err := prfs.FindDescriptorByName(fullName)
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
		return nil, fmt.Errorf("%s", resp.GetErrorMessage())
	}

	fdresp := resp.GetFileDescriptorResponse()

	for _, fdb := range fdresp.FileDescriptorProto {
		fdp := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdb, fdp); err != nil {
			err = fmt.Errorf("failed to unmarshal file descriptor proto: %v", err)
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

	return nil, protoregistry.NotFound
}

func getMdFromFd(fd protoreflect.FileDescriptor, fullname protoreflect.FullName) (md protoreflect.MethodDescriptor, err error) {
	sd := fd.Services().ByName(fullname.Parent().Name())
	if sd == nil {
		return nil, fmt.Errorf("grpc service %q is not found in proto %q", fullname.Parent().Name(), fd.Name())
	}

	md = sd.Methods().ByName(fullname.Name())
	if md == nil {
		return nil, fmt.Errorf("method %q is not found in service %q", fullname.Name(), fullname.Parent().Name())
	}
	return md, nil
}

func splitFullQualifiedName(api string) (protoreflect.FullName, error) {
	qn := regexFullQualifiedName.FindStringSubmatch(api)
	if len(qn) == 0 {
		return "", fmt.Errorf("%q is not a valid gRPC api name", api)
	}
	fn := protoreflect.FullName(strings.Join(qn[2:], "."))
	if !fn.IsValid() {
		return "", fmt.Errorf("%q is not a valid gRPC api name", api)
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

// invokeRPC sends a unary RPC to gRPC server.
func invokeRPC(ctx context.Context, conn grpc.ClientConnInterface, method protoreflect.MethodDescriptor, request *dynamicpb.Message) (
	resp *dynamicpb.Message, err error) {
	resp = dynamicpb.NewMessage(method.Output())
	md, _ := metadata.FromIncomingContext(ctx)
	err = conn.Invoke(ctx, getMethodName(method), request, resp, grpc.Header(&md))
	return
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

func verifyResponsePayload(md protoreflect.MethodDescriptor, caseName string, expect testing.Response, jsonPayload []string) (output any, err error) {
	mapOutput := map[string]any{
		"data": func() []map[string]any {
			r := make([]map[string]any, len(jsonPayload))
			for i := range jsonPayload {
				m := map[string]any{}
				_ = json.Unmarshal([]byte(jsonPayload[i]), &m)
				r[i] = m
			}
			return r
		}(),
	}

	if err = payloadFieldsVerify(md, caseName, expect, jsonPayload); err != nil {
		return
	}

	err = Verify(expect, mapOutput)
	if err != nil {
		return nil, err
	}
	return
}

func payloadFieldsVerify(md protoreflect.MethodDescriptor, caseName string, expect testing.Response, jsonPayload []string) error {
	if expect.Body == "" {
		return nil
	}

	if !gjson.Valid(expect.Body) {
		return fmt.Errorf("case %q: expect body is not a valid JSON", caseName)
	}

	exp, err := parseExpect(md, expect)
	if err != nil {
		return err
	}

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
			return fmt.Errorf("%s", msg)
		}
		return nil
	}
	return fmt.Errorf("case %q: unknown expect content", caseName)
}

func parseExpect(md protoreflect.MethodDescriptor, expect testing.Response) (exps gjson.Result, err error) {
	b := strings.TrimSpace(expect.Body)
	var msgb []byte
	if b[0] == '[' {
		msgpbs, err := getStreamMessagepb(md.Output(), b)
		if err != nil {
			return gjson.Result{}, err
		}
		msgb = append(msgb, '[')
		for i := range msgpbs {
			msg, _ := protojson.Marshal(msgpbs[i])
			msgb = append(msgb, msg...)
			msg = append(msg, ',')
		}
		msgb[len(msgb)-1] = ']'
	} else {
		msgpb, err := getMessagePb(md.Output(), expect.Body)
		if err != nil {
			return gjson.Result{}, err
		}
		msgb, _ = protojson.Marshal(msgpb)
	}
	return gjson.ParseBytes(msgb), nil
}
