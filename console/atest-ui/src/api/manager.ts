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

import axios from "axios";

const service = axios.create({
    baseURL: '/api',
    timeout: 5000
})

// 添加响应拦截器
service.interceptors.response.use((response) => {
    //这里的response就是请求接口以后回调的正确数据
    // console.log(response.data)
    return response;
}, function (error) {
    //这里的error就是请求接口以后回调的错误数据
    // elError(error.response.data || '请求失败')
    return Promise.reject(error.response);
});

// 添加请求拦截器
service.interceptors.request.use((config) => {
    //在请求接口前往header头自动添加token
    const token = sessionStorage.getItem("user-token")
    if (token) {
        config.headers['Authorization'] = "Bearer " + token
    }
    // console.log('看到token',token)
    return config;
}, function (error) {
    //对请求错误做些什么 
    return Promise.reject(error);
});

export default service
