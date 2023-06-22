<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import HelloWorld from './components/HelloWorld.vue'
import { ref } from 'vue'
import type { TabsPaneContext } from 'element-plus'
import * as grpcWeb from 'grpc-web'
import {RunnerClient} from './ServerServiceClientPb'
import {TestTask} from './server_pb'

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
]

const activeName = ref('first')

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event)
}

interface Tree {
  label: string
  children?: Tree[]
}

const handleNodeClick = (data: Tree) => {
  console.log(data)
}

const data: Tree[] = [
  {
    label: 'Suite-1',
    children: [
      {
        label: 'userList',
      },
      {
        label: 'userEdit',
      },
    ],
  }
]

const defaultProps = {
  children: 'children',
  label: 'label',
}

interface User {
  key: string
  value: string
}

let tableData: User[] = [
  {
    key: 'name',
    value: 'Tom',
  },
  {
    key: 'gender',
    value: 'male',
  },
  {
    key: '',
    value: '',
  }
]

function change() {
  let lastItem = tableData[tableData.length-1]
  if (lastItem.key !== '') {
    tableData.push({
      key:'',
      value:''
    })
  }
}
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="200px">
        <el-tree :data="data" :props="defaultProps"
          :default-expanded-keys="['1']"
          @node-click="handleNodeClick" />
      </el-aside>
      

      <el-main>
        <el-header>
            <el-select v-model="value" class="m-2" placeholder="Method" size="large">
              <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
            <el-input v-model="input" placeholder="API Address" />

            <el-button type="primary">Send</el-button>
        </el-header>

        <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
          <el-tab-pane label="Params" name="first">
            
            <el-table :data="tableData" style="width: 100%">
              <el-table-column label="Key" width="180">
                <template #default="scope">
                  <el-input v-model="scope.row.key" placeholder="Key" @change="change"/>  
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
          <el-tab-pane label="Headers" name="second">Config</el-tab-pane>
          <el-tab-pane label="Body" name="third">

            <el-radio-group v-model="radio">
              <el-radio :label="3">None</el-radio>
              <el-radio :label="6">raw</el-radio>
              <el-radio :label="9">form-data</el-radio>
            </el-radio-group>

          </el-tab-pane>
          <el-tab-pane label="Verify" name="fourth">Task</el-tab-pane>
        </el-tabs>
      </el-main>
    </el-container>
  </div>

</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
