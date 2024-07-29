<script setup lang="ts">
import {
  Document,
  Menu as IconMenu,
  Location,
  Share,
  ArrowDown
} from '@element-plus/icons-vue'
import { ref, watch } from 'vue'
import { API } from './views/net'
import { Cache } from './views/cache'
import TestingPanel from './views/TestingPanel.vue'
import MockManager from './views/MockManager.vue'
import StoreManager from './views/StoreManager.vue'
import SecretManager from './views/SecretManager.vue'
import WelcomePage from './views/WelcomePage.vue'
import { useI18n } from 'vue-i18n'

const { t, locale: i18nLocale } = useI18n()

import setAsDarkTheme from './theme'

const asDarkMode = ref(Cache.GetPreference().darkTheme)
setAsDarkTheme(asDarkMode.value)
watch(asDarkMode, Cache.WithDarkTheme)
watch(asDarkMode, () => {
  setAsDarkTheme(asDarkMode.value)
})

const appVersion = ref('')
const appVersionLink = ref('https://github.com/LinuxSuRen/api-testing')
API.GetVersion((d) => {
  appVersion.value = d.version
  const version = d.version.match('^v\\d*.\\d*.\\d*')
  const dirtyVersion = d.version.match('^v\\d*.\\d*.\\d*-\\d*-g')

  if (!version && !dirtyVersion) {
    return
  }

  if (dirtyVersion && dirtyVersion.length > 0) {
    appVersionLink.value = appVersionLink.value + '/commit/' + d.message.replace(dirtyVersion[0], '')
  } else if (version && version.length > 0) {
    appVersionLink.value = appVersionLink.value + '/releases/tag/' + version[0]
  }
})

const sideWidth = ref("width: 200px; display: flex;flex-direction: column;")
const isCollapse = ref(true)
watch(isCollapse, (e) => {
  if (e) {
    sideWidth.value = "width: 80px; display: flex;flex-direction: column;"
  } else {
    sideWidth.value = "width: 200px; display: flex;flex-direction: column;"
  }
})
const lastActiveMenu = window.localStorage.getItem('activeMenu')
const activeMenu = ref(lastActiveMenu === '' ? 'welcome' : lastActiveMenu)
const panelName = ref(activeMenu)
const handleSelect = (key: string) => {
  panelName.value = key
  window.localStorage.setItem('activeMenu', key)
}

const locale = ref(Cache.GetPreference().language)
i18nLocale.value = locale.value

watch(locale, (e) =>{
  Cache.WithLocale(e.value)
  i18nLocale.value = locale.value
})

const handleChangeLan = (command: string) => {
  switch (command) {
    case "chinese":
      locale.value = "zh-CN"
      break;
    case "english":
      locale.value = "en-US"
      break;
  }
};

</script>

<template>
  <el-container style="height: 100%">
    <el-aside :style="sideWidth">
      <el-radio-group v-model="isCollapse">
        <el-radio-button :label="false">+</el-radio-button>
        <el-radio-button :label="true">-</el-radio-button>
      </el-radio-group>
      <el-menu
        style="flex-grow: 1;"
        :default-active="activeMenu"
        :collapse="isCollapse"
        @select="handleSelect"
      >
        <el-menu-item index="welcome">
          <el-icon><share /></el-icon>
          <template #title>{{ t('title.welcome') }}</template>
        </el-menu-item>
        <el-menu-item index="testing" test-id="testing-menu">
          <el-icon><icon-menu /></el-icon>
          <template #title>{{ t('title.testing' )}}</template>
        </el-menu-item>
        <el-menu-item index="mock" test-id="mock-menu">
          <el-icon><icon-menu /></el-icon>
          <template #title>{{ t('title.mock' )}}</template>
        </el-menu-item>
        <el-menu-item index="secret">
          <el-icon><document /></el-icon>
          <template #title>{{ t('title.secrets') }}</template>
        </el-menu-item>
        <el-menu-item index="store">
          <el-icon><location /></el-icon>
          <template #title>{{ t('title.stores') }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main style="padding-top: 5px; padding-bottom: 5px;">
      <div style="position: absolute; top: 10px; right: 20px;">
        <el-col style="display: flex; align-items: center;">
          <el-tag style="font-size: 18px;">{{ t('language') }}</el-tag>
          <el-dropdown trigger="click" @command="(command: string) => handleChangeLan(command)">
            <el-icon><arrow-down /></el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="chinese">中文</el-dropdown-item>
                <el-dropdown-item command="english">English</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-col>
      </div>
      <TestingPanel v-if="panelName === 'testing'" />
      <MockManager v-else-if="panelName === 'mock'" />
      <StoreManager v-else-if="panelName === 'store'" />
      <SecretManager v-else-if="panelName === 'secret'" />
      <WelcomePage v-else />
    </el-main>

    <div style="position: absolute; bottom: 0px; right: 10px;">
      <a :href=appVersionLink target="_blank" rel="noopener">{{appVersion}}</a>
    </div>
  </el-container>
</template>
