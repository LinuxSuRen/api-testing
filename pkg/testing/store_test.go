/**
MIT License

Copyright (c) 2023 Rick

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

package testing

import (
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
				Name: "database",
				URL:  "localhost:7071",
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
		assert.Equal(t, "local", stores[1].Name)
	})

	t.Run("DeleteStore", func(t *testing.T) {
		err := factory.DeleteStore("")
		assert.NoError(t, err)
	})

	t.Run("UpdateStore", func(t *testing.T) {
		err := factory.UpdateStore(Store{})
		assert.NoError(t, err)
	})

	t.Run("no stores.yaml found", func(t *testing.T) {
		factory := NewStoreFactory("testdata-fake")
		stores, err := factory.GetStores()
		assert.NoError(t, err)
		assert.Equal(t, []Store{{
			Name: "local",
		}}, stores)
	})
}

var sampleStoreMap = map[string]string{
	"name":        "test",
	"url":         fooURL,
	"kind.url":    fooURL,
	"kind":        "test",
	"description": "desc",
	"username":    "user",
	"password":    "pass",
	"pro.key":     "val",
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
