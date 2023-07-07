<script setup lang="ts">
import TestCase from './views/TestCase.vue'
import TestSuite from './views/TestSuite.vue'
import { reactive, ref, watch } from 'vue'
import { ElTree } from "element-plus"
import type { FormInstance } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'

interface Tree {
  id: string
  label: string
  parent: string
  children?: Tree[]
}

const testCaseName = ref('')
const testSuite = ref('')
const handleNodeClick = (data: Tree) => {
  if (data.children) {
    viewName.value = "testsuite"
    testSuite.value = data.label

    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            name: data.label,
        })
    };
    fetch('/server.Runner/ListTestCase', requestOptions)
        .then(response => response.json())
        .then(d => {
          if (d.items && d.items.length > 0) {
            data.children=[]
            d.items.forEach((item: any) => {
              data.children?.push({
                id: data.label+item.name,
                label: item.name,
                parent: data.label,
              } as Tree)
            })
          }
        })
  } else {
    testCaseName.value = data.label
    testSuite.value = data.parent
    viewName.value = "testcase"
  }
}

const data = ref([] as Tree[])
const treeRef = ref<InstanceType<typeof ElTree>>()
const currentNodekey = ref('')

function loadTestSuites() {
  const requestOptions = {
      method: 'POST'
  };
  fetch('/server.Runner/GetSuites', requestOptions)
      .then(response => response.json())
      .then(d => {
        data.value = [] as Tree[]
        if (!d.data) {
          return
        }
        Object.keys(d.data).map(k => {
          let suite = {
            id: k,
            label: k,
            children: [] as Tree[],
          } as Tree

          d.data[k].data.forEach((item: any) => {
            suite.children?.push({
              id: k+item,
              label: item,
              parent: k,
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
          
          viewName.value = "testsuite"
          testSuite.value = firstItem.label
        }
      });
}
loadTestSuites()

const dialogVisible = ref(false)
const suiteCreatingLoading = ref(false)
const suiteFormRef = ref<FormInstance>()
const testSuiteForm = reactive({
  name: "",
  api: "",
})

function openTestSuiteCreateDialog() {
  dialogVisible.value = true
}

const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  suiteCreatingLoading.value = true

  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
        name: testSuiteForm.name,
        api: testSuiteForm.api,
    })
  };

  fetch('/server.Runner/CreateTestSuite', requestOptions)
      .then(response => response.json())
      .then(() => {
        suiteCreatingLoading.value = false
        loadTestSuites()
      });
      
  dialogVisible.value = false
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
  <div class="common-layout">
    <el-container style="height: 100vh">
      <el-aside width="200px">
        <el-button type="primary" @click="openTestSuiteCreateDialog" :icon="Edit">New</el-button>
        <el-input v-model="filterText" placeholder="Filter keyword" />

        <el-tree :data="data"
          highlight-current
          :check-on-click-node=true
          :expand-on-click-node=false
          :current-node-key="currentNodekey"
          ref="treeRef"
          node-key="id"
          :filter-node-method="filterTestCases"
          @node-click="handleNodeClick" />
      </el-aside>

      <el-main>
        <TestCase v-if="viewName === 'testcase'" :suite="testSuite" :name="testCaseName" @updated="loadTestSuites"/>
        <TestSuite v-else-if="viewName === 'testsuite'" :name="testSuite" @updated="loadTestSuites"/>
      </el-main>
    </el-container>
  </div>

  <el-dialog v-model="dialogVisible" title="Create Test Suite" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          ref="suiteFormRef"
          status-icon
          label-width="120px"
          class="demo-ruleForm"
        >
          <el-form-item label="Name" prop="name">
            <el-input v-model="testSuiteForm.name" />
          </el-form-item>
          <el-form-item label="API" prop="api">
            <el-input v-model="testSuiteForm.api" placeholder="http://foo" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm(suiteFormRef)" :loading="suiteCreatingLoading">Submit</el-button>
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
