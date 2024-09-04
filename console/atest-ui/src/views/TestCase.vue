<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Edit, Delete, Search, CopyDocument } from '@element-plus/icons-vue'
import JsonViewer from 'vue-json-viewer'
import type { Pair, TestResult, TestCaseWithSuite } from './types'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from './types'
import { Cache } from './cache'
import { API } from './net'
import { UIAPI } from './net-vue'
import type { TestCaseResponse } from './cache'
import { Magic } from './magicKeys'
import { useI18n } from 'vue-i18n'
import { JSONPath } from 'jsonpath-plus'
import { Codemirror } from 'vue-codemirror'
import jsonlint from 'jsonlint-mod'

import CodeMirror from 'codemirror'
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/merge/merge.js'
import 'codemirror/addon/merge/merge.css'

import DiffMatchPatch from 'diff-match-patch';


window.diff_match_patch = DiffMatchPatch;
window.DIFF_DELETE = -1;
window.DIFF_INSERT = 1;
window.DIFF_EQUAL = 0;

const { t } = useI18n()

const props = defineProps({
  name: String,
  suite: String,
  kindName: String,
  historySuiteName: String,
  historyCaseID: String
})
const emit = defineEmits(['updated','toHistoryPanel'])

let querySuggestedAPIs = NewSuggestedAPIsQuery(Cache.GetCurrentStore().name!, props.suite!)
const testResultActiveTab = ref(Cache.GetPreference().responseActiveTab)
watch(testResultActiveTab, Cache.WithResponseActiveTab)
Magic.Keys(() => {
  testResultActiveTab.value = 'output'
}, ['Alt+KeyO'])

const parameters = ref([] as Pair[])
const requestLoading = ref(false)
const testResult = ref({ header: [] as Pair[] } as TestResult)
const sendRequest = async () => {
  if (needUpdate.value) {
    await saveTestCase(false, runTestCase)
    needUpdate.value = false
  } else {
    runTestCase()
  }
}
Magic.Keys(sendRequest, ['Alt+S', 'Alt+ß'])

const runTestCase = () => {
  requestLoading.value = true
  const name = props.name
  const suite = props.suite
  API.RunTestCase({
    suiteName: suite,
    name: name,
    parameters: parameters.value
  }, (e) => {
    handleTestResult(e)
    requestLoading.value = false
  }, (e) => {
    parameters.value = []

    requestLoading.value = false
    UIAPI.ErrorTip(e)
    parseResponseBody(e.body)
  })
}

const parseResponseBody = (body) => {
  if (body === '') {
    return
  }

  try {
    testResult.value.bodyObject = JSON.parse(body)
    testResult.value.originBodyObject = JSON.parse(body)
  } catch {
    testResult.value.bodyText = body
  }
}

const handleTestResult = (e) => {
  testResult.value = e;

  if (!isHistoryTestCase.value) {
    handleTestResultError(e)
  }

  if (e.body !== '') {
    testResult.value.bodyObject = JSON.parse(e.body);
    testResult.value.originBodyObject = JSON.parse(e.body);
  }

    Cache.SetTestCaseResponseCache(suite + '-' + name, {
      body: testResult.value.bodyObject,
      output: e.output,
      statusCode: testResult.value.statusCode
    } as TestCaseResponse)

  parameters.value = [];
}

const handleTestResultError = (e) => {
  if (e.error !== '') {
    ElMessage({
      message: e.error,
      type: 'error'
    });
  } else {
    ElMessage({
      message: 'Pass!',
      type: 'success'
    });
  }
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
      cookie: [],
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

let name
let suite
let historySuiteName
let historyCaseID
const isHistoryTestCase = ref(false)
const HistoryTestCaseCreateTime = ref('')

function load() {
   name = props.name
   suite = props.suite
   historySuiteName = props.historySuiteName
   historyCaseID = props.historyCaseID
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

  if (historySuiteName != '' && historySuiteName != undefined) {
    isHistoryTestCase.value = true
    API.GetHistoryTestCaseWithResult({
      historyCaseID : historyCaseID
    }, (e) => {
      setDefaultValues(e.data)
      determineBodyType(e.data)
      setTestCaseWithSuite(e.data,suite)
      handleTestResult(e.testCaseResult[0])
      HistoryTestCaseCreateTime.value = formatDate(e.createTime)
    })
  } else {
    API.GetTestCase({
      suiteName: suite,
      name: name
    }, (e) => {
      setDefaultValues(e)
      determineBodyType(e)
      setTestCaseWithSuite(e,suite)
    })
  }
}

function formatDate(createTimeStr : string){
  let parts = createTimeStr.split(/[T.Z]/);
    let datePart = parts[0].split("-");
    let timePart = parts[1].split(":");

    let year = datePart[0];
    let month = datePart[1];
    let day = datePart[2];
    let hours = timePart[0];
    let minutes = timePart[1];
    let seconds = timePart[2].split(".")[0];

    let formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    return formattedDate
}

function determineBodyType(e) {
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
}

function setDefaultValues(e) {
  if (e.request.method === '') {
    e.request.method = 'GET'
  }

  e.request.header.push({
    key: '',
    value: ''
  })
  e.request.cookie.push({
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
}

function setTestCaseWithSuite(e, suite) {
  e.suiteName = suite
  testCaseWithSuite.value = {
    suiteName: suite,
    data: e
  } as TestCaseWithSuite;
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
function saveTestCase(tip: boolean = true, callback: (c: any) => void) {
  UIAPI.UpdateTestCase(testCaseWithSuite.value, (e) => {
    if (tip) {
      ElMessage({
        message: 'Saved.',
        type: 'success'
      })
    }

    if (callback) {
      callback()
    }
  }, UIAPI.ErrorTip, saveLoading)
}

function deleteCase() {
  const name = props.name
  const suite = props.suite
  const historyCaseID = props.historyCaseID

  if (isHistoryTestCase.value == true){
    deleteHistoryTestCase(historyCaseID)
  } else {
    deleteTestCase(name, suite)
  }
}

function deleteHistoryTestCase(historyCaseID : string){
  API.DeleteHistoryTestCase({ historyCaseID }, handleDeleteResponse);
}

function deleteTestCase(name : string, suite : string){
  API.DeleteTestCase({ suiteName: suite, name }, handleDeleteResponse);
}

function handleDeleteResponse(e) {
  if (e.ok) {
    emit('updated', 'hello from child');

    ElMessage({
      message: 'Delete.',
      type: 'success'
    });

    // Clean all the values
    testCaseWithSuite.value = emptyTestCaseWithSuite;
  } else {
    UIAPI.ErrorTip(e);
  }
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

const historyDialogOpened = ref(false)
const historyForm = ref({ selectedID: '' })
const historyRecords = ref([]); 
const selectedHistory = ref(null);
const viewHistoryRef = ref(null);
const formatHistoryCase = ref(null);

const rules = {
  selectedID: [
    { required: true, message: 'Please select history TestCase', trigger: 'change' }
  ]
}

function openHistoryDialog(){
  historyDialogOpened.value = true
  name = props.name
  suite = props.suite
  API.GetTestCaseAllHistory({
    suiteName : suite,
    name: name,
  }, (e) => {
    historyRecords.value = e.data
    historyRecords.value.forEach(record => {
      record.createTime = formatDate(record.createTime)
      setDefaultValues(record)
    });
  })
}

function handleHistoryChange(value) {
  selectedHistory.value = historyRecords.value.find(record => record.ID === value);
  const {
  caseName: name,
  suiteName,
  request,
  response
  } = selectedHistory.value;
  formatHistoryCase.value = {
    name,
    suiteName,
    request,
    response
  };
  initCompare(testCaseWithSuite.value.data, formatHistoryCase.value)
}

function initCompare(value, historyValue) {
  const formattedHistoryValue = JSON.stringify(historyValue, null, 2);
  const formattedNewValue = JSON.stringify(value, null, 2);
  const target = document.getElementById('compareView')
  target.innerHTML = ''

  const mergeView = CodeMirror.MergeView(target, {
    value: formattedHistoryValue, 
    origLeft: null,
    orig: formattedNewValue, 
    lineNumbers: true, 
    mode: { name: "javascript", json: true },
    highlightDifferences: true,
    foldGutter:true,
    lineWrapping:true,
    styleActiveLine: true,
    matchBrackets: true, 
    connect: 'align',
    readOnly: true 
  })
}

const caseRevertLoading = ref(false)
const submitForm = async (formEl) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      caseRevertLoading.value = true
      const historyTestCase =  {
        suiteName: props.suite,
        data: formatHistoryCase.value
      } 
      UIAPI.UpdateTestCase(historyTestCase, (e) => {
        if(e.error == ""){
          ElMessage({
            message: 'Saved.',
            type: 'success'
          })
          load()
        }
      }, UIAPI.ErrorTip, saveLoading)
      caseRevertLoading.value = false
      historyDialogOpened.value = false
      historyForm.value.selectedID = ''
      const target = document.getElementById('compareView');
      target.innerHTML = '';
    }
  })
}

const goToHistory = async (formEl) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      caseRevertLoading.value = true
      emit('toHistoryPanel', { ID: selectedHistory.value.ID, panelName: 'history' });
      caseRevertLoading.value = false
      historyDialogOpened.value = false
      historyForm.value.selectedID = ''
      const target = document.getElementById('compareView');
      target.innerHTML = '';
    }
  })
}

const deleteAllHistory = async (formEl) => {
  if (!formEl) return
  caseRevertLoading.value = true
  API.DeleteAllHistoryTestCase(props.suite, props.name, handleDeleteResponse);
  caseRevertLoading.value = false
  historyDialogOpened.value = false
  historyForm.value.selectedID = ''
  const target = document.getElementById('compareView');
  target.innerHTML = '';
}

const options = GetHTTPMethods()
const requestActiveTab = ref(Cache.GetPreference().requestActiveTab)
watch(requestActiveTab, Cache.WatchRequestActiveTab)
Magic.Keys(() => {
  requestActiveTab.value = 'query'
}, ['Alt+KeyQ'])
Magic.Keys(() => {
  requestActiveTab.value = 'header'
}, ['Alt+KeyH'])
Magic.Keys(() => {
  requestActiveTab.value = 'body'
}, ['Alt+KeyB'])

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
function cookieChange(){
  const cookie = testCaseWithSuite.value.data.request.cookie
  let lastItem = cookie[cookie.length - 1]
  if (lastItem.key !== '') {
    testCaseWithSuite.value.data.request.cookie.push({
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

const lintingError = ref('')
function jsonFormat(space) {
  const jsonText = testCaseWithSuite.value.data.request.body
  if (bodyType.value !== 5 || jsonText === '') {
    return
  }

  try {
    const jsonObj = jsonlint.parse(jsonText)
    if (space >= 0) {
      testCaseWithSuite.value.data.request.body = JSON.stringify(jsonObj, null, space)
    }
    lintingError.value = ''
  } catch (e) {
    lintingError.value = e.message
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

const duplicateTestCaseDialog = ref(false)
const targetTestCaseName = ref('')
const openDuplicateTestCaseDialog = () => {
    duplicateTestCaseDialog.value = true
    targetTestCaseName.value = props.name + '-copy'
}
Magic.Keys(openDuplicateTestCaseDialog, ['Alt+KeyD'])
const duplicateTestCase = () => {
    API.DuplicateTestCase(props.suite, props.suite, props.name, targetTestCaseName.value,(d) => {
        duplicateTestCaseDialog.value = false
        ElMessage({
            message: 'Duplicated.',
            type: 'success'
        })
        emit('updated')
    })
}
Magic.Keys(() => {
  if (duplicateTestCaseDialog.value) {
    duplicateTestCase()
  }
}, ['Alt+KeyO'])
</script>

<template>
  <el-container>
    <el-header style="padding-left: 5px;">
      <div style="margin-bottom: 5px">
        <el-button type="primary" @click="saveTestCase" :icon="Edit" v-loading="saveLoading"
          disabled v-if="Cache.GetCurrentStore().readOnly || isHistoryTestCase"
          >{{ t('button.save') }}</el-button>
        <el-button type="primary" @click="saveTestCase" :icon="Edit" v-loading="saveLoading"
          v-if="!Cache.GetCurrentStore().readOnly && !isHistoryTestCase"
          >{{ t('button.save') }}</el-button>
        <el-button type="danger" @click="deleteCase" :icon="Delete">{{ t('button.delete') }}</el-button>
        <el-button type="primary" @click="openDuplicateTestCaseDialog" :icon="CopyDocument" v-if="!isHistoryTestCase">{{ t('button.duplicate') }}</el-button>
        <el-button type="primary" @click="openCodeDialog">{{ t('button.generateCode') }}</el-button>
        <el-button type="primary" v-if="!isHistoryTestCase" @click="openHistoryDialog">{{ t('button.viewHistory') }}</el-button>
        <span v-if="isHistoryTestCase" style="margin-left: 15px;">{{ t('tip.runningAt') }}{{ HistoryTestCaseCreateTime }}</span>
      </div>
      <div style="display: flex;">
        <el-select
          v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'"
          v-model="testCaseWithSuite.data.request.method"
          class="m-2"
          placeholder="Method"
          size="default"
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

        <el-dropdown split-button type="primary" @click="sendRequest" v-loading="requestLoading">
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
        <el-tab-pane name="query" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <template #label>
            <el-badge :value="testCaseWithSuite.data.request.query.length - 1"
              :hidden="testCaseWithSuite.data.request.query.length <=1 " class="item">Query</el-badge>
          </template>
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

        <el-tab-pane name="header">
          <template #label>
            <el-badge :value="testCaseWithSuite.data.request.header.length - 1"
              :hidden="testCaseWithSuite.data.request.header.length <= 1" class="item">Header</el-badge>
          </template>
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

        <el-tab-pane name="cookie">
          <template #label>
            <el-badge :value="testCaseWithSuite.data.request.cookie.length - 1"
              :hidden="testCaseWithSuite.data.request.cookie.length <= 1" class="item">Cookie</el-badge>
          </template>
          <el-table :data="testCaseWithSuite.data.request.cookie" style="width: 100%">
            <el-table-column label="Key">
              <template #default="scope">
                <el-input v-model="scope.row.key" placeholder="Key"
                @change="cookieChange"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="value" />
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane name="body">
          <span style="margin-right: 10px; padding-right: 5px;">
            <el-button type="primary" @click="jsonFormat(4)">Beautify</el-button>
            <el-button type="primary" @click="jsonFormat(0)">Minify</el-button>
            <el-text class="mx-1">Choose the body format</el-text>
          </span>
          <template #label>
            <el-badge :is-dot="testCaseWithSuite.data.request.body !== ''" class="item">Body</el-badge>
          </template>
          <el-radio-group v-model="bodyType" @change="bodyTypeChange">
            <el-radio :label="1">none</el-radio>
            <el-radio :label="2">form-data</el-radio>
            <el-radio :label="3">raw</el-radio>
            <el-radio :label="4">x-www-form-urlencoded</el-radio>
            <el-radio :label="5">JSON</el-radio>
          </el-radio-group>

          <div style="flex-grow: 1;">
            <Codemirror v-if="bodyType === 3 || bodyType === 5"
              @blur="jsonFormat(-1)"
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
          <div v-if="lintingError" style="color: red; margin-top: 10px;">
            {{ lintingError }}
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
          <h4>{{ t('title.codeGenerator') }}</h4>
        </template>
        <template #default>
          <div style="padding-bottom: 10px;">
            <el-select
              v-model="currentCodeGenerator"
              class="m-2"
              style="padding-right: 10px;"
              size="default"
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

      <el-dialog v-model="historyDialogOpened" :title="t('button.viewHistory')" width="50%" draggable>
        <el-form
          ref="viewHistoryRef"
          :model="historyForm"
          status-icon label-width="120px"
          :rules="rules"
        >
          <el-form-item :label="t('title.history')" prop="selectedID">
            <el-row :gutter="20">
              <el-col :span="20">
                <el-select  class="m-2"
                  filterable
                  clearable
                  v-model="historyForm.selectedID"
                  default-first-option
                  placeholder="History Case"
                  size="middle"
                  @change="handleHistoryChange"
                >
                  <el-option
                    v-for="item in historyRecords"
                    :key="item.ID"
                    :label="item.createTime"
                    :value="item.ID"
                  />
                </el-select>
              </el-col>
              <el-col :span="4">
              <div style="display: flex">
                <el-button
                  type="primary"
                  @click="submitForm(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.revert') }}
                </el-button>
                <el-button
                  type="primary"
                  @click="goToHistory(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.goToHistory') }}
                </el-button>
                <el-button
                  type="primary"
                  @click="deleteAllHistory(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.deleteAllHistory') }}
                </el-button>
              </div>
              </el-col>
            </el-row>
          </el-form-item>
        </el-form>
        <div id="compareView"></div>
      </el-dialog>

      <el-drawer v-model="parameterDialogOpened">
        <template #header>
          <h4>{{ t('title.apiRequestParameter') }}</h4>
        </template>
        <template #default>
          <el-table :data="parameters" style="width: 100%"
            :empty-text="t('tip.noParameter')">
            <el-table-column :label="t('field.key')" width="180">
              <template #default="scope">
                <el-input v-model="scope.row.key" placeholder="Key" @change="paramChange"/>
              </template>
            </el-table-column>
            <el-table-column :label="t('field.value')">
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
        <el-tab-pane name="output">
          <template #label>
            <el-badge :is-dot="testResult.output !== ''" class="item">{{ t('title.output') }}</el-badge>
          </template>
          <el-tag class="ml-2" type="success" v-if="testResult.statusCode && testResult.error === ''">{{ t('httpCode.' + testResult.statusCode) }}</el-tag>
          <el-tag class="ml-2" type="danger" v-if="testResult.statusCode && testResult.error !== ''">{{ t('httpCode.' + testResult.statusCode) }}</el-tag>

          <Codemirror v-model="testResult.output"/>
        </el-tab-pane>
        <el-tab-pane label="Body" name="body">
          <div v-if="testResult.bodyObject">
            <el-input :prefix-icon="Search" @change="responseBodyFilter" v-model="responseBodyFilterText"
              clearable placeholder="$.key" />
            <JsonViewer :value="testResult.bodyObject" :expand-depth="2" copyable boxed sort />
          </div>
          <div v-else>
            <Codemirror v-model="testResult.bodyText"/>
          </div>
        </el-tab-pane>
        <el-tab-pane name="response-header">
          <template #label>
            <el-badge :value="testResult.header.length" 
              :hidden="testResult.header.length === 0" class="item">Header</el-badge>
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

    <el-drawer v-model="duplicateTestCaseDialog">
        <template #default>
            New Test Case Name:<el-input v-model="targetTestCaseName" />
        </template>
        <template #footer>
            <el-button type="primary" @click="duplicateTestCase">{{ t('button.ok') }}</el-button>
        </template>
    </el-drawer>
</template>
