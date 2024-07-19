/*
Copyright 2023-2024 API Testing Authors.

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

package testing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreConvert(t *testing.T) {
	t.Run("ToMap", func(t *testing.T) {
		store := sampleStore
		assert.Equal(t, sampleStoreMap, store.ToMap())
	})

	t.Run("MapToStore", func(t *testing.T) {
		store := MapToStore(sampleStoreMap)
		assert.Equal(t, sampleStore, &store)
	})

	t.Run("NewStoreFactory", func(t *testing.T) {
		assert.NotNil(t, NewStoreFactory(""))
	})
}

func TestStoreFactory(t *testing.T) {
	factory := NewStoreFactory("testdata")
	assert.NotNil(t, factory)

	t.Run("GetStoreKinds", func(t *testing.T) {
		_, err := factory.GetStoreKinds()
		assert.NoError(t, err)
	})

	t.Run("GetStore", func(t *testing.T) {
		store, err := factory.GetStore("db")
		assert.Nil(t, err)
		assert.Equal(t, &Store{
			Name: "db",
			Kind: StoreKind{
				Name:    "database",
				URL:     "localhost:7071",
				Enabled: true,
			},
			URL:      "localhost:4000",
			Username: "root",
			Properties: map[string]string{
				"database": "test",
			},
		}, store)
	})

	t.Run("GetAllStores", func(t *testing.T) {
		stores, err := factory.GetStores()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(stores))
		assert.Equal(t, "db", stores[0].Name)
	})

	t.Run("Get stores by owner", func(t *testing.T) {
		stores, err := factory.GetStoresByOwner("rick")
		assert.NoError(t, err)
		assert.Len(t, stores, 1)
		assert.Equal(t, "git", stores[0].Name)

		stores, err = factory.GetStoresByOwner("fake")
		assert.NoError(t, err)
		assert.Len(t, stores, 0)
	})

	t.Run("DeleteStore", func(t *testing.T) {
		err := factory.DeleteStore("")
		assert.NoError(t, err)
	})

	t.Run("UpdateStore", func(t *testing.T) {
		err := factory.UpdateStore(Store{Name: "fake"})
		assert.Error(t, err)
	})

	t.Run("no stores.yaml found", func(t *testing.T) {
		factory := NewStoreFactory("testdata-fake")
		stores, err := factory.GetStores()
		assert.NoError(t, err)
		assert.Nil(t, stores)
	})

	t.Run("CreateStore", func(t *testing.T) {
		dir, err := os.MkdirTemp(os.TempDir(), "store")
		assert.NoError(t, err)
		defer os.RemoveAll(dir)

		factory := NewStoreFactory(dir)
		err = factory.CreateStore(Store{Name: "fake"})
		assert.NoError(t, err)

		// create an existing store
		err = factory.CreateStore(Store{Name: "fake"})
		assert.Error(t, err)

		// update an existing store
		err = factory.UpdateStore(Store{Name: "fake"})
		assert.NoError(t, err)

		// delete an existing store
		err = factory.DeleteStore("fake")
		assert.NoError(t, err)

		// get all stores
		var stores []Store
		stores, err = factory.GetStores()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(stores))
	})
}

var sampleStoreMap = map[string]string{
	"name":        "test",
	"owner":       "",
	"url":         fooURL,
	"kind.url":    fooURL,
	"kind":        "test",
	"description": "desc",
	"disabled":    "false",
	"username":    "user",
	"password":    "pass",
	"pro.key":     "val",
	"readonly":    "false",
}

var sampleStore = &Store{
	Name: "test",
	Kind: StoreKind{
		Name: "test",
		URL:  fooURL,
	},
	URL:         fooURL,
	Description: "desc",
	Username:    "user",
	Password:    "pass",
	Properties: map[string]string{
		"key": "val",
	},
}

const fooURL = "http://foo"
