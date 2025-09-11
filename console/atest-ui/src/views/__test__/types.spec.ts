/*
Copyright 2023-2024 API Testing Authors.

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

import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from '../types'
import type { Pair } from '../types'
import { fetchMock } from '../../test-setup'
import { vi } from 'vitest'

describe('NewSuggestedAPIsQuery', () => {
  test('empty data', () => {
    const func = NewSuggestedAPIsQuery('', '')
    expect(func).not.toBeNull()

    func('xxx', function(e) {
      expect(e.length).toBe(0)
    })
  })
})

describe('CreateFilter', () => {
  const filter = CreateFilter('suite')

  test('ignore letter case', () => {
    expect(filter({ value: 'Suite' } as Pair)).toBeTruthy()
  })

  test('not match', () => {
    expect(filter({ value: 'not match' } as Pair)).not.toBeTruthy()
  })
})

describe('GetHTTPMethods', () => {
  test('HTTP methods', () => {
    const options = GetHTTPMethods()
    expect(options).toHaveLength(7)
    options.forEach((item) => {
      expect(item.key).toBe(item.value)
    })
  })
})

describe('FlattenObject', () => {
  test('simple nested object', () => {
    const obj = {
      a: {
        b: {
          c: 'd'
        }
      }
    }
    const result = FlattenObject(obj);
    expect(result).toEqual(result)
  })
})
