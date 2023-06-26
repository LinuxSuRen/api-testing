package server

import (
	"net"
	"net/http"
)

// HTTPServer is an interface for serving HTTP requests
type HTTPServer interface {
	Serve(lis net.Listener) error
	WithHandler(handler http.Handler)
}

type defaultHTTPServer struct {
	handler http.Handler
}

// NewDefaultHTTPServer creates a default HTTP server
func NewDefaultHTTPServer() HTTPServer {
	return &defaultHTTPServer{}
}

func (s *defaultHTTPServer) Serve(lis net.Listener) (err error) {
	server := &http.Server{Handler: s.handler}
	err = server.Serve(lis)
	return
}

func (s *defaultHTTPServer) WithHandler(h http.Handler) {
	s.handler = h
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
