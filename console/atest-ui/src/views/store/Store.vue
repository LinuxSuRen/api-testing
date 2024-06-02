<template>
  <el-card class="card" shadow="hover">
    <div class="cert-header">
      <div class="index">
        <h3>{{ t('title.storeManager') }}</h3>
        <el-button type="primary" @click="addStore" :icon="Edit">{{ t('button.new') }}</el-button>
        <el-button type="primary" @click="loadStores">{{ t('button.refresh') }}</el-button>
      </div>
      <div>
        Follow
        <el-link
          class="link"
          href="https://linuxsuren.github.io/api-testing/#storage"
          target="_blank"
          >the instructions</el-link
        >
        to configure the storage plugins.
      </div>
    </div>

    <el-card class="tables-container">
      <el-table :data="stores" style="width: 100%" v-loading="storesLoading">
        <el-table-column :label="t('field.name')" width="180">
          <template #default="scope">
            <el-input v-model="scope.row.name" placeholder="Name" />
          </template>
        </el-table-column>
        <el-table-column label="URL">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-input v-model="scope.row.url" placeholder="URL" />
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('field.plugin')">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-input v-model="scope.row.kind.url" placeholder="Plugin" />
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('field.status')" width="100">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <el-text class="mx-1" type="success" v-if="scope.row.ready">Ready</el-text>
              <el-text class="mx-1" type="warning" v-if="!scope.row.ready">Not Ready</el-text>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="t('field.operations')" width="220">
          <template #default="scope">
            <div style="display: flex; align-items: center" v-if="scope.row.name !== 'local'">
              <el-button type="primary" @click="deleteStore(scope.row.name)" :icon="Delete">{{
                t('button.delete')
              }}</el-button>
              <el-button type="primary" @click="editStore(scope.row.name)" :icon="Edit">{{
                t('button.edit')
              }}</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </el-card>

  <el-dialog v-model="dialogVisible" :title="t('title.createStore')" width="30%" draggable>
    <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="storeForm"
          ref="storeFormRef"
          status-icon
          label-width="120px"
        >
          <el-form-item :label="t('field.name')" prop="name">
            <el-input v-model="storeForm.name" test-id="store-form-name" />
          </el-form-item>
          <el-form-item label="URL" prop="url">
            <el-input v-model="storeForm.url" placeholder="http://foo" test-id="store-form-url" />
          </el-form-item>
          <el-form-item :label="t('field.username')" prop="username">
            <el-input v-model="storeForm.username" test-id="store-form-username" />
          </el-form-item>
          <el-form-item :label="t('field.password')" prop="password">
            <el-input v-model="storeForm.password" type="password" test-id="store-form-password" />
          </el-form-item>
          <el-form-item :label="t('field.pluginName')" prop="pluginName">
            <el-select
              v-model="storeForm.kind.name"
              test-id="store-form-plugin-name"
              class="m-2"
              size="middle"
            >
              <el-option
                v-for="item in SupportedExtensions()"
                :key="item.value"
                :label="item.key"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="t('field.pluginURL')" prop="plugin">
            <el-input v-model="storeForm.kind.url" test-id="store-form-plugin" />
          </el-form-item>
          <el-form-item :label="t('field.disabled')" prop="disabled">
            <el-switch v-model="storeForm.disabled" />
          </el-form-item>
          <el-form-item :label="t('field.properties')" prop="properties">
            <el-table :data="storeForm.properties" style="width: 100%">
              <el-table-column label="Key" width="180">
                <template #default="scope">
                  <el-input v-model="scope.row.key" placeholder="Key" @change="updateKeys" />
                </template>
              </el-table-column>
              <el-table-column label="Value">
                <template #default="scope">
                  <div style="display: flex; align-items: center">
                    <el-input v-model="scope.row.value" placeholder="Value" />
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              @click="storeVerify(storeFormRef)"
              test-id="store-form-verify"
              >{{ t('button.verify') }}</el-button
            >
            <el-button
              type="primary"
              @click="submitForm(storeFormRef)"
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
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref, watch } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import { SupportedExtensions } from '../../types/store'
import type { FormInstance, FormRules } from 'element-plus'
import type { Pair } from '../../types/types'
import { GetStores, CreateStore, UpdateStore, DeleteStore, VerifyStore } from '@/api/store/store'

const { t } = useI18n()

const emptyStore = function () {
  return {
    name: '',
    url: '',
    username: '',
    password: '',
    kind: {
      name: '',
      url: ''
    },
    properties: [
      {
        key: '',
        value: ''
      }
    ],
    disabled: false,
    readonly: false
  } as Store
}
const stores = ref([] as Store[])
const dialogVisible = ref(false)
const creatingLoading = ref(false)
const storeFormRef = ref<FormInstance>()
const createAction = ref(true)
const storeForm = reactive(emptyStore())

interface Store {
  name: string
  owner: string
  url: string
  username: string
  password: string
  ready: boolean
  disabled: boolean
  readonly: boolean
  kind: {
    name: string
    url: string
  }
  properties: Pair[]
}

onMounted(() => {
  loadStores()
})

const storesLoading = ref(false)
const loadStores = () => {
  storesLoading.value = true
  GetStores()
    .then((res: any) => {
      stores.value = res.data
    })
    .catch((err: any) => {
      ElMessage({
        type: 'error',
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when fetching store!'
      })
    })
  storesLoading.value = false
}

const deleteStore = (name: string) => {
  DeleteStore({ name: name })
    .then((_: any) => {
      ElMessage({
        showClose: true,
        message: 'Deleted!',
        type: 'success'
      })
      loadStores()
    })
    .catch((err: any) => {
      ElMessage({
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when deleting store!',
        type: 'error'
      })
    })
}

const editStore = (name: string) => {
  dialogVisible.value = true
  stores.value.forEach((e: Store) => {
    if (e.name === name) {
      setStoreForm(e)
      return
    }
  })
  createAction.value = false
}

const setStoreForm = (store: Store) => {
  storeForm.name = store.name
  storeForm.url = store.url
  storeForm.username = store.username
  storeForm.password = store.password
  storeForm.kind = store.kind
  storeForm.disabled = store.disabled
  storeForm.readonly = store.readonly
  storeForm.properties = store.properties
  storeForm.properties.push({
    key: '',
    value: ''
  })
}

const addStore = () => {
  setStoreForm(emptyStore())
  dialogVisible.value = true
  createAction.value = true
}

const rules = reactive<FormRules<Store>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  url: [{ required: true, message: 'URL is required', trigger: 'blur' }],
  'kind.name': [{ required: true, message: 'Plugin is required', trigger: 'blur' }]
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean) => {
    if (valid) {
      CreateStore({
        disabled: storeForm.disabled,
        kind: {
          name: storeForm.kind.name,
          url: storeForm.kind.url
        },
        name: storeForm.name,
        password: storeForm.password,
        properties: {
          ...storeForm.properties
        },
        readonly: storeForm.readonly,
        url: storeForm.url,
        username: storeForm.username
      })
        .then((res: any) => {
          loadStores()
          dialogVisible.value = false
          formEl.resetFields()
        })
        .catch((err: any) => {
          ElMessage({
            showClose: true,
            message: 'Oops, ' + err.message || 'Unknown error when creating store!',
            type: 'error'
          })
          creatingLoading
        })
    }
  })
}

watch(storeForm, (e) => {
  if (e.kind.name === '') {
    if (e.url.startsWith('https://github.com') || e.url.startsWith('https://gitee.com')) {
      e.kind.name = 'atest-store-git'
    }
  }
})

const storeVerify = (formEl: FormInstance | undefined) => {
  if (!formEl) {
    return
  }

  VerifyStore({ name: storeForm.name })
    .then((_: any) => {
      ElMessage({
        showClose: true,
        message: 'Verified!',
        type: 'success'
      })
    })
    .catch((err: any) => {
      ElMessage({
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when verifying store!',
        type: 'error'
      })
    })
}

const updateKeys = () => {
  const props = storeForm.properties
  let lastItem = props[props.length - 1]
  if (lastItem.key !== '') {
    storeForm.properties.push({
      key: '',
      value: ''
    })
  }
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

.tables-container {
  margin-top: 1%;
}
</style>