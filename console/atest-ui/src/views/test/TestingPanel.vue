<template>
  <!-- Index -->
  <div class="index" data-title="Welcome!" data-intro="Welcome to use api-testing! 👋">
    <el-card class="card" shawon="hover" :body-style="{ width: '100%' }">
      <el-container style="height: 100%">
        <el-main style="padding-top: 5px; padding-bottom: 5px">
          <el-container style="height: 100%">
            <el-aside>
              <el-button
                type="primary"
                @click="openTestSuiteCreateDialog"
                data-intro="Click here to create a new test suite"
                test-id="open-new-suite-dialog"
                :icon="Edit"
                >{{ t('button.new') }}</el-button
              >
              <el-button
                type="primary"
                @click="openTestSuiteImportDialog"
                data-intro="Click here to import from Postman"
                test-id="open-import-suite-dialog"
                >{{ t('button.import') }}</el-button
              >
              <el-button type="primary" @click="loadStores" :icon="Refresh">{{
                t('button.refresh')
              }}</el-button>
              <div class="filter-input">
                <el-input v-model="filterText" :placeholder="t('tip.filter')" test-id="search" />
              </div>
              <!-- Test suite tree containers -->
              <el-tree
                v-loading="storesLoading"
                :data="data"
                highlight-current
                :check-on-click-node="true"
                :expand-on-click-node="false"
                :current-node-key="currentNodekey"
                ref="treeRef"
                node-key="id"
                :filter-node-method="filterTestCases"
                @node-click="handleNodeClick"
                data-intro="This is the test suite tree. You can click the test suite to edit it."
              >
                <template #default="{ node, data }">
                  <span>
                    <i :class="getColorClass(data.kind)">{{ data.kind }}</i>
                    <i class="test-suite">{{ node.label }}</i>
                  </span>
                </template>
              </el-tree>
              <TemplateFunctions />
            </el-aside>
            <!-- Test suite and case containers. -->
            <el-main style="padding-top: 0px; padding-right: 0px; padding-bottom: 0px">
              <el-card shadow="hover">
                <TestCase
                  v-if="viewName === 'testcase'"
                  :suite="testSuite"
                  :kindName="testSuiteKind"
                  :name="testCaseName"
                  @updated="loadStores"
                  style="height: 100%"
                  data-intro="This is the test case editor. You can edit the test case here."
                />
                <TestSuite
                  v-else-if="viewName === 'testsuite'"
                  :name="testSuite"
                  @updated="loadStores"
                  data-intro="This is the test suite editor. You can edit the test suite here."
                />
              </el-card>
            </el-main>
          </el-container>
        </el-main>
      </el-container>
    </el-card>
  </div>

  <!-- New Create Btn Dialog -->
  <el-dialog v-model="dialogVisible" :title="t('title.createTestSuite')" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="testSuiteForm"
          ref="suiteFormRef"
          status-icon
          label-width="120px"
        >
          <el-form-item :label="t('field.storageLocation')" prop="store">
            <el-select
              v-model="testSuiteForm.store"
              class="m-2"
              test-id="suite-form-store"
              filterable="true"
              default-first-option="true"
              placeholder="Storage Location"
              size="middle"
            >
              <el-option
                v-for="item in stores"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('field.suiteKind')" prop="kind">
            <el-select
              v-model="testSuiteForm.kind"
              class="m-2"
              filterable="true"
              test-id="suite-form-kind"
              default-first-option="true"
              size="middle"
            >
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
            <el-input
              v-model="testSuiteForm.api"
              placeholder="http://foo"
              test-id="suite-form-api"
            />
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

  <!-- Import Btn Dialog -->
  <el-dialog v-model="importDialogVisible" title="Import Test Suite" width="30%" draggable>
    <span>Supported source URL: Postman collection share link</span>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="importSuiteFormRules"
          :model="importSuiteForm"
          ref="importSuiteFormRef"
          status-icon
          label-width="120px"
        >
          <el-form-item label="Location" prop="store">
            <el-select
              v-model="importSuiteForm.store"
              class="m-2"
              test-id="suite-import-form-store"
              filterable="true"
              default-first-option="true"
              placeholder="Storage Location"
              size="middle"
            >
              <el-option
                v-for="item in stores"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="URL" prop="url">
            <el-input
              v-model="importSuiteForm.url"
              test-id="suite-import-form-api"
              placeholder="https://api.postman.com/collections/xxx"
            />
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

  <!--Login Dialog-->
  <el-dialog v-model="loginDialogVisible" title="You need to login first." width="30%">
    <el-collapse accordion="true">
      <el-collapse-item title="Server in cloud" name="1">
        <a href="/oauth2/token" target="_blank">
          <svg
            height="32"
            aria-hidden="true"
            viewBox="0 0 16 16"
            version="1.1"
            width="32"
            data-view-component="true"
            class="octicon octicon-mark-github v-align-middle color-fg-default"
          >
            <path
              d="M8 0c4.42 0 8 3.58 8 8a8.013 8.013 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38 0-.27.01-1.13.01-2.2 0-.75-.25-1.23-.54-1.48 1.78-.2 3.65-.88 3.65-3.95 0-.88-.31-1.59-.82-2.15.08-.2.36-1.02-.08-2.12 0 0-.67-.22-2.2.82-.64-.18-1.32-.27-2-.27-.68 0-1.36.09-2 .27-1.53-1.03-2.2-.82-2.2-.82-.44 1.1-.16 1.92-.08 2.12-.51.56-.82 1.28-.82 2.15 0 3.06 1.86 3.75 3.64 3.95-.23.2-.44.55-.51 1.07-.46.21-1.61.55-2.33-.66-.15-.24-.6-.83-1.23-.82-.67.01-.27.38.01.53.34.19.73.9.82 1.13.16.45.68 1.31 2.69.94 0 .67.01 1.3.01 1.49 0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8Z"
            ></path>
          </svg>
        </a>
      </el-collapse-item>
      <el-collapse-item title="Server in local" name="2">
        <el-steps :active="deviceAuthActive" finish-status="success">
          <el-step title="Request Device Code" />
          <el-step title="Input Code" />
          <el-step title="Finished" />
        </el-steps>

        <div v-if="deviceAuthActive === 1">
          Open <a :href="deviceAuthResponse.verification_uri" target="_blank">this link</a>, and
          type the code:
          <span>{{ deviceAuthResponse.user_code }}. Then click the next step button.</span>
        </div>
        <el-button style="margin-top: 12px" @click="deviceAuthNext">Next step</el-button>
      </el-collapse-item>
    </el-collapse>
  </el-dialog>
</template>

<script setup lang="ts">
import TestCase from '../../components/test/TestCase.vue'
import TestSuite from '../../components/test/TestSuite.vue'
import TemplateFunctions from '../../components/other/TemplateFunctions.vue'
import { reactive, ref, watch, onMounted } from 'vue'
import { ElTree, ElMessage } from 'element-plus'
import { Edit, Refresh } from '@element-plus/icons-vue'
import { Cache } from '../../utils/cache'
import { useI18n } from 'vue-i18n'
import { DefaultResponseProcess } from '@/api/common'
import { GetStores } from '@/api/store/store'
import type { TestStore, Tree, Suite } from '../../types/types'
import type { FormInstance, FormRules } from 'element-plus'
import { ListTestCase, LoadTestSuite, CreateTestSuite, ImportTestSuite } from '@/api/test/test'
import { da } from 'element-plus/es/locale'

const { t } = useI18n()

const testCaseName = ref('')
const testSuite = ref('')
const testSuiteKind = ref('')
const viewName = ref('')

const handleNodeClick = (data: Tree) => {
  if (data.children) {
    Cache.SetCurrentStore(data.store)
    viewName.value = 'testsuite'
    testSuite.value = data.label
    testSuiteKind.value = data.kind
    Cache.SetCurrentStore(data.store)

    ListTestCase(data.label, data.store)
      .then((res: any) => {
        if (res.items && res.items.length > 0) {
          data.children = []
          res.items.forEach((item: any) => {
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
      .catch((err: any) => {
        ElMessage({
          type: 'error',
          showClose: true,
          message: 'Oops, ' + err.message || 'Unknown error when fetching test case!'
        })
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

const loadTestSuites = async (sn: string) => {
  await LoadTestSuite(sn)
    .then((res: any) => {
      Object.keys(res.data).map((k) => {
        let suite = {
          id: k,
          label: k,
          kind: res.data[k].kind,
          store: sn,
          children: [] as Tree[]
        } as Tree

        res.data[k].data.forEach((item: any) => {
          suite.children?.push({
            id: k + item,
            label: item,
            store: sn,
            kind: suite.kind,
            parent: k,
            parentID: suite.id
          } as Tree)
        })
        data.value.push(suite)
      })
    })
    .catch((err: any) => {
      ElMessage({
        type: 'error',
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when fetching test suite!'
      })
    })
}

const loginDialogVisible = ref(false)
const stores = ref([] as TestStore[])
const storesLoading = ref(false)

const loadStores = async () => {
  storesLoading.value = true
  await GetStores()
    .then(async (res: any) => {
      
      stores.value = res.data
      data.value = [] as Tree[]
      data.value = res.data.slice(1)
      Cache.SetStores(res.data)

      for (const item of res.data) {
        if (item.ready && !item.disabled) {
          await loadTestSuites(item.name)
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
        Cache.SetCurrentStore(targetSuite.store)
        testSuiteKind.value = targetChild.kind
      } else {
        viewName.value = ''
      }
    })
    .catch((err: any) => {
      if (err.message === 'Unauthenticated') {
        loginDialogVisible.value = true
      } else {
        ElMessage({
          type: 'error',
          showClose: true,
          message: 'Oops, ' + err.message || 'Unknown error when fetching test suite!'
        })
      }
    })
    .finally(() => {
      storesLoading.value = false
    })
}

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

const openTestSuiteCreateDialog = () => {
  dialogVisible.value = true
}

const openTestSuiteImportDialog = () => {
  importDialogVisible.value = true
}

const rules = reactive<FormRules<Suite>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  store: [{ required: true, message: 'Location is required', trigger: 'blur' }]
})

// submit new test case.
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) {
    // formEl is undefined or null
    return
  }

  await formEl.validate((valid: boolean) => {
    if (valid) {
      suiteCreatingLoading.value = true
      CreateTestSuite(testSuiteForm)
        .then((_: any) => {
          suiteCreatingLoading.value = false
          loadStores()
          dialogVisible.value = false
          formEl.resetFields()
        })
        .catch((err: any) => {
          suiteCreatingLoading.value = false
          ElMessage({
            type: 'error',
            showClose: true,
            message: 'Oops, ' + err.message || 'Unknown error creating test suite!'
          })
        })
    }
  })
}

onMounted(() => {
  loadStores()
})

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

      ImportTestSuite(importSuiteForm)
        .then((res: any) => {
          if (res.code === 200) {
            loadStores()
            importDialogVisible.value = false
            formEl.resetFields()
          }
        })
        .catch((err: any) => {
          ElMessage({
            type: 'error',
            showClose: true,
            message: 'Oops, ' + err.message || 'Unknown error importing test suite!'
          })
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
      .then(DefaultResponseProcess)
      .then((d) => {
        deviceAuthResponse.value = d
      })
  } else if (deviceAuthActive.value === 2) {
    window.location.href =
      '/oauth2/getUserInfoFromLocalCode?device_code=' + deviceAuthResponse.value.device_code
  }
}

const suiteKinds = [
  {
    name: 'HTTP',
    color: 'blue-text'
  },
  {
    name: 'gRPC',
    color: 'green-text'
  },
  {
    name: 'tRPC',
    color: 'orange-text'
  }
]

const getColorClass = (kind: string) => {
  const suiteKind = suiteKinds.find((suite) => suite.name.toLowerCase() === kind.toLowerCase())
  return suiteKind ? suiteKind.color : 'other-text'
}
</script>

<style scoped>
.index {
  display: flex;
  height: 75vh;
}

@media (max-width: 768px) {
  .index {
    height: 50vh;
  }
}

.blue-text {
  font-size: small;
  font-style: small;
  color: blue;
}

.other-text {
  font-size: small;
  font-style: small;
  color: black;
}

.green-text {
  font-size: small;
  font-style: small;
  color: green;
}

.orange-text {
  font-size: small;
  font-style: small;
  color: orange;
}

.card {
  display: flex;
  margin-top: 1%;
  width: 100%;
  max-width: 1750px;
  height: auto;
  vertical-align: middle;
}

.filter-input {
  vertical-align: middle;
  padding-top: 1vh;
  padding-bottom: 1vh;
}

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

.test-suite {
  padding-left: 10%;
  font-style: normal;
  font-size: small;
}

/* Adjust el-tree style */
::v-deep .el-tree-node__expand-icon {
  color: rgb(64, 158, 255);
  font-size: 20px;
  margin-right: 10px;
  content: '▶';
  transform: rotate(0deg);
}

::v-deep .el-tree-node__expand-icon.is-leaf {
  color: transparent;
  content: '';
}

::v-deep .is-expanded .el-tree-node__expand-icon {
  transform: rotate(90deg);
}
</style>
