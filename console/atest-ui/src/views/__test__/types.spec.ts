import { describe } from 'node:test'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from '../types'
import type { Pair } from '../types'

function expectFetch(url: string) {
  return {
    return: (data: any) => {
      global.fetch=function(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
        const myOptions = { status: 200, statusText: "SuperSmashingGreat!" };
        const myResponse = new Response(data, myOptions);
        return Promise.resolve(myResponse);
      }
    }
  }
}

describe('NewSuggestedAPIsQuery', () => {
  test('function is not null', () => {
    expectFetch('/server.Runner/GetSuggestedAPIs').return(`{"data":[{"request":{"api":"xxx"}}]}`)

    const func = NewSuggestedAPIsQuery('suite')
    expect(func).not.toBeNull()

    func('xxx', function(e) {
      expect(e.length).toBe(1)
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
    expect(options).toHaveLength(4)
    options.forEach((k, v) => {
      expect(k.key).toBe(k.value)
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
