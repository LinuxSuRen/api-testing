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
