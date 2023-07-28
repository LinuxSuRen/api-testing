<script setup lang="ts">
import { ref, watch } from 'vue'
import type { TabsPaneContext } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import JsonViewer from 'vue-json-viewer'
import type { Pair, TestResult, TestCaseWithSuite } from './types'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from './types'

const props = defineProps({
  name: String,
  suite: String,
  store: String
})
const emit = defineEmits(['updated'])

let querySuggestedAPIs = NewSuggestedAPIsQuery(props.suite!)
const testResultActiveTab = ref('output')
const requestLoading = ref(false)
const testResult = ref({ header: [] as Pair[] } as TestResult)
function sendRequest() {
  requestLoading.value = true
  const name = props.name
  const suite = props.suite
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify({
      suite: suite,
      testcase: name
    })
  }
  fetch('/server.Runner/RunTestCase', requestOptions)
    .then((response) => response.json())
    .then((e) => {
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
      }
    })
    .catch((e) => {
      requestLoading.value = false
      ElMessage.error('Oops, ' + e)
      testResult.value.bodyObject = JSON.parse(e.body)
    })
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

  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify({
      suite: suite,
      testcase: name
    })
  }
  fetch('/server.Runner/GetTestCase', requestOptions)
    .then((response) => response.json())
    .then((e) => {
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

const saveLoading = ref(false)
function saveTestCase() {
  saveLoading.value = true

  // remove empty pair
  testCaseWithSuite.value.data.request.header = testCaseWithSuite.value.data.request.header.filter(
    (e) => e.key !== ''
  )
  testCaseWithSuite.value.data.request.query = testCaseWithSuite.value.data.request.query.filter(
    (e) => e.key !== ''
  )
  testCaseWithSuite.value.data.request.form = testCaseWithSuite.value.data.request.form.filter(
    (e) => e.key !== ''
  )
  testCaseWithSuite.value.data.response.header =
    testCaseWithSuite.value.data.response.header.filter((e) => e.key !== '')
  testCaseWithSuite.value.data.response.bodyFieldsExpect =
    testCaseWithSuite.value.data.response.bodyFieldsExpect.filter((e) => e.key !== '')
  testCaseWithSuite.value.data.response.verify =
    testCaseWithSuite.value.data.response.verify.filter((e) => e !== '')

  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify(testCaseWithSuite.value)
  }
  fetch('/server.Runner/UpdateTestCase', requestOptions).then((e) => {
    if (e.ok) {
      ElMessage({
        message: 'Saved.',
        type: 'success'
      })
    } else {
      ElMessage.error('Oops, ' + e.statusText)
    }
    saveLoading.value = false
  })
}

function deleteTestCase() {
  const name = props.name
  const suite = props.suite
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': props.store
    },
    body: JSON.stringify({
      suite: suite,
      testcase: name
    })
  }
  fetch('/server.Runner/DeleteTestCase', requestOptions).then((e) => {
    if (e.ok) {
      emit('updated', 'hello from child')

      ElMessage({
        message: 'Delete.',
        type: 'success'
      })

      // clean all the values
      testCaseWithSuite.value = emptyTestCaseWithSuite
    } else {
      ElMessage.error('Oops, ' + e.statusText)
    }
  })
}

const options = GetHTTPMethods()

const activeName = ref('second')

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
}

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

const bodyType = ref(1)
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
const requestOptions = {
  method: 'POST',
  headers: {
    'X-Store-Name': props.store
  },
}
fetch('/server.Runner/PopularHeaders', requestOptions)
  .then((response) => response.json())
  .then((e) => {
    pupularHeaders.value = e.data
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
</script>

<template>
  <div class="common-layout">
    <el-container style="height: 60vh">
      <el-header style="padding-left: 5px">
        <div style="margin-bottom: 5px">
          <el-button type="primary" @click="saveTestCase" :icon="Edit" :loading="saveLoading"
            >Save</el-button
          >
          <el-button type="primary" @click="deleteTestCase" :icon="Delete">Delete</el-button>
        </div>
        <el-select
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
          style="width: 70%; margin-left: 5px; margin-right: 5px"
        >
          <template #default="{ item }">
            <div class="value">{{ item.request.method }}</div>
            <span class="link">{{ item.request.api }}</span>
          </template>
        </el-autocomplete>

        <el-button type="primary" @click="sendRequest" :loading="requestLoading">Send</el-button>
      </el-header>

      <el-main>
        <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
          <el-tab-pane label="Headers" name="second">
            <el-table :data="testCaseWithSuite.data.request.header" style="width: 100%">
              <el-table-column label="Key" width="180">
                <template #default="scope">
                  <el-autocomplete
                    v-model="scope.row.key"
                    :fetch-suggestions="queryPupularHeaders"
                    placeholder="Key"
                    @change="headerChange"
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

          <el-tab-pane label="Body" name="third">
            <el-radio-group v-model="bodyType" @change="bodyTypeChange">
              <el-radio :label="1">none</el-radio>
              <el-radio :label="2">form-data</el-radio>
              <el-radio :label="3">raw</el-radio>
              <el-radio :label="4">x-www-form-urlencoded</el-radio>
              <el-radio :label="5">JSON</el-radio>
            </el-radio-group>

            <el-input
              v-if="bodyType === 3 || bodyType === 5"
              v-model="testCaseWithSuite.data.request.body"
              :autosize="{ minRows: 4, maxRows: 8 }"
              type="textarea"
              placeholder="Please input"
              @change="jsonForamt"
            />
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
          </el-tab-pane>

          <el-tab-pane label="Expected" name="expected">
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

          <el-tab-pane label="Expected Headers" name="expected-headers">
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

          <el-tab-pane label="BodyFiledExpect" name="fourth">
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

          <el-tab-pane label="Verify" name="fifth">
            <div v-for="verify in testCaseWithSuite.data.response.verify" :key="verify">
              <el-input :value="verify" />
            </div>
          </el-tab-pane>

          <el-tab-pane label="Schema" name="schema">
            <el-input
              v-model="testCaseWithSuite.data.response.schema"
              :autosize="{ minRows: 4, maxRows: 8 }"
              type="textarea"
            />
          </el-tab-pane>
        </el-tabs>
      </el-main>

      <el-footer>
        <el-tabs v-model="testResultActiveTab" class="demo-tabs" @tab-click="handleClick">
          <el-tab-pane label="Output" name="output">
            <el-input
              v-model="testResult.output"
              :autosize="{ minRows: 4, maxRows: 6 }"
              readonly="true"
              type="textarea"
              placeholder="Please input"
            />
          </el-tab-pane>
          <el-tab-pane label="Body" name="body">
            <JsonViewer :value="testResult.bodyObject" :expand-depth="5" copyable boxed sort />
          </el-tab-pane>
          <el-tab-pane name="response-header">
            <template #label>
              <el-badge :value="testResult.header.length" class="item"> Header </el-badge>
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
  </div>
</template>
