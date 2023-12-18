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
	clientv3 "go.etcd.io/etcd/client/v3"
)

type SimpleKV interface {
	clientv3.KV
	Close() error
}

type KVFactory interface {
	New(cfg clientv3.Config) (SimpleKV, error)
}

type realEtcd struct{}

func NewRealEtcd() KVFactory {
	return &realEtcd{}
}

func (r *realEtcd) New(cfg clientv3.Config) (SimpleKV, error) {
	return clientv3.New(cfg)
}
