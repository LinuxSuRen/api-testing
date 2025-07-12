import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import enUS from 'element-plus/dist/locale/en.mjs'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import introJs from 'intro.js'
import 'intro.js/introjs.css'
import { setupI18n } from './i18n'
import en from './locales/en.json'
import zh from './locales/zh.json'
import { name, version } from '../package'
import 'element-plus/theme-chalk/dark/css-vars.css'

const app = createApp(App)

const language = window.navigator.userLanguage || window.navigator.language;
const lang = language.split('-')[0]

const i18n = setupI18n({
  legacy: false,
  locale: lang,
  fallbackLocale: 'en',
  messages: {
    en, zh
  }
})

const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('access_token');
if (token && token !== '') {
  sessionStorage.setItem('token', token)
  window.location.href='/'
}

app.use(ElementPlus, {
  locale: lang === 'zh' ? zhCn : enUS
})

app.use(i18n)
app.mount('#app')

const dontShowAgain = window.location.search.indexOf('newbie') === -1;
introJs().setOptions({
  "dontShowAgain": dontShowAgain,
  "showProgress": true,
}).start();
