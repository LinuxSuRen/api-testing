import { ref } from 'vue'

export function QueryFuncs(filter: string, valRef: typeof ref) {
    const requestOptions = {
      method: 'POST',
      body: JSON.stringify({
        name: filter
      })
    }
    fetch('/server.Runner/FunctionsQuery', requestOptions)
      .then((response) => response.json())
      .then((e) => {
        valRef.value = e.data
      })
}