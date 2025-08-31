<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { reactive, ref } from 'vue'
import { API } from '@/views/net'
import { ElMessage } from 'element-plus'
import type { ImportSource } from '@/views/net'
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
    kind: '',
    data: ''
} as ImportSource)

const importSuiteFormRules = reactive<FormRules<ImportSource>>({
    url: [
        { required: importSuiteForm.kind !== 'native-inline', message: 'URL is required', trigger: 'blur' },
        { type: 'url', message: 'Should be a valid URL value', trigger: 'blur' }
    ],
    data: [{ required: importSuiteForm.kind === 'native-inline', message: 'Data is required', trigger: 'blur' }],
    store: [{ required: true, message: 'Location is required', trigger: 'blur' }],
    kind: [{ required: true, message: 'Kind is required', trigger: 'blur' }]
})

const importSuiteFormSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid: boolean) => {
        if (valid) {
            if (importSuiteForm.kind === 'native-inline') {
                importSuiteForm.kind = 'native'
            }
            API.ImportTestSuite(importSuiteForm, () => {
                emit('created')
                formEl.resetFields()
            }, (e) => {
                ElMessage.error(e)
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
    fetch('/api/v1/stores?kind=store', requestOptions)
        .then(API.DefaultResponseProcess)
        .then(async (d) => {
            stores.value = d.data
        })
}
loadStores()

const importSourceKinds = [{
    "name": "Postman",
    "value": "postman",
    "description": "https://api.postman.com/collections/xxx"
}, {
    "name": "Native",
    "value": "native",
    "description": "http://your-server/api/v1/suites/xxx/yaml?x-store-name=xxx"
}, {
    "name": "Native-Inline",
    "value": "native-inline",
    "description": "Native test suite content in YAML format"
}]
const importSourceDesc = ref("")
const kindChanged = (e) => {
    importSourceKinds.forEach(k => {
        if (k.value === e) {
            importSourceDesc.value = k.description
        }
    });
}
</script>

<template>
    <el-dialog :modelValue="visible" title="Import Test Suite" width="40%"
        draggable destroy-on-close>
        <el-form
            :rules="importSuiteFormRules"
            :model="importSuiteForm"
            ref="importSuiteFormRef"
            status-icon label-width="85px">
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
                    filterable
                    @change="kindChanged"
                    test-id="suite-import-form-kind"
                    default-first-option
                    placeholder="Kind">
                    <el-option
                        v-for="item in importSourceKinds"
                        :key="item.name"
                        :label="item.name"
                        :value="item.value"
                    />
                </el-select>
            </el-form-item>
            <el-form-item label="Data" prop="data" v-if="importSuiteForm.kind === 'native-inline'">
                <el-input v-model="importSuiteForm.data"
                    class="full-width" type="textarea"
                    :placeholder="importSourceDesc" />
            </el-form-item>
            <el-form-item label="URL" prop="url" v-else>
                <el-input v-model="importSuiteForm.url" test-id="suite-import-form-api"
                    class="full-width"
                    :placeholder="importSourceDesc" />
            </el-form-item>
            <el-form-item>
                <el-button
                    type="primary"
                    @click="importSuiteFormSubmit(importSuiteFormRef)"
                    test-id="suite-import-submit"
                >{{ t('button.import') }}</el-button>
            </el-form-item>
        </el-form>
    </el-dialog>
</template>

<style scoped>
.full-width {
    width: 100%;
}
</style>
