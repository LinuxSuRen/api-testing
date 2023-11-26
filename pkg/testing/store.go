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
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type StoreConfig struct {
	Stores  []Store     `yaml:"stores"`
	Plugins []StoreKind `yaml:"plugins"`
}

type Store struct {
	Name        string
	Owner       string
	Kind        StoreKind
	Description string
	URL         string
	Username    string
	Password    string
	ReadOnly    bool
	Disabled    bool
	Properties  map[string]string
}

func (s *Store) ToMap() (result map[string]string) {
	result = map[string]string{
		"name":        s.Name,
		"owner":       s.Owner,
		"kind":        s.Kind.Name,
		"kind.url":    s.Kind.URL,
		"description": s.Description,
		"url":         s.URL,
		"username":    s.Username,
		"password":    s.Password,
		"readonly":    fmt.Sprintf("%t", s.ReadOnly),
		"disabled":    fmt.Sprintf("%t", s.Disabled),
	}
	for key, val := range s.Properties {
		result["pro."+key] = val
	}
	return
}

func MapToStore(data map[string]string) (store Store) {
	store = Store{
		Name:        data["name"],
		Owner:       data["owner"],
		Description: data["description"],
		URL:         data["url"],
		Username:    data["username"],
		Password:    data["password"],
		ReadOnly:    data["readonly"] == "true",
		Disabled:    data["disabled"] == "true",
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
	Name    string
	URL     string
	Enabled bool
}

type StoreGetterAndSetter interface {
	GetStores() (stores []Store, err error)
	GetStoresByOwner(owner string) (stores []Store, err error)
	GetStore(name string) (store *Store, err error)
	DeleteStore(name string) (err error)
	UpdateStore(store Store) (err error)
	CreateStore(store Store) (err error)

	GetStoreKinds() (kinds []StoreKind, err error)
}

type StoreWriterFactory interface {
	NewInstance(store Store) (writer Writer, err error)
}

type storeFactory struct {
	configDir string
}

// NewStoreFactory creates a new store factory
func NewStoreFactory(configDir string) StoreGetterAndSetter {
	return &storeFactory{
		configDir: configDir,
	}
}

// GetStores returns all stores
func (s *storeFactory) GetStores() (stores []Store, err error) {
	var storeConfig *StoreConfig
	if storeConfig, err = s.getStoreConfig(); err == nil {
		stores = storeConfig.Stores
	}
	return
}

func (s *storeFactory) GetStoresByOwner(owner string) (stores []Store, err error) {
	var all []Store
	all, err = s.GetStores()
	if owner == "" {
		return all, err
	}
	if err == nil {
		for _, item := range all {
			if item.Owner != owner {
				continue
			}

			stores = append(stores, item)
		}
	}
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
	var storeConfig *StoreConfig
	if storeConfig, err = s.getStoreConfig(); err == nil {
		for i := range storeConfig.Stores {
			item := storeConfig.Stores[i]
			if item.Name == name {
				storeConfig.Stores = append(storeConfig.Stores[:i], storeConfig.Stores[i+1:]...)
				break
			}
		}
		err = s.save(storeConfig)
	}
	return
}

func (s *storeFactory) UpdateStore(store Store) (err error) {
	var storeConfig *StoreConfig
	if storeConfig, err = s.getStoreConfig(); err == nil {
		exist := false
		for i := range storeConfig.Stores {
			item := storeConfig.Stores[i]
			if item.Name == store.Name {
				storeConfig.Stores[i] = store
				exist = true
				break
			}
		}

		if exist {
			err = s.save(storeConfig)
		} else {
			err = fmt.Errorf("store %s is not exists", store.Name)
		}
	}
	return
}

func (s *storeFactory) CreateStore(store Store) (err error) {
	var storeConfig *StoreConfig
	if storeConfig, err = s.getStoreConfig(); err == nil {
		exist := false
		for i := range storeConfig.Stores {
			item := storeConfig.Stores[i]
			if item.Name == store.Name && item.Owner == store.Owner {
				exist = true
				break
			}
		}

		if !exist {
			storeConfig.Stores = append(storeConfig.Stores, store)
			err = s.save(storeConfig)
		} else {
			err = fmt.Errorf("store %s already exists", store.Name)
		}
	}
	return
}

func (s *storeFactory) save(storeConfig *StoreConfig) (err error) {
	for i, item := range storeConfig.Stores {
		if item.Kind.Name != "" {
			storeConfig.Stores[i].Kind.Enabled = true

			foundPlugin := false
			for j, kind := range storeConfig.Plugins {
				if kind.Name == item.Kind.Name {
					foundPlugin = true
					storeConfig.Plugins[j].Enabled = true
					break
				}
			}

			if !foundPlugin {
				storeConfig.Plugins = append(storeConfig.Plugins, storeConfig.Stores[i].Kind)
			}
		}
	}

	if err = os.MkdirAll(s.configDir, 0755); err == nil {
		var data []byte
		if data, err = yaml.Marshal(storeConfig); err == nil {
			err = os.WriteFile(path.Join(s.configDir, "stores.yaml"), data, 0644)
		}
	}
	return
}

func (s *storeFactory) GetStoreKinds() (kinds []StoreKind, err error) {
	var storeConfig *StoreConfig
	if storeConfig, err = s.getStoreConfig(); err == nil {
		kinds = storeConfig.Plugins
	}
	return
}

func (s *storeFactory) getStoreConfig() (storeConfig *StoreConfig, err error) {
	storeConfig = &StoreConfig{}
	var data []byte
	if data, err = os.ReadFile(path.Join(s.configDir, "stores.yaml")); err == nil {
		err = yaml.Unmarshal(data, storeConfig)
	} else {
		err = nil
	}
	return
}
