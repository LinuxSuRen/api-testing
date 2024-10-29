<script setup lang="ts">
import { ref } from 'vue';
import { Codemirror } from 'vue-codemirror';
import { API } from './net';
import {useI18n} from "vue-i18n";
import EditButton from '../components/EditButton.vue'

const { t } = useI18n()

interface MockConfig {
  Config: string
  Prefix: string
  Port: number
}
const mockConfig = ref({} as MockConfig);
const link = ref('')
API.GetMockConfig((d) => {
  mockConfig.value = d
  link.value = window.location.origin + d.Prefix + "/api.json"
})
const prefixChanged = (p: string) => {
  mockConfig.value.Prefix = p
}
const portChanged = (p: number) => {
  mockConfig.value.Port = p
}
const tabActive = ref('yaml')
const insertSample = () => {
    mockConfig.value.Config = `objects:
  - name: projects
    initCount: 3
    sample: |
      {
        "name": "api-testing",
        "color": "{{ randEnum "blue" "read" "pink" }}"
      }
items:
  - name: base64
    request:
      path: /v1/base64
    response:
      body: aGVsbG8=
      encoder: base64`
}
</script>

<template>
    <div>
        <el-button type="primary" @click="insertSample">{{t('button.insertSample')}}</el-button>
        <el-button type="warning" @click="API.ReloadMockServer(mockConfig)">{{t('button.reload')}}</el-button>
        <el-divider direction="vertical" />
        <el-link target="_blank" :href="link">{{ link }}</el-link> <!-- Noncompliant -->
    </div>
    <div>
      API Prefix:<EditButton :value="mockConfig.Prefix" @changed="prefixChanged"/>
      Port:<EditButton :value="mockConfig.Port" @changed="portChanged"/>
    </div>
    <div>
        <el-tabs v-model="tabActive">
            <el-tab-pane label="YAML" name="yaml">
                <Codemirror v-model="mockConfig.Config" />
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
