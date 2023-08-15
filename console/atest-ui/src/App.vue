<script setup lang="ts">
import TestCase from './views/TestCase.vue'
import TestSuite from './views/TestSuite.vue'
import StoreManager from './views/StoreManager.vue'
import TemplateFunctions from './views/TemplateFunctions.vue'
import { reactive, ref, watch } from 'vue'
import { ElTree } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Edit, Share } from '@element-plus/icons-vue'
import type { Suite } from './types'

interface Tree {
  id: string
  label: string
  parent: string
  store: string
  children?: Tree[]
}

const testCaseName = ref('')
const testSuite = ref('')
const store = ref('')
const handleNodeClick = (data: Tree) => {
  if (data.children) {
    viewName.value = 'testsuite'
    testSuite.value = data.label
    store.value = data.store

    const requestOptions = {
      method: 'POST',
      headers: {
        'X-Store-Name': data.store
      },
      body: JSON.stringify({
        name: data.label
      })
    }
    fetch('/server.Runner/ListTestCase', requestOptions)
      .then((response) => response.json())
      .then((d) => {
        if (d.items && d.items.length > 0) {
          data.children = []
          d.items.forEach((item: any) => {
            data.children?.push({
              id: data.label + item.name,
              label: item.name,
              store: data.store,
              parent: data.label
            } as Tree)
          })
        }
      })
  } else {
    testCaseName.value = data.label
    testSuite.value = data.parent
    store.value = data.store
    viewName.value = 'testcase'
  }
}

const data = ref([] as Tree[])
const treeRef = ref<InstanceType<typeof ElTree>>()
const currentNodekey = ref('')

function loadTestSuites(storeName: string) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': storeName
    },
  }
  fetch('/server.Runner/GetSuites', requestOptions)
    .then((response) => response.json())
    .then((d) => {
      if (!d.data) {
        return
      }
      Object.keys(d.data).map((k) => {
        let suite = {
          id: k,
          label: k,
          store: storeName,
          children: [] as Tree[]
        } as Tree

        d.data[k].data.forEach((item: any) => {
          suite.children?.push({
            id: k + item,
            label: item,
            store: storeName,
            parent: k
          } as Tree)
        })
        data.value.push(suite)
      })

      if (data.value.length > 0) {
        const firstItem = data.value[0]
        if (firstItem.children && firstItem.children.length > 0) {
          const child = firstItem.children[0].id

          currentNodekey.value = child
          treeRef.value!.setCurrentKey(child)
          treeRef.value!.setCheckedKeys([child], false)
        }

        viewName.value = 'testsuite'
        testSuite.value = firstItem.label
        store.value = firstItem.store
      }
    })
}

interface Store {
  name: string,
  description: string,
}

const stores = ref([] as Store[])
function loadStores() {
  const requestOptions = {
    method: 'POST',
  }
  fetch('/server.Runner/GetStores', requestOptions)
    .then((response) => response.json())
    .then((d) => {
      stores.value = d.data
      data.value = [] as Tree[]

      d.data.forEach((item: any) => {
        if (item.ready) {
          loadTestSuites(item.name)
        }
      })
    })
}
loadStores()

const dialogVisible = ref(false)
const suiteCreatingLoading = ref(false)
const suiteFormRef = ref<FormInstance>()
const testSuiteForm = reactive({
  name: '',
  api: '',
  store: ''
})

function openTestSuiteCreateDialog() {
  dialogVisible.value = true
}

const rules = reactive<FormRules<Suite>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  store: [{ required: true, message: 'Location is required', trigger: 'blur' }]
})
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  console.log(formEl)
  await formEl.validate((valid: boolean, fields) => {
    console.log(valid, fields)
    if (valid) {
      suiteCreatingLoading.value = true

      const requestOptions = {
        method: 'POST',
        headers: {
          'X-Store-Name': testSuiteForm.store
        },
        body: JSON.stringify({
          name: testSuiteForm.name,
          api: testSuiteForm.api
        })
      }

      fetch('/server.Runner/CreateTestSuite', requestOptions)
        .then((response) => response.json())
        .then(() => {
          suiteCreatingLoading.value = false
          loadStores()
          dialogVisible.value = false
          formEl.resetFields()
        })
    }
  })
}

const filterText = ref('')
watch(filterText, (val) => {
  treeRef.value!.filter(val)
})
const filterTestCases = (value: string, data: Tree) => {
  if (!value) return true
  return data.label.includes(value)
}

const viewName = ref('testcase')
</script>

<template>
  <div class="common-layout" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-container style="height: 100%">
      <el-header style="height: 30px;justify-content: flex-end;">
        <el-button type="primary" :icon="Share" @click="viewName = ''" />
      </el-header>

      <el-main>
        <el-container style="height: 100%">
          <el-aside width="200px">
            <el-button type="primary" @click="openTestSuiteCreateDialog"
              data-intro="Click here to create a new test suite"
              test-id="open-new-suite-dialog" :icon="Edit">New</el-button>
            <el-input v-model="filterText" placeholder="Filter keyword" test-id="search" />

            <el-tree
              :data=data
              highlight-current
              :check-on-click-node="true"
              :expand-on-click-node="false"
              :current-node-key="currentNodekey"
              ref="treeRef"
              node-key="id"
              :filter-node-method="filterTestCases"
              @node-click="handleNodeClick"
              data-intro="This is the test suite tree. You can click the test suite to edit it."
            />
          </el-aside>

          <el-main>
            <TestCase
              v-if="viewName === 'testcase'"
              :store="store"
              :suite="testSuite"
              :name="testCaseName"
              @updated="loadStores"
              data-intro="This is the test case editor. You can edit the test case here."
            />
            <TestSuite
              v-else-if="viewName === 'testsuite' && testSuite !== ''"
              :name="testSuite"
              :store="store"
              @updated="loadStores"
              data-intro="This is the test suite editor. You can edit the test suite here."
            />
            <StoreManager
            v-else-if="viewName === '' || testSuite === '' || store === ''"
            />
          </el-main>
        </el-container>
      </el-main>
    </el-container>
  </div>

  <el-dialog v-model="dialogVisible" title="Create Test Suite" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="testSuiteForm"
          ref="suiteFormRef"
          status-icon label-width="120px">
          <el-form-item label="Location" prop="store">
            <el-select v-model="testSuiteForm.store" class="m-2"
              test-id="suite-form-store"
              filterable=true
              default-first-option=true
              placeholder="Storage Location" size="middle">
              <el-option
                v-for="item in stores"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="Name" prop="name">
            <el-input v-model="testSuiteForm.name" test-id="suite-form-name" />
          </el-form-item>
          <el-form-item label="API" prop="api">
            <el-input v-model="testSuiteForm.api" placeholder="http://foo" test-id="suite-form-api" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="submitForm(suiteFormRef)"
              :loading="suiteCreatingLoading"
              test-id="suite-form-submit"
              >Submit</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>

  <TemplateFunctions/>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}
.common-layout {
  height: 100%;
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
