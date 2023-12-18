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
package pkg

import (
	"context"
	"log"

	"github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type dbserver struct {
	remote.UnimplementedLoaderServer
}

func NewRemoteServer() remote.LoaderServer {
	return &dbserver{}
}

func (s *dbserver) ListTestSuite(ctx context.Context, _ *server.Empty) (suites *remote.TestSuites, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	suites = &remote.TestSuites{}
	var cur *mongo.Cursor
	if cur, err = client.Find(ctx, bson.D{}); err == nil {
		list := []testing.TestSuite{}

		if err = cur.All(ctx, &list); err == nil {
			for i := range list {
				suites.Data = append(suites.Data, remote.ConvertToGRPCTestSuite(&list[i]))
			}
		}
	}
	return
}

func (s *dbserver) CreateTestSuite(ctx context.Context, testSuite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	_, err = client.InsertOne(ctx, remote.ConvertToNormalTestSuite(testSuite))
	return
}

func (s *dbserver) GetTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	reply = &remote.TestSuite{}
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, bson.M{"name": suite.Name}); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			reply = remote.ConvertToGRPCTestSuite(suite)
		}
	}
	return
}

func (s *dbserver) UpdateTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *remote.TestSuite, err error) {
	reply = &remote.TestSuite{}
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)
	filter := nameFilter(suite)

	reply = &remote.TestSuite{}
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, filter); cur != nil {
		data := &testing.TestSuite{}
		if err = cur.Decode(data); err == nil {
			target := remote.ConvertToNormalTestSuite(suite)
			target.Items = data.Items

			_, err = client.UpdateOne(ctx, filter, bson.M{"$set": target})
		}
	}
	return
}

func (s *dbserver) DeleteTestSuite(ctx context.Context, suite *remote.TestSuite) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	_, err = client.DeleteOne(ctx, bson.M{"name": suite.Name})
	return
}

func (s *dbserver) ListTestCases(ctx context.Context, suite *remote.TestSuite) (reply *server.TestCases, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	reply = &server.TestCases{}
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, bson.M{"name": suite.Name}); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			for i := range suite.Items {
				reply.Data = append(reply.Data, remote.ConvertToGRPCTestCase(suite.Items[i]))
			}
		}
	}
	return
}

func (s *dbserver) CreateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)
	filter := suiteFilter(testcase)

	reply = &server.Empty{}
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, filter); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			suite.Items = append(suite.Items, remote.ConvertToNormalTestCase(testcase))

			_, err = client.UpdateOne(ctx, filter, bson.M{"$set": suite})
		}
	}
	return
}

func (s *dbserver) GetTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, bson.M{"name": testcase.SuiteName}); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			for _, item := range suite.Items {
				if item.Name == testcase.Name {
					reply = remote.ConvertToGRPCTestCase(item)
					break
				}
			}
		}
	}
	return
}

func (s *dbserver) UpdateTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.TestCase, err error) {
	reply = &server.TestCase{}
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)
	filter := suiteFilter(testcase)

	needToUpdate := false
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, filter); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			for i, item := range suite.Items {
				if item.Name == testcase.Name {
					suite.Items[i] = remote.ConvertToNormalTestCase(testcase)
					needToUpdate = true
					break
				}
			}

			if needToUpdate {
				_, err = client.UpdateOne(ctx, filter, bson.M{"$set": suite})
			}
		}
	}
	return
}

func (s *dbserver) DeleteTestCase(ctx context.Context, testcase *server.TestCase) (reply *server.Empty, err error) {
	reply = &server.Empty{}
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)
	filter := suiteFilter(testcase)

	needToUpdate := false
	var cur *mongo.SingleResult
	if cur = client.FindOne(ctx, filter); cur != nil {
		suite := &testing.TestSuite{}
		if err = cur.Decode(suite); err == nil {
			for i, item := range suite.Items {
				if item.Name == testcase.Name {
					suite.Items = append(suite.Items[0:i], suite.Items[i+1:]...)
					needToUpdate = true
					break
				}
			}

			if needToUpdate {
				_, err = client.UpdateOne(ctx, filter, bson.M{"$set": suite})
			}
		}
	}
	return
}

func (s *dbserver) Verify(ctx context.Context, in *server.Empty) (reply *server.ExtensionStatus, err error) {
	var client *mongo.Collection
	if client, err = s.getClient(ctx); err != nil {
		return
	}
	defer client.Database().Client().Disconnect(ctx)

	reply = &server.ExtensionStatus{}
	if pingErr := client.Database().Client().Ping(ctx, readpref.Primary()); pingErr == nil {
		reply.Ready = true
	}
	return
}
func (s *dbserver) PProf(ctx context.Context, in *server.PProfRequest) (data *server.PProfData, err error) {
	log.Println("pprof", in.Name)

	data = &server.PProfData{
		Data: extension.LoadPProf(in.Name),
	}
	return
}

type SuiteNameGetter interface {
	GetSuiteName() string
}

type NameGetter interface {
	GetName() string
}

func suiteFilter(obj SuiteNameGetter) interface{} {
	return bson.M{"name": obj.GetSuiteName()}
}

func nameFilter(obj NameGetter) interface{} {
	return bson.M{"name": obj.GetName()}
}
