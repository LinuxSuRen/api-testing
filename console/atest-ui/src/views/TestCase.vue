<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Edit, Delete, Search } from '@element-plus/icons-vue'
import JsonViewer from 'vue-json-viewer'
import type { Pair, TestResult, TestCaseWithSuite } from './types'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from './types'
import { Cache } from './cache'
import { API } from './net'
import { UIAPI } from './net-vue'
import type { TestCaseResponse } from './cache'
import { useI18n } from 'vue-i18n'
import { JSONPath } from 'jsonpath-plus'
import { Codemirror } from 'vue-codemirror'

const { t } = useI18n()

const props = defineProps({
  name: String,
  suite: String,
  kindName: String,
})
const emit = defineEmits(['updated'])

let querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name!, props.suite!)
const testResultActiveTab = ref(Cache.GetPreference().responseActiveTab)
watch(testResultActiveTab, Cache.WatchResponseActiveTab)

const parameters = ref([] as Pair[])
const requestLoading = ref(false)
const testResult = ref({ header: [] as Pair[] } as TestResult)
const sendRequest = async () => {
  if (needUpdate.value) {
    await saveTestCase(false)
  }

  requestLoading.value = true
  const name = props.name
  const suite = props.suite

  API.RunTestCase({
    suiteName: suite,
    name: name,
      parameters: parameters.value
  }, (e) => {
    testResult.value = e
    requestLoading.value = false

    if (e.error !== '') {
      ElMessage({
        message: e.error,
        type: 'error'
      })
    } else {
      ElMessage({
        message: 'Pass!',
        type: 'success'
      })
    }
    if (e.body !== '') {
      testResult.value.bodyObject = JSON.parse(e.body)
      testResult.value.originBodyObject = JSON.parse(e.body)
    }

    Cache.SetTestCaseResponseCache(suite + '-' + name, {
      body: testResult.value.bodyObject,
      output: e.output,
      statusCode: testResult.value.statusCode
    } as TestCaseResponse)

    parameters.value = []
  }, (e) => {
    parameters.value = []

    requestLoading.value = false
    UIAPI.ErrorTip(e)
    testResult.value.bodyObject = JSON.parse(e.body)
    testResult.value.originBodyObject = JSON.parse(e.body)
  })
}

const responseBodyFilterText = ref('')
function responseBodyFilter() {
  if (responseBodyFilterText.value === '') {
    testResult.value.bodyObject = testResult.value.originBodyObject
  } else {
    const query = JSONPath({
      path: responseBodyFilterText.value,
      json: testResult.value.originBodyObject,
      resultType: 'value'
    })
    testResult.value.bodyObject = query[0]
  }
}

const parameterDialogOpened = ref(false)
function openParameterDialog() {
  API.GetTestSuite(props.suite, (e) => {
      parameters.value = e.param
      parameterDialogOpened.value = true
  }, UIAPI.ErrorTip)
}

function sendRequestWithParameter() {
  parameterDialogOpened.value = false
  sendRequest()
}

function generateCode() {
  const name = props.name
  const suite = props.suite

  API.GenerateCode({
    suiteName: suite,
    name: name,
    generator: currentCodeGenerator.value
  }, (e) => {
      ElMessage({
        message: 'Code generated!',
        type: 'success'
      })
      if (currentCodeGenerator.value === "gRPCPayload") {
        currentCodeContent.value = JSON.stringify(JSON.parse(e.message), null, 4)
      } else {
        currentCodeContent.value = e.message
      }
    }, UIAPI.ErrorTip)
}

function copyCode() {
  navigator.clipboard.writeText(currentCodeContent.value);
}

const queryBodyFields = (queryString: string, cb: any) => {
  if (!testResult.value.bodyObject || !FlattenObject(testResult.value.bodyObject)) {
    cb([])
    return
  }
  const keys = Object.getOwnPropertyNames(FlattenObject(testResult.value.bodyObject))
  if (keys.length <= 0) {
    cb([])
    return
  }

  const pairs = [] as Pair[]
  keys.forEach((e) => {
    pairs.push({
      key: e,
      value: e
    } as Pair)
  })

  const results = queryString ? pairs.filter(CreateFilter(queryString)) : pairs
  // call callback function to return suggestions
  cb(results)
}

const emptyTestCaseWithSuite: TestCaseWithSuite = {
  suiteName: '',
  data: {
    name: '',
    request: {
      api: '',
      method: '',
      header: [],
      query: [],
      form: [],
      body: ''
    },
    response: {
      statusCode: 0,
      body: '',
      header: [],
      bodyFieldsExpect: [],
      verify: [],
      schema: ''
    }
  }
}

const testCaseWithSuite = ref(emptyTestCaseWithSuite)

function load() {
  const name = props.name
  const suite = props.suite
  if (name === '' || suite === '') {
    return
  }

  // load cache
  const cache = Cache.GetTestCaseResponseCache(suite + '-' + name)
  if (cache.body) {
    testResult.value.bodyObject = cache.body
    testResult.value.output = cache.output
    testResult.value.statusCode = cache.statusCode
  } else {
    testResult.value.bodyObject = {}
    testResult.value.output = ''
    testResult.value.statusCode = 0
  }
  testResult.value.originBodyObject = testResult.value.bodyObject

  API.GetTestCase({
    suiteName: suite,
    name: name
  }, (e) => {
      if (e.request.method === '') {
        e.request.method = 'GET'
      }

      e.request.header.push({
        key: '',
        value: ''
      })
      e.request.query.push({
        key: '',
        value: ''
      })
      e.request.form.push({
        key: '',
        value: ''
      })
      e.response.header.push({
        key: '',
        value: ''
      })
      e.response.bodyFieldsExpect.push({
        key: '',
        value: ''
      })
      e.response.verify.push('')
      if (e.response.statusCode === 0) {
        e.response.statusCode = 200
      }

      e.request.header.forEach(item => {
        if (item.key === "Content-Type") {
          switch (item.value) {
            case 'application/x-www-form-urlencoded':
              bodyType.value = 4
              break
            case 'application/json':
              bodyType.value = 5
              break
          }
        }
      });

      testCaseWithSuite.value = {
        suiteName: suite,
        data: e
      } as TestCaseWithSuite
    })
}
load()
watch(props, () => {
  load()
})

const needUpdate = ref(false)
watch(testCaseWithSuite, (after, before) => {
  if (before.data.name !== '' && after.data.name === before.data.name) {
    needUpdate.value = true
  }
}, { deep: true })

const saveLoading = ref(false)
function saveTestCase(tip: boolean = true) {
  UIAPI.UpdateTestCase(testCaseWithSuite.value, (e) => {
    if (tip) {
      ElMessage({
        message: 'Saved.',
        type: 'success'
      })
    }
  }, UIAPI.ErrorTip, saveLoading)
}

function deleteTestCase() {
  const name = props.name
  const suite = props.suite

  API.DeleteTestCase({
    suiteName: suite,
    name: name
  }, (e) => {
    if (e.ok) {
      emit('updated', 'hello from child')

      ElMessage({
        message: 'Delete.',
        type: 'success'
      })

      // clean all the values
      testCaseWithSuite.value = emptyTestCaseWithSuite
    } else {
      UIAPI.ErrorTip(e)
    }
  })
}

const codeDialogOpened = ref(false)
const codeGenerators = ref('')
const currentCodeGenerator = ref('')
const currentCodeContent = ref('')
function openCodeDialog() {
  codeDialogOpened.value = true

  API.ListCodeGenerator((e) => {
    codeGenerators.value = e.data
  })

  if (currentCodeGenerator.value !== '') {
    generateCode()
  }
}
watch(currentCodeGenerator, () => {
  generateCode()
})

const options = GetHTTPMethods()
const requestActiveTab = ref(Cache.GetPreference().requestActiveTab)
watch(requestActiveTab, Cache.WatchRequestActiveTab)

function bodyFiledExpectChange() {
  const data = testCaseWithSuite.value.data.response.bodyFieldsExpect
  let lastItem = data[data.length - 1]
  if (lastItem.key !== '') {
    data.push({
      key: '',
      value: ''
    } as Pair)
  }
}

function queryChange() {
  const query = testCaseWithSuite.value.data.request.query
  let lastItem = query[query.length - 1]
  if (lastItem.key !== '') {
    testCaseWithSuite.value.data.request.query.push({
      key: '',
      value: ''
    } as Pair)
  }
}
function headerChange() {
  const header = testCaseWithSuite.value.data.request.header
  let lastItem = header[header.length - 1]
  if (lastItem.key !== '') {
    testCaseWithSuite.value.data.request.header.push({
      key: '',
      value: ''
    } as Pair)
  }
}

const headerValues = ref([] as Pair[])
const headerSelect = (item: Record<string, any>) => {
  headerValues.value = []
  pupularHeaderPairs.value.filter((v) => {
    if (v.key === item.value) {
      headerValues.value.push({
        key: v.value,
        value: v.value
      } as Pair)
    }
  })
}
function expectedHeaderChange() {
  const header = testCaseWithSuite.value.data.response.header
  let lastItem = header[header.length - 1]
  if (lastItem.key !== '') {
    testCaseWithSuite.value.data.response.header.push({
      key: '',
      value: ''
    } as Pair)
  }
}
function formChange() {
  const form = testCaseWithSuite.value.data.request.form
  let lastItem = form[form.length - 1]
  if (lastItem.key !== '') {
    testCaseWithSuite.value.data.request.form.push({
      key: '',
      value: ''
    } as Pair)
  }
}

const bodyType = ref(5)
function bodyTypeChange(e: number) {
  let contentType = ""
  switch (e) {
    case 4:
      contentType = 'application/x-www-form-urlencoded'
      break;
    case 5:
      contentType = 'application/json'
      break;
  }

  if (contentType !== "") {
    testCaseWithSuite.value.data.request.header = insertOrUpdateIntoMap({
        key: 'Content-Type',
        value: contentType
      } as Pair, testCaseWithSuite.value.data.request.header)
  }
}

function jsonForamt() {
  if (bodyType.value !== 5) {
    return
  }

  try {
    testCaseWithSuite.value.data.request.body = JSON.stringify(JSON.parse(testCaseWithSuite.value.data.request.body), null, 4)
  } catch (e) {
    console.log(e)
  }
}

function insertOrUpdateIntoMap(pair: Pair, pairs: Pair[]) {
  const index = pairs.findIndex((e) => e.key === pair.key)
  if (index === -1) {
    const oldPairs = pairs
    pairs = [pair]
    pairs = pairs.concat(oldPairs)
  } else {
    pairs[index] = pair
  }
  return pairs
}

const pupularHeaders = ref([] as Pair[])
const pupularHeaderPairs = ref([] as Pair[])
API.PopularHeaders((e) => {
  const headerCache = new Map<string, string>();
  for (var i = 0; i < e.data.length; i++) {
    const pair = {
      key: e.data[i].key,
      value: e.data[i].value
    } as Pair

    pupularHeaderPairs.value.push(pair)

    if (!headerCache.get(pair.key)) {
      headerCache.set(pair.key, "index")

      pupularHeaders.value.push({
        key: e.data[i].key,
        value: e.data[i].value
      } as Pair)
    }
  }
})

const queryPupularHeaders = (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
    ? pupularHeaders.value.filter(CreateFilter(queryString))
    : pupularHeaders.value

  results.forEach((e) => {
    e.value = e.key
  })
  cb(results)
}
const queryHeaderValues = (queryString: string, cb: (arg: any) => void) => {
  const results = queryString
    ? headerValues.value.filter(CreateFilter(queryString))
    : headerValues.value

  results.forEach((e) => {
    e.value = e.key
  })
  cb(results)
}
</script>

<template>
  <el-container>
    <el-header style="padding-left: 5px;">
      <div style="margin-bottom: 5px">
        <el-button type="primary" @click="saveTestCase" :icon="Edit" :loading="saveLoading"
          disabled v-if="Cache.GetCurrentStore().readOnly"
          >{{ t('button.save') }}</el-button>
        <el-button type="primary" @click="saveTestCase" :icon="Edit" :loading="saveLoading"
          v-if="!Cache.GetCurrentStore().readOnly"
          >{{ t('button.save') }}</el-button>
        <el-button type="primary" @click="deleteTestCase" :icon="Delete">{{ t('button.delete') }}</el-button>
        <el-button type="primary" @click="openCodeDialog">{{ t('button.generateCode') }}</el-button>
      </div>
      <div style="display: flex;">
        <el-select
          v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'"
          v-model="testCaseWithSuite.data.request.method"
          class="m-2"
          placeholder="Method"
          size="middle"
          test-id="case-editor-method"
        >
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.key"
            :value="item.value"
          />
        </el-select>
        <el-autocomplete
          v-model="testCaseWithSuite.data.request.api"
          :fetch-suggestions="querySuggestedAPIs"
          placeholder="API Address"
          style="width: 50%; margin-left: 5px; margin-right: 5px; flex-grow: 1;"
        >
          <template #default="{ item }">
            <div class="value">{{ item.request.method }}</div>
            <span class="link">{{ item.request.api }}</span>
          </template>
        </el-autocomplete>

        <el-dropdown split-button type="primary" @click="sendRequest" :loading="requestLoading">
          {{ t('button.send') }}
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="openParameterDialog">{{ t('button.sendWithParam') }}</el-dropdown-item>
              </el-dropdown-menu>
            </template>
        </el-dropdown>
      </div>
    </el-header>

    <el-main style="padding-left: 5px;">
      <el-tabs v-model="requestActiveTab">
        <el-tab-pane label="Query" name="query" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-table :data="testCaseWithSuite.data.request.query" style="width: 100%">
            <el-table-column label="Key" width="180">
              <template #default="scope">
                <el-autocomplete
                  v-model="scope.row.key"
                  placeholder="Key"
                  @change="queryChange"
                />
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
        </el-tab-pane>

        <el-tab-pane label="Headers" name="header">
          <el-table :data="testCaseWithSuite.data.request.header" style="width: 100%">
            <el-table-column label="Key" width="180">
              <template #default="scope">
                <el-autocomplete
                  v-model="scope.row.key"
                  :fetch-suggestions="queryPupularHeaders"
                  placeholder="Key"
                  @change="headerChange"
                  @select="headerSelect"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-autocomplete
                    v-model="scope.row.value"
                    :fetch-suggestions="queryHeaderValues"
                    style="width: 100%;"
                  />
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="Body" name="body">
          <el-radio-group v-model="bodyType" @change="bodyTypeChange">
            <el-radio :label="1">none</el-radio>
            <el-radio :label="2">form-data</el-radio>
            <el-radio :label="3">raw</el-radio>
            <el-radio :label="4">x-www-form-urlencoded</el-radio>
            <el-radio :label="5">JSON</el-radio>
          </el-radio-group>

          <div style="flex-grow: 1;">
            <Codemirror v-if="bodyType === 3 || bodyType === 5"
              @change="jsonForamt"
              v-model="testCaseWithSuite.data.request.body"/>
            <el-table :data="testCaseWithSuite.data.request.form" style="width: 100%" v-if="bodyType === 4">
              <el-table-column label="Key" width="180">
                <template #default="scope">
                  <el-input v-model="scope.row.key" placeholder="Key" @change="formChange"/>
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
          </div>
        </el-tab-pane>

        <el-tab-pane label="Expected" name="expected" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-row :gutter="20">
            <span
              class="ml-3 w-50 text-gray-600 inline-flex items-center"
              style="margin-left: 15px; margin-right: 15px"
              >Status Code:</span
            >
            <el-input
              v-model="testCaseWithSuite.data.response.statusCode"
              class="w-50 m-2"
              placeholder="Please input"
              style="width: 200px"
            />
          </el-row>
          <el-input
            v-model="testCaseWithSuite.data.response.body"
            :autosize="{ minRows: 4, maxRows: 8 }"
            type="textarea"
            placeholder="Expected Body"
          />
        </el-tab-pane>

        <el-tab-pane label="Expected Headers" name="expected-headers" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-table :data="testCaseWithSuite.data.response.header" style="width: 100%">
            <el-table-column label="Key" width="180">
              <template #default="scope">
                <el-input
                  v-model="scope.row.key"
                  placeholder="Key"
                  @change="expectedHeaderChange"
                />
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
        </el-tab-pane>

        <el-tab-pane label="BodyFiledExpect" name="bodyFieldExpect" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-table :data="testCaseWithSuite.data.response.bodyFieldsExpect" style="width: 100%">
            <el-table-column label="Key" width="180">
              <template #default="scope">
                <el-autocomplete
                  v-model="scope.row.key"
                  :fetch-suggestions="queryBodyFields"
                  clearable
                  placeholder="Key"
                  @change="bodyFiledExpectChange"
                />
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
        </el-tab-pane>

        <el-tab-pane label="Verify" name="verify" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <div v-for="verify in testCaseWithSuite.data.response.verify" :key="verify">
            <el-input :value="verify" />
          </div>
        </el-tab-pane>

        <el-tab-pane label="Schema" name="schema" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-input
            v-model="testCaseWithSuite.data.response.schema"
            :autosize="{ minRows: 4, maxRows: 20 }"
            type="textarea"
          />
        </el-tab-pane>
      </el-tabs>

      <el-drawer v-model="codeDialogOpened" size="50%">
        <template #header>
          <h4>Code Generator</h4>
        </template>
        <template #default>
          <div style="padding-bottom: 10px;">
            <el-select
              v-model="currentCodeGenerator"
              class="m-2"
              style="padding-right: 10px;"
              size="middle"
            >
              <el-option
                v-for="item in codeGenerators"
                :key="item.key"
                :label="item.key"
                :value="item.key"
              />
            </el-select>
            <el-button type="primary" @click="generateCode">{{ t('button.refresh') }}</el-button>
            <el-button type="primary" @click="copyCode">{{ t('button.copy') }}</el-button>
          </div>
          <Codemirror v-model="currentCodeContent"/>
        </template>
      </el-drawer>

      <el-drawer v-model="parameterDialogOpened">
        <template #header>
          <h4>API Request Parameters</h4>
        </template>
        <template #default>
          <el-table :data="parameters" style="width: 100%">
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

          <el-button type="primary" @click="sendRequestWithParameter">{{ t('button.send') }}</el-button>
        </template>
      </el-drawer>
    </el-main>

    <el-footer style="height: auto;">
      <el-tabs v-model="testResultActiveTab">
        <el-tab-pane label="Output" name="output">
          <el-tag class="ml-2" type="success" v-if="testResult.statusCode && testResult.error === ''">{{ t('httpCode.' + testResult.statusCode) }}</el-tag>
          <el-tag class="ml-2" type="danger" v-if="testResult.statusCode && testResult.error !== ''">{{ t('httpCode.' + testResult.statusCode) }}</el-tag>

          <Codemirror v-model="testResult.output"/>
        </el-tab-pane>
        <el-tab-pane label="Body" name="body">
          <el-input :prefix-icon="Search" @change="responseBodyFilter" v-model="responseBodyFilterText"
            clearable label="dddd" placeholder="$.key" />
          <JsonViewer :value="testResult.bodyObject" :expand-depth="5" copyable boxed sort />
        </el-tab-pane>
        <el-tab-pane name="response-header">
          <template #label>
            <el-badge :value="testResult.header.length" class="item">Header</el-badge>
          </template>
          <el-table :data="testResult.header" style="width: 100%">
            <el-table-column label="Key" width="200">
              <template #default="scope">
                <el-input v-model="scope.row.key" placeholder="Key" readonly="true" />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="Value" readonly="true" />
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-footer>
  </el-container>
</template>
