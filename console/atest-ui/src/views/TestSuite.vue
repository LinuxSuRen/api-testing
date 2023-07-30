<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Suite, TestCase, Pair } from './types'
import { NewSuggestedAPIsQuery } from './types'

const props = defineProps({
  name: String,
  store: String,
})
const emit = defineEmits(['updated'])
let querySuggestedAPIs = NewSuggestedAPIsQuery(props.name!)

const suite = ref({
  name: '',
  api: '',
  param: [] as Pair[],
  spec: {
    kind: '',
    url: ''
  }
} as Suite)
function load() {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify({
      name: props.name
    })
  }
  fetch('/server.Runner/GetTestSuite', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      suite.value = e
      if (suite.value.param.length === 0) {
        suite.value.param.push({
          key: '',
          value: ''
        } as Pair)
      }
    })
    .catch((e) => {
      ElMessage.error('Oops, ' + e)
    })
}
load()
watch(props, () => {
  load()
})

function save() {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify(suite.value)
  }
  fetch('/server.Runner/UpdateTestSuite', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      ElMessage({
        message: 'Updated.',
        type: 'success'
      })
    })
    .catch((e) => {
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
  querySuggestedAPIs = NewSuggestedAPIsQuery(props.name!)
}

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      suiteCreatingLoading.value = true

      const requestOptions = {
        method: 'POST',
        headers: {
          'X-Store-Name': props.store
        },
        body: JSON.stringify({
          suiteName: props.name,
          data: {
            name: testCaseForm.name,
            request: {
              api: testCaseForm.api,
              method: testCaseForm.method
            }
          }
        })
      }

      fetch('/server.Runner/CreateTestCase', requestOptions)
        .then((response) => response.json())
        .then(() => {
          suiteCreatingLoading.value = false
          emit('updated', 'hello from child')
        })

      dialogVisible.value = false
    }
  })
}

function del() {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify({
      name: props.name
    })
  }
  fetch('/server.Runner/DeleteTestSuite', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      ElMessage({
        message: 'Deleted.',
        type: 'success'
      })
      emit('updated')
    })
    .catch((e) => {
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
    <el-input class="mx-1" v-model="suite.api" placeholder="API" test-id="suite-editor-api"></el-input>
    <el-select v-model="suite.spec.kind" class="m-2" placeholder="API Spec Kind" size="middle">
      <el-option
        v-for="item in apiSpecKinds"
        :key="item.value"
        :label="item.label"
        :value="item.value"
      />
    </el-select>
    <el-input class="mx-1" v-model="suite.spec.url" placeholder="API Spec URL"></el-input>

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

    <el-button type="primary" @click="save">Save</el-button>
    <el-button type="primary" @click="del" test-id="suite-del-but">Delete</el-button>

    <el-button type="primary" @click="openNewTestCaseDialog" :icon="Edit" test-id="open-new-case-dialog">New TestCase</el-button>
  </div>

  <el-dialog v-model="dialogVisible" title="Create Test Case" width="40%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="testCaseForm"
          ref="testcaseFormRef"
          status-icon
          label-width="60px"
        >
          <el-form-item label="Name" prop="name">
            <el-input v-model="testCaseForm.name" test-id="case-form-name"/>
          </el-form-item>
          <el-form-item label="Method" prop="method">
            <el-input v-model="testCaseForm.method" test-id="case-form-method" />
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
              >Submit</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
