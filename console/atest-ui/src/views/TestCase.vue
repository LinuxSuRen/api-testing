<script setup lang="ts">
import { ref, watch } from 'vue'
import type { TabsPaneContext } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Delete, Edit, Search, Share, Upload } from '@element-plus/icons-vue'

const props = defineProps({
    name: String,
    suite: String,
})

const testResult = ref('')
function sendRequest() {
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
        });
}

interface Pair{
    key: string,
    value: string
}

const emptyPair: Pair[] = []

const apiAddress = ref('')
const verifyList = ref('')
const requestBody = ref('')
const responseVerifySchema = ref('')
const bodyFieldsExpect = ref(emptyPair)
const headersData = ref(emptyPair)

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
            value.value = e.request.method
            apiAddress.value = e.request.api
            requestBody.value = e.request.body
            verifyList.value = e.response.verify
            responseVerifySchema.value = e.response.schema

            headersData.value = []
            e.request.header.forEach(h => {
                headersData.value.push({
                    key: h.key,
                    value: h.value
                })
            })

            let items: Pair[] = []
            e.response.bodyFieldsExpect.forEach(b => {
                items.push({
                    key: b.key,
                    value: b.value
                })
            })
            bodyFieldsExpect.value = items
        });
})

interface TestCaseWithSuite{
    suiteName: string,
    data: TestCase
}

interface TestCase {
    name: string,
    request: TestCaseRequest,
}

interface TestCaseRequest {
    method: string,
    api: string,
    body: string,
}

function saveTestCase() {
    const p = props
    let testCaseWithSuite: TestCaseWithSuite = {
        suiteName: p.suite,
        data: {
            name: p.name,
            request: {
                method: value.value,
                api: apiAddress.value,
                body: requestBody.value,
            }
        }
    }
    console.log(testCaseWithSuite)

    const requestOptions = {
        method: 'POST',
        body: JSON.stringify(testCaseWithSuite)
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
        })
}

const value = ref('')

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

function change() {
    let lastItem = headersData.value[headersData.value.length - 1]
    if (lastItem.key !== '') {
        headersData.value.push({
            key: '',
            value: ''
        })
    }
}

const radio1 = ref('1')
</script>

<template>
    <div class="common-layout">
        <el-container>
            <el-header style="padding-left: 5px;">
                <div style="margin-bottom: 5px;">
                    <el-button type="primary" @click="saveTestCase" :icon="Edit">Save</el-button>
                </div>
                <el-select v-model="value" class="m-2" placeholder="Method" size="large">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
                <el-input v-model="apiAddress" placeholder="API Address"  style="width: 70%; margin-left: 5px; margin-right: 5px;"/>
                <el-button type="primary" @click="sendRequest">Send</el-button>
            </el-header>

            <el-main>
                <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
                    <el-tab-pane label="Headers" name="second">
                        <el-table :data="headersData" style="width: 100%">
                            <el-table-column label="Key" width="180">
                                <template #default="scope">
                                    <el-input v-model="scope.row.key" placeholder="Key" @change="change" />
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

                        <el-input v-model="requestBody" :autosize="{ minRows: 4, maxRows: 8 }" type="textarea"
                            placeholder="Please input" />
                    </el-tab-pane>

                    <el-tab-pane label="BodyFiledExpect" name="fourth">
                        <el-table :data="bodyFieldsExpect" style="width: 100%">
                            <el-table-column label="Key" width="180">
                                <template #default="scope">
                                    <el-input v-model="scope.row.key" placeholder="Key" @change="change" />
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
                        <div v-for="verify in verifyList" :key="verify">
                            <el-input :value="verify" placeholder="API Address" />
                        </div>
                    </el-tab-pane>

                    <el-tab-pane label="Schema" name="schema">
                        <el-input :value="responseVerifySchema" />
                    </el-tab-pane>
                </el-tabs>
            </el-main>

            <el-footer>
                <div>Test Result:</div>
                <el-input
                    v-model="testResult"
                    :autosize="{ minRows: 4, maxRows: 6 }"
                    readonly="true"
                    type="textarea"
                    placeholder="Please input"
                />
            </el-footer>
        </el-container>
    </div>
</template>
