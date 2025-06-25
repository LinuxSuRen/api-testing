<script setup lang="ts">
import { ref, watch } from 'vue'
import { API } from './net'
import type { QueryObject } from './net'
import type { Store } from './store'
import { Driver, GetDriverName, ExtensionKind } from './store'
import type { Pair } from './types'
import { GetDataManagerPreference, SetDataManagerPreference } from './cache'
import { ElMessage } from 'element-plus'
import { Codemirror } from 'vue-codemirror'
import { sql, StandardSQL, MySQL, PostgreSQL, Cassandra } from "@codemirror/lang-sql"
import type { SQLConfig } from "@codemirror/lang-sql"
import HistoryInput from '../components/HistoryInput.vue'
import type { Ref } from 'vue'
import { Refresh, Document } from '@element-plus/icons-vue'
import { Magic } from './magicKeys'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
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
const globalLoading = ref(false)
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
const storeChangedEvent = (s: string) => {
  kind.value = ''
  SetDataManagerPreference('currentStore', s)
  stores.value.forEach((e: Store) => {
    if (e.name === s) {
      kind.value = e.kind.name
      switch (GetDriverName(e)){
        case Driver.MySQL:
          sqlConfig.value.dialect = MySQL
          break;
        case Driver.Postgres:
          sqlConfig.value.dialect = PostgreSQL
          break;
        case Driver.Cassandra:
          sqlConfig.value.dialect = Cassandra
          break;
        default:
          sqlConfig.value.dialect = StandardSQL
      }
      return
    }
  })

  switch (kind.value) {
    case ExtensionKind.ExtensionKindElasticsearch:
    case ExtensionKind.ExtensionKindEtcd:
    case ExtensionKind.ExtensionKindRedis:
      sqlQuery.value = '*'
      complexEditor.value = false
      break
    default:
      complexEditor.value = true
  }

  queryDataMeta.value.currentDatabase = GetDataManagerPreference().currentDatabase
  executeQuery()
}
watch(store, storeChangedEvent)

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
    case ExtensionKind.ExtensionKindCassandra:
      sqlQuery.value = `@describeTable_${queryDataMeta.value.currentDatabase}:${data.label}`
      break
    default:
      sqlQuery.value = `@describeTable_${data.label}`
  }
  executeQuery()
}
const queryTables = () => {
  executeQuery()
}
watch(kind, (k) => {
  switch (k) {
    case ExtensionKind.ExtensionKindORM:
    case ExtensionKind.ExtensionKindCassandra:
    case ExtensionKind.ExtensionKindIotDB:
      queryTip.value = 'Enter SQL query'
      executeQuery()
      break;
    case ExtensionKind.ExtensionKindEtcd:
    case ExtensionKind.ExtensionKindRedis:
      queryTip.value = 'Enter key'
      break;
    case ExtensionKind.ExtensionKindElasticsearch:
      queryTip.value = 'field:value OR field:other'
      break;
  }
})
watch(sqlQuery, (sql: string) => {
  SetDataManagerPreference('query', sql)
})

API.GetStores((data) => {
  stores.value = data.data

  sqlQuery.value = GetDataManagerPreference().query
  store.value = GetDataManagerPreference().currentStore
  if (store.value) {
    storeChangedEvent(store.value)
  }
}, (e) => {
  ElMessage({
    showClose: true,
    message: e.message,
    type: 'error'
  });
}, () => {
  loadingStores.value = false
})

const showNativeSQL = ref(true)
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
    if (showNativeSQL.value && item.key === '_native_sql') {
      sqlQuery.value = item.value
      return false
    }
    if (item.key === '_columns') {
      columns.value = JSON.parse(item.value)
    }
    return !item.key.startsWith('_')
  })

  SetDataManagerPreference('currentDatabase', data.meta.currentDatabase)
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
  sqlConfig.value.schema = {}
  queryDataMeta.value.tables.forEach((i) => {
    sqlConfig.value.schema[i] = []
    tablesTree.value.push({
      label: i,
    })
  })
}

const keyValueDataHandler = (data: QueryData) => {
  queryResult.value = []
  columns.value = ['key', 'value']
  data.data.forEach(e => {
    queryResult.value.push({
      key: e.key,
      value: e.value
    })
  })
}
const sqlConfig = ref({
  dialect: StandardSQL,
  defaultSchema: queryDataMeta.value.currentDatabase,
  upperCaseKeywords: true,
  schema: {}
} as SQLConfig)

const executeQuery = async () => {
  showNativeSQL.value = true
  if (sqlQuery.value === '') {
    switch (kind.value) {
      case ExtensionKind.ExtensionKindElasticsearch:
      case ExtensionKind.ExtensionKindEtcd:
      case ExtensionKind.ExtensionKindRedis:
        if (sqlQuery.value === '') {
          sqlQuery.value = '*'
        }
        break
    }
  }
  return executeWithQuery(sqlQuery.value)
}
const executeWithQuery = async (sql: string) => {
  let success = false
  query.value.store = store.value
  query.value.key = queryDataMeta.value.currentDatabase
  query.value.sql = sql

  try {
    globalLoading.value = true
    const data = await API.DataQueryAsync(query.value, () => {
      globalLoading.value = false
    });
    switch (kind.value) {
      case ExtensionKind.ExtensionKindORM:
      case ExtensionKind.ExtensionKindCassandra:
      case ExtensionKind.ExtensionKindIotDB:
      case ExtensionKind.ExtensionKindOpengeMini:
      case ExtensionKind.ExtensionKindElasticsearch:
        ormDataHandler(data)
        success = true
        break;
      case ExtensionKind.ExtensionKindEtcd:
      case ExtensionKind.ExtensionKindRedis:
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
  query.value.offset = Number(query.value.limit) + Number(query.value.offset)
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

Magic.AdvancedKeys([{
  Keys: ['Ctrl+E', 'Ctrl+Enter'],
  Func: executeQuery,
  Description: 'Execute query'
}, {
  Keys: ['Ctrl+Shift+O'],
  Func: () => {
    showNativeSQL.value = false
    executeWithQuery(sqlQuery.value)
  },
  Description: 'Execute query without showing native SQL'
}])
</script>

<template>
  <div>
    <div class="page-header">
      <span class="page-title">{{t('title.dataManager')}}</span>
    </div>
    <el-container style="height: calc(100vh - 80px);" v-loading="globalLoading">
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
            <el-row :gutter="10" v-if="kind === 'atest-store-elasticsearch' || kind === 'atest-store-redis'">
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
                <el-button type="primary" @click="nextPage">{{ t("button.next-page") }}</el-button>
              </el-col>
            </el-row>
          </el-form>
            <Codemirror
            v-model="sqlQuery"
            v-if="complexEditor"
            style="height: var(--sql-editor-height);"
            :extensions="[sql(sqlConfig)]"
            />
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
