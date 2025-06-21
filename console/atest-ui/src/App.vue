<script setup lang="ts">
import {
  Document,
  Menu as IconMenu,
  Histogram,
  Location,
  Share,
  ArrowDown,
  Guide,
  DataAnalysis,
  Help,
  Setting
} from '@element-plus/icons-vue'
import { ref, watch } from 'vue'
import { API } from './views/net'
import { Cache } from './views/cache'
import TestingPanel from './views/TestingPanel.vue'
import TestingHistoryPanel from './views/TestingHistoryPanel.vue'
import MockManager from './views/MockManager.vue'
import StoreManager from './views/StoreManager.vue'
import SecretManager from './views/SecretManager.vue'
import WelcomePage from './views/WelcomePage.vue'
import DataManager from './views/DataManager.vue'
import { useI18n } from 'vue-i18n'
import { setAsDarkTheme, getThemes, setTheme, getTheme } from './theme'

const { t, locale: i18nLocale } = useI18n()

const allThemes = ref(getThemes())
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

  if (!version && !dirtyVersion) return

  if (dirtyVersion && dirtyVersion.length > 0) {
    appVersionLink.value += '/commit/' + d.message.replace(dirtyVersion[0], '')
  } else if (version && version.length > 0) {
    appVersionLink.value += '/releases/tag/' + version[0]
  }
})

const isCollapse = ref(true)
watch(isCollapse, (v: boolean) => {
  window.localStorage.setItem('button.style', v ? 'simple' : '')
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
watch(locale, (e: string) => {
  Cache.WithLocale(e)
  i18nLocale.value = e
})

const handleChangeLan = (command: string) => {
  switch (command) {
    case 'chinese':
      locale.value = 'zh-CN'
      break
    case 'english':
      locale.value = 'en-US'
      break
  }
}

const ID = ref(null)
const toHistoryPanel = ({ ID: selectID, panelName: historyPanelName }) => {
  ID.value = selectID
  panelName.value = historyPanelName
}

const settingDialogVisible = ref(false)
watch(settingDialogVisible, (v: boolean) => {
  if (v) {
    allThemes.value = getThemes()
  }
})

const theme = ref(getTheme())
watch(theme, (e) => {
  setTheme(e)
})
</script>

<template>
  <el-container class="h-screen">
    <!-- Sidebar -->
    <el-aside class="flex flex-col border-r border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900" :width="isCollapse ? '80px' : '200px'">
      <!-- Collapse radio buttons, centered -->
      <el-radio-group v-model="isCollapse" class="p-3 flex justify-center space-x-3">
        <el-radio-button :value="false" class="px-3">+</el-radio-button>
        <el-radio-button :value="true" class="px-3">-</el-radio-button>
      </el-radio-group>

      <!-- Menu -->
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        @select="handleSelect"
        class="el-menu-vertical flex-1 overflow-auto border-r-0"
      >
        <el-menu-item index="welcome">
          <el-icon><Share /></el-icon>
          <template #title>{{ t('title.welcome') }}</template>
        </el-menu-item>
        <el-menu-item index="testing" test-id="testing-menu">
          <el-icon><IconMenu /></el-icon>
          <template #title>{{ t('title.testing') }}</template>
        </el-menu-item>
        <el-menu-item index="history" test-id="history-menu">
          <el-icon><Histogram /></el-icon>
          <template #title>{{ t('title.history') }}</template>
        </el-menu-item>
        <el-menu-item index="mock" test-id="mock-menu">
          <el-icon><Guide /></el-icon>
          <template #title>{{ t('title.mock') }}</template>
        </el-menu-item>
        <el-menu-item index="data" test-id="data-menu">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>{{ t('title.data') }}</template>
        </el-menu-item>
        <el-menu-item index="secret">
          <el-icon><Document /></el-icon>
          <template #title>{{ t('title.secrets') }}</template>
        </el-menu-item>
        <el-menu-item index="store">
          <el-icon><Location /></el-icon>
          <template #title>{{ t('title.stores') }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- Main content -->
    <el-main class="relative flex-1 p-4 bg-white dark:bg-[#1e1e1e] dark:text-white overflow-auto">
      

      <!-- Panels -->
      <div class="pr-12">
        <TestingPanel v-if="panelName === 'testing'" @toHistoryPanel="toHistoryPanel" />
        <TestingHistoryPanel v-else-if="panelName === 'history'" :ID="ID" />
        <DataManager v-else-if="panelName === 'data'" />
        <MockManager v-else-if="panelName === 'mock'" />
        <StoreManager v-else-if="panelName === 'store'" />
        <SecretManager v-else-if="panelName === 'secret'" />
        <WelcomePage v-else />
      </div>
    </el-main>
    <!-- Settings button top-right -->
      <div class="absolute top-4 right-4 cursor-pointer z-10" @click="settingDialogVisible = true" title="Settings">
        <el-icon size="20"><Setting /></el-icon>
      </div>

    <!-- Footer version link -->
    <div class="absolute bottom-2 right-4 text-sm z-10">
      <a
        :href="appVersionLink"
        target="_blank"
        rel="noopener"
        class="text-emerald-500 hover:bg-emerald-100 dark:hover:bg-emerald-900 transition rounded px-2 py-1 select-none"
      >
        {{ appVersion }}
      </a>
    </div>
  </el-container>

  <!-- Settings Dialog -->
  <el-dialog
    v-model="settingDialogVisible"
    :title="t('title.setting')"
    width="50%"
    draggable
    destroy-on-close
    class="!p-6"
  >
    <el-row class="mb-6 items-center">
      <el-col :span="4" class="font-semibold text-gray-700 dark:text-gray-300">Theme:</el-col>
      <el-col :span="18">
        <el-select v-model="theme" placeholder="Select a theme" class="w-full">
          <el-option v-for="item in allThemes" :key="item" :label="item" :value="item" />
        </el-select>
        <el-icon class="ml-3 align-middle">
          <el-link
            href="https://github.com/LinuxSuRen/atest-ext-data-swagger/tree/master/data/theme"
            target="_blank"
            class="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
            > <Help /> </el-link
          >
        </el-icon>
      </el-col>
    </el-row>

    <el-row class="mb-6 items-center">
      <el-col :span="4" class="font-semibold text-gray-700 dark:text-gray-300">Language:</el-col>
      <el-col :span="18" class="flex items-center gap-3">
        <el-tag class="text-base">{{ t('language') }}</el-tag>
        <el-dropdown trigger="click" @command="handleChangeLan">
          <el-icon><ArrowDown /></el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="chinese">中文</el-dropdown-item>
              <el-dropdown-item command="english">English</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-col>
    </el-row>

    <el-row class="items-center">
      <el-col :span="4" class="font-semibold text-gray-700 dark:text-gray-300">Dark Mode:</el-col>
      <el-col :span="18">
        <el-switch type="primary" v-model="asDarkMode" />
      </el-col>
    </el-row>
  </el-dialog>
</template>

<style scoped>
.el-menu-vertical:not(.el-menu--collapse) {
  width: 200px;
}

.el-menu-vertical {
  border-right: none !important;
}
</style>