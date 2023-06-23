<script setup lang="ts">
import { ref, watch } from 'vue'
import type { TabsPaneContext } from 'element-plus'

const props = defineProps({
    name: String,
})

interface Pair {
    key: string
    value: string
}

const verifyList = ref('')
const requestBody = ref('')
const bodyFieldsExpect = ref('')
const headersData = ref('')
watch(props, (p) => {
    const name = p.name
    const requestOptions = {
        method: 'POST'
    };
    fetch('/server.Runner/GetSuite', requestOptions)
        .then(response => response.json())
        .then(d => {
            d.items.forEach(e => {
                if (e.name === name) {
                    if (e.request.method === "") {
                        e.request.method = "GET"
                    }
                    value.value = e.request.method
                    input.value = d.api + e.request.api
                    requestBody.value = e.request.body
                    verifyList.value = e.response.verify

                    headersData.value = []
                    e.request.header.forEach(h => {
                        headersData.value.push({
                            key: h.key,
                            value: h.value
                        })
                    })

                    let items = []
                    e.response.bodyFieldsExpect.forEach(b => {
                        items.push({
                            key: b.key,
                            value: b.value
                        })
                    })
                    bodyFieldsExpect.value = items
                    console.log(items)
                }
            });
        });
})

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

const defaultProps = {
    children: 'children',
    label: 'label',
}

const input = ref('')

function change() {
    let lastItem = tableData[tableData.length - 1]
    if (lastItem.key !== '') {
        tableData.push({
            key: '',
            value: ''
        })
    }
}
</script>

<template>
    <div class="common-layout">
        <el-container>
            <el-header style="padding-left: 5px;">
                <el-select v-model="value" class="m-2" placeholder="Method" size="large">
                    <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
                </el-select>
                <el-input v-model="input" placeholder="API Address"  style="width: 70%; margin-left: 5px; margin-right: 5px;"/>
                <el-button type="primary">Send</el-button>
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
                        <el-radio-group v-model="radio">
                            <el-radio :label="3">none</el-radio>
                            <el-radio :label="9">form-data</el-radio>
                            <el-radio :label="6">raw</el-radio>
                            <el-radio :label="9">x-www-form-urlencoded</el-radio>
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
                </el-tabs>
            </el-main>
        </el-container>
    </div>
</template>
