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
