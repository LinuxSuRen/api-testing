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
import {API, TestCase} from '../net'

const fakeFetch: { [key:string]:string; } = {};

global.fetch = jest.fn((key: string) =>
    Promise.resolve({
        json: () => {
            if (fakeFetch[key] === undefined) {
                return Promise.resolve({})
            }
            return Promise.resolve(JSON.parse(fakeFetch[key]))
        },
    }),
) as jest.Mock;

fakeFetch['/api/v1/version'] = '{"version":"0.0.1"}'

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

describe('net', () => {
    test('GetVersion', () => {
        API.GetVersion()
    })

    test('CreateTestSuite', () => {
        API.CreateTestSuite({
            store: 'store',
            name: 'name',
            api: 'api',
            kind: 'kind',
        }, (d) => {
            expect(d).toEqual({})
        })
    })

    test('UpdateTestSuite', () => {
        API.UpdateTestSuite({}, (d) => {})
    })

    test('GetTestSuite', () => {
        API.GetTestSuite('fake', (d) => {})
    })

    test('DeleteTestSuite', () => {
        API.DeleteTestSuite('fake', (d) => {})
    })

    test('ConvertTestSuite', () => {
        API.ConvertTestSuite('fake', 'generator', (d) => {})
    })

    test('DuplicateTestSuite', () => {
        API.DuplicateTestSuite('source', 'target', (d) => {})
    })

    test('GetTestSuiteYaml', () => {
        API.GetTestSuiteYaml('fake', (d) => {})
    })

    test('CreateTestCase', () => {
        API.CreateTestCase({
            suiteName: 'store',
            name: 'name'
        } as TestCase, (d) => {
            expect(d).toEqual({})
        })
    })
})
