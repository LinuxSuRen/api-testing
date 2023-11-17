<script setup lang="ts">
import WelcomePage from './views/WelcomePage.vue'
import TestCase from './views/TestCase.vue'
import TestSuite from './views/TestSuite.vue'
import StoreManager from './views/StoreManager.vue'
import SecretManager from './views/SecretManager.vue'
import TemplateFunctions from './views/TemplateFunctions.vue'
import { reactive, ref, watch } from 'vue'
import { ElTree, ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Edit, Share } from '@element-plus/icons-vue'
import type { Suite } from './types'
import { API } from './views/net'
import { Cache } from './views/cache'
import { useI18n } from 'vue-i18n'
import ClientMonitor from 'skywalking-client-js'
import { name, version } from '../package'

import setAsDarkTheme from './theme'

const { t } = useI18n()

const asDarkMode = ref(Cache.GetPreference().darkTheme)
setAsDarkTheme(asDarkMode.value)
watch(asDarkMode, Cache.WatchDarkTheme)
watch(asDarkMode, () => {
  setAsDarkTheme(asDarkMode.value)
})

interface Tree {
  id: string
  label: string
  parent: string
  parentID: string
  store: string
  kind: string
  children?: Tree[]
}

const testCaseName = ref('')
const testSuite = ref('')
const testSuiteKind = ref('')
const handleNodeClick = (data: Tree) => {
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
  return async () => {
    await fetch('/server.Runner/GetSuites', requestOptions)
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
              id: k + item,
              label: item,
              store: storeName,
              kind: suite.kind,
              parent: k,
              parentID: suite.id
            } as Tree)
          })
          data.value.push(suite)
        })
      })
  }
}

interface Store {
  name: string,
  description: string,
}

const stores = ref([] as Store[])
const storesLoading = ref(false)
function loadStores() {
  storesLoading.value = true
  const requestOptions = {
    method: 'POST',
  }
  fetch('/server.Runner/GetStores', requestOptions)
    .then((response) => response.json())
    .then(async (d) => {
      storesLoading.value = false
      stores.value = d.data
      data.value = [] as Tree[]
      Cache.SetStores(d.data)

      for (const item of d.data) {
        if (item.ready && !item.disabled) {
          await loadTestSuites(item.name)()
        }
      }

      if (data.value.length > 0) {
        const key = Cache.GetLastTestCaseLocation()

        let targetSuite = {} as Tree
        let targetChild = {} as Tree
        if (key.suite !== '' && key.testcase !== '') {
          for (var i = 0; i < data.value.length; i++) {
            const item = data.value[i]
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
          targetSuite = data.value[0]
          if (targetSuite.children && targetSuite.children.length > 0) {
            targetChild = targetSuite.children[0]
          }
        }

        viewName.value = 'testsuite'
        currentNodekey.value = targetChild.id
        treeRef.value!.setCurrentKey(targetChild.id)
        treeRef.value!.setCheckedKeys([targetChild.id], false)
        testSuite.value = targetSuite.label
        Cache.SetCurrentStore(targetSuite.store )
        testSuiteKind.value = targetChild.kind
      } else {
        viewName.value = ""
      }
    })
}
loadStores()

const dialogVisible = ref(false)
const importDialogVisible = ref(false)
const suiteCreatingLoading = ref(false)
const suiteFormRef = ref<FormInstance>()
const testSuiteForm = reactive({
  name: '',
  api: '',
  store: '',
  kind: ''
})
const importSuiteFormRef = ref<FormInstance>()
const importSuiteForm = reactive({
  url: '',
  store: ''
})

function openTestSuiteCreateDialog() {
  dialogVisible.value = true
}

function openTestSuiteImportDialog() {
  importDialogVisible.value = true
}

const rules = reactive<FormRules<Suite>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  store: [{ required: true, message: 'Location is required', trigger: 'blur' }]
})
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      suiteCreatingLoading.value = true

      API.CreateTestSuite(testSuiteForm, (e) => {
          suiteCreatingLoading.value = false
          if (e.error !== "") {
            ElMessage.error('Oops, ' + e.error)
          } else {
            loadStores()
            dialogVisible.value = false
            formEl.resetFields()
          }
        })
        .catch((e) => {
          suiteCreatingLoading.value = false
          ElMessage.error('Oops, ' + e)
        })
    }
  })
}

const importSuiteFormRules = reactive<FormRules<Suite>>({
  url: [
    { required: true, message: 'URL is required', trigger: 'blur' },
    { type: 'url', message: 'Should be a valid URL value', trigger: 'blur' }
  ],
  store: [{ required: true, message: 'Location is required', trigger: 'blur' }]
})
const importSuiteFormSubmit = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      suiteCreatingLoading.value = true

      API.ImportTestSuite(importSuiteForm, () => {
        loadStores()
        importDialogVisible.value = false
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

const viewName = ref('')
watch(viewName, (val) => {
  ClientMonitor.setPerformance({
    service: name,
    serviceVersion: version,
    pagePath: val,
    useFmp: true,
    enableSPA: true,
    customTags: [{
      key: 'theme', value: asDarkMode.value ? 'dark' : 'light'
    }, {
      key: 'store', value: Cache.GetCurrentStore().name
    }]
  });
})

const suiteKinds = [{
  "name": "HTTP",
}, {
  "name": "gRPC",
}, {
  "name": "tRPC",
}]

const appVersion = ref('')
const appVersionLink = ref('https://github.com/LinuxSuRen/api-testing')
API.GetVersion((d) => {
  appVersion.value = d.message
  const version = d.message.match('^v\\d*.\\d*.\\d*')
  const dirtyVersion = d.message.match('^v\\d*.\\d*.\\d*-\\d*-g')

  if (!version && !dirtyVersion) {
    return
  }

  console.log(dirtyVersion)
  if (dirtyVersion && dirtyVersion.length > 0) {
    appVersionLink.value = appVersionLink.value + '/commit/' + d.message.replace(dirtyVersion[0], '')
  } else if (version && version.length > 0) {
    appVersionLink.value = appVersionLink.value + '/releases/tag/' + version[0]
  }
})
</script>

<template>
  <div class="common-layout" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-container style="height: 100%">
      <el-header style="height: 30px;justify-content: flex-end;">
        <el-button type="primary" :icon="Edit" @click="viewName = 'secret'" data-intro="Manage the secrets."/>
        <el-button type="primary" :icon="Share" @click="viewName = 'store'" data-intro="Manage the store backends." />
        <el-form-item label="Dark Mode" style="margin-left:20px;">
          <el-switch type="primary" data-intro="Switch light and dark modes" v-model="asDarkMode" />
        </el-form-item>
      </el-header>

      <el-main>
        <el-container style="height: 100%">
          <el-aside>
            <el-button type="primary" @click="openTestSuiteCreateDialog"
              data-intro="Click here to create a new test suite"
              test-id="open-new-suite-dialog" :icon="Edit">{{ t('button.new') }}</el-button>
            <el-button type="primary" @click="openTestSuiteImportDialog"
              data-intro="Click here to import from Postman"
              test-id="open-import-suite-dialog">{{ t('button.import') }}</el-button>
            <el-input v-model="filterText" placeholder="Filter keyword" test-id="search" />

            <el-tree
              v-loading="storesLoading"
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
            <WelcomePage
              v-if="viewName === ''"
            />
            <TestCase
              v-else-if="viewName === 'testcase'"
              :suite="testSuite"
              :kindName="testSuiteKind"
              :name="testCaseName"
              @updated="loadStores"
              data-intro="This is the test case editor. You can edit the test case here."
            />
            <TestSuite
              v-else-if="viewName === 'testsuite'"
              :name="testSuite"
              @updated="loadStores"
              data-intro="This is the test suite editor. You can edit the test suite here."
            />
            <StoreManager
            v-else-if="viewName === 'store'"
            />
            <SecretManager
            v-else-if="viewName === 'secret'"
            />
          </el-main>
        </el-container>
      </el-main>
    <div style="position: absolute; bottom: 0px; right: 10px;">
      <a :href=appVersionLink target="_blank" rel="noopener">{{appVersion}}</a>
    </div>
    <TemplateFunctions/>
    </el-container>
  </div>

  <el-dialog v-model="dialogVisible" :title="t('title.createTestSuite')" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="testSuiteForm"
          ref="suiteFormRef"
          status-icon label-width="120px">
          <el-form-item :label="t('field.storageLocation')" prop="store">
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
          <el-form-item :label="t('field.suiteKind')" prop="kind">
            <el-select v-model="testSuiteForm.kind" class="m-2"
              filterable=true
              default-first-option=true
              size="middle">
              <el-option
                v-for="item in suiteKinds"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('field.name')" prop="name">
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
              >{{ t('button.submit') }}</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>

  <el-dialog v-model="importDialogVisible" title="Import Test Suite" width="30%" draggable>
    <span>Supported source URL: Postman collection share link</span>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="importSuiteFormRules"
          :model="importSuiteForm"
          ref="importSuiteFormRef"
          status-icon label-width="120px">
          <el-form-item label="Location" prop="store">
            <el-select v-model="importSuiteForm.store" class="m-2"
              test-id="suite-import-form-store"
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
          <el-form-item label="URL" prop="url">
            <el-input v-model="importSuiteForm.url" test-id="suite-import-form-api" placeholder="https://api.postman.com/collections/xxx" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="importSuiteFormSubmit(importSuiteFormRef)"
              test-id="suite-import-submit"
              >{{ t('button.import') }}</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
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
