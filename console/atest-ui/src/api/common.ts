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

export const DefaultResponseProcess = (response: any) => {

  if (!response.ok) {
    switch (response.status) {
      case 401:
        throw new Error('Unauthenticated')
    }
    throw new Error(response.statusText)
  } else {
    return response.json()
  }
}

export const safeToggleFunc = (toggle?: (e: boolean) => void) => {

  if (!toggle) {
    return (e: boolean) => {}
  }
  
  return toggle
}

export interface AppVersion {

  message: string
}

export interface TestSuite {

  store: string
  name: string
  api: string
  kind: string
}

export interface ImportSource {
 
  store: string
  url: string
}

export interface TestCase {
 
  suiteName: string
  name: string
  api: string
  method: string
}

export interface GenerateRequest {
 
  suiteName: string
  name: string
  generator: string
}

export interface RunTestCaseRequest {
 
  suiteName: string
  name: string
  parameters: any
}
