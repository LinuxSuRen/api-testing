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
