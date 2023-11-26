/**
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

package oauth

import (
	"context"

	"errors"

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

func GetUserFromContext(ctx context.Context) (user *UserInfo) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		data := md.Get("auth")
		if len(data) > 0 {
			user = GetUser(data[0])
		}
	}
	return
}
