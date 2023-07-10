package main

import (
	"context"
	"fmt"
	"time"

	"github.com/linuxsuren/api-testing/extensions/store-etcd/remote"
	"go.etcd.io/etcd/clientv3"
)

type server struct {
	remote.UnimplementedLoaderServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer() remote.LoaderServer {
	return &server{}
}

func (s *server) ListTestSuite(context.Context, *remote.Empty) (suites *remote.TestSuites, err error) {
	suites = &remote.TestSuites{
		Data: []*remote.TestSuite{{
			Name: "fake",
		}, {
			Name: "fake2",
		}},
	}
	return
}
func (s *server) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *remote.Empty, err error) {
	reply = &remote.Empty{}
	fmt.Println(*testSuite)

	var cli *clientv3.Client
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer cli.Close()

	_, err = cli.Put(ctx, testSuite.Name, "test")
	return
}
func (s *server) GetTestSuite(context.Context, *remote.TestSuite) (*remote.TestSuite, error) {
	return nil, nil
}
func (s *server) UpdateTestSuite(context.Context, *remote.TestSuite) (*remote.TestSuite, error) {
	return nil, nil
}
func (s *server) DeleteTestSuite(context.Context, *remote.TestSuite) (*remote.Empty, error) {
	return nil, nil
}
func (s *server) ListTestCases(context.Context, *remote.TestSuite) (*remote.TestCases, error) {
	return nil, nil
}
func (s *server) CreateTestCase(context.Context, *remote.TestCase) (*remote.Empty, error) {
	return nil, nil
}
func (s *server) GetTestCase(context.Context, *remote.TestCase) (*remote.TestCase, error) {
	return nil, nil
}
func (s *server) UpdateTestCase(context.Context, *remote.TestCase) (*remote.TestCase, error) {
	return nil, nil
}
func (s *server) DeleteTestCase(context.Context, *remote.TestCase) (*remote.Empty, error) {
	return nil, nil
}
