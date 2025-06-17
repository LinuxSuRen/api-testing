export default function setAsDarkTheme(darkMode: boolean) {
  document.documentElement.className = darkMode ? 'dark' : 'light'
}
