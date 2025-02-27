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

interface Tree {
    label: string
    children?: Tree[]
}
const tablesTree: Tree[] = []
watch(tables, (t) => {
    // clear tablesTree
    tablesTree.splice(0, tablesTree.length)
    t.forEach((i) => {
        tablesTree.push({
            label: i,
        })
    })
    console.log(tablesTree)
})
watch(store, (s) => {
    stores.value.forEach((e: Store) => {
        if (e.name === s) {
            kind.value = e.kind.name
            return
        }
    })
})
watch(kind, (k) => {
    switch (k) {
        case 'atest-store-orm':
            sqlQuery.value = 'show tables'
            queryTip.value = 'Enter SQL query'
            break;
        case 'atest-store-etcd':
            sqlQuery.value = ''
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
    queryResult.value = result
    columns.value = Array.from(cols).sort((a, b) => {
        if (a === 'id') return -1;
        if (b === 'id') return 1;
        return a.localeCompare(b);
    })
}

const keyValueDataHandler = (data) => {
    queryResult.value = []
    data.Pairs.forEach(e => {
        const obj = {}
        obj['key'] = e.key
        obj['value'] = e.value
        queryResult.value.push(obj)

        columns.value = ['key', 'value']
    })
}

const executeQuery = async () => {
    API.DataQuery(store.value, kind.value, sqlQuery.value, (data) => {
        switch (kind.value) {
            case 'atest-store-orm':
                ormDataHandler(data)
                break;
            case 'atest-store-etcd':
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
    <el-container>
      <el-aside width="200px">
          <el-select v-model="currentDatabase" placeholder="Select database">
              <el-option v-for="item in databases" :key="item" :label="item"
                         :value="item"></el-option>
          </el-select>
          <el-tree :data="tablesTree" />
      </el-aside>
      <el-container>
          <el-header>
              <el-form @submit.prevent="executeQuery">
                  <el-row :gutter="10">
                      <el-col :span="2">
                          <el-form-item>
                              <el-select v-model="store" placeholder="Select store">
                                  <el-option v-for="item in stores" :key="item.name" :label="item.name"
                                             :value="item.name" :disabled="!item.ready" :kind="item.kind.name"></el-option>
                              </el-select>
                          </el-form-item>
                      </el-col>
                      <el-col :span="18">
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
