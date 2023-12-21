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
	"errors"
	"fmt"

	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"github.com/linuxsuren/api-testing/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *dbserver) getClient(ctx context.Context) (collection *mongo.Collection, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to mongodb")
		return
	}

	if store.Properties == nil {
		store.Properties = map[string]string{}
	}

	databaseName := util.EmptyThenDefault(store.Properties["database"], "testing")
	collectionName := util.EmptyThenDefault(store.Properties["collection"], "atest")

	address := fmt.Sprintf("mongodb://%s:%s@%s", store.Username, store.Password, store.URL)
	var client *mongo.Client
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(address))
	if err != nil {
		return
	}

	collection = client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		if err = client.Database(databaseName).CreateCollection(ctx, collectionName); err == nil {
			collection = client.Database(databaseName).Collection(collectionName)
		}
	}
	return
}
