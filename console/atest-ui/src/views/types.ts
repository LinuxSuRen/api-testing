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
import { ref } from 'vue'
import _ from 'lodash'
import { API } from './net'

export interface Suite {
  name: string
  api: string
  param: Pair[]
  spec: {
    kind: string
    url: string
    secure: {
        insecure: boolean
    }
  }
  proxy: {
    http: string
    https: string
    no: string
  }
}

export interface TestResult {
  body: string
  bodyObject: {}
  bodyText: string
  bodyLength: number
  output: string
  error: string
  statusCode: number
  header: Pair[]

  // inner fields
  originBodyObject:{}
}

export interface Pair {
  key: string
  value: string
  type: string
  defaultValue: string
  description: string
}

export interface TestCaseWithSuite {
  suiteName: string
  data: TestCase
}

export interface TestCase {
  name: string
  server: string
  request: TestCaseRequest
  response: TestCaseResponse
}

export interface TestCaseRequest {
  api: string
  method: string
  header: Pair[]
  cookie: Pair[]
  query: Pair[]
  form: Pair[]
  body: string
  filepath: string
}

export interface TestCaseResponse {
  statusCode: number
  body: string
  header: Pair[]
  bodyFieldsExpect: Pair[]
  verify: string[]
  schema: string
}

// Suggested APIs query
const localCache = ref({} as TestCaseWithValue[])
export function NewSuggestedAPIsQuery(store: string, suite: string) {
  return function (queryString: string, cb: (arg: any) => void) {
    loadCache(store, suite, function () {
      const results = queryString
        ? localCache.value.filter(CreateFilter(queryString))
        : localCache.value

      cb(results.slice(0, 10))
    })
  }
}
export function CreateFilter(queryString: string) {
  return (v: Pair) => {
    return v.value.toLowerCase().indexOf(queryString.toLowerCase()) !== -1
  }
}
function loadCache(store: string, suite: string, callback: Function) {
  if (localCache.value.length > 0) {
    callback()
    return
  }

  if (suite === '') {
    return
  }

  API.GetSuggestedAPIs(suite, (e) => {
    localCache.value = e.data
    localCache.value.forEach((v: TestCaseWithValue) => {
      v.value = v.request.api
    })
    callback()
  })
}

interface TestCaseWithValue extends TestCase, Pair {}

export function GetHTTPMethods() {
  return [
    {
      value: 'GET',
      key: 'GET',
      type: ''
    },
    {
      value: 'POST',
      key: 'POST',
      type: 'success'
    },
    {
      value: 'DELETE',
      key: 'DELETE',
      type: 'danger'
    },
    {
      value: 'PUT',
      key: 'PUT',
      type: 'warning'
    },
    {
      value: 'HEAD',
      key: 'HEAD',
      type: ''
    },
    {
      value: 'PATCH',
      key: 'PATCH',
      type: ''
    },
    {
      value: 'OPTIONS',
      key: 'OPTIONS',
      type: ''
    }
  ] as Pair[]
}

export function GetHTTPMethod(name: string) {
  for (const method of GetHTTPMethods()) {
    if (method.key === name) {
      return method
    }
  }
  return {} as Pair
}

export function FlattenObject(obj: any): any {
  function _flattenPairs(obj: any, prefix: string): [string, any][] {
    if (!_.isObject(obj)) {
      return [prefix, obj]
    }

    return _.toPairs(obj).reduce((final: [string, any][], nPair: [string, any]) => {
      const flattened = _flattenPairs(nPair[1], `${prefix}.${nPair[0]}`)
      if (flattened.length === 2 && !_.isObject(flattened[0]) && !_.isObject(flattened[1])) {
        return final.concat([flattened as [string, any]])
      } else {
        return final.concat(flattened)
      }
    }, [])
  }

  if (!_.isObject(obj)) {
    return JSON.stringify(obj)
  }

  const pairs: [string, any][] = _.toPairs(obj).reduce(
    (final: [string, any][], pair: [string, any]) => {
      const flattened = _flattenPairs(pair[1], pair[0])
      if (flattened.length === 2 && !_.isObject(flattened[0]) && !_.isObject(flattened[1])) {
        return final.concat([flattened as [string, any]])
      } else {
        return final.concat(flattened)
      }
    },
    []
  )

  return pairs.reduce((acc: any, pair: [string, any]) => {
    acc[pair[0]] = pair[1]
    return acc
  }, {})
}
