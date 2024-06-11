export const getToken = () => {

    const token = sessionStorage.getItem('token')
    if (!token) {
      return ''
    }
    return token
  }
  