<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Edit, Delete, Search, CopyDocument, Help } from '@element-plus/icons-vue'
import JsonViewer from 'vue-json-viewer'
import type { Pair, TestResult, TestCaseWithSuite, TestCase } from './types'
import { NewSuggestedAPIsQuery, CreateFilter, GetHTTPMethods, FlattenObject } from './types'
import Button from '../components/Button.vue'
import { Cache } from './cache'
import { API } from './net'
import EditButton from '../components/EditButton.vue'
import type { RunTestCaseRequest } from './net'
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

const runTestCaseResultHandler = (e: any) => {
  requestLoading.value = false
  handleTestResult(e)
}
const runTestCase = () => {
  requestLoading.value = true
  const name = props.name
  const suite = props.suite
  const request = {
    suiteName: suite,
    name: name,
    parameters: parameters.value
  } as RunTestCaseRequest

  if (batchRunMode.value) {
    API.BatchRunTestCase({
      count: batchRunCount.value,
      interval: batchRunInterval.value,
      request: request
    }, runTestCaseResultHandler, (e) => {
      parameters.value = []

      requestLoading.value = false
      UIAPI.ErrorTip(e)
      parseResponseBody(e.body)
    })
  } else {
    API.RunTestCase(request, runTestCaseResultHandler, (e) => {
      parameters.value = []

      requestLoading.value = false
      UIAPI.ErrorTip(e)
      parseResponseBody(e.body)
    })
  }
}

const parseResponseBody = (body: any) => {
  if (body === '') {
    return
  }

  try {
    testResult.value.bodyLength = body.length
    testResult.value.bodyObject = JSON.parse(body)
    testResult.value.originBodyObject = JSON.parse(body)
  } catch {
    testResult.value.bodyText = body
  }
}

/**
 * Handles test result data from API response
 *
 * Processes the test response with proper error handling and content type detection:
 * - For JSON responses: Parses and makes it available for filtering/display
 * - For plain text responses: Displays as raw text without JSON parsing
 * - For file responses: Handles as downloadable content
 *
 * @param e The test result data from API
 */
const handleTestResult = (e: any): void => {
  testResult.value = e;

  if (!isHistoryTestCase.value) {
    handleTestResultError(e)
  }
  const isFilePath = e.body.startsWith("isFilePath-")

  if(isFilePath){
    isResponseFile.value = true
  } else if(e.body !== ''){
    testResult.value.bodyLength = e.body.length
    try {
      // Try to parse as JSON, fallback to plain text if parsing fails
      testResult.value.bodyObject = JSON.parse(e.body);
      testResult.value.originBodyObject = JSON.parse(e.body);
      responseBodyFilter()
    } catch (error) {
      // JSON parsing failed, display as plain text
      testResult.value.bodyText = e.body;
      testResult.value.bodyObject = null;
      testResult.value.originBodyObject = null;
    }
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
      showClose: true,
      message: e.error,
      type: 'error'
    });
  } else {
    ElMessage({
      showClose: true,
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
    testResult.value.bodyObject = query
  }
}

const parameterDialogOpened = ref(false)
const batchRunMode = ref(false)
const batchRunCount = ref(1)
const batchRunInterval = ref('1s')
const openBatchRunDialog = () => {
  batchRunMode.value = true
  openParameterDialog()
}
function openParameterDialog() {
  API.GetTestSuite(props.suite, (e) => {
    parameters.value = e.param
    parameterDialogOpened.value = true
  }, UIAPI.ErrorTip)
}

function sendRequestWithParameter() {
  parameterDialogOpened.value = false
  sendRequest()
  batchRunMode.value = false
}

function generateCode() {
  const name = props.name
  const suite = props.suite
  const ID = props.historyCaseID
  if (isHistoryTestCase.value == true){
    API.HistoryGenerateCode({
      id: ID,
      generator: currentCodeGenerator.value
    }, (e) => {
      ElMessage({
        showClose: true,
        message: 'Code generated!',
        type: 'success'
      })
      if (currentCodeGenerator.value === "gRPCPayload") {
        currentCodeContent.value = JSON.stringify(JSON.parse(e.message), null, 4)
      } else {
        currentCodeContent.value = e.message
      }
    }, UIAPI.ErrorTip)
  } else{
    API.GenerateCode({
      suiteName: suite,
      name: name,
      generator: currentCodeGenerator.value
    }, (e) => {
      ElMessage({
        showClose: true,
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
      body: '',
      filepath: ''
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

function determineBodyType(e: TestCase) {
  e.request.header.forEach(item => {
    if (item.key === "Content-Type") {
      switch (item.value) {
        case 'application/x-www-form-urlencoded':
          bodyType.value = 4
          break
        case 'application/json':
          bodyType.value = 5
          break
        case 'multipart/form-data':
          bodyType.value = 6

          e.request.form.forEach(fItem => {
            if (fItem.key !== '' && fItem.key !== '') {
              e.request.filepath = fItem.key + "=" + fItem.value
            }
          })
          break
      }
    }
  });
}

function base64ToBinary(base64: string): Uint8Array {
    const binaryString = atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes;
}

const isResponseFile = ref(false)
function downloadResponseFile(){
  API.DownloadResponseFile({
      body: testResult.value.body
    }, (e) => {
      if (e && e.data) {
        try {
        const bytes = base64ToBinary(e.data);
        const blob = new Blob([bytes], { type: 'mimeType' });
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        if (e.filename.indexOf('isFilePath-') === -1) {
            link.download = e.filename;
        } else {
            link.download = e.filename.substring("isFilePath-".length);
        }

        console.log(e.filename);
        document.body.appendChild(link);
        link.click();

        window.URL.revokeObjectURL(link.href);
        document.body.removeChild(link);
        } catch (error) {
          console.error('Error during file download:', error);
        }
      } else {
        console.error('No data to download.');
      }
    })
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
    if (isHistoryTestCase.value == true){ 
      testCaseWithSuite.value.data.request.api = `${testCaseWithSuite.value.data.suiteApi}${testCaseWithSuite.value.data.request.api}` 
    }
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
        showClose: true,
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
      showClose: true,
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

function handleDialogClose(){
  caseRevertLoading.value = false
  historyDialogOpened.value = false
  historyForm.value.selectedID = ''
  const target = document.getElementById('compareView');
  target.innerHTML = ''
}

function handleHistoryChange(value) {
  selectedHistory.value = historyRecords.value.find(record => record.ID === value);
  const {
  caseName: name,
  suiteName,
  request,
  response,
  historyHeader,
  } = selectedHistory.value;
  request.header = historyHeader
  request.header.push({
        key: '',
        value: ''
  })
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
    collapseIdentical: true,
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
            showClose: true,
            message: 'Saved.',
            type: 'success'
          })
          load()
        }
      }, UIAPI.ErrorTip, saveLoading)
      handleDialogClose()
    }
  })
}

const goToHistory = async (formEl) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      caseRevertLoading.value = true
      emit('toHistoryPanel', { ID: selectedHistory.value.ID, panelName: 'history' })
      handleDialogClose()
    }
  })
}

const deleteAllHistory = async (formEl) => {
  if (!formEl) return
  caseRevertLoading.value = true
  API.DeleteAllHistoryTestCase(props.suite, props.name, handleDeleteResponse)
  handleDialogClose()
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

const filepathChange = () => {
  const items = testCaseWithSuite.value.data.request.filepath.split("=")
      if (items && items.length > 1) {
        testCaseWithSuite.value.data.request.form = [{
          key: items[0],
          value: items[1]
        } as Pair]
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
    case 6:
      contentType = 'multipart/form-data'
      filepathChange()
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
function jsonFormat(space: number) {
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
            showClose: true,
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

const renameTestCase = (name: string) => {
  const suiteName = props.suite
  API.RenameTestCase(suiteName, suiteName, props.name, name, (d) => {
    emit('updated', suiteName, name)
  })
}
</script>

<template>
  <el-container style="height: 100%;">
    <el-header style="padding-left: 5px;">
      <div style="margin-bottom: 5px">
        <Button type="primary" @click="saveTestCase" :icon="Edit" v-loading="saveLoading"
          disabled v-if="Cache.GetCurrentStore().readOnly || isHistoryTestCase"
          >{{ t('button.save') }}</Button>
        <Button type="primary" @click="saveTestCase" :icon="Edit" v-loading="saveLoading"
          v-if="!Cache.GetCurrentStore().readOnly && !isHistoryTestCase"
          >{{ t('button.save') }}</Button>
        <Button type="danger" @click="deleteCase" :icon="Delete">{{ t('button.delete') }}</Button>
        <Button type="primary" @click="openDuplicateTestCaseDialog" :icon="CopyDocument" v-if="!isHistoryTestCase">{{ t('button.duplicate') }}</Button>
        <Button type="primary" @click="openCodeDialog">{{ t('button.generateCode') }}</Button>
        <Button type="primary" v-if="!isHistoryTestCase && Cache.GetCurrentStore().kind.name == 'atest-store-orm'" @click="openHistoryDialog">{{ t('button.viewHistory') }}</Button>
        <span v-if="isHistoryTestCase" style="margin-left: 15px;">{{ t('tip.runningAt') }}{{ HistoryTestCaseCreateTime }}</span>
        <EditButton :value="props.name" @changed="renameTestCase"/>
      </div>
      <div>
        <el-row justify="space-between" gutter="10">
          <el-col :span="3">
            <el-select
              v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'"
              v-model="testCaseWithSuite.data.request.method"
              class="m-2"
              placeholder="Method"
              size="default"
              test-id="case-editor-method"
              :disabled="isHistoryTestCase"
            >
              <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.key"
                :value="item.value"
              >
              <el-text class="mx-1" :type="item.type">{{ item.key }}</el-text>
             </el-option>
            </el-select>
          </el-col>
          <el-col :span="19">
            <el-autocomplete
              v-model="testCaseWithSuite.data.request.api"
              style="width: 100%"
              :fetch-suggestions="querySuggestedAPIs"
              :readonly="isHistoryTestCase">
              <template #default="{ item }">
                <div class="value">{{ item.request.method }}</div>
                <span class="link">{{ item.request.api }}</span>
              </template>
              <template #prefix v-if="!testCaseWithSuite.data.request.api.startsWith('http://') && !testCaseWithSuite.data.request.api.startsWith('https://')">
                {{ testCaseWithSuite.data.server }}
              </template>
            </el-autocomplete>
          </el-col>
          <el-col :span="2" style="text-align-last: right;">
            <el-dropdown split-button type="primary"
              @click="sendRequest"
              v-loading="requestLoading"
              v-if="!isHistoryTestCase">
              {{ t('button.send') }}
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="openParameterDialog">{{ t('button.sendWithParam') }}</el-dropdown-item>
                  <el-dropdown-item @click="openBatchRunDialog">Batch Send</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </el-col>
        </el-row>
      </div>
    </el-header>

    <el-main style="padding-left: 5px; min-height: 280px;">
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
                  :readonly="isHistoryTestCase"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="Value" :readonly="isHistoryTestCase"/>
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
                  :readonly="isHistoryTestCase"
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
                    :readonly="isHistoryTestCase"
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
                :readonly="isHistoryTestCase"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="value" :readonly="isHistoryTestCase"/>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane name="body">
          <span style="margin-right: 10px; padding-right: 5px;">
            <Button type="primary" @click="jsonFormat(4)">Beautify</Button>
            <Button type="primary" @click="jsonFormat(0)">Minify</Button>
            <el-text class="mx-1">Choose the body format</el-text>
          </span>
          <template #label>
            <el-badge :is-dot="testCaseWithSuite.data.request.body !== ''" class="item">Body</el-badge>
          </template>
          <el-radio-group v-model="bodyType" @change="bodyTypeChange">
            <el-radio :value="1">none</el-radio>
            <el-radio :value="2">form-data</el-radio>
            <el-radio :value="3">raw</el-radio>
            <el-radio :value="4">x-www-form-urlencoded</el-radio>
            <el-radio :value="5">JSON</el-radio>
            <el-radio :value="6">EmbedFile</el-radio>
          </el-radio-group>

          <div style="flex-grow: 1;">
            <div v-if="bodyType === 6">
              <el-row>
                <el-col :span="4">Filename:</el-col>
                <el-col :span="20">
                  <el-input v-model="testCaseWithSuite.data.request.filepath" placeholder="file=sample.txt" @change="filepathChange" />
                </el-col>
              </el-row>
            </div>
            <Codemirror v-if="bodyType === 3 || bodyType === 5 || bodyType === 6"
              @blur="jsonFormat(-1)"
              v-model="testCaseWithSuite.data.request.body"
              :disabled="isHistoryTestCase"/>
            <el-table :data="testCaseWithSuite.data.request.form" style="width: 100%" v-if="bodyType === 4">
              <el-table-column label="Key" width="180">
                <template #default="scope">
                  <el-input v-model="scope.row.key" placeholder="Key" @change="formChange" :readonly="isHistoryTestCase"/>
                </template>
              </el-table-column>
              <el-table-column label="Value">
                <template #default="scope">
                  <div style="display: flex; align-items: center">
                    <el-input v-model="scope.row.value" placeholder="Value" :readonly="isHistoryTestCase"/>
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
          <el-row>
            <el-col :span="4">
              Status Code:
            </el-col>
            <el-col :span="20">
              <el-input
                v-model="testCaseWithSuite.data.response.statusCode"
                class="w-50 m-2"
                placeholder="Please input"
                :readonly="isHistoryTestCase">
                <template #append>
                  {{  t('httpCode.' + testCaseWithSuite.data.response.statusCode) }}
                </template>
              </el-input>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="4">Body:</el-col>
            <el-col :span="20">
              <el-input
                v-model="testCaseWithSuite.data.response.body"
                :autosize="{ minRows: 4, maxRows: 8 }"
                type="textarea"
                placeholder="Expected Body"
                :readonly="isHistoryTestCase"
              />
            </el-col>
          </el-row>
        </el-tab-pane>

        <el-tab-pane label="Expected Headers" name="expected-headers" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-table :data="testCaseWithSuite.data.response.header" style="width: 100%">
            <el-table-column label="Key" width="180">
              <template #default="scope">
                <el-input
                  v-model="scope.row.key"
                  placeholder="Key"
                  @change="expectedHeaderChange"
                  :readonly="isHistoryTestCase"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="Value" :readonly="isHistoryTestCase"/>
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
                  :readonly="isHistoryTestCase"
                />
              </template>
            </el-table-column>
            <el-table-column label="Value">
              <template #default="scope">
                <div style="display: flex; align-items: center">
                  <el-input v-model="scope.row.value" placeholder="Value" :readonly="isHistoryTestCase"/>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="Verify" name="verify" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <div v-for="verify in testCaseWithSuite.data.response.verify" :key="verify">
            <el-input :value="verify" :readonly="isHistoryTestCase"/>
          </div>
        </el-tab-pane>

        <el-tab-pane label="Schema" name="schema" v-if="props.kindName !== 'tRPC' && props.kindName !== 'gRPC'">
          <el-input
            v-model="testCaseWithSuite.data.response.schema"
            :autosize="{ minRows: 4, maxRows: 20 }"
            type="textarea"
            :readonly="isHistoryTestCase"
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
            <Button type="primary" @click="generateCode">{{ t('button.refresh') }}</Button>
            <Button type="primary" @click="copyCode">{{ t('button.copy') }}</Button>
          </div>
          <Codemirror v-model="currentCodeContent"/>
        </template>
      </el-drawer>

      <el-dialog @close="handleDialogClose" v-model="historyDialogOpened" :title="t('button.viewHistory')" width="60%" draggable>
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
              <div style="display: flex;flex-wrap: nowrap;justify-content: flex-end;">
                <Button
                  type="primary"
                  @click="submitForm(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.revert') }}
                </Button>
                <Button
                  type="primary"
                  @click="goToHistory(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.goToHistory') }}
                </Button>
                <Button
                  type="primary"
                  @click="deleteAllHistory(viewHistoryRef)"
                  :loading="caseRevertLoading"
                  >{{ t('button.deleteAllHistory') }}
                </Button>
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
          <div v-if="batchRunMode">
            <el-row>
              <el-col :span="6">
                Count:
              </el-col>
              <el-col :span="18">
                <el-input v-model="batchRunCount" type="number" min="1" max="100"/>
              </el-col>
            </el-row>
            <el-row>
              <el-col :span="6">
                Interval:
              </el-col>
              <el-col :span="18">
                <el-input v-model="batchRunInterval" />
              </el-col>
            </el-row>
          </div>
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

          <Button type="primary" @click="sendRequestWithParameter">{{ t('button.send') }}</Button>
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
              clearable placeholder="$.data[?(@.status==='SUCCEED')]">
              <template #prepend v-if="testResult.bodyLength > 0">Body Size: {{testResult.bodyLength}}</template>
              <template #suffix>
                <a href="https://www.npmjs.com/package/jsonpath-plus" target="_blank"><el-icon><Help /></el-icon></a>
              </template>
            </el-input>
            <JsonViewer :value="testResult.bodyObject" :expand-depth="5" copyable boxed sort />
          </div>
          <div v-else>
            <Codemirror v-if="!isResponseFile" v-model="testResult.bodyText"/>
            <div v-if="isResponseFile" style="padding-top: 10px;">
            <el-row>
              <el-col :span="10">
                <div>Response body is too large, please download to view.</div>
              </el-col>
              <el-col :span="2">
                <Button type="primary" @click="downloadResponseFile">Download</Button>
              </el-col>
            </el-row>
          </div>
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
            <Button type="primary" @click="duplicateTestCase">{{ t('button.ok') }}</Button>
        </template>
    </el-drawer>
</template>
