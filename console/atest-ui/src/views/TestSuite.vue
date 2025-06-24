<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit, CopyDocument, Delete, View } from '@element-plus/icons-vue'
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
let querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name, props.name!)
const querySwaggers = SwaggerSuggestion()

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
  querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name!, props.name!)
}

Magic.AdvancedKeys([{
  Keys: ['Alt+N', 'Alt+dead'],
  Func: openNewTestCaseDialog,
  Description: 'Open new test case dialog',
}, {
  Keys: ['Alt+S', 'Alt+ß'],
  Func: saveTestSuite,
  Description: 'Save test suite',
}])

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
        }
      )

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
  <div class="common-layout">
    <el-form :rules="rules"
      ref="testSuiteFormRef"
      :model="suite"
      label-width="auto">
      {{ t('tip.testsuite') }}<EditButton :value="suite.name" @changed="renameTestSuite"/>

      <el-form-item :label="t('tip.apiAddress')" prop="api">
        <HistoryInput placeholder="API" v-model="suite.api" group="apiAddress" />
      </el-form-item>
      <table class="full-width">
        <tbody>
        <tr>
          <td>
            <el-select
              v-model="suite.spec.kind"
              class="m-2"
              placeholder="API Spec Kind"
              size="default"
            >
              <el-option
                v-for="item in apiSpecKinds"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </td>
          <td>
              <el-autocomplete
                  v-model="suite.spec.url"
                  :fetch-suggestions="querySwaggers"
              />
          </td>
        </tr>
       </tbody>
      </table>

      <el-collapse>
        <el-collapse-item :title="t('title.parameter')">
          <el-table :data="suite.param" class="full-width">
            <el-table-column :label="t('field.key')" width="180">
              <template #default="scope">
                <el-input v-model="scope.row.key" :placeholder="t('field.key')" @change="paramChange" />
              </template>
            </el-table-column>
            <el-table-column :label="t('field.value')">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" :placeholder="t('field.value')" />
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-collapse-item>
          <el-collapse-item v-if="suite.spec.rpc">
              <div>
                  <span>{{ t('title.refelction') }}</span>
                  <el-switch v-model="suite.spec.rpc.serverReflection" />
              </div>
              <div>
                  <span>{{ t('title.protoContent') }}</span>
                  <el-input
                      v-model="suite.spec.rpc.raw"
                      :autosize="{ minRows: 4, maxRows: 8 }"
                      type="textarea"
                  />
              </div>
              <div>
                  <span>{{ t('title.protoImport') }}</span>
                  <el-input class="mx-1" v-model="suite.spec.rpc.import"></el-input>
              </div>
              <div>
                  <span>{{ t('title.protoFile') }}</span>
                  <el-input class="mx-1" v-model="suite.spec.rpc.protofile"></el-input>
              </div>
          </el-collapse-item>
          <el-collapse-item :title="t('title.secure')">
              <el-switch v-model="suite.spec.secure.insecure" active-text="Insecure" inactive-text="Secure" inline-prompt/>
          </el-collapse-item>
        <el-collapse-item :title="t('title.proxy')">
          <div>
            <el-form-item :label="t('proxy.http')" prop="proxy.http">
              <el-input class="mx-1" v-model="suite.proxy.http" placeholder="HTTP Proxy"></el-input>
            </el-form-item>
          </div>
          <div>
            <el-form-item :label="t('proxy.https')" prop="proxy.http">
              <el-input class="mx-1" v-model="suite.proxy.https" placeholder="HTTPS Proxy"></el-input>
            </el-form-item>
          </div>
          <div>
            <el-form-item :label="t('proxy.no')">
              <el-input class="mx-1" v-model="suite.proxy.no" placeholder="No Proxy"></el-input>
            </el-form-item>
          </div>
        </el-collapse-item>
      </el-collapse>

      <div class="button-container">
          Share link: <el-input readonly v-model="shareLink" style="width: 80%" />
      </div>
      <div class="button-container">
        <Button type="primary" @click="updateTestSuiteForm(testSuiteFormRef)" v-if="!Cache.GetCurrentStore().readOnly">{{
          t('button.save')
        }}</Button>
        <Button type="primary" @click="updateTestSuiteForm(testSuiteFormRef)" disabled v-if="Cache.GetCurrentStore().readOnly">{{
          t('button.save')
        }}</Button>
        <Button type="danger" @click="del" :icon="Delete" test-id="suite-del-but">{{
          t('button.delete')
        }}</Button>
        <Button type="primary" @click="convert" test-id="convert">{{
          t('button.export')
        }}</Button>
        <Button
          type="primary"
          @click="openDuplicateDialog"
          :icon="CopyDocument"
          test-id="duplicate"
          >{{ t('button.duplicate') }}</Button
        >
        <Button type="primary" @click="viewYaml" :icon="View" test-id="view-yaml">{{
          t('button.viewYaml')
        }}</Button>
      </div>
      <div class="button-container">
        <Button
          type="primary"
          @click="openNewTestCaseDialog"
          :icon="Edit"
          test-id="open-new-case-dialog"
          >{{ t('button.newtestcase') }}</Button
        >
      </div>
    </el-form>
  </div>

  <el-dialog v-model="dialogVisible" :title="t('title.createTestCase')" width="40%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="testCaseForm"
          ref="testcaseFormRef"
          status-icon
          label-width="60px"
        >
          <el-form-item :label="t('field.name')" prop="name">
            <el-input v-model="testCaseForm.name" test-id="case-form-name" />
          </el-form-item>
          <el-form-item
            label="Method"
            prop="method"
            v-if="suite.spec.kind !== 'tRPC' && suite.spec.kind !== 'gRPC'"
          >
            <el-select
              v-model="testCaseForm.request.method"
              class="m-2"
              placeholder="Method"
              size="middle"
              test-id="case-form-method"
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
              :fetch-suggestions="querySuggestedAPIs"
              @select="handleAPISelect"
              placeholder="API Address"
              style="width: 100%; margin-left: 5px; margin-right: 5px"
              test-id="case-form-api"
            >
              <template #default="{ item }">
                <div class="value">{{ item.request.method }}</div>
                <span class="link">{{ item.request.api }}</span>
              </template>
            </el-autocomplete>
          </el-form-item>
          <el-form-item>
            <Button
              type="primary"
              @click="submitTestCaseForm(testcaseFormRef)"
              v-loading="suiteCreatingLoading"
              test-id="case-form-submit"
              >{{ t('button.submit') }}</Button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>

  <el-dialog
    v-model="yamlDialogVisible"
    :title="t('button.viewYaml')"
    :fullscreen="isFullScreen"
    width="40%"
    draggable
  >
    <Button type="primary" @click="isFullScreen = !isFullScreen" style="margin-bottom: 10px">
      <p>{{ isFullScreen ? t('button.cancelFullScreen') : t('button.fullScreen') }}</p>
    </Button>
    <el-scrollbar>
      <Codemirror v-model="yamlFormat" />
    </el-scrollbar>
  </el-dialog>

  <el-drawer v-model="testSuiteDuplicateDialog">
    <template #default>
      New Test Suite Name:<el-input v-model="targetSuiteDuplicateName" />
    </template>
    <template #footer>
      <Button type="primary" @click="duplicateTestSuite">{{ t('button.ok') }}</Button>
    </template>
  </el-drawer>
</template>

<style scoped>
.button-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 8px;
}

.button-container > .Button + .Button {
  margin-left: 0px;
}
</style>
