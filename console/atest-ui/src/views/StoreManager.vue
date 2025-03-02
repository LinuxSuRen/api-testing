<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { reactive, ref, watch } from 'vue'
import { Edit, Delete, QuestionFilled, Help } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { Pair } from './types'
import { API } from './net'
import { UIAPI } from './net-vue'
import { SupportedExtensions, SupportedExtension } from './store'
import { useI18n } from 'vue-i18n'
import { Magic } from './magicKeys'

const { t } = useI18n()

const emptyStore = function() {
  return {
    name: '',
    url: '',
    username: '',
    password: '',
    kind: {
        name: '',
        url: ''
    },
    properties: [{
        key: '',
        value: '',
        description: '',
    }],
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

const storesLoading = ref(false)
function loadStores() {
  storesLoading.value = true
  API.GetStores((e) => {
    stores.value = e.data
  }, (e) => {
    ElMessage.error('Oops, ' + e)
  }, () => {
    storesLoading.value = false
  })
}
loadStores()
Magic.Keys(loadStores, ['Alt+KeyR'])

function deleteStore(name: string) {
  API.DeleteStore(name, (e) => {
    ElMessage({
      message: 'Deleted.',
      type: 'success'
    })
    loadStores()
  }, (e) => {
    ElMessage.error('Oops, ' + e)
  })
}

function editStore(name: string) {
    dialogVisible.value = true
    stores.value.forEach((e: Store) => {
        if (e.name === name) {
            setStoreForm(e)
            return
        }
    })
    createAction.value = false
}

function setStoreForm(store: Store) {
    storeForm.name = store.name
    storeForm.url = store.url
    storeForm.username = store.username
    storeForm.password = store.password
    storeForm.kind = store.kind
    storeForm.disabled = store.disabled
    storeForm.readonly = store.readonly
    storeForm.properties = store.properties
}

function addStore() {
    setStoreForm(emptyStore())
    dialogVisible.value = true
    createAction.value = true
}
Magic.Keys(addStore, ['Alt+KeyN'])

const rules = reactive<FormRules<Store>>({
  name: [{ required: true, message: 'Name is required', trigger: 'blur' }],
  url: [{ required: true, message: 'URL is required', trigger: 'blur' }]
})
const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid: boolean) => {
    if (valid) {
      UIAPI.CreateOrUpdateStore(storeForm, createAction.value, () => {
          loadStores()
          dialogVisible.value = false
          formEl.resetFields()
      }, creatingLoading)
    }
  })
}

watch(() => storeForm.kind.name, (name) => {
  const ext = SupportedExtension(name)
  if (ext) {
    storeExtLink.value = ext.link
    let pro = storeForm.properties.slice()

    for (var i = 0; i < pro.length;) {
      // remove it if the value or key is empty
      if (pro[i].key === '' || pro[i].value === '') {
        pro.splice(i, 1)
      } else {
        i++
      }
    }

    // add extension related params
    ext.params.forEach(p => {
      const index = pro.findIndex(e => e.key === p.key)
      if (index === -1) {
        pro.push({
          key: p.key,
          value: '',
          defaultValue: p.defaultValue,
          description: p.description,
          type: p.type,
          enum: p.enum
        } as Pair)
      } else {
        pro[index].description = p.description
      }
    })

    // make sure there is always a empty pair for letting users input
    pro.push({
      key: '',
      value: ''
    } as Pair)
    storeForm.properties = pro
  }
})
watch(storeForm, (e) => {
  if (e.kind.name === '') {
    if (e.url.startsWith('https://github.com') || e.url.startsWith('https://gitee.com')) {
      e.kind.name = 'atest-store-git'
    }
  }
})

function storeVerify(formEl: FormInstance | undefined) {
  if (!formEl) return
  
  API.VerifyStore(storeForm.name, (e) => {
    if (e.ready) {
      ElMessage({
        message: 'Verified!',
        type: 'success'
      })
    } else {
      ElMessage.error(e.message)
    }
  }, (e) => {
    ElMessage.error(e.message)
  })
}

function updateKeys() {
  const props = storeForm.properties
  if (props.findIndex(p => p.key === '') === -1) {
    storeForm.properties.push({
      key: '',
      value: ''
    } as Pair)
  }
}

const storeExtLink = ref('')
</script>

<template>
    <div>Store Manager</div>
    <div>
        <el-button type="primary" @click="addStore" :icon="Edit">{{t('button.new')}}</el-button>
        <el-button type="primary" @click="loadStores">{{t('button.refresh')}}</el-button>
    </div>
    <el-table :data="stores" style="width: 100%" v-loading=storesLoading>
      <el-table-column :label="t('field.name')" width="180">
        <template #default="scope">
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column label="URL">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            {{ scope.row.url }}
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('field.plugin')">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            {{ scope.row.kind.name }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="Socket">
        <template #default="scope">
          <div style="display: flex; align-items: center">
            {{ scope.row.kind.url }}
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('field.status')" width="100">
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
      <el-table-column :label="t('field.operations')" width="220">
        <template #default="scope">
          <div style="display: flex; align-items: center" v-if="scope.row.name !== 'local'">
            <el-button type="primary" @click="deleteStore(scope.row.name)" :icon="Delete">{{t('button.delete')}}</el-button>
            <el-button type="primary" @click="editStore(scope.row.name)" :icon="Edit">{{t('button.edit')}}</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top: 20px; margin-bottom: 20px; position: absolute; bottom: 0px;">
      Follow <el-link href="https://linuxsuren.github.io/api-testing/#storage" target="_blank">the instructions</el-link> to configure the storage plugins.
    </div>

    <el-dialog v-model="dialogVisible" :title="t('title.createStore')" width="30%" draggable>
      <template #footer>
      <span class="dialog-footer">
        <el-form
          :rules="rules"
          :model="storeForm"
          ref="storeFormRef"
          status-icon label-width="120px">
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
            >
              <el-option
                v-for="item in SupportedExtensions()"
                :key="item.name"
                :label="item.name"
                :value="item.name"
              />
            </el-select>
            <el-icon v-if="storeExtLink && storeExtLink !== ''">
              <el-link :href="storeExtLink" target="_blank">
                <Help />
              </el-link>
            </el-icon>
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
                        <el-input v-model="scope.row.key" placeholder="Key" @change="updateKeys"/>
                    </template>
                </el-table-column>
                <el-table-column label="Value">
                    <template #default="scope">
                      <div style="display: flex; align-items: center">
                          <el-select v-model="scope.row.value" v-if="scope.row.enum">
                            <el-option
                                v-for="item in scope.row.enum"
                                :key="item"
                                :label="item"
                                :value="item"
                            />
                          </el-select>
                          <el-input-number v-model="scope.row.value" v-else-if="scope.row.type === 'number'"/>
                            <el-input v-model="scope.row.value" :placeholder="scope.row.defaultValue" v-else/>
                            <el-tooltip :content="scope.row.description" v-if="scope.row.description">
                              <el-icon>
                                <QuestionFilled/>
                              </el-icon>
                            </el-tooltip>
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
              >{{t('button.verify')}}</el-button
            >
            <el-button
              type="primary"
              @click="submitForm(storeFormRef)"
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
