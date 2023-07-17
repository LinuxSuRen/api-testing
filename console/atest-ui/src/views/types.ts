import { ref } from 'vue'

export interface Suite {
  name: string
  api: string
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
export function NewSuggestedAPIsQuery(suite: string) {
  return function (queryString: string, cb: (arg: any) => void) {
    loadCache(suite, function () {
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
function loadCache(suite: string, callback: Function) {
  if (localCache.value.length > 0) {
    callback()
    return
  }

  if (suite === '') {
    return
  }

  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
      name: suite
    })
  }
  fetch('/server.Runner/GetSuggestedAPIs', requestOptions)
    .then((response) => response.json())
    .then((e) => {
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
    }
  ] as Pair[]
}
