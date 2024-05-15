/*
Copyright 2024 API Testing Authors.

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
	"log"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type grpcResultWriter struct {
	targetUrl string
}

// NewGRPCResultWriter creates a new grpcResultWriter
func NewGRPCResultWriter(url string) ReportResultWriter {
	return &grpcResultWriter{
		targetUrl: url,
	}
}

// Output writes the JSON base report to target writer
func (w *grpcResultWriter) Output(result []ReportResult) (err error) {
	server := getHost(w.targetUrl, "127.0.0.1")
	log.Println("will send report to:" + server)
	conn, err := getConnection(server)
	if err != nil {
		log.Println("Error when connecting to grpc server", err)
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	md, err := w.getMethodDescriptor(ctx, conn)
	if err != nil {
		if err == protoregistry.NotFound {
			return fmt.Errorf("api %q is not found on grpc server", w.targetUrl)
		}
		return err
	}
	jsonPayload, _ := json.Marshal(
		map[string][]ReportResult{
			"data": result,
		})
	payload := string(jsonPayload)
	resp, err := invokeRequest(ctx, md, payload, conn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("getting response back:", resp)
	return
}

// use server reflection to get the method descriptor
func (w *grpcResultWriter) getMethodDescriptor(ctx context.Context, conn *grpc.ClientConn) (protoreflect.MethodDescriptor, error) {
	fullName, err := splitFullQualifiedName(w.targetUrl)
	if err != nil {
		return nil, err
	}
	var dp protoreflect.Descriptor

	dp, err = getByReflect(ctx, nil, fullName, conn)
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

// get connection with gRPC server
func getConnection(host string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

// WithAPIConverage sets the api coverage
func (w *grpcResultWriter) WithAPIConverage(apiConverage apispec.APIConverage) ReportResultWriter {
	return w
}

func (w *grpcResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}
