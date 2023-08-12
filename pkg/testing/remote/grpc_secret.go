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

package remote

import (
	context "context"

	"github.com/linuxsuren/api-testing/pkg/secret"
	"github.com/linuxsuren/api-testing/pkg/server"
	"google.golang.org/grpc"
)

type gRPCSecret struct {
	UnimplementedSecretServiceServer
	client SecretServiceClient
}

func NewGRPCSecretFrom(address string) (server SecretServiceServer, err error) {
	server = &gRPCSecret{}
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(address, grpc.WithInsecure()); err == nil {
		server = &gRPCSecret{
			client: NewSecretServiceClient(conn),
		}
	}
	return
}

func (s *gRPCSecret) GetSecret(ctx context.Context, in *server.Secret) (reply *server.Secret, err error) {
	return s.client.GetSecret(ctx, in)
}

func (s *gRPCSecret) GetSecrets(ctx context.Context, in *server.Empty) (reply *server.Secrets, err error) {
	return s.client.GetSecrets(ctx, in)
}

func (s *gRPCSecret) CreateSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	return s.client.CreateSecret(ctx, in)
}

func (s *gRPCSecret) DeleteSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	return s.client.DeleteSecret(ctx, in)
}

func (s *gRPCSecret) UpdateSecret(ctx context.Context, in *server.Secret) (reply *server.CommonResult, err error) {
	return s.client.UpdateSecret(ctx, in)
}

// make sure gRPCSecret implemented SecretServiceClient
var _ SecretServiceServer = &gRPCSecret{}

type grpcSecretGetter struct {
	remoteServer server.SecertServiceGetable
}

func NewGRPCSecretGetter(remoteServer server.SecertServiceGetable) secret.SecretGetter {
	return &grpcSecretGetter{
		remoteServer: remoteServer,
	}
}

func (s *grpcSecretGetter) GetSecret(name string) (reply secret.Secret, err error) {
	var result *server.Secret
	if result, err = s.remoteServer.GetSecret(context.Background(),
		&server.Secret{Name: name}); err == nil && result != nil {
		reply.Name = result.Name
		reply.Value = result.Value
	}
	return
}
