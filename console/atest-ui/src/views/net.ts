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

import { Cache } from './cache'

function DefaultResponseProcess(response: any) {
  if (!response.ok) {
    switch (response.status) {
      case 401:
        throw new Error("Unauthenticated")
    }
    throw new Error(response.statusText)
  } else {
    return response.json()
  }
}

interface AppVersion {
  message: string
}

function safeToggleFunc(toggle?: (e: boolean) => void) {
  if (!toggle) {
    return (e: boolean) => {}
  }
  return toggle
}

function GetVersion(callback: (v: AppVersion) => void) {
  const requestOptions = {
    method: 'POST',
  }
  fetch('/server.Runner/GetVersion', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback)
}

interface TestSuite {
  store: string
  name: string
  api: string
  kind: string
}

function CreateTestSuite(suite: TestSuite,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': suite.store,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: suite.name,
      api: suite.api,
      kind: suite.kind
    })
  }

  fetch('/server.Runner/CreateTestSuite', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

interface ImportSource {
  store: string
  url: string
}

function UpdateTestSuite(suite: any,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify(suite)
  }
  fetch('/server.Runner/UpdateTestSuite', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function GetTestSuite(name: string,
  callback: (d: any) => void, errHandle: (e: any) => void) {
  const store = Cache.GetCurrentStore()
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': store.name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: name
    })
  }
  fetch('/server.Runner/GetTestSuite', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteTestSuite(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: name
    })
  }
  fetch('/server.Runner/DeleteTestSuite', requestOptions)
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

function ConvertTestSuite(suiteName: string, genertor: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      Generator: genertor,
      TestSuite: suiteName
    })
  }
  fetch('/server.Runner/ConvertTestSuite', requestOptions)
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

function ImportTestSuite(source: ImportSource, callback: (d: any) => void) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': source.store,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      url: source.url
    })
  }

  fetch('/server.Runner/ImportTestSuite', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback)
}

interface TestCase {
  suiteName: string
  name: string
  api: string
  method: string
}

function CreateTestCase(testcase: TestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      suiteName: testcase.suiteName,
      data: {
        name: testcase.name,
        request: {
          api: testcase.api,
          method: testcase.method
        }
      }
    })
  }

  fetch('/server.Runner/CreateTestCase', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function UpdateTestCase(testcase: any,
  callback: (d: any) => void, errHandle?: (e: any) => void | null,
  toggle?: (e: boolean) => void) {
    const requestOptions = {
      method: 'POST',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify(testcase)
    }
    safeToggleFunc(toggle)(true)
    fetch('/server.Runner/UpdateTestCase', requestOptions)
      .then(DefaultResponseProcess)
      .then(callback).catch(errHandle)
      .finally(() => {
        safeToggleFunc(toggle)(false)
      })
}

function GetTestCase(req: TestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      suite: req.suiteName,
      testcase: req.name
    })
  }
  fetch('/server.Runner/GetTestCase', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function ListTestCase(suite: string, store: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': store,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: suite
    })
  }
  fetch('/server.Runner/ListTestCase', requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteTestCase(testcase: TestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    const requestOptions = {
      method: 'POST',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify({
        suite: testcase.suiteName,
        testcase: testcase.name
      })
    }
    fetch('/server.Runner/DeleteTestCase', requestOptions)
      .then(callback).catch(errHandle)
}

interface RunTestCaseRequest {
  suiteName: string
  name: string
  parameters: any
}

function RunTestCase(request: RunTestCaseRequest,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      suite: request.suiteName,
      testcase: request.name,
      parameters: request.parameters
    })
  }
  fetch('/server.Runner/RunTestCase', requestOptions)
  .then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

interface GenerateRequest {
  suiteName: string
  name: string
  generator: string
}

function GenerateCode(request: GenerateRequest,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      TestSuite: request.suiteName,
      TestCase: request.name,
      Generator: request.generator
    })
  }
  fetch('/server.Runner/GenerateCode', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function ListCodeGenerator(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  fetch('/server.Runner/ListCodeGenerator', {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
  }).then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

function PopularHeaders(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
  }
  fetch('/server.Runner/PopularHeaders', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function GetStores(callback: (d: any) => void,
  errHandle?: (e: any) => void | null, final?: () => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
  }
  fetch('/server.Runner/GetStores', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle).finally(final)
}

function DeleteStore(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: name
    })
  }
  fetch('/server.Runner/DeleteStore', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function VerifyStore(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: name
    })
  }
  
  let api = '/server.Runner/VerifyStore'
  fetch(api, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function getToken() {
  const token = sessionStorage.getItem('token')
  if (!token) {
    return ''
  }
  return token
}

export const API = {
  DefaultResponseProcess,
  GetVersion,
  CreateTestSuite, UpdateTestSuite, ImportTestSuite, GetTestSuite, DeleteTestSuite, ConvertTestSuite,
  CreateTestCase, UpdateTestCase, GetTestCase, ListTestCase, DeleteTestCase, RunTestCase,
  GenerateCode, ListCodeGenerator,
  PopularHeaders,
  GetStores, DeleteStore, VerifyStore,
  getToken
}
