/*
Copyright 2023-2025 API Testing Authors.

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
import { API } from './net'

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
const ExtensionKindAI = "atest-ext-ai"

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
    ExtensionKindOpengeMini,
    ExtensionKindAI
}

const MySQL = "mysql";
const Postgres = "postgres";
const SQLite = "sqlite";
const TDengine = "tdengine";
const Cassandra = "cassandra";

export const GetDriverName = (store: Store): string => {
    switch (store.kind.name) {
    case 'atest-store-orm':
        const properties = store.properties as Pair[];
        for (let i = 0; i < properties.length; i++) {
            if (properties[i].key === 'driver') {
                return properties[i].value || MySQL;
            }
        }
        return MySQL;
    case 'atest-store-cassandra':
        return Cassandra;
    }
    return ""
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
    },
    {
        name: ExtensionKindAI,
        params: [{
            key: 'model_provider',
            defaultValue: 'ollama',
            enum: ['ollama', 'openai'],
            description: 'AI model provider used in GenerateContentRequest.parameters for content generation'
        }, {
            key: 'model_name',
            defaultValue: 'llama3.2',
            description: 'AI model name passed to GenerateContentRequest.parameters["model_name"] for specific model selection'
        }, {
            key: 'api_key',
            defaultValue: '',
            description: 'API key for external providers, used in GenerateContentRequest.parameters["api_key"] for authentication'
        }, {
            key: 'base_url',
            defaultValue: 'http://localhost:11434',
            description: 'Base URL for AI service, used in GenerateContentRequest.parameters["base_url"] to specify service endpoint'
        }, {
            key: 'temperature',
            defaultValue: '0.7',
            description: 'Controls randomness in AI responses (0.0-1.0), passed via GenerateContentRequest.parameters["temperature"]'
        }, {
            key: 'max_tokens',
            defaultValue: '2048',
            description: 'Maximum tokens in AI response, used in GenerateContentRequest.parameters["max_tokens"] to limit output length'
        }],
        link: 'https://github.com/LinuxSuRen/atest-ext-ai'
    }
]

export const Driver = {
    MySQL, Postgres, SQLite, TDengine, Cassandra
}

export async function SupportedExtensions() {
    const kinds = await API.GetStoreKinds()
    if (kinds) {
        return kinds.data
    }
    return []
}

export async function SupportedExtension(name: string) {
    const storeExtensions = await SupportedExtensions()
    return storeExtensions.find((e: Store) => e.name === name)
}
