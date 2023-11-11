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
