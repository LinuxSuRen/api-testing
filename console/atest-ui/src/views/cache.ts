export interface TestCaseResponse {
    output: string
    body: {},
    statusCode: number
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
