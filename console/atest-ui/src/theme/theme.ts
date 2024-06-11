const lightTheme: { [k: string]: string } = {
  '--color-background': 'var(--vt-c-white)',
  '--color-background-soft': 'var(--vt-c-white-soft)',
  '--color-background-mute': 'var(--vt-c-white-mute)',
  '--color-border': 'var(--vt-c-divider-light-2)',
  '--color-border-hover': 'var(--vt-c-divider-light-1)',
  '--color-heading': 'var(--vt-c-text-light-1)',
  '--color-text': 'var(--vt-c-text-light-1)'
}

const darkTheme: { [k: string]: string } = {
  '--color-background': 'var(--vt-c-black)',
  '--color-background-soft': 'var(--vt-c-black-soft)',
  '--color-background-mute': 'var(--vt-c-black-mute)',
  '--color-border': 'var(--vt-c-divider-dark-2)',
  '--color-border-hover': 'var(--vt-c-divider-dark-1)',
  '--color-heading': 'var(--vt-c-text-dark-1)',
  '--color-text': 'var(--vt-c-text-dark-2)'
}


export const setAsDarkTheme = (darkMode: boolean) => {

  const theme = darkMode ? darkTheme : lightTheme
  
  Object.keys(theme).forEach((key) => {    
    document.documentElement.style.setProperty(key, theme[key])
  })
}
