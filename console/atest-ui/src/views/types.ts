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
  }
}

export interface TestResult {
  body: string
  bodyObject: {}
  output: string
  error: string
  statusCode: number
  header: Pair[]

  // inner fileds
  originBodyObject:{}
}

export interface Pair {
  key: string
  value: string
}

export interface TestCaseWithSuite {
  suiteName: string
  data: TestCase
}

export interface TestCase {
  name: string
  request: TestCaseRequest
  response: TestCaseResponse
}

export interface TestCaseRequest {
  api: string
  method: string
  header: Pair[]
  query: Pair[]
  form: Pair[]
  body: string
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
      key: 'GET'
    },
    {
      value: 'POST',
      key: 'POST'
    },
    {
      value: 'DELETE',
      key: 'DELETE'
    },
    {
      value: 'PUT',
      key: 'PUT'
    },
    {
      value: 'HEAD',
      key: 'HEAD'
    },
    {
      value: 'PATCH',
      key: 'PATCH'
    },
    {
      value: 'OPTIONS',
      key: 'OPTIONS'
    }
  ] as Pair[]
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