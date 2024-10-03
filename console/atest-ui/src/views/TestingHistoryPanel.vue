<script setup lang="ts">
import TestCase from './TestCase.vue'
import TestSuite from './TestSuite.vue'
import TemplateFunctions from './TemplateFunctions.vue'
import { reactive, ref, watch } from 'vue'
import { ElTree, ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Edit, Refresh } from '@element-plus/icons-vue'
import type { Suite } from './types'
import { API } from './net'
import { Cache } from './cache'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

interface Tree {
  id: string
  label: string
  parent: string
  parentID: string
  store: string
  kind: string
  suiteLabel: string
  children?: Tree[]
}

const props = defineProps({
  ID: String,
})
const testCaseName = ref('')
const testSuite = ref('')
const testKind = ref('')
const historySuiteName = ref('')
const historyCaseID = ref('')

const handleNodeClick = (data: Tree) => {
  if (data.children) {
    Cache.SetCurrentStore(data.store)
    viewName.value = 'testsuite'
    historySuiteName.value = data.label
    Cache.SetCurrentStore(data.store)

    API.ListTestCase(data.label, data.store, (d) => {
      if (d.items && d.items.length > 0) {
        data.children = []
        d.items.forEach((item: any) => {
          data.children?.push({
            id: data.label,
            label: item.name,
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
    historySuiteName.value = data.parent
    testSuite.value = data.suiteLabel
    testCaseName.value = data.label
    historyCaseID.value = data.id
    testKind.value = data.kind
    viewName.value = 'testcase'
  }
}

const treeData = ref([] as Tree[])
const treeRef = ref<InstanceType<typeof ElTree>>()
const currentNodekey = ref('')

function loadHistoryTestSuites(storeName: string) {
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Store-Name': storeName,
      'X-Auth': API.getToken()
    },
  }
  return async () => {
    await fetch('/api/v1/historySuites', requestOptions)
      .then((response) => response.json())
      .then((d) => {
        if (!d.data) {
          return
        }
        const sortedKeys = Object.keys(d.data).sort((a, b) => new Date(a) - new Date(b));
        sortedKeys.map((k) => {
          let suite = {
            id: k,
            label: k,
            store: storeName,
            children: [] as Tree[]
          } as Tree

          d.data[k].data.forEach((item: any) => {
            suite.children?.push({
              id: item.ID,
              label: item.testcase,
              suiteLabel: item.suite,
              store: storeName,
              kind: item.kind,
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
          await loadHistoryTestSuites(item.name)()
        }
      }

      if (treeData.value.length > 0) {
        const key = Cache.GetLastTestCaseLocation()

        let targetSuite = {} as Tree
        let targetChild = {} as Tree

        const targetID = props.ID
        if (targetID && targetID !== '') {
          for (const suite of treeData.value) {
            if (suite.children) {
              const foundChild = suite.children.find(child => child.id === targetID)
              if (foundChild) {
                targetSuite = suite
                targetChild = foundChild
                handleNodeClick(targetChild)
                updateTreeSelection(targetSuite, targetChild)
                return
              }
            }
          }
        } else {
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
        }
        if (!targetChild.id || targetChild.id === '') {
          targetSuite = treeData.value[0]
          if (targetSuite.children && targetSuite.children.length > 0) {
            targetChild = targetSuite.children[0]
          }
        }

        viewName.value = 'testsuite'
        updateTreeSelection(targetSuite, targetChild)
      } else {
        viewName.value = ""
      }
    }).catch((e) => {
      if (e.message === "Unauthenticated") {
        loginDialogVisible.value = true
      } else {
        ElMessage.error('Oops, ' + e)
      }
    }).finally(() => {
      storesLoading.value = false
    })
}
loadStores()

function updateTreeSelection(targetSuite: Tree, targetChild: Tree) {
    currentNodekey.value = targetChild.id

    treeRef.value!.setCurrentKey(targetChild.id)
    treeRef.value!.setCheckedKeys([targetChild.id], false)

    testSuite.value = targetSuite.label
    Cache.SetCurrentStore(targetSuite.store)
    testKind.value = targetChild.kind
}

const filterText = ref('')
watch(filterText, (val) => {
  treeRef.value!.filter(val)
})
const filterTestCases = (value: string, data: Tree) => {
  if (!value) return true
  return data.label.includes(value)
}

const viewName = ref('')
</script>

<template>
  <div class="common-layout" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-container style="height: 100%">
      <el-main style="padding-top: 5px; padding-bottom: 5px;">
        <el-container style="height: 100%">
          <el-aside>
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
              @node-click="handleNodeClick"
              data-intro="This is the test history tree. You can click the history test to browse it."
            />
            <TemplateFunctions />
          </el-aside>

          <el-main style="padding-top: 0px; padding-right: 0px; padding-bottom: 0px;">
            <TestCase v-if="viewName === 'testcase'" :suite="testSuite" :kindName="testKind" :name="testCaseName"
              :historySuiteName="historySuiteName" :historyCaseID="historyCaseID" @updated="loadStores" style="height: 100%;"
              data-intro="This is the test case editor. You can edit the test case here." />
          </el-main>
        </el-container>
      </el-main>
    </el-container>
  </div>
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

.demo-tabs>.el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
