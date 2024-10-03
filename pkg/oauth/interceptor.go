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
package oauth

import (
	"context"

	"errors"

	"github.com/linuxsuren/oauth-hub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewAuthInterceptor(groups []string) grpc.ServerOption {
	inter := &authInter{
		groups: groups,
	}
	return grpc.UnaryInterceptor(inter.authInterceptor)
}

type authInter struct {
	groups []string
}

func (a *authInter) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (
	resp interface{}, err error) {
	if info.FullMethod == "/server.Runner/GetVersion" {
		return handler(ctx, req)
	}

	if user := GetUserFromContext(ctx); user != nil {
		approve := len(a.groups) == 0
		for _, g := range a.groups {
			for _, gg := range user.Groups {
				if g == gg {
					approve = true
					break
				}
			}
			if approve {
				break
			}
		}

		if approve {
			resp, err = handler(ctx, req)
		} else {
			err = errors.New("invalid group")
		}
		return
	}

	msg := "no auth found"
	if err != nil {
		msg = err.Error()
	}

	sta := status.New(codes.Unauthenticated, msg)
	resp = sta
	err = sta.Err()
	return
}

func GetUserFromContext(ctx context.Context) (user *oauth.UserInfo) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		data := md.Get("auth")
		if len(data) > 0 {
			user = oauth.GetUser(data[0])
		}
	}
	return
}
