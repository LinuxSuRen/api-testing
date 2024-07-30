<script setup lang="ts">
import TestCase from './TestCase.vue'
import TestSuite from './TestSuite.vue'
import TemplateFunctions from './TemplateFunctions.vue'
import { ref, watch } from 'vue'
import { ElTree, ElMessage } from 'element-plus'
import { Edit, Refresh } from '@element-plus/icons-vue'
import { API } from './net'
import { Cache } from './cache'
import { useI18n } from 'vue-i18n'
import { Magic } from './magicKeys'
import TestSuiteCreationDialog from '../components/TestSuiteCreationDialog.vue'
import TestSuiteImportDialog from '../components/TestSuiteImportDialog.vue'
import LoginDialog from '../components/LoginDialog.vue'

const { t } = useI18n()

interface Tree {
  id: string
  label: string
  method: string
  parent: string
  parentID: string
  store: string
  kind: string
  children?: Tree[]
}

const testCaseName = ref('')
const testSuite = ref('')
const testSuiteKind = ref('')
const handleTreeClick = (data: Tree) => {
  if (data.children) {
    Cache.SetCurrentStore(data.store)
    viewName.value = 'testsuite'
    testSuite.value = data.label
    testSuiteKind.value = data.kind
    Cache.SetCurrentStore(data.store)

    API.ListTestCase(data.label, data.store, (d) => {
        if (d.items && d.items.length > 0) {
          data.children = []
          d.items.forEach((item: any) => {
            data.children?.push({
              id: data.label,
              label: item.name,
              method: item.request.method,
              kind: data.kind,
              store: data.store,
              parent: data.label,
              parentID: data.id
            } as Tree)
          })
        }
    })
  } else {
    Cache.SetCurrentStore(data.store)
    Cache.SetLastTestCaseLocation(data.parentID, data.id)
    testCaseName.value = data.label
    testSuite.value = data.parent
    testSuiteKind.value = data.kind
    viewName.value = 'testcase'
  }
}

Magic.Keys((k) => {
  const currentKey = currentNodekey.value

  if (treeRef.value) {
    treeRef.value.data.forEach((n) => {
      if (n.children) {
        n.children.forEach((c, index) => {
          if (c.id === currentKey) {
            var nextIndex = -1
            if (k.endsWith('Up')) {
              if (index > 0) {
                nextIndex = index - 1
              }
            } else {
              if (index < n.children.length - 1) {
                nextIndex = index + 1
              }
            }

            if (nextIndex >= 0 < n.children.length) {
              const next = n.children[nextIndex]
              currentNodekey.value = next.id
              treeRef.value!.setCurrentKey(next.id)
              treeRef.value!.setCheckedKeys([next.id], false)
            }
            return
          }
        })
      }
    })
  }
}, ['Alt+ArrowUp', 'Alt+ArrowDown'])

const treeData = ref([] as Tree[])
const treeRef = ref<InstanceType<typeof ElTree>>()
const currentNodekey = ref('')

function loadTestSuites(storeName: string) {
  const requestOptions = {
    method: 'GET',
    headers: {
      'X-Store-Name': storeName,
      'X-Auth': API.getToken()
    },
  }
  return async () => {
    await fetch('/api/v1/suites', requestOptions)
      .then((response) => response.json())
      .then((d) => {
        if (!d.data) {
          return
        }
        Object.keys(d.data).map((k) => {
          let suite = {
            id: k,
            label: k,
            kind: d.data[k].kind,
            store: storeName,
            children: [] as Tree[]
          } as Tree

          d.data[k].data.forEach((item: any) => {
            suite.children?.push({
              id: generateTestCaseID(k, item),
              label: item,
              store: storeName,
              kind: suite.kind,
              parent: k,
              parentID: suite.id
            } as Tree)
          })
          treeData.value.push(suite)
        })
      })
  }
}

function generateTestCaseID(suiteName: string, caseName: string) {
  return suiteName + caseName
}

interface Store {
  name: string,
  description: string,
}

const loginDialogVisible = ref(false)
const stores = ref([] as Store[])
const storesLoading = ref(false)
function loadStores(lastSuitName?: string, lastCaseName?: string) {
  if (lastSuitName && lastCaseName && lastSuitName !== '' && lastCaseName !== '') {
    // get data from emit event
    Cache.SetLastTestCaseLocation(lastSuitName, generateTestCaseID(lastSuitName, lastCaseName))
  }

  storesLoading.value = true
  const requestOptions = {
    headers: {
      'X-Auth': API.getToken()
    }
  }
  fetch('/api/v1/stores', requestOptions)
    .then(API.DefaultResponseProcess)
    .then(async (d) => {
      stores.value = d.data
      treeData.value = [] as Tree[]
      Cache.SetStores(d.data)

      for (const item of d.data) {
        if (item.ready && !item.disabled) {
          await loadTestSuites(item.name)()
        }
      }

      if (treeData.value.length > 0) {
        const key = Cache.GetLastTestCaseLocation()

        let targetSuite = {} as Tree
        let targetChild = {} as Tree
        if (key.suite !== '' && key.testcase !== '') {
          for (var i = 0; i < treeData.value.length; i++) {
            const item = treeData.value[i]
            if (item.id === key.suite && item.children) {
              for (var j = 0; j < item.children.length; j++) {
                const child = item.children[j]
                if (child.id === key.testcase) {
                  targetSuite = item
                  targetChild = child
                  break
                }
              }
              break
            }
          }
        }
        
        if (!targetChild.id || targetChild.id === '') {
          targetSuite = treeData.value[0]
          if (targetSuite.children && targetSuite.children.length > 0) {
            targetChild = targetSuite.children[0]
          }
        }

        viewName.value = 'testsuite'
        currentNodekey.value = targetChild.id

        treeRef.value!.setCurrentKey(targetChild.id)
        treeRef.value!.setCheckedKeys([targetChild.id], false)

        testSuite.value = targetSuite.label
        Cache.SetCurrentStore(targetSuite.store)
        testSuiteKind.value = targetChild.kind
      } else {
        viewName.value = ""
      }
    }).catch((e) => {
      if(e.message === "Unauthenticated") {
        loginDialogVisible.value = true
      } else {
        ElMessage.error('Oops, ' + e)
      }
    }).finally(() => {
      storesLoading.value = false
    })
}
loadStores()

const dialogVisible = ref(false)
const importDialogVisible = ref(false)

function openTestSuiteCreateDialog() {
  dialogVisible.value = true
}

function openTestSuiteImportDialog() {
  importDialogVisible.value = true
}

const filterText = ref('')
watch(filterText, (val) => {
  treeRef.value!.filter(val)
})
const filterTestCases = (value: string, data: Tree) => {
  if (!value) return true
  return data.label.toLocaleLowerCase().includes(value.toLocaleLowerCase())
}

const viewName = ref('')

</script>

<template>
  <div class="common-layout" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-container style="height: 100%">
      <el-main style="padding-top: 5px; padding-bottom: 5px;">
        <el-container style="height: 100%">
          <el-aside>
            <el-button type="primary" @click="openTestSuiteCreateDialog"
              data-intro="Click here to create a new test suite"
              test-id="open-new-suite-dialog" :icon="Edit">{{ t('button.new') }}</el-button>
            <el-button type="primary" @click="openTestSuiteImportDialog"
              data-intro="Click here to import from Postman"
              test-id="open-import-suite-dialog">{{ t('button.import') }}</el-button>
            <el-button type="primary" @click="loadStores" :icon="Refresh">{{ t('button.refresh') }}</el-button>
            <el-input v-model="filterText" :placeholder="t('tip.filter')" test-id="search" style="padding: 5px;" />

            <el-tree
              v-loading="storesLoading"
              :data=treeData
              highlight-current
              :check-on-click-node="true"
              :expand-on-click-node="false"
              :current-node-key="currentNodekey"
              ref="treeRef"
              node-key="id"
              :filter-node-method="filterTestCases"
              @current-change="handleTreeClick"
              data-intro="This is the test suite tree. You can click the test suite to edit it."
            >
              <template #default="{ node, data }">
                <span class="custom-tree-node">
                  <el-text class="mx-1" v-if="data.method === 'POST'" type="success">{{ node.label }}</el-text>
                  <el-text class="mx-1" v-else-if="data.method === 'PUT'" type="warning">{{ node.label }}</el-text>
                  <el-text class="mx-1" v-else-if="data.method === 'DELETE'" type="danger">{{ node.label }}</el-text>
                  <el-text class="mx-1" v-else>{{ node.label }}</el-text>
                </span>
              </template>
            </el-tree>
            <TemplateFunctions/>
          </el-aside>

          <el-main style="padding-top: 0px; padding-right: 0px; padding-bottom: 0px;">
            <TestCase
              v-if="viewName === 'testcase'"
              :suite="testSuite"
              :kindName="testSuiteKind"
              :name="testCaseName"
              @updated="loadStores"
              style="height: 100%;"
              data-intro="This is the test case editor. You can edit the test case here."
            />
            <TestSuite
              v-else-if="viewName === 'testsuite'"
              :name="testSuite"
              @updated="loadStores"
              data-intro="This is the test suite editor. You can edit the test suite here."
            />
          </el-main>
        </el-container>
      </el-main>
    </el-container>
  </div>

    <TestSuiteCreationDialog
        :visible="dialogVisible"
        @created="dialogVisible=false; loadStores()"/>

    <TestSuiteImportDialog
        :visible="importDialogVisible"
        @created="importDialogVisible=false; loadStores()"/>

    <LoginDialog
        :visible="loginDialogVisible"/>
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
