<script setup lang="ts">
import { ref } from 'vue'
import type { Pair } from './types'
import { API } from './net'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const dialogVisible = ref(false)
const query = ref('')
const funcs = ref([] as Pair[])

function queryFuncs() {
    API.FunctionsQuery(query.value, (d) => {
        funcs.value = d.data
    })
}
</script>

<template>
    <el-affix position="bottom" :offset="20">
        <el-button type="primary" @click="dialogVisible = !dialogVisible"
            data-intro="You can search your desired template functions.">{{ t('button.toolbox') }}</el-button>
    </el-affix>

    <el-dialog v-model="dialogVisible" title="Template Functions Query" width="40%" draggable destroy-on-close>
        <template #footer>
            <el-input v-model="query" placeholder="Query after enter" v-on:keyup.enter="queryFuncs" />
            <span class="dialog-footer">
                <el-table :data="funcs" style="width: 100%">
                    <el-table-column label="Key" width="250">
                        <template #default="scope">
                            <el-input v-model="scope.row.key" placeholder="Value" />
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
            </span>
        </template>
    </el-dialog>
</template>