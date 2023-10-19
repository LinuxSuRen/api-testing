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

package server

import "github.com/linuxsuren/api-testing/pkg/testing"

// ToGRPCStore convert the normal store to GRPC store
func ToGRPCStore(store testing.Store) (result *Store) {
	result = &Store{
		Name: store.Name,
		Kind: &StoreKind{
			Name:    store.Kind.Name,
			Url:     store.Kind.URL,
			Enabled: store.Kind.Enabled,
		},
		Description: store.Description,
		Url:         store.URL,
		Username:    store.Username,
		Password:    store.Password,
		Properties:  mapToPair(store.Properties),
	}
	return
}

// ToNormalStore convert the GRPC store to normal store
func ToNormalStore(store *Store) (result testing.Store) {
	result = testing.Store{
		Name:        store.Name,
		Description: store.Description,
		URL:         store.Url,
		Username:    store.Username,
		Password:    store.Password,
		Properties:  pairToMap(store.Properties),
	}
	if store.Kind != nil {
		result.Kind = testing.StoreKind{
			Name: store.Kind.Name,
			URL:  store.Kind.Url,
		}
	}
	return
}
