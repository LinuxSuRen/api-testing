<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { reactive, ref } from 'vue'
import type { Suite } from '@/views/types'
import { API } from '@/views/net'
import type { FormInstance, FormRules } from 'element-plus'

const { t } = useI18n()
const emit = defineEmits(['created'])
const props = defineProps({
    visible: Boolean,
})

const importSuiteFormRef = ref<FormInstance>()
const importSuiteForm = reactive({
    url: '',
    store: '',
    kind: ''
})

const importSuiteFormRules = reactive<FormRules<Suite>>({
    url: [
        { required: true, message: 'URL is required', trigger: 'blur' },
        { type: 'url', message: 'Should be a valid URL value', trigger: 'blur' }
    ],
    store: [{ required: true, message: 'Location is required', trigger: 'blur' }],
    kind: [{ required: true, message: 'Kind is required', trigger: 'blur' }]
})

const importSuiteFormSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid: boolean) => {
        if (valid) {
            API.ImportTestSuite(importSuiteForm, () => {
                emit('created')
                formEl.resetFields()
            })
        }
    })
}

interface Store {
    name: string,
    description: string,
}
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

const importSourceKinds = [{
  "name": "Postman",
  "value": "postman"
}, {
  "name": "Native",
  "value": "native"
}]
</script>

<template>
    <el-dialog :modelValue="visible" title="Import Test Suite" width="30%" draggable>
        <span>Supported source URL: Postman collection share link</span>
        <template #footer>
          <span class="dialog-footer">
            <el-form
                :rules="importSuiteFormRules"
                :model="importSuiteForm"
                ref="importSuiteFormRef"
                status-icon label-width="120px">
                <el-form-item label="Location" prop="store">
                    <el-select v-model="importSuiteForm.store" class="m-2"
                                test-id="suite-import-form-store"
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
                <el-form-item label="Kind" prop="kind">
                    <el-select v-model="importSuiteForm.kind" class="m-2"
                    filterable=true
                    test-id="suite-import-form-kind"
                    default-first-option=true
                    placeholder="Kind" size="middle">
                    <el-option
                        v-for="item in importSourceKinds"
                        :key="item.name"
                        :label="item.name"
                        :value="item.value"
                    />
                    </el-select>
                </el-form-item>
                <el-form-item label="URL" prop="url">
                    <el-input v-model="importSuiteForm.url" test-id="suite-import-form-api" placeholder="https://api.postman.com/collections/xxx" />
                </el-form-item>
                <el-form-item>
                    <el-button
                        type="primary"
                        @click="importSuiteFormSubmit(importSuiteFormRef)"
                        test-id="suite-import-submit"
                    >{{ t('button.import') }}</el-button>
                </el-form-item>
            </el-form>
          </span>
        </template>
    </el-dialog>
</template>
