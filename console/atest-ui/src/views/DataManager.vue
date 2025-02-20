<script setup lang="ts">
import { ref } from 'vue'
import { API } from './net'
import { Cache } from './cache'

const sqlQuery = ref('')
const queryResult = ref([])

const executeQuery = async () => {
  Cache.SetCurrentStore('mysql');

    API.DataQuery(sqlQuery.value, (data) => {
        queryResult.value = data.data
    })
}
</script>

<template>
  <div>
    <el-form @submit.prevent="executeQuery">
      <el-form-item>
        <el-input v-model="sqlQuery" placeholder="Enter SQL query"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="executeQuery">Execute</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="queryResult">
      <el-table-column v-for="(value, key) in queryResult[0]" :key="key" :prop="key" :label="key"></el-table-column>
    </el-table>
  </div>
</template>