/*
Copyright 2024 API Testing Authors.

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

import { post } from '../axios'
import type { TestSuite, ImportSource, TestCase, RunTestCaseRequest } from '../common'
import { getToken } from '../../utils/auth/token'
import { Cache } from '../../utils/cache'

const stroeName = Cache.GetCurrentStore().name

export const CreateTestSuite = (params: TestSuite) =>
  post(
    '/server.Runner/CreateTestSuite',
    {
      name: params.name,
      api: params.api,
      kind: params.kind
    },
    {
      'X-Store-Name': params.store,
      'X-Auth': getToken()
    }
  )

export const UpdateTestSuite = (suite: any) =>
  post('/server.Runner/UpdateTestSuite', suite, {
    'X-Store-Name': stroeName,
    'X-Auth': getToken()
  })

export const GetTestSuite = (name: string) =>
  post(
    '/server.Runner/GetTestSuite',
    { name: name },
    {
      'X-Store-Name': stroeName,
      'X-Auth': getToken()
    }
  )

export const DeleteTestSuite = (name: string) =>
  post('/server.Runner/DeleteTestSuite', name, {
    'X-Store-Name': stroeName,
    'X-Auth': getToken()
  })

export const ConvertTestSuite = (params: any) =>
  post('/server.Runner/ConvertTestSuite', params, {
    'X-Store-Name': Cache.GetCurrentStore().name,
    'X-Auth': getToken()
  })

export const ImportTestSuite = (params: ImportSource) =>
  post('/server.Runner/ImportTestSuite', params, {
    'X-Store-Name': params.store,
    'X-Auth': getToken()
  })

export const LoadTestSuite = (params: any) =>
  post('/server.Runner/GetSuites', null, {
    'X-Store-Name': params,
    'X-Auth': getToken()
  })

export const CreateTestCase = (params: any) =>
  post('/server.Runner/CreateTestCase', params, {
    'X-Store-Name': Cache.GetCurrentStore().name,
    'X-Auth': getToken()
  })

export const UpdateTestCase = (params: any) =>
  post('/server.Runner/UpdateTestCase', params, {
    'X-Store-Name': stroeName,
    'X-Auth': getToken()
  })

export const GetTestCase = (params: TestCase) =>
  post('/server.Runner/GetTestCase', params, {
    'X-Store-Name': Cache.GetCurrentStore().name,
    'X-Auth': getToken()
  })

export const ListTestCase = (params: any, store: any) =>
  post(
    '/server.Runner/ListTestCase',
    {
      name: params
    },
    {
      'X-Store-Name': store,
      'X-Auth': getToken()
    }
  )

export const DeleteTestCase = (params: TestCase) =>
  post('/server.Runner/DeleteTestCase', params, {
    'X-Store-Name': stroeName,
    'X-Auth': getToken()
  })

export const RunTestCase = (params: RunTestCaseRequest) =>
  post('/server.Runner/RunTestCase', params, {
    'X-Store-Name': stroeName,
    'X-Auth': getToken()
  })
