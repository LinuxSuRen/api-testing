<script setup lang="ts">
import { ref, watch } from 'vue'
import { API } from './net'
import type { Store } from './store'
import type { Pair } from './types'
import { ElMessage } from 'element-plus'
import { Codemirror } from 'vue-codemirror'
import HistoryInput from '../components/HistoryInput.vue'

const stores = ref([] as Store[])
const kind = ref('')
const store = ref('')
const sqlQuery = ref('')
const queryResult = ref([] as any[])
const queryResultAsJSON= ref('')
const columns = ref([] as string[])
const queryTip = ref('')
const loadingStores = ref(true)
const dataFormat = ref('table')
const dataFormatOptions = ['table', 'json']
const queryDataMeta = ref({} as QueryDataMeta)

const tablesTree = ref([])
watch(store, (s) => {
    kind.value = ''
    stores.value.forEach((e: Store) => {
        if (e.name === s) {
            kind.value = e.kind.name
            return
        }
    })
    queryDataMeta.value.currentDatabase = ''
    sqlQuery.value = ''
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
    executeQuery()
}
const queryTables = () => {
    sqlQuery.value = ``
    executeQuery()
}
watch(kind, (k) => {
    switch (k) {
        case 'atest-store-orm':
            queryTip.value = 'Enter SQL query'
            executeQuery()
            break;
        case 'atest-store-etcd':
        case 'atest-store-redis':
            queryTip.value = 'Enter key'
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

    data.meta.labels = data.meta.labels.filter((item) => {
        if (item.key === '_native_sql') {
            sqlQuery.value = item.value
            return false
        }
        return !item.key.startsWith('_')
    })

    queryDataMeta.value = data.meta
    queryResult.value = result
    queryResultAsJSON.value = JSON.stringify(result, null, 2)
    columns.value = Array.from(cols).sort((a, b) => {
        if (a === 'id') return -1;
        if (b === 'id') return 1;
        return a.localeCompare(b);
    })

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
    switch (kind.value) {
        case 'atest-store-etcd':
            sqlQuery.value = '*'
            break;
    }

    let success = false
    try {
        const data = await API.DataQueryAsync(store.value, kind.value, queryDataMeta.value.currentDatabase, sqlQuery.value);
        switch (kind.value) {
            case 'atest-store-orm':
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
</script>

<template>
  <div>
    <el-container style="height: calc(100vh - 50px);">
      <el-aside v-if="kind === 'atest-store-orm'">
          <el-scrollbar>
              <el-select v-model="queryDataMeta.currentDatabase" placeholder="Select database" @change="queryTables" filterable>
                  <el-option v-for="item in queryDataMeta.databases" :key="item" :label="item"
                             :value="item"></el-option>
              </el-select>
              <el-tree :data="tablesTree" node-key="label" @node-click="queryDataFromTable" highlight-current draggable/>
          </el-scrollbar>
      </el-aside>
      <el-container>
          <el-header>
              <el-form @submit.prevent="executeQuery">
                  <el-row :gutter="10">
                      <el-col :span="4">
                          <el-form-item>
                              <el-select v-model="store" placeholder="Select store" filterable :loading="loadingStores">
                                  <el-option v-for="item in stores" :key="item.name" :label="item.name"
                                         :value="item.name" :disabled="!item.ready" :kind="item.kind.name"></el-option>
                              </el-select>
                          </el-form-item>
                      </el-col>
                      <el-col :span="16">
                          <el-form-item>
                              <HistoryInput :placeholder="queryTip" :callback="executeQuery" v-model="sqlQuery"/>
                          </el-form-item>
                      </el-col>
                      <el-col :span="2">
                          <el-form-item>
                              <el-button type="primary" @click="executeQuery">Execute</el-button>
                          </el-form-item>
                      </el-col>
                      <el-col :span="2">
                        <el-select v-model="dataFormat" placeholder="Select data format">
                            <el-option v-for="item in dataFormatOptions" :key="item" :label="item" :value="item"></el-option>
                        </el-select>
                      </el-col>
                  </el-row>
              </el-form>
          </el-header>
          <el-main>
              <div style="display: flex; gap: 8px;">
                <el-tag type="primary" v-if="queryResult.length > 0">{{ queryResult.length }} rows</el-tag>
                <el-tag type="primary" v-if="queryDataMeta.duration">{{  queryDataMeta.duration }}</el-tag>
                <el-tag type="primary" v-for="label in queryDataMeta.labels">{{  label.value }}</el-tag>
              </div>
              <el-table :data="queryResult" stripe v-if="dataFormat === 'table'">
                  <el-table-column v-for="col in columns" :key="col" :prop="col" :label="col" sortable/>
              </el-table>
              <Codemirror v-else-if="dataFormat === 'json'"
                v-model="queryResultAsJSON"/>
          </el-main>
      </el-container>
    </el-container>
  </div>
</template>
