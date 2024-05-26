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

function DuplicateTestSuite(sourceSuiteName: string, targetSuiteName: string,
    callback: (d: any) => void, errHandle?: ((reason: any) => PromiseLike<never>) | undefined | null ) {
    const requestOptions = {
      method: 'POST',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify({
          sourceSuiteName: sourceSuiteName,
          targetSuiteName: targetSuiteName,
      })
    }
    fetch('/server.Runner/DuplicateTestSuite', requestOptions)
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

function DuplicateTestCase(sourceSuiteName: string, targetSuiteName: string,
                            sourceTestCaseName: string, targetTestCaseName: string,
                            callback: (d: any) => void, errHandle?: ((reason: any) => PromiseLike<never>) | undefined | null ) {
    const requestOptions = {
        method: 'POST',
        headers: {
            'X-Store-Name': Cache.GetCurrentStore().name,
            'X-Auth': getToken()
        },
        body: JSON.stringify({
            sourceSuiteName: sourceSuiteName,
            targetSuiteName: targetSuiteName,
            sourceCaseName: sourceTestCaseName,
            targetCaseName: targetTestCaseName,
        })
    }
    fetch('/server.Runner/DuplicateTestCase', requestOptions)
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

function CreateOrUpdateStore(payload: any, create: boolean,
  callback: (d: any) => void, errHandle?: (e: any) => void | null,
  toggle?: (e: boolean) => void) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify(payload)
  }
  
  let api = '/server.Runner/CreateStore'
  if (!create) {
    api = '/server.Runner/UpdateStore'
  }

  safeToggleFunc(toggle)(true)
  fetch(api, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle).finally(() => {
      safeToggleFunc(toggle)(false)
    })
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

function GetSecrets(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
  }
  fetch('/server.Runner/GetSecrets', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

function FunctionsQuery(filter: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      name: filter
    })
  }
  fetch('/server.Runner/FunctionsQuery', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteSecret(name: string,
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
  fetch('/server.Runner/DeleteSecret', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

function CreateOrUpdateSecret(payload: any, create: boolean,
  callback: (d: any) => void, errHandle?: (e: any) => void | null,
  toggle?: (e: boolean) => void) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify(payload)
  }
  
  let api = '/server.Runner/CreateSecret'
  if (!create) {
    api = '/server.Runner/UpdateSecret'
  }

  safeToggleFunc(toggle)(true)
  fetch(api, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle).finally(() => {
      safeToggleFunc(toggle)(false)
    })
}

function GetSuggestedAPIs(name: string,
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
  fetch('/server.Runner/GetSuggestedAPIs', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
}

function ReloadMockServer(config: string) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      Config: config
    })
  }
  fetch('/server.Mock/Reload', requestOptions)
      .then(DefaultResponseProcess)
}

function GetMockConfig(callback: (d: any) => void) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    }
  }
  fetch('/server.Mock/GetConfig', requestOptions)
      .then(DefaultResponseProcess)
      .then(callback)
}

function getToken() {
  const token = sessionStorage.getItem('token')
  if (!token) {
    return ''
  }
  return token
}

const GetTestSuiteYaml = (suite: string, store: string, callback: (d: any) => void, errHandle?: (e: any) => void | null) => {
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
  fetch('/server.Runner/GetTestSuiteYaml', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

export const API = {
    DefaultResponseProcess,
    GetVersion,
    CreateTestSuite, UpdateTestSuite, ImportTestSuite, GetTestSuite, DeleteTestSuite, ConvertTestSuite, GetTestSuiteYaml,
    DuplicateTestSuite,
    CreateTestCase, UpdateTestCase, GetTestCase, ListTestCase, DeleteTestCase, RunTestCase, DuplicateTestCase,
    GenerateCode, ListCodeGenerator,
    PopularHeaders,
    CreateOrUpdateStore, GetStores, DeleteStore, VerifyStore,
    FunctionsQuery,
    GetSecrets, DeleteSecret, CreateOrUpdateSecret,
    GetSuggestedAPIs,
    ReloadMockServer, GetMockConfig,
    getToken
}
