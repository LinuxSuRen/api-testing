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

const path = require('node:path')

exports.control = function(okCallback, errorCallback) {
    fetch(getHealthzUrl()).
        then(okCallback).catch(errorCallback)
}

function getPort() {
    // TODO support set this value
    return 7788
}

function getHomePage() {
    return 'http://localhost:' + getPort()
}

function getHealthzUrl() {
    return 'http://localhost:' + getPort() + '/healthz'
}

function getHomeDir() {
    const homedir = require('os').homedir();
    return path.join(homedir, ".config", 'atest')
}

function getLogfile() {
    return path.join(getHomeDir(), 'log.log')
}

exports.getPort = getPort
exports.getHomePage = getHomePage
exports.getHomeDir = getHomeDir
exports.getLogfile = getLogfile
exports.getHealthzUrl = getHealthzUrl
