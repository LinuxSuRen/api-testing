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

package pkg

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/version"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type remoteserver struct {
	kvFactory KVFactory
	remote.UnimplementedLoaderServer
}

// NewRemoteServer creates a remote server instance
func NewRemoteServer(kvFactory KVFactory) remote.LoaderServer {
	return &remoteserver{
		kvFactory: kvFactory,
	}
}

const keyPrefix = "atest-"

func (s *remoteserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	suites = &remote.TestSuites{}

	var resp *clientv3.GetResponse
	var testsuite *testing.TestSuite
	if resp, err = cli.Get(ctx, keyPrefix, clientv3.WithLimit(100), clientv3.WithPrefix()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				suites.Data = append(suites.Data, remote.ConvertToGRPCTestSuite(testsuite))
			}
		}
	}
	return
}
func (s *remoteserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}

	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	err = saveTestSuite(ctx, cli, remote.ConvertToNormalTestSuite(testSuite))
	return
}
func (s *remoteserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &remote.TestSuite{}

	var resp *clientv3.GetResponse
	var testsuite *testing.TestSuite
	if resp, err = cli.Get(ctx, getKey(suite.Name), clientv3.WithFromKey()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				reply = remote.ConvertToGRPCTestSuite(testsuite)
				return
			}
		}
	}
	return
}
func (s *remoteserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &remote.TestSuite{}

	var resp *clientv3.GetResponse
	var testsuite *testing.TestSuite
	if resp, err = cli.Get(ctx, getKey(suite.Name), clientv3.WithFromKey()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				oldItems := testsuite.Items
				testsuite = remote.ConvertToNormalTestSuite(suite)
				testsuite.Items = oldItems
				err = saveTestSuite(ctx, cli, testsuite)
				return
			}
		}
	}
	return
}
func (s *remoteserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()
	reply = &server.Empty{}

	_, err = cli.Delete(ctx, getKey(suite.Name))
	return
}
func (s *remoteserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (reply *server.TestCases, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &server.TestCases{}

	var resp *clientv3.GetResponse
	var testsuite *testing.TestSuite
	if resp, err = cli.Get(ctx, getKey(suite.Name), clientv3.WithFromKey()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				reply.Data = remote.ConvertToGRPCTestSuite(testsuite).Items
				return
			}
		}
	}
	return
}
func (s *remoteserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &server.Empty{}

	var resp *clientv3.GetResponse
	var testsuite *testing.TestSuite
	if resp, err = cli.Get(ctx, getKey(testcase.SuiteName), clientv3.WithFromKey()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				suite := remote.ConvertToGRPCTestSuite(testsuite)
				suite.Items = append(suite.Items, testcase)

				err = saveTestSuite(ctx, cli, remote.ConvertToNormalTestSuite(suite))
				return
			}
		}
	}
	return
}
func (s *remoteserver) GetTestCase(ctx context.Context, input *server.TestCase) (reply *server.TestCase, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &server.TestCase{}

	var testcase *testing.TestCase
	var index int
	if testcase, index, _, err = getTestCase(ctx, cli, input.SuiteName, input.Name); err == nil && index != NotFound {
		reply = remote.ConvertToGRPCTestCase(*testcase)
	}
	return
}
func (s *remoteserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &server.TestCase{}

	var testsuite *testing.TestSuite
	var index int
	if _, index, testsuite, err = getTestCase(ctx, cli, testcase.SuiteName, testcase.Name); err == nil && index != NotFound {
		testsuite.Items[index] = remote.ConvertToNormalTestCase(testcase)
		err = saveTestSuite(ctx, cli, testsuite)
	}
	return
}
func (s *remoteserver) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		return
	}
	defer cli.Close()

	reply = &server.Empty{}

	var testsuite *testing.TestSuite
	var index int
	if _, index, testsuite, err = getTestCase(ctx, cli, testcase.SuiteName, testcase.Name); err == nil && index != NotFound {
		testsuite.Items = append(testsuite.Items[:index], testsuite.Items[index+1:]...)
		err = saveTestSuite(ctx, cli, testsuite)
	}
	return
}
func (s *remoteserver) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	reply = &server.ExtensionStatus{
		Version: version.GetVersion(),
	}

	var cli SimpleKV
	cli, err = s.getClient(ctx)
	if err != nil {
		reply.Message = err.Error()
		return
	}

	defer cli.Close()
	// try to connect
	if _, err = cli.Get(ctx, keyPrefix, connectTestOption()...); err == nil {
		reply.Ready = true
	}
	return
}

func connectTestOption() []clientv3.OpOption {
	return []clientv3.OpOption{clientv3.WithLimit(1), clientv3.WithPrefix(),
		clientv3.WithCountOnly(), clientv3.WithKeysOnly()}
}

func getTestCase(ctx context.Context, cli SimpleKV, suiteName, caseName string) (testcase *testing.TestCase, index int, testsuite *testing.TestSuite, err error) {
	index = NotFound
	var resp *clientv3.GetResponse
	if resp, err = cli.Get(ctx, getKey(suiteName), clientv3.WithFromKey()); err == nil {
		for _, kv := range resp.Kvs {
			if testsuite, err = testing.ParseFromData(kv.Value); err == nil && testsuite != nil {
				for i, item := range testsuite.Items {
					if item.Name == caseName {
						testcase = &item
						index = i
						return
					}
				}
			}
		}
	}
	return
}
func saveTestSuite(ctx context.Context, cli SimpleKV, suite *testing.TestSuite) (err error) {
	var data []byte
	if data, err = testing.ToYAML(suite); err == nil {
		_, err = cli.Put(ctx, getKey(suite.Name), string(data))
		log.Println("save to etcd", err)
	}
	return
}

const NotFound = -1

func (s *remoteserver) getClient(ctx context.Context) (cli SimpleKV, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to etcd server")
	} else {
		cli, err = s.kvFactory.New(clientv3.Config{
			Endpoints:   []string{store.URL},
			DialTimeout: 5 * time.Second,
			Username:    store.Username,
			Password:    store.Password,
			Context:     ctx,
		})
	}
	return
}
func getKey(name string) string {
	return keyPrefix + name
}
