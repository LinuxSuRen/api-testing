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

import service from './manager'

export const get = (url: string, params: any, header: any) => {
    
    return new Promise((resolve, reject) => {
        service.get(url, {
            params: params,
            headers: header,
        }).then(res => {
            resolve(res.data);
        }).catch(err => {
            reject(err.data);
        });
    });
}

export const del = (url: string, params: any, header: any) => {

    return new Promise((resolve, reject) => {
        service.delete(
            url, 
            {
                headers: header,
            }
        ).then(res => {
            resolve(res.data);
        }).catch(err => {
            reject(err.data);
        });
    });
}

export const put = (url: string, data: any, header: any) => {

    return new Promise((resolve, reject) => {
        service.put(
            url, 
            data, 
            {
                headers: header,
            }
        ).then(res => {
            resolve(res.data);
        }).catch(err => {
            reject(err.data);
        });
    });
}

export const post = (url: string, params: any, header: any) => {

    return new Promise((resolve, reject) => {
        service.post(
            url, 
            params, 
            {
                headers: header,
            }
        ).then(res => {
            resolve(res.data);
        }).catch(err => {
            reject(err.data);
        });
    });
}
