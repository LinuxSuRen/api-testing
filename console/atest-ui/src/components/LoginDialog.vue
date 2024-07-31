<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { reactive, ref } from 'vue'
import { API } from '@/views/net'

const { t } = useI18n()
const props = defineProps({
    visible: Boolean,
})

const deviceAuthActive = ref(0)
const deviceAuthResponse = ref({
    user_code: '',
    verification_uri: '',
    device_code: ''
})
const deviceAuthNext = () => {
    if (deviceAuthActive.value++ > 2) {
        return
    }

    if (deviceAuthActive.value === 1) {
        fetch('/oauth2/getLocalCode')
            .then(API.DefaultResponseProcess)
            .then((d) => {
                deviceAuthResponse.value = d
            })
    } else if (deviceAuthActive.value === 2) {
        window.location.href = '/oauth2/getUserInfoFromLocalCode?device_code=' + deviceAuthResponse.value.device_code
    }
}
</script>

<template>
    <el-dialog
        :modelValue="visible"
        title="You need to login first."
        width="30%"
    >
        <el-collapse accordion>
            <el-collapse-item title="Server in cloud" name="1">
                <a href="/oauth2/token" target="_blank">
                    <svg height="32" aria-hidden="true" viewBox="0 0 16 16" version="1.1" width="32" data-view-component="true" class="octicon octicon-mark-github v-align-middle color-fg-default">
                        <path d="M8 0c4.42 0 8 3.58 8 8a8.013 8.013 0 0 1-5.45 7.59c-.4.08-.55-.17-.55-.38 0-.27.01-1.13.01-2.2 0-.75-.25-1.23-.54-1.48 1.78-.2 3.65-.88 3.65-3.95 0-.88-.31-1.59-.82-2.15.08-.2.36-1.02-.08-2.12 0 0-.67-.22-2.2.82-.64-.18-1.32-.27-2-.27-.68 0-1.36.09-2 .27-1.53-1.03-2.2-.82-2.2-.82-.44 1.1-.16 1.92-.08 2.12-.51.56-.82 1.28-.82 2.15 0 3.06 1.86 3.75 3.64 3.95-.23.2-.44.55-.51 1.07-.46.21-1.61.55-2.33-.66-.15-.24-.6-.83-1.23-.82-.67.01-.27.38.01.53.34.19.73.9.82 1.13.16.45.68 1.31 2.69.94 0 .67.01 1.3.01 1.49 0 .21-.15.45-.55.38A7.995 7.995 0 0 1 0 8c0-4.42 3.58-8 8-8Z"></path>
                    </svg>
                </a>
            </el-collapse-item>
            <el-collapse-item title="Server in local" name="2">
                <el-steps :active="deviceAuthActive" finish-status="success">
                    <el-step title="Request Device Code" />
                    <el-step title="Input Code"/>
                    <el-step title="Finished" />
                </el-steps>

                <div v-if="deviceAuthActive===1">
                    Open <a :href="deviceAuthResponse.verification_uri" target="_blank">this link</a>, and type the code: <span>{{ deviceAuthResponse.user_code }}. Then click the next step button.</span>
                </div>
                <el-button style="margin-top: 12px" @click="deviceAuthNext">Next step</el-button>
            </el-collapse-item>
        </el-collapse>
    </el-dialog>
</template>
