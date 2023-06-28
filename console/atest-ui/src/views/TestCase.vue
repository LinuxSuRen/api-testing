<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import type { TabsPaneContext } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Delete, Edit, Search, Share, Upload } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'

const props = defineProps({
    name: String,
    suite: String,
})

const requestLoading = ref(false)
const testResult = ref('')
function sendRequest() {
    requestLoading.value = true
    const name = props.name
    const suite = props.suite
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            suite: suite,
            testcase: name,
        })
    };
    fetch('/server.Runner/RunTestCase', requestOptions)
        .then(response => response.json())
        .then(e => {
            testResult.value = e.body
            requestLoading.value = false

            if (e.error !== "") {
                ElMessage({
                    message: e.error,
                    type: 'error'
                })
            } else {
                ElMessage({
                    message: 'Pass!',
                    type: 'success',
                })
            }
        });
}

interface Pair {
    key: string,
    value: string
}

const emptyTestCaseWithSuite: TestCaseWithSuite = {
    suiteName: "",
    data: {
        name: "",
        request: {
            api: "",
            method: "",
            header: [],
            query: [],
            form: [],
            body: "",
        },
        response: {
            statusCode: 0,
            body: "",
            header: [],
            bodyFieldsExpect: [],
            verify: [],
            schema: "",
        },
    }
}

const testCaseWithSuite = ref(emptyTestCaseWithSuite)

watch(props, (p) => {
    const name = p.name
    const suite = p.suite
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            suite: suite,
            testcase: name,
        })
    };
    fetch('/server.Runner/GetTestCase', requestOptions)
        .then(response => response.json())
        .then(e => {
            if (e.request.method === "") {
                e.request.method = "GET"
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

            testCaseWithSuite.value = {
                suiteName: suite,
                data: e
            }
        });
})

interface TestCaseWithSuite{
    suiteName: string,
    data: TestCase
}

interface TestCase {
    name: string,
    request: TestCaseRequest,
    response: TestCaseResponse,
}

interface TestCaseRequest {
    api: string,
    method: string,
    header: Pair[],
    query: Pair[],
    form: Pair[],
    body: string,
}

interface TestCaseResponse {
    statusCode: number,
    body: string,
    header: Pair[],
    bodyFieldsExpect: Pair[],
    verify: string[],
    schema: string,
}

const saveLoading = ref(false)
function saveTestCase() {
    saveLoading.value = true

    // remove empty pair
    testCaseWithSuite.value.data.request.header = testCaseWithSuite.value.data.request.header.filter(e => e.key !== '')
    testCaseWithSuite.value.data.request.query = testCaseWithSuite.value.data.request.query.filter(e => e.key !== '')
    testCaseWithSuite.value.data.request.form = testCaseWithSuite.value.data.request.form.filter(e => e.key !== '')
    testCaseWithSuite.value.data.response.header = testCaseWithSuite.value.data.response.header.filter(e => e.key !== '')
    testCaseWithSuite.value.data.response.bodyFieldsExpect = testCaseWithSuite.value.data.response.bodyFieldsExpect.filter(e => e.key !== '')
    testCaseWithSuite.value.data.response.verify = testCaseWithSuite.value.data.response.verify.filter(e => e !== '')

    const requestOptions = {
        method: 'POST',
        body: JSON.stringify(testCaseWithSuite.value)
    };
    fetch('/server.Runner/UpdateTestCase', requestOptions)
        .then(e => {
            if (e.ok) {
                ElMessage({
                    message: 'Saved.',
                    type: 'success',
                })
            } else {
                ElMessage.error('Oops, ' + e.statusText)
            }
            saveLoading.value = false
        })
}

const options = [
    {
        value: 'GET',
        label: 'GET',
    },
    {
        value: 'POST',
        label: 'POST',
    },
    {
        value: 'DELETE',
        label: 'DELETE',
    },
    {
        value: 'PUT',
        label: 'PUT',
    },
]

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
        })
    }
}

function headerChange() {
    const header = testCaseWithSuite.value.data.request.header
    let lastItem = header[header.length - 1]
    if (lastItem.key !== '') {
        header.push({
            key: '',
            value: ''
        })
    }
}

const radio1 = ref('1')

const dialogVisible = ref(false)
const testcaseFormRef = ref<FormInstance>()
const testCaseForm = reactive({
    suiteName: "",
    name: "",
    api: "",
})
function openNewTestCaseDialog() {
    loadTestSuites()
    dialogVisible.value = true
}

const suiteCreatingLoading = ref(false)
const testSuiteList = ref([])
function loadTestSuites() {
  const requestOptions = {
      method: 'POST'
  };
  fetch('/server.Runner/GetSuites', requestOptions)
      .then(response => response.json())
      .then(d => {
        Object.keys(d.data).map(k => {
          testSuiteList.value.push(k)
        })
      });
}
const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  suiteCreatingLoading.value = true

  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
        suiteName: testCaseForm.suiteName,
        data: {
            name: testCaseForm.name,
            request: {
                api: testCaseForm.api,
                method: "GET",
            }
        },
    })
  };

  fetch('/server.Runner/UpdateTestCase', requestOptions)
      .then(response => response.json())
      .then(() => {
        suiteCreatingLoading.value = false
      });
      
  dialogVisible.value = false
}
</script>

<template>
    <div class="common-layout">
        <el-container style="height: 60vh">
            <el-header style="padding-left: 5px;">
                <div style="margin-bottom: 5px;">
                    <el-button type="primary" @click="saveTestCase" :icon="Edit" :loading="saveLoading">Save</el-button>
                    <el-button type="primary" @click="openNewTestCaseDialog" :icon="Edit">New</el-button>
                    <el-text class="mx-1" type="primary">{{props.name}}</el-text>
                </div>
                <el-select v-model="testCaseWithSuite.data.request.method" class="m-2" placeholder="Method" size="middle">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
                <el-input v-model="testCaseWithSuite.data.request.api" placeholder="API Address"  style="width: 70%; margin-left: 5px; margin-right: 5px;"/>
                <el-button type="primary" @click="sendRequest" :loading="requestLoading">Send</el-button>
            </el-header>

            <el-main>
                <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
                    <el-tab-pane label="Headers" name="second">
                        <el-table :data="testCaseWithSuite.data.request.header" style="width: 100%">
                            <el-table-column label="Key" width="180">
                                <template #default="scope">
                                    <el-input v-model="scope.row.key" placeholder="Key" @change="headerChange" />
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
                        <el-radio-group v-model="radio1">
                            <el-radio :label="1">none</el-radio>
                            <el-radio :label="2">form-data</el-radio>
                            <el-radio :label="3">raw</el-radio>
                            <el-radio :label="4">x-www-form-urlencoded</el-radio>
                        </el-radio-group>

                        <el-input v-model="testCaseWithSuite.data.request.body" :autosize="{ minRows: 4, maxRows: 8 }" type="textarea"
                            placeholder="Please input" />
                    </el-tab-pane>

                    <el-tab-pane label="BodyFiledExpect" name="fourth">
                        <el-table :data="testCaseWithSuite.data.response.bodyFieldsExpect" style="width: 100%">
                            <el-table-column label="Key" width="180">
                                <template #default="scope">
                                    <el-input v-model="scope.row.key" placeholder="Key" @change="bodyFiledExpectChange" />
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
                        <el-input v-model="testCaseWithSuite.data.response.schema"
                            :autosize="{ minRows: 4, maxRows: 8 }" type="textarea" />
                    </el-tab-pane>
                </el-tabs>
            </el-main>

            <el-footer>
                <div>Test Result:</div>
                <el-input
                    v-model="testResult"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    readonly=true
                    type="textarea"
                    placeholder="Please input"
                />
            </el-footer>
        </el-container>
    </div>

  <el-dialog v-model="dialogVisible" title="Create Test Case" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          ref="testcaseFormRef"
          status-icon
          label-width="120px"
          class="demo-ruleForm"
        >
          <el-form-item label="Suite" prop="suite">
            <el-select class="m-2" v-model="testCaseForm.suiteName" placeholder="Select" size="large">
                <el-option
                v-for="item in testSuiteList"
                :key="item"
                :label="item"
                :value="item"
                />
            </el-select>
          </el-form-item>
          <el-form-item label="Name" prop="name">
            <el-input v-model="testCaseForm.name" />
          </el-form-item>
          <el-form-item label="API" prop="api">
            <el-input v-model="testCaseForm.api" placeholder="http://foo" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm(testcaseFormRef)" :loading="suiteCreatingLoading">Submit</el-button>
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
