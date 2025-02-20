<script setup lang="ts">
import { ref } from 'vue'
import { API } from './net'
import { Cache } from './cache'
import { ElMessage } from 'element-plus'

const sqlQuery = ref('select * from t_sys_global_config')
const queryResult = ref([])
const columns = ref([])

const executeQuery = async () => {
  Cache.SetCurrentStore('mysql');

  API.DataQuery(sqlQuery.value, (data) => {
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

    queryResult.value = result
    columns.value = Array.from(cols)
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
    <el-form @submit.prevent="executeQuery">
        <el-row :gutter="10">
          <el-col :span="20">
            <el-form-item>
              <el-input v-model="sqlQuery" placeholder="Enter SQL query"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item>
              <el-button type="primary" @click="executeQuery">Execute</el-button>
            </el-form-item>
          </el-col>
        </el-row>
    </el-form>
    <el-table :data="queryResult">
      <el-table-column v-for="col in columns" :key="col" :prop="col" :label="col"></el-table-column>
    </el-table>
  </div>
</template>