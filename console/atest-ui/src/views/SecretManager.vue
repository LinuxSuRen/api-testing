<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const secrets = ref([] as Secret[])
const dialogVisible = ref(false)
const creatingLoading = ref(false)
const secretFormRef = ref<FormInstance>()
const secret = ref({} as Secret)
const createAction = ref(true)
const secretForm = reactive(secret)

interface Secret {
  Name: string
  Value: string
}

function loadStores() {
  const requestOptions = {
    method: 'POST',
  }
  fetch('/server.Runner/GetSecrets', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      secrets.value = e.data
    })
    .catch((e) => {
      ElMessage.error('Oops, ' + e)
    })
}
loadStores()

function deleteSecret(name: string) {
  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
      name: name
    })
  }
  fetch('/server.Runner/DeleteSecret', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      ElMessage({
        message: 'Deleted.',
        type: 'success'
      })
      loadStores()
    })
    .catch((e) => {
      ElMessage.error('Oops, ' + e)
    })
}

function editSecret(name: string) {
    dialogVisible.value = true
    secrets.value.forEach((e: Secret) => {
        if (e.Name === name) {
            secret.value = e
        }
    })
    createAction.value = false
}

function addSecret() {
    dialogVisible.value = true
    createAction.value = true
}

const rules = reactive<FormRules<Secret>>({
  Name: [{ required: true, message: 'Name is required', trigger: 'blur' }]
})
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      creatingLoading.value = true

      const requestOptions = {
        method: 'POST',
        body: JSON.stringify(secret.value)
      }
      
      let api = '/server.Runner/CreateSecret'
      if (!createAction.value) {
        api = '/server.Runner/UpdateSecret'
      }

      fetch(api, requestOptions)
        .then((response) => response.json())
        .then(() => {
          creatingLoading.value = false
          loadStores()
          dialogVisible.value = false
          formEl.resetFields()
        })
    }
  })
}

</script>

<template>
    <div>Secret Manager</div>
    <div>
        <el-button type="primary" @click="addSecret" :icon="Edit">New</el-button>
    </div>
    <el-table :data="secrets" style="width: 100%">
      <el-table-column label="Name" width="180">
        <template #default="scope">
          <el-text class="mx-1">{{ scope.row.Name }}</el-text>
        </template>
      </el-table-column>
      <el-table-column label="Operations" width="220">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <el-button type="primary" @click="deleteSecret(scope.row.Name)" :icon="Delete">Delete</el-button>
            <el-button type="primary" @click="editSecret(scope.row.Name)" :icon="Edit">Edit</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="Create Secret" width="30%" draggable>
      <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="secretForm"
          ref="secretFormRef"
          status-icon label-width="120px">
          <el-form-item label="Name" prop="Name">
            <el-input v-model="secretForm.Name" test-id="secret-form-name" />
          </el-form-item>
          <el-form-item label="Password" prop="Value">
            <el-input v-model="secretForm.Value" type="password" test-id="secret-form-password" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="submitForm(secretFormRef)"
              :loading="creatingLoading"
              test-id="store-form-submit"
              >Submit</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
