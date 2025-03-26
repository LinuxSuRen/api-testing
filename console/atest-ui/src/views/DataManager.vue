<script setup lang="ts">
import { ref, watch } from 'vue'
import { API } from './net'
import { ElMessage } from 'element-plus'

const stores = ref([])
const kind = ref('')
const store = ref('')
const sqlQuery = ref('')
const queryResult = ref([])
const columns = ref([])
const queryTip = ref('')
const databases = ref([])
const tables = ref([])
const currentDatabase = ref('')
const loadingStores = ref(true)

const tablesTree = ref([])
watch(store, (s) => {
    kind.value = ''
    stores.value.forEach((e: Store) => {
        if (e.name === s) {
            kind.value = e.kind.name
            return
        }
    })
    currentDatabase.value = ''
    sqlQuery.value = ''
    executeQuery()
})
const queryDataFromTable = (data) => {
    sqlQuery.value = `select * from ${data.label} limit 10`
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

const ormDataHandler = (data) => {
    const result = []
    const cols = new Set()

    data.items.forEach(e => {
        const obj = {}
        e.data.forEach(item => {
            obj[item.key] = item.value
            cols.add(item.key)
        })
        result.push(obj)
    })

    databases.value = data.meta.databases
    tables.value = data.meta.tables
    currentDatabase.value = data.meta.currentDatabase
    queryResult.value = result
    columns.value = Array.from(cols).sort((a, b) => {
        if (a === 'id') return -1;
        if (b === 'id') return 1;
        return a.localeCompare(b);
    })

    tablesTree.value = []
    tables.value.forEach((i) => {
        tablesTree.value.push({
            label: i,
        })
    })
}

const keyValueDataHandler = (data) => {
    queryResult.value = []
    data.data.forEach(e => {
        const obj = {}
        obj['key'] = e.key
        obj['value'] = e.value
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

    API.DataQuery(store.value, kind.value, currentDatabase.value, sqlQuery.value, (data) => {
        switch (kind.value) {
            case 'atest-store-orm':
                ormDataHandler(data)
                break;
            case 'atest-store-iotdb':
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
    }, (e) => {
        ElMessage({
            showClose: true,
            message: e.message,
            type: 'error'
        });
    })
}
</script>

<template>
  <div>
    <el-container style="height: calc(100vh - 45px);">
      <el-aside v-if="kind === 'atest-store-orm'">
          <el-scrollbar>
              <el-select v-model="currentDatabase" placeholder="Select database" @change="queryTables" filterable>
                  <el-option v-for="item in databases" :key="item" :label="item"
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
                      <el-col :span="17">
                          <el-form-item>
                              <el-input v-model="sqlQuery" :placeholder="queryTip" @keyup.enter="executeQuery"></el-input>
                          </el-form-item>
                      </el-col>
                      <el-col :span="2">
                          <el-form-item>
                              <el-button type="primary" @click="executeQuery">Execute</el-button>
                          </el-form-item>
                      </el-col>
                  </el-row>
              </el-form>
          </el-header>
          <el-main>
              <el-table :data="queryResult">
                  <el-table-column v-for="col in columns" :key="col" :prop="col" :label="col"></el-table-column>
              </el-table>
          </el-main>
      </el-container>
    </el-container>
  </div>
</template>
