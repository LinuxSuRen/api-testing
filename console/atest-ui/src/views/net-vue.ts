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
import { API } from './net'
import type { Ref } from 'vue'
import { ElMessage } from 'element-plus'

function UpdateTestCase(testcase: any,
    callback: (d: any) => void, errHandle?: (e: any) => void | null,
    loadingRef?: Ref<Boolean>) {
    API.UpdateTestCase(testcase, callback, errHandle, (e: boolean) => {
        if (loadingRef) {
            loadingRef.value = e
        }
    })
}

function CreateOrUpdateSecret(payload: any, create: boolean,
    callback: (d: any) => void,
    loadingRef?: Ref<Boolean>) {
    API.CreateOrUpdateSecret(payload, create, callback, ErrorTip, (e: boolean) => {
        if (loadingRef) {
            loadingRef.value = e
        }
    })
}

function CreateOrUpdateStore(payload: any, create: boolean,
    callback: (d: any) => void,
    loadingRef?: Ref<Boolean>) {
    API.CreateOrUpdateStore(payload, create, callback, ErrorTip, (e: boolean) => {
        if (loadingRef) {
            loadingRef.value = e
        }
    })
}

function ErrorTip(e: {
    statusText:''
}) {
    ElMessage.error('Oops, ' + e.statusText)
}

export const UIAPI = {
    UpdateTestCase, CreateOrUpdateSecret,
    CreateOrUpdateStore,
    ErrorTip
}
