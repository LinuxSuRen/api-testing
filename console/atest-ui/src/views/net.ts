export function DefaultResponseProcess() {
    return (response: any) => {
        if (!response.ok) {
          throw new Error(response.statusText)
        } else {
          return response.json()
        }
      }
}