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

	mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type fakeKV struct {
	data map[string]string
}

func (f *fakeKV) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	r, err := f.Do(ctx, clientv3.OpPut(key, val, opts...))
	return r.Put(), err
}

func (f *fakeKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (resp *clientv3.GetResponse, err error) {
	Kvs := []*mvccpb.KeyValue{}
	for k, v := range f.data {
		Kvs = append(Kvs, &mvccpb.KeyValue{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	resp = &clientv3.GetResponse{
		Kvs: Kvs,
	}
	if len(Kvs) == 0 {
		err = errors.New("not found")
	}
	return
}

func (f *fakeKV) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	r, err := f.Do(ctx, clientv3.OpDelete(key, opts...))
	return r.Del(), err
}

func (f *fakeKV) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, nil
}

func (f *fakeKV) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	switch {
	case op.IsPut():
		f.data[string(op.KeyBytes())] = string(op.ValueBytes())
	case op.IsDelete():
		delete(f.data, string(op.KeyBytes()))
	}
	return clientv3.OpResponse{}, nil
}

func (f *fakeKV) Txn(ctx context.Context) clientv3.Txn {
	return nil
}

func (f *fakeKV) Close() error {
	return nil
}

func (f *fakeKV) New(cfg clientv3.Config) (SimpleKV, error) {
	return f, nil
}
