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
import { Cache } from './cache'

async function DefaultResponseProcess(response: any) {
  if (!response.ok) {
    switch (response.status) {
      case 401:
        throw new Error("Unauthenticated")
    }

    const message = await response.json().then((data :any) => data.message)
    throw new Error(message)
  } else {
    return response.json()
  }
}

interface AppVersion {
  version: string
  commit: string
  date: string
}

function safeToggleFunc(toggle?: (e: boolean) => void) {
  if (!toggle) {
    return (e: boolean) => {}
  }
  return toggle
}

function GetVersion(callback: (v: AppVersion) => void) {
  const requestOptions = {
    method: 'GET',
  }
  fetch('/api/v1/version', requestOptions)
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

  fetch('/api/v1/suites', requestOptions)
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
    method: 'PUT',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify(suite)
  }
  fetch(`/api/v1/suites/${suite.name}`, requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function GetTestSuite(name: string,
  callback: (d: any) => void, errHandle: (e: any) => void) {
  const store = Cache.GetCurrentStore()
  const requestOptions = {
    headers: {
      'X-Store-Name': store.name,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suites/${name}`, requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteTestSuite(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'DELETE',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suites/${name}`, requestOptions)
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
  fetch(`/api/v1/converters/convert`, requestOptions)
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
    fetch(`/api/v1/suites/${sourceSuiteName}/duplicate`, requestOptions)
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

  fetch(`/api/v1/suites/import`, requestOptions)
  .then(DefaultResponseProcess)
    .then(callback)
}

interface TestCase {
  suiteName: string
  name: string
  request: any
}

interface HistoryTestCase {
  historyCaseID : string,
  suiteName: string
  caseName: string
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
        request: testcase.request
      }
    })
  }

  fetch(`/api/v1/suites/${testcase.suiteName}/cases`, requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function UpdateTestCase(testcase: any,
  callback: (d: any) => void, errHandle?: (e: any) => void | null,
  toggle?: (e: boolean) => void) {
    const requestOptions = {
      method: 'PUT',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify(testcase)
    }
    safeToggleFunc(toggle)(true)
    fetch(`/api/v1/suites/${testcase.suiteName}/cases/${testcase.data.name}`, requestOptions)
      .then(DefaultResponseProcess)
      .then(callback).catch(errHandle)
      .finally(() => {
        safeToggleFunc(toggle)(false)
      })
}

function GetTestCase(req: TestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suites/${req.suiteName}/cases/${req.name}`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function ListTestCase(suite: string, store: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    headers: {
      'X-Store-Name': store,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suites/${suite}/cases`, requestOptions)
  .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteTestCase(testcase: TestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    const requestOptions = {
      method: 'DELETE',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify({
        suite: testcase.suiteName,
        testcase: testcase.name
      })
    }
    fetch(`/api/v1/suites/${testcase.suiteName}/cases/${testcase.name}`, requestOptions)
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
  fetch(`/api/v1/suites/${request.suiteName}/cases/${request.name}/run`, requestOptions)
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
    fetch(`/api/v1/suites/${sourceSuiteName}/cases/${sourceTestCaseName}/duplicate`, requestOptions)
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
  fetch(`/api/v1/codeGenerators/generate`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function ListCodeGenerator(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  fetch('/api/v1/codeGenerators', {
    headers: {
      'X-Auth': getToken()
    },
  }).then(DefaultResponseProcess)
  .then(callback).catch(errHandle)
}

function PopularHeaders(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
  }
  fetch(`/api/v1/popularHeaders`, requestOptions)
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
  
  let api = '/api/v1/stores'
  if (!create) {
    api = `/api/v1/stores/${payload.name}`
    requestOptions.method = "PUT"
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
    headers: {
      'X-Auth': getToken()
    },
  }
  fetch('/api/v1/stores', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle).finally(final)
}

function DeleteStore(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'DELETE',
    headers: {
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/stores/${name}`, requestOptions)
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
  
  const api = `/api/v1/stores/verify`
  fetch(api, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

export interface Secret {
  Name: string
  Value: string
}

function GetSecrets(callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    headers: {
      'X-Auth': getToken()
    },
  }
  fetch(`/api/v1/secrets`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

function FunctionsQuery(filter: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    headers: {
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/functions?name=${filter}`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteSecret(name: string,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: "DELETE",
    headers: {
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/secrets/${name}`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

function CreateOrUpdateSecret(payload: Secret, create: boolean,
  callback: (d: any) => void, errHandle?: (e: any) => void | null,
  toggle?: (e: boolean) => void) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': getToken()
    },
    body: JSON.stringify(payload)
  }
  
  let api = `/api/v1/secrets`
  if (!create) {
    api = `/api/v1/secrets/${payload.Name}`
    requestOptions.method = "PUT"
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
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suggestedAPIs?name=${name}`, requestOptions)
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
  fetch(`/api/v1/mock/reload`, requestOptions)
      .then(DefaultResponseProcess)
}

function GetMockConfig(callback: (d: any) => void) {
  const requestOptions = {
    headers: {
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/mock/config`, requestOptions)
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

const GetTestSuiteYaml = (suite: string, callback: (d: any) => void, errHandle?: (e: any) => void | null) => {
  const requestOptions = {
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    }
  }
  fetch(`/api/v1/suites/${suite}/yaml`, requestOptions)
    .then(DefaultResponseProcess)
    .then(callback)
    .catch(errHandle)
}

function GetHistoryTestCaseWithResult(req: HistoryTestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      ID : req.historyCaseID
    })
  }
  fetch('/server.Runner/GetHistoryTestCaseWithResult', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function GetHistoryTestCase(req: HistoryTestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': Cache.GetCurrentStore().name,
      'X-Auth': getToken()
    },
    body: JSON.stringify({
      ID : req.historyCaseID
    })
  }
  fetch('/server.Runner/GetHistoryTestCase', requestOptions)
    .then(DefaultResponseProcess)
    .then(callback).catch(errHandle)
}

function DeleteHistoryTestCase(req: HistoryTestCase,
  callback: (d: any) => void, errHandle?: (e: any) => void | null) {
    const requestOptions = {
      method: 'POST',
      headers: {
        'X-Store-Name': Cache.GetCurrentStore().name,
        'X-Auth': getToken()
      },
      body: JSON.stringify({
        ID : req.historyCaseID
      })
    }
    fetch('/server.Runner/DeleteHistoryTestCase', requestOptions)
      .then(callback).catch(errHandle)
}

export const API = {
  DefaultResponseProcess,
  GetVersion,
  CreateTestSuite, UpdateTestSuite, ImportTestSuite, GetTestSuite, DeleteTestSuite, ConvertTestSuite,GetTestSuiteYaml,
  CreateTestCase, UpdateTestCase, GetTestCase, ListTestCase, DeleteTestCase, RunTestCase,
  GetHistoryTestCaseWithResult, DeleteHistoryTestCase, GetHistoryTestCase, GetTestCaseAllHistory,
  GenerateCode, ListCodeGenerator,
  PopularHeaders,
  CreateOrUpdateStore, GetStores, DeleteStore, VerifyStore,
  FunctionsQuery,
  GetSecrets, DeleteSecret, CreateOrUpdateSecret,
  GetSuggestedAPIs,
  ReloadMockServer, GetMockConfig,
  getToken
}
