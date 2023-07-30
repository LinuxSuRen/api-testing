import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
// import router from './router'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import introJs from 'intro.js'
import 'intro.js/introjs.css'

const app = createApp(App)

app.use(ElementPlus, {
  locale: zhCn
})
// app.use(router)

app.mount('#app')

const dontShowAgain = window.location.search.indexOf('newbie') === -1;
introJs().setOptions({
  "dontShowAgain": dontShowAgain,
  "showProgress": true,
}).start();
