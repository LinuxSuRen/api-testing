import { API } from './views/net'

const themes: { [k: string]: any } = {}

export function getThemes() {
    return Object.keys(themes)
}

API.GetThemes().then(data => {
    data.data.forEach((theme) => {
        const key = theme.key
        API.GetTheme(key).then((data: any) => {
            themes[key] = JSON.parse(data.message)

            const theme = getTheme()
            if (theme) {
                setTheme(theme)
            }
        })
    })
})

export function setTheme(theme: string) {
    const themeObj = themes[theme]
    if (themeObj) {
        applyTheme(themeObj)
        window.localStorage.setItem('theme', theme)
    }
}

export function getTheme() {
    return window.localStorage.getItem('theme')
}

const applyTheme = function (data: any) {
    if (data instanceof Object) {
        Object.keys(data).forEach((key) => {
            if (data[key] instanceof Object) {
                applyTheme(data[key])
            } else {
                document.documentElement.style.setProperty(key, data[key])
            }
        })
    }
}

export function setAsDarkTheme(darkMode: boolean) {
  document.documentElement.className = darkMode ? 'dark' : 'light'
}
