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
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/linuxsuren/api-testing/pkg/apispec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type grpcResultWriter struct {
	context   context.Context
	targetUrl string
}

// NewGRPCResultWriter creates a new grpcResultWriter
func NewGRPCResultWriter(ctx context.Context, url string) ReportResultWriter {
	return &grpcResultWriter{
		context:   ctx,
		targetUrl: url,
	}
}

// Output writes the JSON base report to target writer
func (w *grpcResultWriter) Output(result []ReportResult) (err error) {
	server, err := w.getHost()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("will send report to:" + server)
	conn, err := getConnection(server)
	if err != nil {
		log.Println("Error when connecting to grpc server", err)
		return err
	}
	defer conn.Close()
	md, err := w.getMethodDescriptor(conn)
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
	resp, err := invokeRequest(w.context, md, payload, conn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("getting response back:", resp)
	return
}

// use server reflection to get the method descriptor
func (w *grpcResultWriter) getMethodDescriptor(conn *grpc.ClientConn) (protoreflect.MethodDescriptor, error) {
	fullName, err := splitFullQualifiedName(w.targetUrl)
	if err != nil {
		return nil, err
	}
	var dp protoreflect.Descriptor

	dp, err = getByReflect(w.context, nil, fullName, conn)
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
func (w *grpcResultWriter) getHost() (host string, err error) {
	qn := regexFullQualifiedName.FindStringSubmatch(w.targetUrl)
	if len(qn) == 0 {
		return "", errors.New("can not get host from url")
	}
	return qn[1], nil
}

// get connection with gRPC server
func getConnection(host string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

// WithAPIConverage sets the api coverage
func (w *grpcResultWriter) WithAPICoverage(apiConverage apispec.APICoverage) ReportResultWriter {
	return w
}

func (w *grpcResultWriter) WithResourceUsage([]ResourceUsage) ReportResultWriter {
	return w
}

func (w *grpcResultWriter) GetWriter() io.Writer {
	return nil
}
