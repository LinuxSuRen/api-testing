// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: pkg/server/server.proto

package server

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RunnerClient is the client API for Runner service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RunnerClient interface {
	Run(ctx context.Context, in *TestTask, opts ...grpc.CallOption) (*TestResult, error)
	Sample(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error)
	GetVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error)
	GetSuites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Suites, error)
	CreateTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*HelloReply, error)
	GetTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*TestSuite, error)
	UpdateTestSuite(ctx context.Context, in *TestSuite, opts ...grpc.CallOption) (*HelloReply, error)
	DeleteTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*HelloReply, error)
	ListTestCase(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*Suite, error)
	RunTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*TestCaseResult, error)
	GetTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*TestCase, error)
	CreateTestCase(ctx context.Context, in *TestCaseWithSuite, opts ...grpc.CallOption) (*HelloReply, error)
	UpdateTestCase(ctx context.Context, in *TestCaseWithSuite, opts ...grpc.CallOption) (*HelloReply, error)
	DeleteTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*HelloReply, error)
	PopularHeaders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Pairs, error)
}

type runnerClient struct {
	cc grpc.ClientConnInterface
}

func NewRunnerClient(cc grpc.ClientConnInterface) RunnerClient {
	return &runnerClient{cc}
}

func (c *runnerClient) Run(ctx context.Context, in *TestTask, opts ...grpc.CallOption) (*TestResult, error) {
	out := new(TestResult)
	err := c.cc.Invoke(ctx, "/server.Runner/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) Sample(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/Sample", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) GetVersion(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) GetSuites(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Suites, error) {
	out := new(Suites)
	err := c.cc.Invoke(ctx, "/server.Runner/GetSuites", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) CreateTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/CreateTestSuite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) GetTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*TestSuite, error) {
	out := new(TestSuite)
	err := c.cc.Invoke(ctx, "/server.Runner/GetTestSuite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) UpdateTestSuite(ctx context.Context, in *TestSuite, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/UpdateTestSuite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) DeleteTestSuite(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/DeleteTestSuite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) ListTestCase(ctx context.Context, in *TestSuiteIdentity, opts ...grpc.CallOption) (*Suite, error) {
	out := new(Suite)
	err := c.cc.Invoke(ctx, "/server.Runner/ListTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) RunTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*TestCaseResult, error) {
	out := new(TestCaseResult)
	err := c.cc.Invoke(ctx, "/server.Runner/RunTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) GetTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*TestCase, error) {
	out := new(TestCase)
	err := c.cc.Invoke(ctx, "/server.Runner/GetTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) CreateTestCase(ctx context.Context, in *TestCaseWithSuite, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/CreateTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) UpdateTestCase(ctx context.Context, in *TestCaseWithSuite, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/UpdateTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) DeleteTestCase(ctx context.Context, in *TestCaseIdentity, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/server.Runner/DeleteTestCase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runnerClient) PopularHeaders(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Pairs, error) {
	out := new(Pairs)
	err := c.cc.Invoke(ctx, "/server.Runner/PopularHeaders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RunnerServer is the server API for Runner service.
// All implementations must embed UnimplementedRunnerServer
// for forward compatibility
type RunnerServer interface {
	Run(context.Context, *TestTask) (*TestResult, error)
	Sample(context.Context, *Empty) (*HelloReply, error)
	GetVersion(context.Context, *Empty) (*HelloReply, error)
	GetSuites(context.Context, *Empty) (*Suites, error)
	CreateTestSuite(context.Context, *TestSuiteIdentity) (*HelloReply, error)
	GetTestSuite(context.Context, *TestSuiteIdentity) (*TestSuite, error)
	UpdateTestSuite(context.Context, *TestSuite) (*HelloReply, error)
	DeleteTestSuite(context.Context, *TestSuiteIdentity) (*HelloReply, error)
	ListTestCase(context.Context, *TestSuiteIdentity) (*Suite, error)
	RunTestCase(context.Context, *TestCaseIdentity) (*TestCaseResult, error)
	GetTestCase(context.Context, *TestCaseIdentity) (*TestCase, error)
	CreateTestCase(context.Context, *TestCaseWithSuite) (*HelloReply, error)
	UpdateTestCase(context.Context, *TestCaseWithSuite) (*HelloReply, error)
	DeleteTestCase(context.Context, *TestCaseIdentity) (*HelloReply, error)
	PopularHeaders(context.Context, *Empty) (*Pairs, error)
	mustEmbedUnimplementedRunnerServer()
}

// UnimplementedRunnerServer must be embedded to have forward compatible implementations.
type UnimplementedRunnerServer struct {
}

func (UnimplementedRunnerServer) Run(context.Context, *TestTask) (*TestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (UnimplementedRunnerServer) Sample(context.Context, *Empty) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sample not implemented")
}
func (UnimplementedRunnerServer) GetVersion(context.Context, *Empty) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedRunnerServer) GetSuites(context.Context, *Empty) (*Suites, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuites not implemented")
}
func (UnimplementedRunnerServer) CreateTestSuite(context.Context, *TestSuiteIdentity) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTestSuite not implemented")
}
func (UnimplementedRunnerServer) GetTestSuite(context.Context, *TestSuiteIdentity) (*TestSuite, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestSuite not implemented")
}
func (UnimplementedRunnerServer) UpdateTestSuite(context.Context, *TestSuite) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTestSuite not implemented")
}
func (UnimplementedRunnerServer) DeleteTestSuite(context.Context, *TestSuiteIdentity) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTestSuite not implemented")
}
func (UnimplementedRunnerServer) ListTestCase(context.Context, *TestSuiteIdentity) (*Suite, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTestCase not implemented")
}
func (UnimplementedRunnerServer) RunTestCase(context.Context, *TestCaseIdentity) (*TestCaseResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunTestCase not implemented")
}
func (UnimplementedRunnerServer) GetTestCase(context.Context, *TestCaseIdentity) (*TestCase, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestCase not implemented")
}
func (UnimplementedRunnerServer) CreateTestCase(context.Context, *TestCaseWithSuite) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTestCase not implemented")
}
func (UnimplementedRunnerServer) UpdateTestCase(context.Context, *TestCaseWithSuite) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTestCase not implemented")
}
func (UnimplementedRunnerServer) DeleteTestCase(context.Context, *TestCaseIdentity) (*HelloReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTestCase not implemented")
}
func (UnimplementedRunnerServer) PopularHeaders(context.Context, *Empty) (*Pairs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PopularHeaders not implemented")
}
func (UnimplementedRunnerServer) mustEmbedUnimplementedRunnerServer() {}

// UnsafeRunnerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RunnerServer will
// result in compilation errors.
type UnsafeRunnerServer interface {
	mustEmbedUnimplementedRunnerServer()
}

func RegisterRunnerServer(s grpc.ServiceRegistrar, srv RunnerServer) {
	s.RegisterService(&Runner_ServiceDesc, srv)
}

func _Runner_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).Run(ctx, req.(*TestTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_Sample_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).Sample(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/Sample",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).Sample(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).GetVersion(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_GetSuites_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).GetSuites(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/GetSuites",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).GetSuites(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_CreateTestSuite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSuiteIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).CreateTestSuite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/CreateTestSuite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).CreateTestSuite(ctx, req.(*TestSuiteIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_GetTestSuite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSuiteIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).GetTestSuite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/GetTestSuite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).GetTestSuite(ctx, req.(*TestSuiteIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_UpdateTestSuite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSuite)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).UpdateTestSuite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/UpdateTestSuite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).UpdateTestSuite(ctx, req.(*TestSuite))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_DeleteTestSuite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSuiteIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).DeleteTestSuite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/DeleteTestSuite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).DeleteTestSuite(ctx, req.(*TestSuiteIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_ListTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestSuiteIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).ListTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/ListTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).ListTestCase(ctx, req.(*TestSuiteIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_RunTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCaseIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).RunTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/RunTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).RunTestCase(ctx, req.(*TestCaseIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_GetTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCaseIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).GetTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/GetTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).GetTestCase(ctx, req.(*TestCaseIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_CreateTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCaseWithSuite)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).CreateTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/CreateTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).CreateTestCase(ctx, req.(*TestCaseWithSuite))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_UpdateTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCaseWithSuite)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).UpdateTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/UpdateTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).UpdateTestCase(ctx, req.(*TestCaseWithSuite))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_DeleteTestCase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestCaseIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).DeleteTestCase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/DeleteTestCase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).DeleteTestCase(ctx, req.(*TestCaseIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Runner_PopularHeaders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RunnerServer).PopularHeaders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/server.Runner/PopularHeaders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RunnerServer).PopularHeaders(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Runner_ServiceDesc is the grpc.ServiceDesc for Runner service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Runner_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "server.Runner",
	HandlerType: (*RunnerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Run",
			Handler:    _Runner_Run_Handler,
		},
		{
			MethodName: "Sample",
			Handler:    _Runner_Sample_Handler,
		},
		{
			MethodName: "GetVersion",
			Handler:    _Runner_GetVersion_Handler,
		},
		{
			MethodName: "GetSuites",
			Handler:    _Runner_GetSuites_Handler,
		},
		{
			MethodName: "CreateTestSuite",
			Handler:    _Runner_CreateTestSuite_Handler,
		},
		{
			MethodName: "GetTestSuite",
			Handler:    _Runner_GetTestSuite_Handler,
		},
		{
			MethodName: "UpdateTestSuite",
			Handler:    _Runner_UpdateTestSuite_Handler,
		},
		{
			MethodName: "DeleteTestSuite",
			Handler:    _Runner_DeleteTestSuite_Handler,
		},
		{
			MethodName: "ListTestCase",
			Handler:    _Runner_ListTestCase_Handler,
		},
		{
			MethodName: "RunTestCase",
			Handler:    _Runner_RunTestCase_Handler,
		},
		{
			MethodName: "GetTestCase",
			Handler:    _Runner_GetTestCase_Handler,
		},
		{
			MethodName: "CreateTestCase",
			Handler:    _Runner_CreateTestCase_Handler,
		},
		{
			MethodName: "UpdateTestCase",
			Handler:    _Runner_UpdateTestCase_Handler,
		},
		{
			MethodName: "DeleteTestCase",
			Handler:    _Runner_DeleteTestCase_Handler,
		},
		{
			MethodName: "PopularHeaders",
			Handler:    _Runner_PopularHeaders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/server/server.proto",
}
