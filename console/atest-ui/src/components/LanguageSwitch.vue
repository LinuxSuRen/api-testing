<template>
    <el-col style="display: flex; align-items: center; vertical-align: middle;">
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
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { Cache } from '../utils/cache'

const { t, locale: i18nLocale } = useI18n()

const locale = ref(Cache.GetPreference().language)
i18nLocale.value = locale.value

watch(locale, (value) =>{
Cache.WatchLocale(value)
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
