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
	"fmt"

	"github.com/linuxsuren/api-testing/pkg/testing/remote"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *dbserver) getClient(ctx context.Context) (collection *mongo.Collection, err error) {
	store := remote.GetStoreFromContext(ctx)
	if store == nil {
		err = errors.New("no connect to mongodb")
	} else {
		databaseName := "testing"
		collectionName := "atest"

		if store.Properties != nil {
			if v, ok := store.Properties["database"]; ok && v != "" {
				databaseName = v
			}
			if v, ok := store.Properties["collection"]; ok && v != "" {
				collectionName = v
			}
		}

		address := fmt.Sprintf("mongodb://%s:%s@%s", store.Username, store.Password, store.URL)
		var client *mongo.Client
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(address))
		// TODO need a way to maintain the connection of it
		if err == nil {
			collection = client.Database(databaseName).Collection(collectionName)
			if collection == nil {
				if err = client.Database(databaseName).CreateCollection(ctx, collectionName); err == nil {
					collection = client.Database(databaseName).Collection(collectionName)
				}
			}
		}
	}
	return
}
