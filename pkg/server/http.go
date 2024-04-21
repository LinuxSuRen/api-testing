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
package server

import (
	context "context"
	"net"
	"net/http"
	"strings"
)

// HTTPServer is an interface for serving HTTP requests
type HTTPServer interface {
	Serve(lis net.Listener) error
	WithHandler(handler http.Handler)
	Shutdown(ctx context.Context) error
}

type CombineHandler interface {
	PutHandler(string, http.Handler)
	GetHandler() http.Handler
}

type defaultHTTPServer struct {
	server  *http.Server
	handler http.Handler
}

// NewDefaultHTTPServer creates a default HTTP server
func NewDefaultHTTPServer() HTTPServer {
	return &defaultHTTPServer{}
}

func (s *defaultHTTPServer) Serve(lis net.Listener) (err error) {
	s.server = &http.Server{Handler: s.handler}
	err = s.server.Serve(lis)
	return
}

func (s *defaultHTTPServer) WithHandler(h http.Handler) {
	s.handler = h
}

func (s *defaultHTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

type defaultCombineHandler struct {
	handlerMapping map[string]http.Handler
	defaultHandler http.Handler
}

func NewDefaultCombineHandler() CombineHandler {
	return &defaultCombineHandler{
		handlerMapping: make(map[string]http.Handler),
	}
}

func (s *defaultCombineHandler) PutHandler(name string, handler http.Handler) {
	if name == "" {
		s.defaultHandler = handler
	} else {
		s.handlerMapping[name] = handler
	}
}

func (s *defaultCombineHandler) GetHandler() http.Handler {
	if len(s.handlerMapping) == 0 {
		return s.defaultHandler
	}
	return s
}

func (s *defaultCombineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for prefix, handler := range s.handlerMapping {
		if strings.HasPrefix(r.URL.Path, prefix) {
			handler.ServeHTTP(w, r)
			return
		}
	}
	s.defaultHandler.ServeHTTP(w, r)
}

type fakeHandler struct{}

// NewFakeHTTPServer creates a fake HTTP server
func NewFakeHTTPServer() HTTPServer {
	return &fakeHandler{}
}

func (s *fakeHandler) Serve(lis net.Listener) (err error) {
	// do nothing due to this is a fake method
	return
}

func (s *fakeHandler) WithHandler(h http.Handler) {
	// do nothing due to this is a fake method
}

func (s *fakeHandler) Shutdown(ctx context.Context) error {
	// do nothing due to this is a fake method
	return nil
}
