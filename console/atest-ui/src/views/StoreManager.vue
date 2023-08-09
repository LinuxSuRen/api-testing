<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Pair } from './types'

const stores = ref([] as Store[])
const dialogVisible = ref(false)
const creatingLoading = ref(false)
const storeFormRef = ref<FormInstance>()
const store = ref({
    kind: {},
    properties: [{
        key: '',
        value: ''
    }]
} as Store)
const createAction = ref(true)
const storeForm = reactive(store)

interface Store {
  name: string
  url: string
  username: string
  password: string
  ready: boolean
  kind: {
    url: string
  }
  properties: Pair[]
}

function loadStores() {
  const requestOptions = {
    method: 'POST',
  }
  fetch('/server.Runner/GetStores', requestOptions)
    .then((response) => response.json())
    .then((e) => {
      stores.value = e.data
    })
    .catch((e) => {
      ElMessage.error('Oops, ' + e)
    })
}
loadStores()

function deleteStore(name: string) {
  const requestOptions = {
    method: 'POST',
    body: JSON.stringify({
      name: name
    })
  }
  fetch('/server.Runner/DeleteStore', requestOptions)
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

function editStore(name: string) {
    dialogVisible.value = true
    stores.value.forEach((e: Store) => {
        if (e.name === name) {
            store.value = e
        }
    })
    createAction.value = false
}

function addStore() {
    dialogVisible.value = true
    createAction.value = true
}

const rules = reactive<FormRules<Store>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }]
})
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean, fields) => {
    if (valid) {
      creatingLoading.value = true

       // remove empty pair
      store.value.properties = store.value.properties.filter(
       (e) => e.key !== ''
      )

      const requestOptions = {
        method: 'POST',
        body: JSON.stringify(store.value)
      }
      
      let api = '/server.Runner/CreateStore'
      if (!createAction.value) {
        api = '/server.Runner/UpdateStore'
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

function updateKeys() {
  const props = store.value.properties
  let lastItem = props[props.length - 1]
  if (lastItem.key !== '') {
    store.value.properties.push({
      key: '',
      value: ''
    })
  }
}
</script>

<template>
    <div>Store Manager</div>
    <div>
        <el-button type="primary" @click="addStore" :icon="Edit">New</el-button>
    </div>
    <el-table :data="stores" style="width: 100%">
      <el-table-column label="Name" width="180">
        <template #default="scope">
          <el-input v-model="scope.row.name" placeholder="Name"/>
        </template>
      </el-table-column>
      <el-table-column label="URL">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <el-input v-model="scope.row.url" placeholder="URL" />
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Plugin">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <el-input v-model="scope.row.kind.url" placeholder="Plugin" />
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Status" width="100">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            <el-text class="mx-1"
            type="success"
            v-if="scope.row.ready"
            >Ready</el-text>
            <el-text class="mx-1"
            type="warning"
            v-if="!scope.row.ready"
            >Not Ready</el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Operations" width="220">
        <template #default="scope">
          <div style="display: flex; align-items: center" v-if="scope.row.name !== 'local'">
            <el-button type="primary" @click="deleteStore(scope.row.name)" :icon="Delete">Delete</el-button>
            <el-button type="primary" @click="editStore(scope.row.name)" :icon="Edit">Edit</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="Create Store" width="30%" draggable>
      <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="storeForm"
          ref="storeFormRef"
          status-icon label-width="120px">
          <el-form-item label="Name" prop="name">
            <el-input v-model="storeForm.name" test-id="store-form-name" />
          </el-form-item>
          <el-form-item label="URL" prop="url">
            <el-input v-model="storeForm.url" placeholder="http://foo" test-id="store-form-url" />
          </el-form-item>
          <el-form-item label="Username" prop="username">
            <el-input v-model="storeForm.username" test-id="store-form-username" />
          </el-form-item>
          <el-form-item label="Password" prop="password">
            <el-input v-model="storeForm.password" type="password" test-id="store-form-password" />
          </el-form-item>
          <el-form-item label="Plugin" prop="plugin">
            <el-input v-model="storeForm.kind.url" test-id="store-form-plugin" />
          </el-form-item>
          <el-form-item label="Properties" prop="properties">
            <el-table :data="storeForm.properties" style="width: 100%">
                <el-table-column label="Key" width="180">
                    <template #default="scope">
                        <el-input v-model="scope.row.key" placeholder="Key" @change="updateKeys"/>
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
              @click="submitForm(storeFormRef)"
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
