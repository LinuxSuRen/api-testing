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
import {API} from '../net'
import { type TestCase } from '../net'
import { SetupStorage } from './common'
import fetchMock from "jest-fetch-mock";

fetchMock.enableMocks();
SetupStorage()

beforeEach(() => {
    fetchMock.resetMocks();
});

describe('net', () => {
    test('GetVersion', () => {
        fetchMock.mockResponseOnce(`{"version":"v0.0.1"}`)
        API.GetVersion((d) => {
            expect(d.version).toEqual('v0.0.2')
        })
    })

    test('CreateTestSuite', () => {
        fetchMock.mockResponseOnce(`{"version":"v0.0.1"}`)
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
        API.UpdateTestSuite({}, () => {})
    })

    test('GetTestSuite', () => {
        API.GetTestSuite('fake', () => {})
    })

    test('DeleteTestSuite', () => {
        API.DeleteTestSuite('fake', () => {})
    })

    test('ConvertTestSuite', () => {
        API.ConvertTestSuite('fake', 'generator', () => {})
    })

    test('DuplicateTestSuite', () => {
        API.DuplicateTestSuite('source', 'target', () => {})
    })

    test('GetTestSuiteYaml', () => {
        API.GetTestSuiteYaml('fake', () => {})
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
