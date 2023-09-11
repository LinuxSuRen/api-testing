package server

import (
	context "context"
	"net"
	"net/http"
)

// HTTPServer is an interface for serving HTTP requests
type HTTPServer interface {
	Serve(lis net.Listener) error
	WithHandler(handler http.Handler)
	Shutdown(ctx context.Context) error
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
