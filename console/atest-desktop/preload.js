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

const { contextBridge, ipcRenderer } = require('electron')

// All the Node.js APIs are available in the preload process.
// It has the same sandbox as a Chrome extension.
window.addEventListener('DOMContentLoaded', () => {
  const replaceText = (selector, text) => {
    const element = document.getElementById(selector)
    if (element) element.innerText = text
  }

  for (const dependency of ['chrome', 'node', 'electron']) {
    replaceText(`${dependency}-version`, process.versions[dependency])
  }
})

contextBridge.exposeInMainWorld('electronAPI', {
  openLogDir: () => ipcRenderer.send('openLogDir'),
  openWithExternalBrowser: (address) => ipcRenderer.invoke('openWithExternalBrowser', address),
  startServer: () => ipcRenderer.send('startServer'),
  stopServer: () => ipcRenderer.send('stopServer'),
  control: (okCallback, errCallback) => ipcRenderer.send('control', okCallback, errCallback),
  getHomePage: () => ipcRenderer.invoke('getHomePage'),
  getPort: () => ipcRenderer.invoke('getPort'),
  setPort: (port) => ipcRenderer.invoke('setPort', port),
  setExtensionRegistry: (registry) => ipcRenderer.invoke('setExtensionRegistry', registry),
  getExtensionRegistry: () => ipcRenderer.invoke('getExtensionRegistry'),
  getDownloadTimeout: () => ipcRenderer.invoke('getDownloadTimeout'),
  setDownloadTimeout: (timeout) => ipcRenderer.invoke('setDownloadTimeout', timeout),
  getHealthzUrl: () => ipcRenderer.invoke('getHealthzUrl'),
})
