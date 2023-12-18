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
package server

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestToGRPCStore(t *testing.T) {
	assert.Equal(t, &Store{
		Name:        "test",
		Owner:       "rick",
		Description: "desc",
		Kind: &StoreKind{
			Name: "test",
			Url:  urlFoo,
		},
		Url:      urlFoo,
		Username: "user",
		Password: "pass",
		Disabled: true,
		Properties: []*Pair{{
			Key: "foo", Value: "bar",
		}},
	}, ToGRPCStore(atest.Store{
		Name:        "test",
		Owner:       "rick",
		Description: "desc",
		Kind: atest.StoreKind{
			Name: "test",
			URL:  urlFoo,
		},
		URL:      urlFoo,
		Username: "user",
		Password: "pass",
		Disabled: true,
		Properties: map[string]string{
			"foo": "bar",
		},
	}))
}

func TestToNormalStore(t *testing.T) {
	assert.Equal(t, atest.Store{
		Name:        "test",
		Description: "desc",
		Kind: atest.StoreKind{
			Name: "test",
			URL:  urlFoo,
		},
		URL:      urlFoo,
		Username: "user",
		Password: "pass",
		Properties: map[string]string{
			"foo": "bar",
		},
	}, ToNormalStore(&Store{
		Name:        "test",
		Description: "desc",
		Kind: &StoreKind{
			Name: "test",
			Url:  urlFoo,
		},
		Url:      urlFoo,
		Username: "user",
		Password: "pass",
		Properties: []*Pair{{
			Key: "foo", Value: "bar",
		}},
	}))
}

func TestToGRPCSuite(t *testing.T) {
	assert.Equal(t, &TestSuite{
		Name: "test",
		Api:  "api",
		Param: []*Pair{{
			Key: "foo", Value: "bar",
		}},
		Spec: &APISpec{
			Secure: &Secure{
				Insecure: true,
			},
			Rpc: &RPC{
				Raw: "raw",
			},
		},
	}, ToGRPCSuite(&atest.TestSuite{
		Name: "test",
		API:  "api",
		Param: map[string]string{
			"foo": "bar",
		},
		Spec: atest.APISpec{
			Secure: &atest.Secure{
				Insecure: true,
			},
			RPC: &atest.RPCDesc{
				Raw: "raw",
			},
		},
	}))
}

func TestToNormalSuite(t *testing.T) {
	assert.Equal(t, &atest.TestSuite{
		Name: "test",
		API:  "api",
		Param: map[string]string{
			"foo": "bar",
		},
		Spec: atest.APISpec{
			Secure: &atest.Secure{
				Insecure: true,
			},
			RPC: &atest.RPCDesc{
				Raw: "raw",
			},
		},
	}, ToNormalSuite(&TestSuite{
		Name: "test",
		Api:  "api",
		Param: []*Pair{{
			Key: "foo", Value: "bar",
		}},
		Spec: &APISpec{
			Secure: &Secure{
				Insecure: true,
			},
			Rpc: &RPC{
				Raw: "raw",
			},
		},
	}))
}
