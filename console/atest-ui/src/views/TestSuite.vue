<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit, CopyDocument, Delete, View, Setting, Link, Plus, Check, Download } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Suite, TestCase, Pair } from './types'
import { NewSuggestedAPIsQuery, GetHTTPMethods, SwaggerSuggestion } from './types'
import EditButton from '../components/EditButton.vue'
import HistoryInput from '../components/HistoryInput.vue'
import Button from '../components/Button.vue'
import { Cache } from './cache'
import { useI18n } from 'vue-i18n'
import { API } from './net'
import { Magic } from './magicKeys'
import { Codemirror } from 'vue-codemirror'
import yaml from 'js-yaml'

const { t } = useI18n()

const props = defineProps({
  name: String
})
const emit = defineEmits(['updated'])

const suite = ref({
  name: '',
  api: '',
  param: [] as Pair[],
  spec: {
    kind: '',
    url: '',
    rpc: {
      raw: '',
      protofile: '',
      serverReflection: false
    },
    secure: {
      insecure: true
    }
  },
  proxy: {
    http: '',
    https: '',
    no: ''
  }
} as Suite)
const shareLink = ref('')
const querySuggestedAPIs = ref(NewSuggestedAPIsQuery(Cache.GetCurrentStore().name, props.name!))
const querySwaggers = ref(SwaggerSuggestion())

function load() {
  const store = Cache.GetCurrentStore()
  if (!props.name || store.name === '') return

  API.GetTestSuite(
    props.name,
    (e) => {
      suite.value = e
      if (suite.value.param.length === 0) {
        suite.value.param.push({
          key: '',
          value: ''
        } as Pair)
      }
      if (!suite.value.proxy) {
        suite.value.proxy = {
          http: '',
          https: '',
          no: ''
        }
      }
      if (!suite.value.spec.secure) {
        suite.value.spec.secure = {
          insecure: false
        }
      }

      shareLink.value = `${window.location.href}api/v1/suites/${e.name}/yaml?x-store-name=${store.name}`
    },
    (e) => {
      ElMessage.error('Oops, ' + e)
    }
  )
}
load()
watch(props, () => {
  load()
})

const testSuiteFormRef = ref<FormInstance>()
const updateTestSuiteForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      saveTestSuite()
    }
  })
}
const saveTestSuite = () => {
  let oldImportPath = ''
  let hasImport = false
  if (suite.value.spec && suite.value.spec.rpc) {
    oldImportPath = suite.value.spec.rpc.import
    hasImport = true
    if (typeof oldImportPath === 'string' && oldImportPath !== '') {
      const importPath = oldImportPath.split(',')
      suite.value.spec.rpc.import = importPath
    }
  }

  API.UpdateTestSuite(
    suite.value,
    (e) => {
      if (e.error === '') {
        ElMessage({
          message: 'Updated.',
          type: 'success'
        })
      } else {
        ElMessage.error('Oops, ' + e.message)
      }

      if (hasImport) {
        suite.value.spec.rpc.import = oldImportPath
      }
    },
    (e) => {
      if (hasImport) {
        suite.value.spec.rpc.import = oldImportPath
      }
      ElMessage.error('Oops, ' + e)
    }
  )
}
Magic.Keys(saveTestSuite, ['Alt+S', 'Alt+ß'])

const isFullScreen = ref(false)
const dialogVisible = ref(false)
const testcaseFormRef = ref<FormInstance>()
const testCaseForm = reactive({
  suiteName: '',
  name: '',
  request: {
    api: '',
    method: 'GET'
  }
})
const rules = reactive<FormRules<Suite>>({
  name: [{ required: true, message: 'Please input TestCase name', trigger: 'blur' }],
  'proxy.http': [{ type: 'url', message: 'Please input a valid URL', trigger: 'blur' }],
})

function openNewTestCaseDialog() {
  dialogVisible.value = true
  querySuggestedAPIs.value = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name!, props.name!)
}
Magic.Keys(openNewTestCaseDialog, ['Alt+N', 'Alt+dead'])

const submitTestCaseForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean) => {
    if (valid) {
      suiteCreatingLoading.value = true

      API.CreateTestCase({
        suiteName: props.name,
        name: testCaseForm.name,
        request: testCaseForm.request
      }, () => {
        suiteCreatingLoading.value = false
        emit('updated', props.name, testCaseForm.name)
      }, (e) => {
        suiteCreatingLoading.value = false
        ElMessage.error('Oops, ' + e)
      })

      dialogVisible.value = false
    }
  })
}

function del() {
  API.DeleteTestSuite(
    props.name,
    () => {
      ElMessage({
        message: 'Deleted.',
        type: 'success'
      })
      emit('updated')
    },
    (e) => {
      ElMessage.error('Oops, ' + e)
    }
  )
}

function convert() {
  API.ConvertTestSuite(
    props.name,
    'jmeter',
    (e) => {
      const blob = new Blob([e.message], { type: `text/xml;charset=utf-8;` })
      const link = document.createElement('a')
      if (link.download !== undefined) {
        const url = URL.createObjectURL(blob)
        link.setAttribute('href', url)
        link.setAttribute('download', `jmeter.jmx`)
        link.style.visibility = 'hidden'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      }

      ElMessage({
        message: 'Converted.',
        type: 'success'
      })
      emit('updated')
    },
    (e) => {
      ElMessage.error('Oops, ' + e)
    }
  )
}

const suiteCreatingLoading = ref(false)

const apiSpecKinds = [
  {
    value: 'swagger',
    label: 'Swagger'
  },
  {
    value: 'openapi',
    label: 'OpenAPI'
  }
]

const handleAPISelect = (item: TestCase) => {
  testCaseForm.method = item.request.method
  if (testCaseForm.name === '') {
    testCaseForm.name = item.name
  }
  testCaseForm.request = item.request
}

function paramChange() {
  const form = suite.value.param
  let lastItem = form[form.length - 1]
  if (lastItem.key !== '') {
    suite.value.param.push({
      key: '',
      value: ''
    } as Pair)
  }
}

const yamlFormat = ref('')
const yamlDialogVisible = ref(false)

function viewYaml() {
  yamlDialogVisible.value = true
  API.GetTestSuiteYaml(props.name, (d) => {
    try {
      yamlFormat.value = yaml.dump(yaml.load(atob(d.data)))
    } catch (e) {
      ElMessage.error('Oops, ' + e)
    }
  })
}

const openDuplicateDialog = () => {
  testSuiteDuplicateDialog.value = true
  targetSuiteDuplicateName.value = props.name + '-copy'
}
const duplicateTestSuite = () => {
  API.DuplicateTestSuite(props.name, targetSuiteDuplicateName.value, (d) => {
    testSuiteDuplicateDialog.value = false
    ElMessage({
      message: 'Duplicated.',
      type: 'success'
    })
    emit('updated')
  })
}
const testSuiteDuplicateDialog = ref(false)
const targetSuiteDuplicateName = ref('')

const renameTestSuite = (name: string) => {
  API.RenameTestSuite(props.name, name, (d) => {
    emit('updated', name)
  })
}
</script>

<template>
  <div class="test-suite-container">
    <!-- Header Section -->
    <div class="header-section">
      <div class="suite-title">
        <h2 class="title-text">{{ t('tip.testsuite') }}</h2>
        <EditButton :value="suite.name" @changed="renameTestSuite" class="edit-button"/>
      </div>
      <div class="suite-actions">
        <el-button 
          type="primary" 
          :icon="Plus" 
          @click="openNewTestCaseDialog"
          test-id="open-new-case-dialog"
          class="primary-action"
        >
          {{ t('button.newtestcase') }}
        </el-button>
      </div>
    </div>

    <!-- Main Form -->
    <el-form 
      :rules="rules"
      ref="testSuiteFormRef"
      :model="suite"
      label-width="140px"
      class="main-form"
    >
      <!-- API Configuration Card -->
      <div class="config-card">
        <div class="card-header">
          <h3 class="card-title">
            <el-icon><Setting /></el-icon>
            API Configuration
          </h3>
        </div>
        <div class="card-content">
          <el-form-item :label="t('tip.apiAddress')" prop="api" class="form-item">
            <HistoryInput 
              placeholder="Enter API address" 
              v-model="suite.api" 
              group="apiAddress"
              class="full-width-input"
            />
          </el-form-item>
          
          <div class="spec-row">
            <el-form-item label="API Spec" class="spec-kind">
              <el-select
                v-model="suite.spec.kind"
                placeholder="Select API Spec Kind"
                size="default"
                class="spec-select"
              >
                <el-option
                  v-for="item in apiSpecKinds"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="Spec URL" class="spec-url">
              <el-autocomplete
                v-model="suite.spec.url"
                :fetch-suggestions="querySwaggers.value"
                placeholder="Enter specification URL"
                class="full-width-input"
              />
            </el-form-item>
          </div>
        </div>
      </div>

      <!-- Advanced Configuration -->
      <div class="config-card">
        <div class="card-header">
          <h3 class="card-title">Advanced Configuration</h3>
        </div>
        <div class="card-content">
          <el-collapse class="custom-collapse">
            <el-collapse-item name="parameters">
              <template #title>
                <div class="collapse-title">
                  <el-icon><Edit /></el-icon>
                  {{ t('title.parameter') }}
                </div>
              </template>
              <div class="parameters-section">
                <el-table :data="suite.param" class="params-table">
                  <el-table-column :label="t('field.key')" width="200">
                    <template #default="scope">
                      <el-input 
                        v-model="scope.row.key" 
                        :placeholder="t('field.key')" 
                        @change="paramChange"
                        class="param-input"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column :label="t('field.value')">
                    <template #default="scope">
                      <el-input 
                        v-model="scope.row.value" 
                        :placeholder="t('field.value')"
                        class="param-input"
                      />
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </el-collapse-item>

            <el-collapse-item v-if="suite.spec.rpc" name="rpc">
              <template #title>
                <div class="collapse-title">
                  <el-icon><Setting /></el-icon>
                  RPC Configuration
                </div>
              </template>
              <div class="rpc-section">
                <div class="rpc-item">
                  <label class="rpc-label">{{ t('title.refelction') }}</label>
                  <el-switch v-model="suite.spec.rpc.serverReflection" />
                </div>
                <div class="rpc-item">
                  <label class="rpc-label">{{ t('title.protoContent') }}</label>
                  <el-input
                    v-model="suite.spec.rpc.raw"
                    :autosize="{ minRows: 4, maxRows: 8 }"
                    type="textarea"
                    class="proto-textarea"
                  />
                </div>
                <div class="rpc-item">
                  <label class="rpc-label">{{ t('title.protoImport') }}</label>
                  <el-input v-model="suite.spec.rpc.import" class="proto-input" />
                </div>
                <div class="rpc-item">
                  <label class="rpc-label">{{ t('title.protoFile') }}</label>
                  <el-input v-model="suite.spec.rpc.protofile" class="proto-input" />
                </div>
              </div>
            </el-collapse-item>

            <el-collapse-item name="security">
              <template #title>
                <div class="collapse-title">
                  <el-icon><View /></el-icon>
                  {{ t('title.secure') }}
                </div>
              </template>
              <div class="security-section">
                <el-switch 
                  v-model="suite.spec.secure.insecure" 
                  active-text="Insecure" 
                  inactive-text="Secure" 
                  inline-prompt
                  class="security-switch"
                />
              </div>
            </el-collapse-item>

            <el-collapse-item name="proxy">
              <template #title>
                <div class="collapse-title">
                  <el-icon><Link /></el-icon>
                  {{ t('title.proxy') }}
                </div>
              </template>
              <div class="proxy-section">
                <el-form-item :label="t('proxy.http')" prop="proxy.http">
                  <el-input 
                    v-model="suite.proxy.http" 
                    placeholder="HTTP Proxy URL"
                    class="proxy-input"
                  />
                </el-form-item>
                <el-form-item :label="t('proxy.https')" prop="proxy.https">
                  <el-input 
                    v-model="suite.proxy.https" 
                    placeholder="HTTPS Proxy URL"
                    class="proxy-input"
                  />
                </el-form-item>
                <el-form-item :label="t('proxy.no')">
                  <el-input 
                    v-model="suite.proxy.no" 
                    placeholder="No Proxy URLs"
                    class="proxy-input"
                  />
                </el-form-item>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
      </div>

      <!-- Share Link Section -->
      <div class="config-card">
        <div class="card-header">
          <h3 class="card-title">
            <el-icon><Link /></el-icon>
            Share Link
          </h3>
        </div>
        <div class="card-content">
          <el-input 
            readonly 
            v-model="shareLink" 
            class="share-input"
            placeholder="Share link will appear here"
          >
            <template #append>
              <el-button :icon="CopyDocument" @click="navigator.clipboard.writeText(shareLink)">
                Copy
              </el-button>
            </template>
          </el-input>
        </div>
      </div>

      <!-- Action Buttons -->
     <div class="actions-section">
  <div class="primary-actions">
    <Button 
      type="primary" 
      @click="updateTestSuiteForm(testSuiteFormRef)" 
      :disabled="Cache.GetCurrentStore().readOnly"
      class="save-button"
    >
      {{ t('button.save') }}
    </Button>
  </div>
        
        <div class="secondary-actions">
    <Button 
      type="info" 
      @click="convert" 
      test-id="convert"
    >
      Export to JMeter
    </Button>
    <Button
      type="info"
      @click="openDuplicateDialog"
      test-id="duplicate"
    >
      Duplicate Suite
    </Button>
    <Button 
      type="info" 
      @click="viewYaml" 
      test-id="view-yaml"
    >
      View YAML
    </Button>
    <Button 
      type="danger" 
      @click="del" 
      test-id="suite-del-but"
      class="delete-button"
    >
      Delete Suite
    </Button>
        </div>
      </div>
    </el-form>
  </div>

  <!-- Create Test Case Dialog -->
  <el-dialog 
    v-model="dialogVisible" 
    :title="t('title.createTestCase')" 
    width="500px" 
    draggable
    class="create-dialog"
  >
    <el-form
      :rules="rules"
      :model="testCaseForm"
      ref="testcaseFormRef"
      status-icon
      label-width="80px"
      class="dialog-form"
    >
      <el-form-item :label="t('field.name')" prop="name">
        <el-input 
          v-model="testCaseForm.name" 
          test-id="case-form-name"
          placeholder="Enter test case name"
        />
      </el-form-item>
      <el-form-item
        label="Method"
        prop="method"
        v-if="suite.spec.kind !== 'tRPC' && suite.spec.kind !== 'gRPC'"
      >
        <el-select
          v-model="testCaseForm.request.method"
          placeholder="Select HTTP Method"
          test-id="case-form-method"
          class="method-select"
        >
          <el-option
            v-for="item in GetHTTPMethods()"
            :key="item.value"
            :label="item.key"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="API" prop="api">
        <el-autocomplete
          v-model="testCaseForm.request.api"
          :fetch-suggestions="querySuggestedAPIs.value"
          @select="handleAPISelect"
          placeholder="Enter API endpoint"
          test-id="case-form-api"
          class="api-autocomplete"
        >
          <template #default="{ item }">
            <div class="api-suggestion">
              <span class="method-badge">{{ item.request.method }}</span>
              <span class="api-path">{{ item.request.api }}</span>
            </div>
          </template>
        </el-autocomplete>
      </el-form-item>
    </el-form>
    
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button
          type="primary"
          @click="submitTestCaseForm(testcaseFormRef)"
          :loading="suiteCreatingLoading"
          test-id="case-form-submit"
        >
          {{ t('button.submit') }}
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- YAML Viewer Dialog -->
  <el-dialog
    v-model="yamlDialogVisible"
    :title="t('button.viewYaml')"
    :fullscreen="isFullScreen"
    width="70%"
    draggable
    class="yaml-dialog"
  >
    <div class="yaml-controls">
      <el-button 
        type="primary" 
        @click="isFullScreen = !isFullScreen"
        :icon="isFullScreen ? 'Minus' : 'Plus'"
      >
        {{ isFullScreen ? t('button.cancelFullScreen') : t('button.fullScreen') }}
      </el-button>
    </div>
    <div class="yaml-content">
      <Codemirror v-model="yamlFormat" class="yaml-editor" />
    </div>
  </el-dialog>

  <!-- Duplicate Suite Drawer -->
  <el-drawer v-model="testSuiteDuplicateDialog" title="Duplicate Test Suite" size="400px">
    <div class="duplicate-content">
      <el-form label-width="120px">
        <el-form-item label="New Name:">
          <el-input 
            v-model="targetSuiteDuplicateName" 
            placeholder="Enter new suite name"
          />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <div class="drawer-footer">
        <el-button @click="testSuiteDuplicateDialog = false">Cancel</el-button>
        <el-button type="primary" @click="duplicateTestSuite">
          {{ t('button.ok') }}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<style scoped>
.test-suite-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
 
  min-height: 100vh;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.suite-title {
  display: flex;
  align-items: center;
  gap: 16px;
}

.title-text {
  margin: 0;
  font-size: 28px;
  font-weight: 700;

  background-clip: text;
}

.edit-button {
  margin-left: 12px;
}

.primary-action {
  padding: 12px 24px;
  font-weight: 600;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}

.primary-action:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
}

.main-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.config-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: all 0.3s ease;
}

.config-card:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.card-header {
  padding: 20px 24px;
  background: #409eff;
  color: white;
}

.card-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-content {
  padding: 24px;
}

.form-item {
  margin-bottom: 20px;
}

.full-width-input {
  width: 100%;
}

.spec-row {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: 20px;
  align-items: end;
}

.spec-select {
  width: 100%;
}

.custom-collapse {
  border: none;
  background: transparent;
}

.custom-collapse :deep(.el-collapse-item__header) {
  background: #f8f9fa;
  border: none;
  border-radius: 8px;
  margin-bottom: 8px;
  padding: 16px;
  font-weight: 600;
}

.custom-collapse :deep(.el-collapse-item__content) {
  padding: 20px 16px;
  background: #fafbfc;
  border-radius: 8px;
  margin-bottom: 16px;
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #2c3e50;
}

.parameters-section {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.params-table {
  width: 100%;
}

.param-input {
  border-radius: 6px;
}

.rpc-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rpc-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rpc-label {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.proto-textarea, .proto-input {
  border-radius: 6px;
}

.security-section {
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.security-switch {
  font-size: 16px;
}

.proxy-section {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.proxy-input {
  border-radius: 6px;
}

.share-input {
  border-radius: 8px;
}

.share-input :deep(.el-input__wrapper) {
  border-radius: 8px 0 0 8px;
}

.actions-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.primary-actions {
  display: flex;
  justify-content: center;
}

.save-button {
  padding: 12px 32px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}

.save-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
}

.secondary-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

.secondary-actions .el-button {
  border-radius: 6px;
  transition: all 0.3s ease;
}

.secondary-actions .el-button:hover {
  transform: translateY(-1px);
}

.delete-button {
  box-shadow: 0 4px 12px rgba(245, 108, 108, 0.3);
}

.delete-button:hover {
  box-shadow: 0 6px 20px rgba(245, 108, 108, 0.4);
}

.create-dialog :deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

.create-dialog :deep(.el-dialog__header) {

  color: white;
  padding: 20px 24px;
}

.dialog-form {
  padding: 20px 0;
}

.method-select, .api-autocomplete {
  width: 100%;
}

.api-suggestion {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 0;
}

.method-badge {
  background: #409eff;
  color: white;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  min-width: 60px;
  text-align: center;
}

.api-path {
  color: #606266;
  font-family: monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.yaml-dialog :deep(.el-dialog) {
  border-radius: 16px;
}

.yaml-controls {
  margin-bottom: 16px;
}

.yaml-content {
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
}

.yaml-editor {
  min-height: 400px;
}

.duplicate-content {
  padding: 20px 0;
}

.drawer-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 0;
}

/* Responsive Design */
@media (max-width: 768px) {
  .test-suite-container {
    padding: 16px;
  }
  
  .header-section {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .spec-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .secondary-actions {
    flex-direction: column;
  }
  
  .create-dialog {
    width: 90% !important;
  }
}

/* Animation for smooth transitions */
.config-card {
  animation: slideInUp 0.6s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Custom scrollbar */
:deep(.el-scrollbar__wrap) {
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 transparent;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar) {
  width: 6px;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar-track) {
  background: transparent;
}

:deep(.el-scrollbar__wrap::-webkit-scrollbar-thumb) {
  background-color: #c1c1c1;
  border-radius: 3px;
}
</style>