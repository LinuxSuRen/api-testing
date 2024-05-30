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
import { getToken } from '../../utils/auth/token'

export const GetSecrets = () =>
  post('/server.Runner/GetSecrets', null, {
    'X-Auth': getToken()
  })

export const FunctionsQuery = (params: string) =>
  post('/server.Runner/FunctionsQuery', params, {
    'X-Auth': getToken()
  })

export const DeleteSecret = (name: any) =>
  post('/server.Runner/DeleteSecret', name, {
    'X-Auth': getToken()
  })

export const CreateSecret = (params: any) =>
  post('/server.Runner/CreateSecret', params, {
    'X-Auth': getToken()
  })

export const UpdateSecret = (params: any) =>
  post('/server.Runner/UpdateSecret', params, {
    'X-Auth': getToken()
  })
