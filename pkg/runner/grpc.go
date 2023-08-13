package runner

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bufbuild/protocompile"
	"github.com/linuxsuren/api-testing/pkg/compare"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type gRPCTestCaseRunner struct {
	UnimplementedRunner
	host  string
	proto testing.GRPCDesc
}

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

	md, err := getMethodDescriptor(ctx, r, testcase)
	if err != nil {
		return nil, err
	}

	if err = runJob(testcase.Before); err != nil {
		return
	}

	r.log.Info("start to send request to %s\n", testcase.Request.API)
	conn, err := grpc.Dial(r.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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
		result := gjson.Parse(payload)
		if !result.IsArray() {
			return nil, fmt.Errorf("payload is not a JSON array")
		}

		reqs := make([]*dynamicpb.Message, len(result.Array()))
		for i, v := range result.Array() {
			req := dynamicpb.NewMessage(md.Input())
			err := protojson.Unmarshal([]byte(v.Raw), req)
			if err != nil {
				return nil, err
			}
			reqs[i] = req
		}

		resps, err = invokeRPCStream(ctx, conn, md, reqs)
		if err != nil {
			return nil, err
		}

	} else {
		request := dynamicpb.NewMessage(md.Input())
		if payload != "" {
			err = protojson.Unmarshal([]byte(payload), request)
			if err != nil {
				return nil, err
			}
		}

		resp, err := invokeRPC(ctx, conn, md, request)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}

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

func getMethodDescriptor(ctx context.Context, r *gRPCTestCaseRunner, testcase *testing.TestCase) (protoreflect.MethodDescriptor, error) {
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

	fd, err := linker.AsResolver().FindFileByPath(r.proto.ProtoFile)
	if err != nil {
		return nil, err
	}

	api := splitServiceAndMethod(testcase.Request.API)
	if len(api) != 2 {
		return nil, fmt.Errorf("%s is not a valid gRPC api name", testcase.Request.API)
	}

	sd := fd.Services().ByName(protoreflect.Name(api[0]))
	if sd == nil {
		return nil, fmt.Errorf("grpc service %s is not found in proto %s", api[0], fd.Name())
	}

	md := sd.Methods().ByName(protoreflect.Name(api[1]))
	if md == nil {
		return nil, fmt.Errorf("method %s is not found in service %s", api[1], api[0])
	}
	return md, nil
}

func splitServiceAndMethod(api string) []string {
	return strings.Split(api, ".")
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
