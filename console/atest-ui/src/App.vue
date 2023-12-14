<script setup lang="ts">
import {
  Document,
  Menu as IconMenu,
  Location,
  Share,
} from '@element-plus/icons-vue'
import { ref, watch } from 'vue'
import { API } from './views/net'
import { Cache } from './views/cache'
import TestingPanel from './views/TestingPanel.vue'
import StoreManager from './views/StoreManager.vue'
import SecretManager from './views/SecretManager.vue'
import WelcomePage from './views/WelcomePage.vue'

import setAsDarkTheme from './theme'

const asDarkMode = ref(Cache.GetPreference().darkTheme)
setAsDarkTheme(asDarkMode.value)
watch(asDarkMode, Cache.WatchDarkTheme)
watch(asDarkMode, () => {
  setAsDarkTheme(asDarkMode.value)
})

const appVersion = ref('')
const appVersionLink = ref('https://github.com/LinuxSuRen/api-testing')
API.GetVersion((d) => {
  appVersion.value = d.message
  const version = d.message.match('^v\\d*.\\d*.\\d*')
  const dirtyVersion = d.message.match('^v\\d*.\\d*.\\d*-\\d*-g')

  if (!version && !dirtyVersion) {
    return
  }

  if (dirtyVersion && dirtyVersion.length > 0) {
    appVersionLink.value = appVersionLink.value + '/commit/' + d.message.replace(dirtyVersion[0], '')
  } else if (version && version.length > 0) {
    appVersionLink.value = appVersionLink.value + '/releases/tag/' + version[0]
  }
})

const panelName = ref('')
const sideWidth = ref("width: 200px; display: flex;flex-direction: column;")
const isCollapse = ref(false)
watch(isCollapse, (e) => {
  if (e) {
    sideWidth.value = "width: 80px; display: flex;flex-direction: column;"
  } else {
    sideWidth.value = "width: 200px; display: flex;flex-direction: column;"
  }
})
const handleSelect = (key: string) => {
  panelName.value = key
}
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
        default-active="welcome"
        :collapse="isCollapse"
        @select="handleSelect"
      >
        <el-menu-item index="welcome">
          <el-icon><share /></el-icon>
          <template #title>Welcome</template>
        </el-menu-item>
        <el-menu-item index="testing" test-id="testing-menu">
          <el-icon><icon-menu /></el-icon>
          <template #title>Testing</template>
        </el-menu-item>
        <el-menu-item index="secret">
          <el-icon><document /></el-icon>
          <template #title>Secrets</template>
        </el-menu-item>
        <el-menu-item index="store">
          <el-icon><location /></el-icon>
          <template #title>Stores</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main style="padding-top: 5px; padding-bottom: 5px;">
      <TestingPanel v-if="panelName === 'testing'" />
      <StoreManager v-else-if="panelName === 'store'" />
      <SecretManager v-else-if="panelName === 'secret'" />
      <WelcomePage v-else />
    </el-main>

    <div style="position: absolute; bottom: 0px; right: 10px;">
      <a :href=appVersionLink target="_blank" rel="noopener">{{appVersion}}</a>
    </div>
  </el-container>
</template>
