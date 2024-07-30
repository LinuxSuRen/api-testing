<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { reactive, ref } from 'vue'
import type { Suite } from '@/views/types'
import { ElMessage } from 'element-plus'
import { API } from '@/views/net'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()

interface Store {
    name: string,
    description: string,
}

const emit = defineEmits(['created'])
const props = defineProps({
    visible: Boolean,
})

const stores = ref([] as Store[])
function loadStores() {
    const requestOptions = {
        headers: {
            'X-Auth': API.getToken()
        }
    }
    fetch('/api/v1/stores', requestOptions)
        .then(API.DefaultResponseProcess)
        .then(async (d) => {
            stores.value = d.data
        })
}
loadStores()

const suiteCreatingLoading = ref(false)
const suiteFormRef = ref<FormInstance>()
const testSuiteForm = reactive({
    name: '',
    api: '',
    store: '',
    kind: ''
})
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
                    formEl.resetFields()
                }

                emit('created')
            }, (e) => {
                suiteCreatingLoading.value = false
                ElMessage.error('Oops, ' + e)
            })
        }
    })
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
    <el-dialog :modelValue="visible" :title="t('title.createTestSuite')" width="30%" draggable>
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
                           filterable
                           default-first-option
                           placeholder="Storage Location">
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
                           filterable
                           test-id="suite-form-kind"
                           default-first-option>
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
</template>
