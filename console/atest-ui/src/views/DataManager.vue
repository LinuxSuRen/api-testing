<script setup lang="ts">
import { ref, watch } from 'vue'
import { API } from './net'
import type { QueryObject } from './net'
import type { Store } from './store'
import type { Pair } from './types'
import { ElMessage } from 'element-plus'
import { Codemirror } from 'vue-codemirror'
import HistoryInput from '../components/HistoryInput.vue'
import type { Ref } from 'vue'
import { Refresh, Document } from '@element-plus/icons-vue'

const stores: Ref<Store[]> = ref([])
const kind = ref('')
const store = ref('')
const query = ref({
    offset: 0,
    limit: 10
} as QueryObject)
const currentTable = ref('')
const sqlQuery = ref('')
const queryResult = ref([] as any[])
const queryResultAsJSON = ref('')
const columns = ref([] as string[])
const queryTip = ref('')
const loadingStores = ref(true)
const showOverflowTooltip = ref(true)
const complexEditor = ref(false)
const dataFormat = ref('table')
const dataFormatOptions = ['table', 'json']
const queryDataMeta = ref({} as QueryDataMeta)
const largeContent = ref('')
const largeContentDialogVisible = ref(false)

interface TreeItem {
    label: string
}
const tablesTree = ref([] as TreeItem[])
watch(store, (s) => {
    kind.value = ''
    stores.value.forEach((e: Store) => {
        if (e.name === s) {
            kind.value = e.kind.name
            return
        }
    })

    switch (kind.value) {
        case 'atest-store-elasticsearch':
        case 'atest-store-etcd':
            sqlQuery.value = '*'
            complexEditor.value = false
            break
        default:
            complexEditor.value = true
            queryDataMeta.value.currentDatabase = ''
            sqlQuery.value = ''
    }

    executeQuery()
})

interface QueryDataMeta {
    databases: string[]
    tables: string[]
    currentDatabase: string
    duration: string
    labels: Pair[]
}

interface QueryData {
    items: any[]
    data: any[]
    label: string
    meta: QueryDataMeta
}

const queryDataFromTable = (data: QueryData) => {
    sqlQuery.value = `@selectTableLImit100_${data.label}`
    currentTable.value = data.label
    executeQuery()
}
const describeTable = (data: QueryData) => {
    switch (kind.value) {
        case 'atest-store-cassandra':
            sqlQuery.value = `@describeTable_${queryDataMeta.value.currentDatabase}:${data.label}`
            break
            break
        default:
            sqlQuery.value = `@describeTable_${data.label}`
    }
    executeQuery()
}
const queryTables = () => {
    switch (kind.value) {
        case 'atest-store-elasticsearch':
            if (sqlQuery.value === '') {
                sqlQuery.value = '*'
            }
            break
        default:
            sqlQuery.value = ``
    }
    executeQuery()
}
watch(kind, (k) => {
    switch (k) {
        case 'atest-store-orm':
        case 'atest-store-cassandra':
        case 'atest-store-iotdb':
            queryTip.value = 'Enter SQL query'
            executeQuery()
            break;
        case 'atest-store-etcd':
        case 'atest-store-redis':
            queryTip.value = 'Enter key'
            break;
        case 'atest-store-elasticsearch':
            queryTip.value = 'field:value OR field:other'
            break;
    }
})

API.GetStores((data) => {
    stores.value = data.data
}, (e) => {
    ElMessage({
        showClose: true,
        message: e.message,
        type: 'error'
    });
}, () => {
    loadingStores.value = false
})

const ormDataHandler = (data: QueryData) => {
    const result = [] as any[]
    const cols = new Set<string>()

    data.items.forEach(e => {
        const obj = {}
        e.data.forEach((item: Pair) => {
            obj[item.key] = item.value
            cols.add(item.key)
        })
        result.push(obj)
    })

    columns.value = [] as string[]
    data.meta.labels = data.meta.labels.filter((item) => {
        if (item.key === '_native_sql') {
            sqlQuery.value = item.value
            return false
        }
        if (item.key === '_columns') {
            columns.value = JSON.parse(item.value)
        }
        return !item.key.startsWith('_')
    })

    queryDataMeta.value = data.meta
    queryResult.value = result
    queryResultAsJSON.value = JSON.stringify(result, null, 2)
    if (columns.value.length == 0) {
        columns.value = Array.from(cols).sort((a, b) => {
            if (a === 'id') return -1;
            if (b === 'id') return 1;
            return a.localeCompare(b);
        })
    }

    tablesTree.value = []
    queryDataMeta.value.tables.forEach((i) => {
        tablesTree.value.push({
            label: i,
        })
    })
}

const keyValueDataHandler = (data: QueryData) => {
    queryResult.value = []
    data.data.forEach(e => {
        const obj = new Map<string, string>();
        obj.set('key', e.key)
        obj.set('value', e.value)
        queryResult.value.push(obj)

        columns.value = ['key', 'value']
    })
}

const executeQuery = async () => {
    return executeWithQuery(sqlQuery.value)
}
const executeWithQuery = async (sql: string) => {
    switch (kind.value) {
        case 'atest-store-etcd':
            sqlQuery.value = '*'
            break;
        case '':
            return;
    }

    let success = false
    query.value.store = store.value
    query.value.key = queryDataMeta.value.currentDatabase
    query.value.sql = sql

    try {
        const data = await API.DataQueryAsync(query.value);
        switch (kind.value) {
            case 'atest-store-orm':
            case 'atest-store-cassandra':
            case 'atest-store-iotdb':
            case 'atest-store-opengemini':
            case 'atest-store-elasticsearch':
                ormDataHandler(data)
                success = true
                break;
            case 'atest-store-etcd':
                keyValueDataHandler(data)
                break;
            case 'atest-store-redis':
                keyValueDataHandler(data)
                break;
            default:
                ElMessage({
                    showClose: true,
                    message: 'Unsupported store kind',
                    type: 'error'
                });
        }
    } catch (e: any) {
        ElMessage({
            showClose: true,
            message: e.message,
            type: 'error'
        });
    }
    return success
}
const nextPage = () => {
    query.value.offset += query.value.limit
    executeQuery()
}
const overflowChange = () => {
    showOverflowTooltip.value = !showOverflowTooltip.value
}
const tryShowPrettyJSON = (row: any, column: any, cell: HTMLTableCellElement, event: Event) => {
    const cellValue = row[column.rawColumnKey]
    const obj = JSON.parse(cellValue)
    if (obj) {
        largeContent.value = JSON.stringify(obj, null, 2)
    }
}
watch(largeContent, (e) => {
    largeContentDialogVisible.value = e !== ''
})
</script>

<template>
    <div>
        <el-container style="height: calc(100vh - 80px);">
            <el-aside v-if="kind === 'atest-store-orm' || kind === 'atest-store-iotdb' || kind === 'atest-store-cassandra' || kind === 'atest-store-elasticsearch' || kind === 'atest-store-opengemini'">
                <el-scrollbar>
                    <el-select v-model="queryDataMeta.currentDatabase" placeholder="Select database"
                        @change="queryTables" filterable>
                        <template #header>
                            <el-button type="primary" :icon="Refresh" @click="executeWithQuery('')"></el-button>
                        </template>
                        <el-option v-for="item in queryDataMeta.databases" :key="item" :label="item"
                            :value="item"></el-option>
                    </el-select>
                    <el-tree :data="tablesTree" node-key="label" highlight-current
                        draggable>
                        <template #default="{node, data}">
                            <div class="space-between">
                                <span @click="queryDataFromTable(data)">
                                    {{ node.label }}
                                </span>
                                <el-icon style="margin-left: 6px;" @click="describeTable(data)" v-if="kind === 'atest-store-orm' || kind === 'atest-store-cassandra'"><Document /></el-icon>
                            </div>
                        </template>
                    </el-tree>
                </el-scrollbar>
            </el-aside>
            <el-container>
                <el-header style="height: auto">
                    <el-form @submit.prevent="executeQuery">
                        <el-row :gutter="10" justify="center">
                            <el-col :span="4">
                                <el-form-item>
                                    <el-select v-model="store" placeholder="Select store" filterable
                                        :loading="loadingStores">
                                        <el-option v-for="item in stores" :key="item.name" :label="item.name"
                                            :value="item.name" :disabled="!item.ready"
                                            :kind="item.kind.name"></el-option>
                                    </el-select>
                                </el-form-item>
                            </el-col>
                            <el-col :span="16">
                                <el-form-item v-if="!complexEditor">
                                    <HistoryInput :placeholder="queryTip" :callback="executeQuery" v-model="sqlQuery" />
                                </el-form-item>
                            </el-col>
                            <el-col :span="2">
                                <el-form-item>
                                    <el-button type="primary" @click="executeQuery" :disabled="kind === ''">Execute</el-button>
                                </el-form-item>
                            </el-col>
                            <el-col :span="2">
                                <el-select v-model="dataFormat" placeholder="Select data format">
                                    <el-option v-for="item in dataFormatOptions" :key="item" :label="item"
                                        :value="item"></el-option>
                                </el-select>
                            </el-col>
                        </el-row>
                        <el-row :gutter="10" v-if="kind === 'atest-store-elasticsearch'">
                            <el-col :span="10">
                                <el-input type="number" v-model="query.offset">
                                    <template #prepend>Offset</template>
                                </el-input>
                            </el-col>
                            <el-col :span="10">
                                <el-input type="number" v-model="query.limit">
                                    <template #prepend>Limit</template>
                                </el-input>
                            </el-col>
                            <el-col :span="2">
                                <el-button type="primary" @click="nextPage">Next</el-button>
                            </el-col>
                        </el-row>
                    </el-form>
                    <Codemirror v-model="sqlQuery" v-if="complexEditor" style="height: 180px"/>
                </el-header>
                <el-main>
                    <div style="display: flex; gap: 8px;">
                        <el-tag type="primary" v-if="queryResult.length > 0">{{ queryResult.length }} rows</el-tag>
                        <el-tag type="primary" v-if="queryDataMeta.duration">{{ queryDataMeta.duration }}</el-tag>
                        <el-tag type="primary" v-for="label in queryDataMeta.labels">{{ label.value }}</el-tag>
                        <el-check-tag type="primary" :checked="showOverflowTooltip" @change="overflowChange" v-if="queryResult.length > 0">overflow</el-check-tag>
                    </div>
                    <el-table :data="queryResult" stripe v-if="dataFormat === 'table'" height="calc(100vh - 380px)" @cell-dblclick="tryShowPrettyJSON">
                        <el-table-column v-for="col in columns" :key="col" :prop="col" :label="col" sortable :show-overflow-tooltip="showOverflowTooltip" />
                    </el-table>
                    <Codemirror v-else-if="dataFormat === 'json'" v-model="queryResultAsJSON" />
                </el-main>
            </el-container>
        </el-container>
    </div>

    <el-dialog
        v-model="largeContentDialogVisible"
        destroy-on-close
        @closed="largeContent=''"
        center
    >
        <Codemirror v-model="largeContent" />
    </el-dialog>
</template>

<style>
.space-between {
    justify-content: space-between;
    display: flex;
    width: 100%;
    padding-right: 8px;
}
</style>
