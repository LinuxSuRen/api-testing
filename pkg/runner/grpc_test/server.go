package grpc_test

import (
	"context"
	"io"
)

type TestServer struct {
	UnimplementedMainServer
}

func (s *TestServer) Unary(context.Context, *Empty) (*HelloReply, error) {
	return &HelloReply{
		Message: "Hello!",
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
