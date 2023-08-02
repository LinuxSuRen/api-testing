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
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type Store struct {
	Name        string
	Kind        StoreKind
	Description string
	URL         string
	Username    string
	Password    string
	Properties  map[string]string
}

func (s *Store) ToMap() (result map[string]string) {
	result = map[string]string{
		"name":        s.Name,
		"kind":        s.Kind.Name,
		"kind.url":    s.Kind.URL,
		"description": s.Description,
		"url":         s.URL,
		"username":    s.Username,
		"password":    s.Password,
	}
	for key, val := range s.Properties {
		result["pro."+key] = val
	}
	return
}

func (s *Store) IsLocal() bool {
	return s.Name == "local"
}

func MapToStore(data map[string]string) (store Store) {
	store = Store{
		Name:        data["name"],
		Description: data["description"],
		URL:         data["url"],
		Username:    data["username"],
		Password:    data["password"],
		Kind: StoreKind{
			Name: data["kind"],
			URL:  data["kind.url"],
		},
		Properties: make(map[string]string),
	}
	for key, val := range data {
		if strings.HasPrefix(key, "pro.") {
			store.Properties[strings.TrimPrefix(key, "pro.")] = val
		}
	}
	return
}

// StoreKind represents a gRPC-based store
type StoreKind struct {
	Name string
	URL  string
}

type StoreGetterAndSetter interface {
	GetStores() (stores []Store, err error)
	GetStore(name string) (store *Store, err error)
	DeleteStore(name string) (err error)
	UpdateStore(store Store) (err error)

	GetStoreKinds() (kinds []StoreKind, err error)
}

type StoreWriterFactory interface {
	NewInstance(store Store) (writer Writer, err error)
}

type storeFactory struct {
	configDir string
}

func NewStoreFactory(configDir string) StoreGetterAndSetter {
	return &storeFactory{
		configDir: configDir,
	}
}

func (s *storeFactory) GetStores() (stores []Store, err error) {
	var data []byte
	if data, err = os.ReadFile(path.Join(s.configDir, "stores.yaml")); err == nil {
		err = yaml.Unmarshal(data, &stores)
	} else {
		err = nil
	}
	stores = append(stores, Store{Name: "local"})
	return
}

func (s *storeFactory) GetStore(name string) (store *Store, err error) {
	var stores []Store
	if stores, err = s.GetStores(); err == nil {
		for i := range stores {
			item := stores[i]
			if item.Name == name {
				store = &item
			}
		}
	}
	return
}

func (s *storeFactory) DeleteStore(name string) (err error) {
	return
}

func (s *storeFactory) UpdateStore(store Store) (err error) {
	return
}

func (s *storeFactory) GetStoreKinds() (kinds []StoreKind, err error) {
	return
}
