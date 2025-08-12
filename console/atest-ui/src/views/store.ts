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

export async function SupportedExtensions() {
    const kinds = await API.GetStoreKinds()
    return kinds.data
}

export async function SupportedExtension(name: string) {
    const storeExtensions = await SupportedExtensions()
    return storeExtensions.find(e => e.name === name)
}
