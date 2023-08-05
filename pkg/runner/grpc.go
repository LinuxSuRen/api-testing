package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/bufbuild/protocompile"
	"github.com/linuxsuren/api-testing/pkg/runner/kubernetes"
	"github.com/linuxsuren/api-testing/pkg/testing"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type gRPCTestCaseRunner struct {
	simpleTestCaseRunner
	host  string
	proto testing.GRPCDesc
}

func NewGRPCTestCaseRunner(host string, proto testing.GRPCDesc) TestCaseRunner {
	runner := &gRPCTestCaseRunner{
		simpleTestCaseRunner: simpleTestCaseRunner{},
		host:                 host,
		proto:                proto,
	}
	runner.WithOutputWriter(io.Discard).
		WithWriteLevel("info").
		WithTestReporter(NewDiscardTestReporter()).
		WithExecer(fakeruntime.DefaultExecer{})
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

	if err = runJob(testcase.Before); err != nil {
		return
	}

	payload := testcase.Request.Payload

	r.log.Info("start to send request to %s\n", testcase.Request.API)
	conn, err := grpc.Dial(r.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resps := make([]*dynamicpb.Message, 0)

	if md.IsStreamingClient() || md.IsStreamingServer() {
		result := gjson.Parse(payload)
		if !result.IsArray() {
			return nil, fmt.Errorf("payload is not a json array")
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
	record.Body = strings.Join(respsStr, ",")
	r.log.Debug("response body: %s\n", record.Body)

	output, err = verifyResponsePayload(testcase.Name, testcase.Expect, respsStr)
	if err != nil {
		return nil, err
	}

	return output, nil
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

	for _, verify := range expect.Verify {
		var program *vm.Program
		if program, err = expr.Compile(verify, expr.Env(mapOutput),
			expr.AsBool(), kubernetes.PodValidatorFunc(),
			kubernetes.KubernetesValidatorFunc()); err != nil {
			return
		}

		var result interface{}
		if result, err = expr.Run(program, mapOutput); err != nil {
			return
		}

		if !result.(bool) {
			err = fmt.Errorf("failed to verify: %s", verify)
			fmt.Println(err)
			break
		}
	}
	return
}

func payloadFieldsVerify(caseName string, expect testing.Response, jsonPayload []string) error {
	if expect.Payload == "" {
		return nil
	}
	if !gjson.Valid(expect.Payload) {
		return fmt.Errorf("case %s: expect payload is not a valid json", caseName)
	}

	result := gjson.Parse(expect.Payload)
	gjsonPayload := make([]gjson.Result, len(jsonPayload))
	for i := range jsonPayload {
		gjsonPayload[i] = gjson.Parse(jsonPayload[i])
	}

	if result.IsArray() {
		return compareArr(caseName, result.Array(), gjsonPayload)
	}
	if result.IsObject() {
		var errlist error
		for i := range jsonPayload {
			err := compareObj(fmt.Sprintf("%v[%v]", caseName, i),
				result.Map(), gjsonPayload[i].Map())
			if err != nil {
				errlist = errors.Join(err)
			}
		}
		return errlist
	}
	return fmt.Errorf("case %s: unknown expect content", caseName)
}

func compareObj(field string, expect, actul map[string]gjson.Result) error {
	var errlist error
	for k, ev := range expect {
		av, ok := actul[k]
		if !ok {
			errors.Join(errlist, fmt.Errorf("field %s: field %s is not exist", field, k))
			continue
		}

		err := compareElement(k, ev, av)
		if err != nil {
			errlist = errors.Join(errlist, fmt.Errorf("field %s: fail at %s", field, err.Error()))
		}
	}

	return errlist
}

func compareArr(field string, expect, actul []gjson.Result) error {
	var errlist error
	if l1, l2 := len(expect), len(actul); l1 != l2 {
		return fmt.Errorf("field %s: expect %v fields but got %v", field, l1, l2)
	}

	for i := range expect {
		err := compareElement(strconv.Itoa(i), expect[i], actul[i])
		if err != nil {
			errlist = errors.Join(errlist, fmt.Errorf("field %s: fail at %s", field, err.Error()))
		}
	}
	return errlist
}

func compareElement(field string, expect, actul gjson.Result) error {
	if expect.Type != actul.Type {
		return fmt.Errorf("field %s: expect type %s but got %v", field, expect.Type.String(), actul.Type.String())
	}

	if expect.IsObject() {
		return compareObj(field, expect.Map(), actul.Map())
	}
	if expect.IsArray() {
		return compareArr(field, expect.Array(), actul.Array())
	}

	if !reflect.DeepEqual(expect.Value(), actul.Value()) {
		return fmt.Errorf("field %s: expect %v but got %v", field, expect.Value(), actul.Value())
	}

	return nil
}
