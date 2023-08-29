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
