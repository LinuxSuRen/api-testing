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

import { describe } from 'node:test'
import {Cache, SetPreference, TestCaseResponse, Store, Stores} from '../cache'

const localStorageMock = (() => {
    let store = {};

    return {
        getItem(key) {
            return store[key] || null;
        },
        setItem(key, value) {
            store[key] = value.toString();
        },
        removeItem(key) {
            delete store[key];
        },
        clear() {
            store = {};
        }
    };
})();

Object.defineProperty(global, 'sessionStorage', {
    value: localStorageMock
});
Object.defineProperty(global, 'localStorage', {
    value: localStorageMock
});
Object.defineProperty(global, 'navigator', {
    value: {
        language: 'en'
    }
});

describe('TestCaseResponseCache', () => {
    test('should set and get test case response cache', () => {
        const id = 'test-case-id'
        const resp = {
            output: 'test-body',
            body: {},
            statusCode: 200,
        } as TestCaseResponse
        Cache.SetTestCaseResponseCache(id, resp)
        const result = Cache.GetTestCaseResponseCache(id)
        expect(result).toEqual(resp)
    })

    test('get a non-existent test case response cache', () => {
        expect(Cache.GetTestCaseResponseCache('non-existent-id')).toEqual({})
    })
})

describe('LastTestCaseLocation', () => {
    test('should get empty object when no last test case location', () => {
        expect(Cache.GetLastTestCaseLocation()).toEqual({})
    })

    test('should set and get last test case location', () => {
        const suite = 'test-suite'
        const testcase = 'test-case'
        Cache.SetLastTestCaseLocation(suite, testcase)
        const result = Cache.GetLastTestCaseLocation()
        expect(result).toEqual({ suite, testcase })
    })
})

describe('Preference', () => {
    test('get the default preference', () => {
        expect(Cache.GetPreference()).toEqual({
            darkTheme: false,
            requestActiveTab: 'body',
            responseActiveTab: 'body',
            language: 'en',
        })
    })

    test('set and get preference', () => {
        const preference = {
            darkTheme: true,
            requestActiveTab: 'header',
            responseActiveTab: 'header',
            language: 'zh-cn',
        }
        SetPreference(preference)
        expect(Cache.GetPreference()).toEqual(preference)
    })

    test('set and get dark theme', () => {
        Cache.WithDarkTheme(true)
        expect(Cache.GetPreference().darkTheme).toEqual(true)
    })

    test('set and get request active tab', () => {
        Cache.WithRequestActiveTab('request')
        expect(Cache.GetPreference().requestActiveTab).toEqual('request')
    })

    test('set and get response active tab', () => {
        Cache.WithResponseActiveTab('response')
        expect(Cache.GetPreference().responseActiveTab).toEqual('response')
    })

    it('set and get language', () => {
        Cache.WithLocale('zh-cn')
        expect(Cache.GetPreference().language).toEqual('zh-cn')
    })
})

describe('stores', () => {
    test('should get empty object when no stores', () => {
        expect(Cache.GetCurrentStore()).toEqual({})
    })

    test('should set and get stores', () => {
        const stores = {
            current: 'test-store',
            items: [
                {
                    name: 'test-store',
                    readOnly: false,
                } as Store,
                {
                    name: 'read-only-store',
                    readOnly: true,
                } as Store,
            ],
        }
        Cache.SetStores(stores)
        expect(Cache.GetCurrentStore()).toEqual({
            name: 'test-store',
            readOnly: false,
        })

        Cache.SetCurrentStore('read-only-store')
        expect(Cache.GetCurrentStore()).toEqual({
            name: 'read-only-store',
            readOnly: true,
        })

        Cache.SetStores({} as Stores)
    })
})
