import { describe } from 'node:test'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from '../types'
import type { Pair } from '../types'

const fakeFetch: { [key:string]:string; } = {};
function matchFake(url: string, data: string) {
  fakeFetch[url] = data
}

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

describe('NewSuggestedAPIsQuery', () => {
  test('empty data', () => {
    const func = NewSuggestedAPIsQuery('', '')
    expect(func).not.toBeNull()

    func('xxx', function(e) {
      expect(e.length).toBe(0)
    })
  })

  test('have data', () => {
    matchFake('/server.Runner/GetSuggestedAPIs', `{"data":[{"request":{"api":"xxx"}}]}`)

    const func = NewSuggestedAPIsQuery('', 'suite')
    expect(func).not.toBeNull()

    func('xxx', function(e) {
      expect(e.length).toBe(1)
    })

    func('not-match', function(e) {
      expect(e.length).toBe(0)
    })

    const anotherFunc = NewSuggestedAPIsQuery('', '')
    expect(anotherFunc).not.toBeNull()
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
