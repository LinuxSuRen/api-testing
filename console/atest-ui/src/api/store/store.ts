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

export const CreateStore = (params: any) =>
  post('/server.Runner/CreateStore', params, {
    'X-Auth': getToken()
  })

export const UpdateStore = (params: any) =>
  post('/server.Runner/UpdateStore', params, {
    'X-Auth': getToken()
  })

export const GetStores = () =>
  post('/server.Runner/GetStores', null, {
    'X-Auth': getToken()
  })

export const DeleteStore = (params: any) =>
  post('/server.Runner/DeleteStore', params, {
    'X-Auth': getToken()
  })

export const VerifyStore = (params: any) =>
  post('/server.Runner/VerifyStore', params, {
    'X-Auth': getToken()
  })
