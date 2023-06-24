// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: pkg/server/server.proto

/*
Package server is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package server

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Runner_Run_0(ctx context.Context, marshaler runtime.Marshaler, client RunnerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq TestTask
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.Run(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Runner_Run_0(ctx context.Context, marshaler runtime.Marshaler, server RunnerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq TestTask
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.Run(ctx, &protoReq)
	return msg, metadata, err

}

func request_Runner_Sample_0(ctx context.Context, marshaler runtime.Marshaler, client RunnerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.Sample(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Runner_Sample_0(ctx context.Context, marshaler runtime.Marshaler, server RunnerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.Sample(ctx, &protoReq)
	return msg, metadata, err

}

func request_Runner_GetVersion_0(ctx context.Context, marshaler runtime.Marshaler, client RunnerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetVersion(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Runner_GetVersion_0(ctx context.Context, marshaler runtime.Marshaler, server RunnerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GetVersion(ctx, &protoReq)
	return msg, metadata, err

}

func request_Runner_GetSuites_0(ctx context.Context, marshaler runtime.Marshaler, client RunnerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetSuites(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Runner_GetSuites_0(ctx context.Context, marshaler runtime.Marshaler, server RunnerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Empty
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GetSuites(ctx, &protoReq)
	return msg, metadata, err

}

func request_Runner_GetTestCase_0(ctx context.Context, marshaler runtime.Marshaler, client RunnerClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq TestCaseIdentity
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetTestCase(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Runner_GetTestCase_0(ctx context.Context, marshaler runtime.Marshaler, server RunnerServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq TestCaseIdentity
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GetTestCase(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterRunnerHandlerServer registers the http handlers for service Runner to "mux".
// UnaryRPC     :call RunnerServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterRunnerHandlerFromEndpoint instead.
func RegisterRunnerHandlerServer(ctx context.Context, mux *runtime.ServeMux, server RunnerServer) error {

	mux.Handle("POST", pattern_Runner_Run_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/server.Runner/Run", runtime.WithHTTPPathPattern("/server.Runner/Run"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Runner_Run_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_Run_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_Sample_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/server.Runner/Sample", runtime.WithHTTPPathPattern("/server.Runner/Sample"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Runner_Sample_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_Sample_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/server.Runner/GetVersion", runtime.WithHTTPPathPattern("/server.Runner/GetVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Runner_GetVersion_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetSuites_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/server.Runner/GetSuites", runtime.WithHTTPPathPattern("/server.Runner/GetSuites"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Runner_GetSuites_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetSuites_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetTestCase_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/server.Runner/GetTestCase", runtime.WithHTTPPathPattern("/server.Runner/GetTestCase"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Runner_GetTestCase_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetTestCase_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterRunnerHandlerFromEndpoint is same as RegisterRunnerHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterRunnerHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterRunnerHandler(ctx, mux, conn)
}

// RegisterRunnerHandler registers the http handlers for service Runner to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterRunnerHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterRunnerHandlerClient(ctx, mux, NewRunnerClient(conn))
}

// RegisterRunnerHandlerClient registers the http handlers for service Runner
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "RunnerClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "RunnerClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "RunnerClient" to call the correct interceptors.
func RegisterRunnerHandlerClient(ctx context.Context, mux *runtime.ServeMux, client RunnerClient) error {

	mux.Handle("POST", pattern_Runner_Run_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/server.Runner/Run", runtime.WithHTTPPathPattern("/server.Runner/Run"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Runner_Run_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_Run_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_Sample_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/server.Runner/Sample", runtime.WithHTTPPathPattern("/server.Runner/Sample"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Runner_Sample_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_Sample_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetVersion_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/server.Runner/GetVersion", runtime.WithHTTPPathPattern("/server.Runner/GetVersion"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Runner_GetVersion_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetVersion_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetSuites_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/server.Runner/GetSuites", runtime.WithHTTPPathPattern("/server.Runner/GetSuites"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Runner_GetSuites_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetSuites_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("POST", pattern_Runner_GetTestCase_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/server.Runner/GetTestCase", runtime.WithHTTPPathPattern("/server.Runner/GetTestCase"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Runner_GetTestCase_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Runner_GetTestCase_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_Runner_Run_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"server.Runner", "Run"}, ""))

	pattern_Runner_Sample_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"server.Runner", "Sample"}, ""))

	pattern_Runner_GetVersion_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"server.Runner", "GetVersion"}, ""))

	pattern_Runner_GetSuites_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"server.Runner", "GetSuites"}, ""))

	pattern_Runner_GetTestCase_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"server.Runner", "GetTestCase"}, ""))
)

var (
	forward_Runner_Run_0 = runtime.ForwardResponseMessage

	forward_Runner_Sample_0 = runtime.ForwardResponseMessage

	forward_Runner_GetVersion_0 = runtime.ForwardResponseMessage

	forward_Runner_GetSuites_0 = runtime.ForwardResponseMessage

	forward_Runner_GetTestCase_0 = runtime.ForwardResponseMessage
)
