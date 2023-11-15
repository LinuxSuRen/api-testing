/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

export interface TestCaseResponse {
    output: string
    body: {},
    statusCode: number
}

interface Store {
    name: string
    readOnly: boolean
}

interface Stores {
    current: string
    items: Store[]
}

export interface Preference {
    darkTheme: boolean
    requestActiveTab: string,
    responseActiveTab: string
}

export function GetTestCaseResponseCache(id: string) {
    const val = sessionStorage.getItem(id)
    if (val && val !== '') {
        return JSON.parse(val)
    } else {
        return {} as TestCaseResponse
    }
}

export function SetTestCaseResponseCache(id: string, resp: TestCaseResponse) {
    sessionStorage.setItem(id, JSON.stringify(resp))
}

const lastTestCaseLocationKey = "api-testing-case-location"
export function GetLastTestCaseLocation() {
    const val = localStorage.getItem(lastTestCaseLocationKey)
    if (val && val !== '') {
        return JSON.parse(val)
    } else {
        return {}
    }
}

export function SetLastTestCaseLocation(suite: string, testcase: string) {
    localStorage.setItem(lastTestCaseLocationKey, JSON.stringify({
        suite: suite,
        testcase: testcase
    }))
    return
}

const preferenceKey = "api-testing-preference"
export function GetPreference() {
    const val = localStorage.getItem(preferenceKey)
    if (val && val !== '') {
        return JSON.parse(val)
    } else {
        return {
            darkTheme: false,
            requestActiveTab: "body",
            responseActiveTab: "body"
        } as Preference
    }
}

export function SetPreference(preference: Preference) {
    localStorage.setItem(preferenceKey, JSON.stringify(preference))
    return
}

export function WatchRequestActiveTab(tab: string) {
    const preference = GetPreference()
    preference.requestActiveTab = tab
    SetPreference(preference)
}

function WatchResponseActiveTab(tab: string) {
    const preference = GetPreference()
    preference.responseActiveTab = tab
    SetPreference(preference)
}

function WatchDarkTheme(darkTheme: boolean) {
    const preference = GetPreference()
    preference.darkTheme = darkTheme
    SetPreference(preference)
}

const storeKey = "stores"
function GetCurrentStore() {
    const val = sessionStorage.getItem(storeKey)
    if (val && val !== '') {
        const stores = JSON.parse(val)
        for (var i = 0; i < stores.items.length; i++) {
            if (stores.items[i].name === stores.current) {
                return stores.items[i]
            }
        }
    }
    return {}
}
function SetCurrentStore(name: string) {
    const val = sessionStorage.getItem(storeKey)
    if (val && val !== '') {
        const stores = JSON.parse(val)
        stores.current = name
        SetStores(stores)
    }
}
function SetStores(stores: Stores | Store[]) {
    if ('current' in stores) {
        sessionStorage.setItem(storeKey, JSON.stringify(stores))
    } else {
        sessionStorage.setItem(storeKey, JSON.stringify({
            items: stores
        }))
    }
    return
}

export const Cache = {
    GetTestCaseResponseCache,
    SetTestCaseResponseCache,
    GetLastTestCaseLocation,
    SetLastTestCaseLocation,
    GetPreference,
    WatchRequestActiveTab,
    WatchResponseActiveTab,
    WatchDarkTheme,
    GetCurrentStore, SetStores, SetCurrentStore
}
