<template>
  <div class="pb-6">
    <el-breadcrumb separator="/">
      <el-breadcrumb-item v-for="(item, index) in $route.matched" :key="index">{{ t(item.meta.title) }}</el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>
