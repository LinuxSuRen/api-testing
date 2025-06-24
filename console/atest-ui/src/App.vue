<script setup lang="ts">
import {
    Document,
    Menu as IconMenu,
    Histogram,
    Location,
    Share,
    ArrowDown,
    Guide,
    DataAnalysis, Help, Setting
} from '@element-plus/icons-vue'
import { ref, watch, getCurrentInstance} from 'vue'
import { API } from './views/net'
import { Cache } from './views/cache'
import TestingPanel from './views/TestingPanel.vue'
import TestingHistoryPanel from './views/TestingHistoryPanel.vue'
import MockManager from './views/MockManager.vue'
import StoreManager from './views/StoreManager.vue'
import SecretManager from './views/SecretManager.vue'
import WelcomePage from './views/WelcomePage.vue'
import DataManager from './views/DataManager.vue'
import MagicKey from './components/MagicKey.vue'
import { useI18n } from 'vue-i18n'
import ElementPlus from 'element-plus'; 
import zhCn from 'element-plus/dist/locale/zh-cn.mjs' 
import enUS from 'element-plus/dist/locale/en.mjs'   

const { t, locale: i18nLocale } = useI18n()
const app = getCurrentInstance()?.appContext.app;

import { setAsDarkTheme, getThemes, setTheme, getTheme } from './theme'

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

  if (!version && !dirtyVersion) {
    return
  }

  if (dirtyVersion && dirtyVersion.length > 0) {
    appVersionLink.value = appVersionLink.value + '/commit/' + d.message.replace(dirtyVersion[0], '')
  } else if (version && version.length > 0) {
    appVersionLink.value = appVersionLink.value + '/releases/tag/' + version[0]
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

watch(locale, (newLocale: string) => {
  updateLocale(newLocale);
});
const updateLocale = (newLocale: string) => {
  locale.value = newLocale;
  i18nLocale.value = newLocale;
  Cache.WithLocale(newLocale);

  const elementLocale = newLocale === 'zh-CN' ? zhCn : enUS;
  app?.use(ElementPlus, { locale: elementLocale }); 
};
const handleChangeLan = (command: string) => {
  switch (command) {
    case "chinese":
    locale.value = 'zh-CN';
    break;
    case "english":
    locale.value = 'en-US';
    break;
  }
};

const ID = ref(null);
const toHistoryPanel = ({ ID: selectID, panelName: historyPanelName }) => {
  ID.value = selectID;
  panelName.value = historyPanelName;
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
  <el-container class="full-height">
    <el-aside width="auto" style="display: flex; flex-direction: column;">
      <el-radio-group v-model="isCollapse" class="el-menu">
        <el-radio-button :value="false">+</el-radio-button>
        <el-radio-button :value="true">-</el-radio-button>
      </el-radio-group>
      <el-menu
        class="el-menu-vertical full-height"
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
        <el-menu-item index="history" test-id="history-menu">
          <el-icon><histogram /></el-icon>
          <template #title>{{ t('title.history' )}}</template>
        </el-menu-item>
        <el-menu-item index="mock" test-id="mock-menu">
          <el-icon><Guide /></el-icon>
          <template #title>{{ t('title.mock' )}}</template>
        </el-menu-item>
        <el-menu-item index="data" test-id="data-menu">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>{{ t('title.data' )}}</template>
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

    <el-main class="center-zone">
      <div class="top-menu">
        <el-col style="display: flex; align-items: center;">
          <el-icon @click="settingDialogVisible=true" size="20"><Setting /></el-icon>
        </el-col>
      </div>
      <TestingPanel v-if="panelName === 'testing'" @toHistoryPanel="toHistoryPanel"/>
      <TestingHistoryPanel v-else-if="panelName === 'history'" :ID="ID"/>
      <DataManager v-else-if="panelName === 'data'" />
      <MockManager v-else-if="panelName === 'mock'" />
      <StoreManager v-else-if="panelName === 'store'" />
      <SecretManager v-else-if="panelName === 'secret'" />
      <WelcomePage v-else />
    </el-main>

    <div style="position: absolute; bottom: 0px; right: 10px;">
      <a :href=appVersionLink target="_blank" rel="noopener">{{appVersion}}</a>
    </div>
  </el-container>

    <el-dialog v-model="settingDialogVisible" :title="t('title.setting' )" width="50%" draggable destroy-on-close>
        <el-row>
            <el-col :span="4">
              Theme:
            </el-col>
            <el-col :span="18">
              <el-select v-model="theme" placeholder="Select a theme">
                  <el-option
                      v-for="item in allThemes"
                      :key="item"
                      :label="item"
                      :value="item"
                  />
              </el-select>
              <el-icon>
                  <el-link href="https://github.com/LinuxSuRen/atest-ext-data-swagger/tree/master/data/theme" target="_blank">
                      <Help />
                  </el-link>
              </el-icon>
            </el-col>
        </el-row>

        <el-row>
            <el-col :span="4">
              Language:
            </el-col>
            <el-col :span="18">
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
        </el-row>
        <el-row>
            <el-col :span="4">
              Dark Mode:
            </el-col>
            <el-col :span="18">
              <el-switch type="primary" data-intro="Switch light and dark modes" v-model="asDarkMode"/>
            </el-col>
        </el-row>
    </el-dialog>

    <MagicKey />
</template>

<style>
.el-menu-vertical:not(.el-menu--collapse) {
  width: 200px;
}
.el-menu-vertical:is(.el-menu--collapse) {
  width: 80px;
}
</style>
