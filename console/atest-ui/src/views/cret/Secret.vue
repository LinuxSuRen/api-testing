<template>
  <el-card class="card" shadow="hover">
    <div class="cert-header">
      <div class="index">
        <h3>{{ t('title.secretManager') }}</h3>
        <el-button type="primary" @click="addSecret" :icon="Edit">{{ t('button.new') }}</el-button>
      </div>
      <div>
        Follow
        <el-link
          class="link"
          href="https://linuxsuren.github.io/api-testing/#secret-server"
          target="_blank"
          >the instructions</el-link
        >
        to configure the secret server.
      </div>
    </div>

    <el-card class="tables-container">
      <el-table :data="secrets" style="width: 100%">
        <el-table-column :label="t('field.name')" width="180">
          <template #default="scope">
            <el-text class="mx-1">{{ scope.row.Name }}</el-text>
          </template>
        </el-table-column>
        <el-table-column :label="t('field.operations')" width="220">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-button type="primary" @click="deleteSecret(scope.row.Name)" :icon="Delete">{{
                t('button.delete')
              }}</el-button>
              <el-button type="primary" @click="editSecret(scope.row.Name)" :icon="Edit">{{
                t('button.edit')
              }}</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </el-card>

  <el-dialog v-model="dialogVisible" :title="t('title.createSecret')" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="secretForm"
          ref="secretFormRef"
          status-icon
          label-width="120px"
        >
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
              :loading="creatingLoading"
              test-id="store-form-submit"
              >{{ t('button.submit') }}</el-button
            >
          </el-form-item>
        </el-form>
      </span>
    </template>
  </el-dialog>
</template>
    
<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { GetSecrets, DeleteSecret, CreateSecret, UpdateSecret } from '@/api/cert/cert'

const { t } = useI18n()

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

onMounted(() => {
  loadSecrets()
})

const loadSecrets = () => {
  GetSecrets()
    .then((res: any) => {
      secrets.value = res.data
    })
    .catch((err: any) => {
      ElMessage({
        type: 'error',
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when fetching secret!'
      })
    })
}

const deleteSecret = (name: string) => {
  DeleteSecret({ name: name })
    .then((res: any) => {
      ElMessage({
        showClose: true,
        message: 'Deleted.',
        type: 'success'
      })
      loadSecrets()
    })
    .catch((err: any) => {
      ElMessage({
        type: 'error',
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when deleting secret!'
      })
    })
}

const editSecret = (name: string) => {
  dialogVisible.value = true
  secrets.value.forEach((e: Secret) => {
    if (e.Name === name) {
      secret.value = e
    }
  })
  createAction.value = false
}

const addSecret = () => {
  dialogVisible.value = true
  createAction.value = true
}

const rules = reactive<FormRules<Secret>>({
  Name: [{ required: true, message: 'Name is required', trigger: 'blur' }]
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean) => {
    if (valid) {
      CreateSecret({
        payload: secret.value,
        create: createAction.value
      })
        .then((res: any) => {
          loadSecrets()
          dialogVisible.value = false
          formEl.resetFields()
        })
        .catch((err: any) => {
          ElMessage({
            type: 'error',
            showClose: true,
            message: 'Oops, ' + err.message || 'Unknown error creating secret!'
          })
        })
      creatingLoading
    }
  })
}
</script>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  margin-top: 1%;
  width: 100%;
  max-width: 1750px;
  height: auto;
  vertical-align: middle;
}

.tables-container {
  margin-top: 1%;
}

h3 {
  display: inline-flex;
  margin-right: 2%;
  vertical-align: middle;
}

.index {
  display: flex;
  width: 30%;
}

.cert-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.link {
  color: #409eff;
  font-style: italic;
}
</style>  
