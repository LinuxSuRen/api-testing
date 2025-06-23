<script setup lang="ts">
import { ref } from 'vue'
import type { Pair } from './types'
import { API } from './net'
import { useI18n } from 'vue-i18n'
import { Magic } from './magicKeys'

const { t } = useI18n()

const functionKind = ref('template')
const dialogVisible = ref(false)
const query = ref('')
const funcs = ref([] as Pair[])

function queryFuncs() {
    API.FunctionsQuery(query.value, functionKind.value, (d) => {
        funcs.value = d.data
    })
}

Magic.AdvancedKeys([{
    Keys: ['Alt+KeyT'],
    Func: () => {
        dialogVisible.value = true
    },
    Description: 'Open template functions dialog',
}])
</script>

<template>
    <el-affix position="bottom" :offset="20" class="tool-box">
        <el-button type="primary" @click="dialogVisible = !dialogVisible"
            data-intro="You can search your desired template functions.">{{ t('button.toolbox') }}</el-button>
    </el-affix>

    <el-dialog v-model="dialogVisible" :title="t('title.functionQuery')" width="50%" draggable destroy-on-close>
        <el-input
            v-model="query" placeholder="Query after enter" v-on:keyup.enter="queryFuncs">
            <template #append v-if="funcs.length > 0">
                {{ funcs.length }}
            </template>
            <template #prepend>
                <el-select
                    v-model="functionKind"
                >
                    <el-option label="Template" value="template" />
                    <el-option label="Verify" value="verify" />
                </el-select>
            </template>
        </el-input>
        <span class="dialog-footer">
            <el-table :data="funcs">
                <el-table-column label="Name" width="250">
                    <template #default="scope">
                        {{ scope.row.key }}
                    </template>
                </el-table-column>
                <el-table-column label="Function">
                    <template #default="scope">
                        {{ scope.row.value }}
                    </template>
                </el-table-column>
                <el-table-column label="Usage">
                    <template #default="scope">
                        <div style="display: flex; align-items: center">
                            <el-input v-model="scope.row.description" readonly />
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </span>
        <div>
            Powered by <a href="https://masterminds.github.io/sprig/" target="_blank">Sprig</a> and <a href="https://pkg.go.dev/text/template" target="_blank">built-in templates</a>.
        </div>
    </el-dialog>
</template>
