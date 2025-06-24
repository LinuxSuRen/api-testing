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

export interface Store {
  name: string;
  link: string;
  ready: boolean;
  kind: {
    name: string;
    description: string;
  };
  properties: Pair[];
  params: Pair[];
}

const MySQL = "mysql";
const Postgres = "postgres";
const SQLite = "sqlite";
const TDengine = "tdengine";
const Cassandra = "cassandra";

export const GetDriverName = (store: Store): string => {
    switch (store.kind.name) {
    case 'atest-store-orm':
        return store.properties.find((p: Pair) => p.key === 'driver')?.value || MySQL;
    case 'atest-store-cassandra':
        return Cassandra;
    }
    return ""
}

export const Driver = {
    MySQL, Postgres, SQLite, TDengine, Cassandra
}

const ExtensionKindGit = "atest-store-git"
const ExtensionKindS3 = "atest-store-s3"
const ExtensionKindORM = "atest-store-orm"
const ExtensionKindIotDB = "atest-store-iotdb"
const ExtensionKindCassandra = "atest-store-cassandra"
const ExtensionKindEtcd = "atest-store-etcd"
const ExtensionKindRedis = "atest-store-redis"
const ExtensionKindMongoDB = "atest-store-mongodb"
const ExtensionKindElasticsearch = "atest-store-elasticsearch"
const ExtensionKindOpengeMini = "atest-store-opengemini"

export const ExtensionKind = {
    ExtensionKindGit,
    ExtensionKindS3,
    ExtensionKindORM,
    ExtensionKindIotDB,
    ExtensionKindCassandra,
    ExtensionKindEtcd,
    ExtensionKindRedis,
    ExtensionKindMongoDB,
    ExtensionKindElasticsearch,
    ExtensionKindOpengeMini
}

const storeExtensions = [
    {
        name: ExtensionKindGit,
        params: [{
            key: 'insecure'
        }, {
            key: 'timeout'
        }, {
            key: 'targetpath'
        }, {
            key: 'branch'
        }, {
            key: 'email',
            description: 'See also: git config --local user.email xxx@xxx.com'
        }, {
            key: 'name',
            description: 'See also: git config --local user.name xxx'
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-git'
    },
    {
        name: ExtensionKindS3,
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
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-s3'
    },
    {
        name: ExtensionKindORM,
        params: [{
            key: 'driver',
            defaultValue: 'mysql',
            enum: [MySQL, Postgres, SQLite, TDengine],
            description: 'Supported: mysql, postgres, sqlite, tdengine'
        }, {
            key: 'database',
            defaultValue: 'atest'
        }, {
            key: 'historyLimit',
            defaultValue: '',
            // type: 'number',
            description: 'Set the limit of the history record count'
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-orm'
    },
    {
        name: ExtensionKindIotDB,
        params: [],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-iotdb'
    },
    {
        name: ExtensionKindCassandra,
        params: [{
            key: 'keyspace',
            defaultValue: ''
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-cassandra'
    },
    {
        name: ExtensionKindEtcd,
        params: [],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-etcd'
    },
    {
        name: ExtensionKindRedis,
        params: [],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-redis'
    },
    {
        name: ExtensionKindMongoDB,
        params: [{
            key: 'collection'
        }, {
            key: 'database',
            defaultValue: 'atest'
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-mongodb'
    },
    {
        name: ExtensionKindElasticsearch,
        params: [],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-elasticsearch'
    },
    {
        name: ExtensionKindOpengeMini,
        params: [],
        link: 'https://github.com/LinuxSuRen/atest-ext-store-opengemini'
    }
] as Store[]

export function SupportedExtensions() {
    return storeExtensions
}

export function SupportedExtension(name: string) {
    return storeExtensions.find(e => e.name === name)
}
