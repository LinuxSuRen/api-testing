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
import type { Pair } from './types'

interface Store {
  name: string;
  params: Pair[];
}

const storeExtensions = [
  {
    name: 'atest-store-git',
    params: [{
      key: 'insecure'
    }, {
      key: 'timeout'
    }, {
      key: 'targetpath'
    }, {
      key: 'url'
    }, {
      key: 'email'
    }, {
      key: 'name'
    }]
  },
  {
    name: 'atest-store-s3',
    params: [{
      key: 'accesskeyid'
    }, {
      key: 'secretaccesskey'
    }, {
      key: 'sessiontoken'
    }, {
      key: 'region'
    }, {
      key: 'disablessl'
    }, {
      key: 'forcepathstyle'
    }, {
      key: 'bucket'
    }]
  },
  {
    name: 'atest-store-orm',
    params: [{
      key: 'driver'
    }, {
      key: 'database'
    }]
  },
  {
    name: 'atest-store-etcd',
    params: []
  },
  {
    name: 'atest-store-mongodb',
    params: [{
      key: 'collection'
    }, {
      key: 'database'
    }]
  }
] as Store[]

export function SupportedExtensions() {
    return storeExtensions
}

export function SupportedExtension(name: string) {
    return storeExtensions.find(e => e.name === name)
}
