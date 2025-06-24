<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { Magic } from '@/views/magicKeys'

const { t } = useI18n()
const keyBindingsDialogVisible = ref(false)
const keyBindings = ref([])
const showKeyBindingsDialog = (keys: any) => {
  keyBindingsDialogVisible.value = true
  keyBindings.value = keys.detail || []
}

onMounted(() => {
  window.addEventListener(Magic.MagicKeyEventName, showKeyBindingsDialog)
})

onBeforeUnmount(() => {
  window.removeEventListener(Magic.MagicKeyEventName, showKeyBindingsDialog)
})
</script>

<template>
    <el-drawer v-model="keyBindingsDialogVisible" size="50%">
      <template #header>
        <h4>{{ t('title.keyBindings') }}</h4>
      </template>
      <template #default>
        <el-table :data="keyBindings">
          <el-table-column prop="keys" :label="t('field.shortcut')" />
          <el-table-column prop="description" :label="t('field.description')" />
        </el-table>
      </template>
    </el-drawer>
</template>
