<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'

const props = defineProps({
    name: String,
})
const emit = defineEmits(['updated'])

interface Suite {
    name: string;
    api: string;
}

const suite = ref({} as Suite)
function load() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            name: props.name,
        })
    };
    fetch('/server.Runner/GetTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            suite.value = {
                name: e.name,
                api: e.api,
            } as Suite
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}
load()
watch(props, () => {
    load()
})

function save() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify(suite.value),
    };
    fetch('/server.Runner/UpdateTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            ElMessage({
                    message: 'Updated.',
                    type: 'success',
                })
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}

const dialogVisible = ref(false)
const testcaseFormRef = ref<FormInstance>()
const testCaseForm = reactive({
    suiteName: "",
    name: "",
    api: "",
})
function openNewTestCaseDialog() {
    loadTestSuites()
    dialogVisible.value = true
}

function loadTestSuites() {
  const requestOptions = {
      method: 'POST'
  };
  fetch('/server.Runner/GetSuites', requestOptions)
      .then(response => response.json())
      .then(d => {
        Object.keys(d.data).map(k => {
          testSuiteList.value.push(k)
        })
      });
}

const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  suiteCreatingLoading.value = true

  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
        suiteName: testCaseForm.suiteName,
        data: {
            name: testCaseForm.name,
            request: {
                api: testCaseForm.api,
                method: "GET",
            }
        },
    })
  };

  fetch('/server.Runner/CreateTestCase', requestOptions)
      .then(response => response.json())
      .then(() => {
        suiteCreatingLoading.value = false
        emit('updated', 'hello from child')
      });
      
  dialogVisible.value = false
}

function del() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            name: props.name,
        })
    };
    fetch('/server.Runner/DeleteTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            ElMessage({
                    message: 'Deleted.',
                    type: 'success',
                })
            emit('updated')
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}

const suiteCreatingLoading = ref(false)
const testSuiteList = ref([])
</script>

<template>
    <div class="common-layout">
        <el-text class="mx-1" type="primary">{{suite.name}}</el-text>
        <el-input class="mx-1" v-model="suite.api" placeholder="API"></el-input>

        <el-button type="primary" @click="save">Save</el-button>
        <el-button type="primary" @click="del">Delete</el-button>

        <el-button type="primary" @click="openNewTestCaseDialog" :icon="Edit">New TestCase</el-button>
    </div>

  <el-dialog v-model="dialogVisible" title="Create Test Case" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          ref="testcaseFormRef"
          status-icon
          label-width="120px"
          class="demo-ruleForm"
        >
          <el-form-item label="Suite" prop="suite">
            <el-select class="m-2" v-model="testCaseForm.suiteName" placeholder="Select" size="large">
                <el-option
                v-for="item in testSuiteList"
                :key="item"
                :label="item"
                :value="item"
                />
            </el-select>
          </el-form-item>
          <el-form-item label="Name" prop="name">
            <el-input v-model="testCaseForm.name" />
          </el-form-item>
          <el-form-item label="API" prop="api">
            <el-input v-model="testCaseForm.api" placeholder="http://foo" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm(testcaseFormRef)" :loading="suiteCreatingLoading">Submit</el-button>
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
