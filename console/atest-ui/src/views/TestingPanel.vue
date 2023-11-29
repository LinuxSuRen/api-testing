<script setup lang="ts">
import WelcomePage from './WelcomePage.vue'
import TestCase from './TestCase.vue'
import TestSuite from './TestSuite.vue'
import TemplateFunctions from './TemplateFunctions.vue'
import { reactive, ref, watch } from 'vue'
import { ElTree, ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'
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
      'X-Store-Name': storeName,
      'X-Auth': API.getToken()
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

const loginDialogVisible = ref(false)
const stores = ref([] as Store[])
const storesLoading = ref(false)
function loadStores() {
  storesLoading.value = true
  const requestOptions = {
    method: 'POST',
    headers: {
      'X-Auth': API.getToken()
    }
  }
  fetch('/server.Runner/GetStores', requestOptions)
    .then(API.DefaultResponseProcess)
    .then(async (d) => {
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
  await formEl.validate((valid: boolean) => {
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
        }, (e) => {
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
  await formEl.validate((valid: boolean) => {
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

const deviceAuthActive = ref(0)
const deviceAuthResponse = ref({
  user_code: '',
  verification_uri: '',
  device_code: ''
})
const deviceAuthNext = () => {
  if (deviceAuthActive.value++ > 2) {
    return
  }

  if (deviceAuthActive.value === 1) {
    fetch('/oauth2/getLocalCode')
    .then(API.DefaultResponseProcess)
    .then((d) => {
      deviceAuthResponse.value = d
    })
  } else if (deviceAuthActive.value === 2) {
    window.location.href = '/oauth2/getUserInfoFromLocalCode?device_code=' + deviceAuthResponse.value.device_code
  }
}

const suiteKinds = [{
  "name": "HTTP",
}, {
  "name": "gRPC",
}, {
  "name": "tRPC",
}]

</script>

<template>
  <div class="common-layout" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-container style="height: 100%">
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
            <TestCase
              v-if="viewName === 'testcase'"
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
          </el-main>
        </el-container>
      </el-main>

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
              test-id="suite-form-kind"
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

  <el-dialog
    v-model="loginDialogVisible"
    title="You need to login first."
    width="30%"
  >
    <el-collapse accordion="true">
      <el-collapse-item title="Server in cloud" name="1">
        <a href="/oauth2/token" target="_blank">
          <svg height="32" aria-hidden="true" viewBox="0 0 16 16" version="1.1" width="32" data-view-component="true" class="octicon octicon-mark-github v-align-middle color-fg-default">
              <path d="M8 0c4.42 0 8 3.58 8 8a8.013 8.013 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38 0-.27.01-1.13.01-2.2 0-.75-.25-1.23-.54-1.48 1.78-.2 3.65-.88 3.65-3.95 0-.88-.31-1.59-.82-2.15.08-.2.36-1.02-.08-2.12 0 0-.67-.22-2.2.82-.64-.18-1.32-.27-2-.27-.68 0-1.36.09-2 .27-1.53-1.03-2.2-.82-2.2-.82-.44 1.1-.16 1.92-.08 2.12-.51.56-.82 1.28-.82 2.15 0 3.06 1.86 3.75 3.64 3.95-.23.2-.44.55-.51 1.07-.46.21-1.61.55-2.33-.66-.15-.24-.6-.83-1.23-.82-.67.01-.27.38.01.53.34.19.73.9.82 1.13.16.45.68 1.31 2.69.94 0 .67.01 1.3.01 1.49 0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8Z"></path>
          </svg>
        </a>
      </el-collapse-item>
      <el-collapse-item title="Server in local" name="2">
        <el-steps :active="deviceAuthActive" finish-status="success">
          <el-step title="Request Device Code" />
          <el-step title="Input Code"/>
          <el-step title="Finished" />
        </el-steps>

        <div v-if="deviceAuthActive===1">
          Open <a :href="deviceAuthResponse.verification_uri" target="_blank">this link</a>, and type the code: <span>{{ deviceAuthResponse.user_code }}. Then click the next step button.</span>
        </div>
        <el-button style="margin-top: 12px" @click="deviceAuthNext">Next step</el-button>
      </el-collapse-item>
    </el-collapse>
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
