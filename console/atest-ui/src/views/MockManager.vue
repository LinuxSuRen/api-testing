<script setup lang="ts">
import { ref } from 'vue';
import { Codemirror } from 'vue-codemirror';
import { API } from './net';
import {useI18n} from "vue-i18n";

const { t } = useI18n()

const mockConfig = ref('');
const link = ref('')
API.GetMockConfig((d) => {
  mockConfig.value = d.Config
  link.value = window.location.origin + d.Prefix + "/api.json"
})
const tabActive = ref('yaml')
const insertSample = () => {
    mockConfig.value = `objects:
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
        <el-tabs v-model="tabActive">
            <el-tab-pane label="YAML" name="yaml">
                <Codemirror v-model="mockConfig" />
            </el-tab-pane>
        </el-tabs>
    </div>
</template>
