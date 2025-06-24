<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { API } from './net'
import type { Secret } from './net'
import { UIAPI } from './net-vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const secrets = ref([] as Secret[])
const dialogVisible = ref(false)
const creatingLoading = ref(false)
const secretFormRef = ref<FormInstance>()
const secret = ref({} as Secret)
const createAction = ref(true)
const secretForm = reactive(secret)

function loadSecrets() {
  API.GetSecrets((e) => {
    secrets.value = e.data
  }, UIAPI.ErrorTip)
}
loadSecrets()

function deleteSecret(name: string) {
  API.DeleteSecret(name, () => {
      ElMessage({
        message: 'Deleted.',
        type: 'success'
      })
      loadSecrets()
    }, UIAPI.ErrorTip)
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

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean) => {
    if (valid) {
      UIAPI.CreateOrUpdateSecret(secretForm.value, createAction.value, () => {
          loadSecrets()
          dialogVisible.value = false
          formEl.resetFields()
      }, creatingLoading)
    }
  })
}

</script>

<template>
    <div class="page-header">
      <span class="page-title">{{t('title.secretManager')}}</span>
        <el-button type="primary" @click="addSecret" :icon="Edit">{{t('button.new')}}</el-button>
    </div>
    <el-table :data="secrets" style="width: 100%">
      <el-table-column :label="t('field.name')">
        <template #default="scope">
          <el-text class="mx-1">{{ scope.row.Name }}</el-text>
        </template>
      </el-table-column>
      <el-table-column :label="t('field.operations')" width="220">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <el-button type="primary" @click="deleteSecret(scope.row.Name)" :icon="Delete">{{t('button.delete')}}</el-button>
            <el-button type="primary" @click="editSecret(scope.row.Name)" :icon="Edit">{{t('button.edit')}}</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 20px; margin-bottom: 20px; position: absolute; bottom: 0px;">
      Follow <el-link href="https://linuxsuren.github.io/api-testing/#secret-server" target="_blank">the instructions</el-link> to configure the secret server.
    </div>

    <el-dialog v-model="dialogVisible" :title="t('title.createSecret')" width="30%" draggable>
      <template #footer>
      <span class="dialog-footer">
        <el-form
          :model="secretForm"
          ref="secretFormRef"
          status-icon label-width="120px">
          <el-form-item :label="t('field.name')" prop="Name">
            <el-input v-model="secretForm.Name" test-id="secret-form-name" />
          </el-form-item>
          <el-form-item :label="t('field.password')" prop="Value">
            <el-input v-model="secretForm.Value" type="password" test-id="secret-form-password" />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="submitForm(secretFormRef)"
              v-loading="creatingLoading"
              test-id="store-form-submit"
              >{{t('button.submit')}}</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
