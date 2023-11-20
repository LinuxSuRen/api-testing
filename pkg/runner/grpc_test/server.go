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

package grpc_test

import (
	"context"
	"google.golang.org/grpc/metadata"
	"io"
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
