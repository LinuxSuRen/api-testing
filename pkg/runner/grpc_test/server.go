/*
Copyright 2023 API Testing Authors.

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

package grpc_test

import (
	"context"
	"io"

	"google.golang.org/grpc/metadata"
)

type TestServer struct {
	UnimplementedMainServer
}

func (s *TestServer) Unary(ctx context.Context, _ *Empty) (*HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	var msg string
	if ok {
		// just for test purpose
		if items := md.Get("message"); len(items) > 0 {
			msg = items[0]
		}
	}

	if msg == "" {
		msg = "Hello!"
	}

	return &HelloReply{
		Message: msg,
	}, nil
}

func (s *TestServer) ClientStream(stream Main_ClientStreamServer) error {
	msgs := make([]*StreamMessage, 0)
	for {
		v, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		msgs = append(msgs, v)
	}
	return stream.SendAndClose(&StreamMessageRepeated{Data: msgs})
}

func (s *TestServer) ServerStream(msg *StreamMessageRepeated, stream Main_ServerStreamServer) error {
	var err error
	for i := range msg.Data {
		err = stream.Send(msg.Data[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) BidStream(stream Main_BidStreamServer) error {
	msgs := make([]*StreamMessage, 0)
	for {
		v, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		msgs = append(msgs, v)
	}

	var err error
	for i := range msgs {
		err = stream.Send(msgs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) TestBasicType(ctx context.Context, baiscType *BasicType) (*BasicType, error) {
	return baiscType, nil
}

func (s *TestServer) TestAdvancedType(ctx context.Context, advancedType *AdvancedType) (*AdvancedType, error) {
	return advancedType, nil
}
