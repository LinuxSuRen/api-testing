<script setup lang="ts">
import { ref, watch } from 'vue';
import { Codemirror } from 'vue-codemirror';
import yaml from 'js-yaml';
import { jsonSchema } from "codemirror-json-schema";
import { NewTemplateLangComplete, NewHeaderLangComplete } from './languageComplete'
import { jsonLanguage } from "@codemirror/lang-json"
import { API } from './net';
import {useI18n} from "vue-i18n";
import EditButton from '../components/EditButton.vue'

const { t } = useI18n()

const mockschema = ref({}); // Type assertion to any for JSON schema
const jsonComplete = NewTemplateLangComplete(jsonLanguage)
const headerComplete = NewHeaderLangComplete(jsonLanguage)

API.GetSchema('mock').then((schema) => {
    if (schema.success && schema.message !== '') {
        mockschema.value = JSON.parse(schema.message);
    }
});

const logOutput = ref('');
API.GetStream((stream) => {
    try {
        const data = JSON.parse(stream);
        logOutput.value += `${data.result.message}`;
        console.log('Received schema:', logOutput.value);
    } catch (e) {}
});

interface MockConfig {
  Config: string
  ConfigAsJSON: string
  Prefix: string
  Port: number
}

const tabActive = ref('yaml')
const mockConfig = ref({} as MockConfig);
watch(mockConfig, (newValue) => {
    if (tabActive.value === 'json') {
        // convert JSON to YAML string
        try {
            newValue.Config = jsonToYaml(newValue.ConfigAsJSON);
        } catch (e) {
        }
    } else {
        newValue.ConfigAsJSON = JSON.stringify(yaml.load(newValue.Config), null, 2);
    }
}, { deep: true });

function jsonToYaml(jsonData: object | string): string {
  const data = typeof jsonData === 'string' 
    ? JSON.parse(jsonData) 
    : jsonData;
  
  return yaml.dump(data, {
    indent: 2,
    skipInvalid: true,
    noRefs: true,
    lineWidth: -1,
  });
}

const link = ref('')
API.GetMockConfig((d) => {
  mockConfig.value = d
  link.value = `http://${window.location.hostname}:${d.Port}${d.Prefix}/api.json`
})
const prefixChanged = (p: string) => {
  mockConfig.value.Prefix = p
}
const portChanged = (p: number) => {
  mockConfig.value.Port = p
}
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
    <div class="config">
      API Prefix:<EditButton :value="mockConfig.Prefix" @changed="prefixChanged"/>
      Port:<EditButton :value="mockConfig.Port" @changed="portChanged"/>
    </div>
    <div>
        <el-tabs v-model="tabActive">
            <el-tab-pane label="YAML" name="yaml">
                <Codemirror v-model="mockConfig.Config"
                    :extensions="[jsonComplete, headerComplete]" />
            </el-tab-pane>
            <el-tab-pane label="JSON" name="json">
                <Codemirror v-model="mockConfig.ConfigAsJSON"
                    :extensions="[jsonSchema(mockschema), jsonComplete, headerComplete]" />
            </el-tab-pane>
        </el-tabs>
        <el-card class="log-output" shadow="hover">
          <template #header>
            <span>{{ t('label.logs') }}</span>
          </template>
          <el-scrollbar height="200px" ref="logScrollbar">
            <pre style="white-space: pre-wrap; word-break: break-all;">{{ logOutput }}</pre>
          </el-scrollbar>
        </el-card>
    </div>
</template>

<style>
.config {
  margin: 6px 0; 
  display: flex;
  align-items: center; 
  gap: 8px; 
}
</style>
