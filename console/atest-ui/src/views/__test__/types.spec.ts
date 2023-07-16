import { describe } from 'node:test'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods } from '../types'
import type { Pair } from '../types'

describe('NewSuggestedAPIsQuery', () => {
  test('function is not null', () => {
    const func = NewSuggestedAPIsQuery('suite')
    expect(func).not.toBeNull()
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
    expect(options).toHaveLength(4)
  })
})
