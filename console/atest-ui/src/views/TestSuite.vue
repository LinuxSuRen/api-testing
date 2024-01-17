<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Suite, TestCase, Pair } from './types'
import { NewSuggestedAPIsQuery, GetHTTPMethods } from './types'
import { Cache } from './cache'
import { useI18n } from 'vue-i18n'
import { API } from './net'

const { t } = useI18n()

const props = defineProps({
  name: String,
})
const emit = defineEmits(['updated'])
let querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name, props.name!)

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
    }
  }
} as Suite)
function load() {
  const store = Cache.GetCurrentStore()
  if (!props.name || store.name === "") return

  API.GetTestSuite(props.name, (e) => {
      suite.value = e
      if (suite.value.param.length === 0) {
        suite.value.param.push({
          key: '',
          value: ''
        } as Pair)
      }
    }, (e) => {
      ElMessage.error('Oops, ' + e)
    })
}
load()
watch(props, () => {
  load()
})

function save() {
  let oldImportPath = ""
  let hasImport = false
  if (suite.value.spec && suite.value.spec.rpc) {
    oldImportPath = suite.value.spec.rpc.import
    hasImport = true
    if (typeof oldImportPath === 'string' && oldImportPath !== '') {
      const importPath = oldImportPath.split(',')
      suite.value.spec.rpc.import = importPath
    }
  }

  API.UpdateTestSuite(suite.value, (e) => {
      if (e.error === "") {
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
    }, (e) => {
      if (hasImport) {
        suite.value.spec.rpc.import = oldImportPath
      }
      ElMessage.error('Oops, ' + e)
    })
}

const dialogVisible = ref(false)
const testcaseFormRef = ref<FormInstance>()
const testCaseForm = reactive({
  suiteName: '',
  name: '',
  api: '',
  method: 'GET'
})
const rules = reactive<FormRules<Suite>>({
  name: [{ required: true, message: 'Please input TestCase name', trigger: 'blur' }]
})

function openNewTestCaseDialog() {
  dialogVisible.value = true
  querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name!, props.name!)
}

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      suiteCreatingLoading.value = true

      API.CreateTestCase({
        suiteName: props.name,
        name: testCaseForm.name,
        api: testCaseForm.api,
        method: testCaseForm.method
      }, () => {
        suiteCreatingLoading.value = false
        emit('updated', 'hello from child')
      })

      dialogVisible.value = false
    }
  })
}

function del() {
  API.DeleteTestSuite(props.name, () => {
      ElMessage({
        message: 'Deleted.',
        type: 'success'
      })
      emit('updated')
    }, (e) => {
      ElMessage.error('Oops, ' + e)
    })
}

function convert() {
  API.ConvertTestSuite(props.name, 'jmeter', (e) => {
      const blob = new Blob([e.message], { type: `text/xml;charset=utf-8;` });
      const link = document.createElement('a');
      if (link.download !== undefined) {
        const url = URL.createObjectURL(blob);
        link.setAttribute('href', url);
        link.setAttribute('download', `jmeter.jmx`);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      }

      ElMessage({
        message: 'Converted.',
        type: 'success'
      })
      emit('updated')
    }, (e) => {
      ElMessage.error('Oops, ' + e)
    })
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
  if (testCaseForm.method === '') {
    testCaseForm.method = item.request.method
  }
  if (testCaseForm.name === '') {
    testCaseForm.name = item.name
  }
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
</script>

<template>
  <div class="common-layout">
    <el-text class="mx-1" type="primary">{{ suite.name }}</el-text>

    <el-input class="w-50 m-2" v-model="suite.api" placeholder="API" test-id="suite-editor-api"></el-input>

    <el-select v-model="suite.spec.kind" class="m-2" placeholder="API Spec Kind" size="middle">
      <el-option
        v-for="item in apiSpecKinds"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
    <el-input class="mx-1" v-model="suite.spec.url" placeholder="API Spec URL"></el-input>

    <div>
      <span>{{ t('title.parameter') }}</span>
      <el-table :data="suite.param" style="width: 100%">
        <el-table-column label="Key" width="180">
          <template #default="scope">
            <el-input v-model="scope.row.key" placeholder="Key" @change="paramChange"/>
          </template>
        </el-table-column>
        <el-table-column label="Value">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-input v-model="scope.row.value" placeholder="Value" />
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-divider />
    </div>

    <div v-if="suite.spec.rpc">
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
      <el-divider />
    </div>

    <el-button type="primary" @click="save" v-if="!Cache.GetCurrentStore().readOnly">{{ t('button.save') }}</el-button>
    <el-button type="primary" @click="save" disabled v-if="Cache.GetCurrentStore().readOnly">{{ t('button.save') }}</el-button>
    <el-button type="primary" @click="del" test-id="suite-del-but">{{ t('button.delete') }}</el-button>
    <el-button type="primary" @click="openNewTestCaseDialog" :icon="Edit" test-id="open-new-case-dialog">{{ t('button.newtestcase') }}</el-button>
    <el-button type="primary" @click="convert" test-id="convert">{{ t('button.export') }}</el-button>
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
            <el-input v-model="testCaseForm.name" test-id="case-form-name"/>
          </el-form-item>
          <el-form-item label="Method" prop="method" v-if="suite.spec.kind !== 'tRPC' && suite.spec.kind !== 'gRPC'">
            <el-select
              v-model="testCaseForm.method"
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
              v-model="testCaseForm.api"
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
            <el-button
              type="primary"
              @click="submitForm(testcaseFormRef)"
              :loading="suiteCreatingLoading"
              test-id="case-form-submit"
              >{{ t('button.submit') }}</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
