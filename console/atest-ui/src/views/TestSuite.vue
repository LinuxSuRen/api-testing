<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
    name: String,
})
const emit = defineEmits(['updated'])

interface Suite {
    name: string;
    api: string;
}

const suite = ref({} as Suite)
function load() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            name: props.name,
        })
    };
    fetch('/server.Runner/GetTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            suite.value = {
                name: e.name,
                api: e.api,
            } as Suite
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}
load()
watch(props, () => {
    load()
})

function save() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify(suite.value),
    };
    fetch('/server.Runner/UpdateTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            ElMessage({
                    message: 'Updated.',
                    type: 'success',
                })
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}

function del() {
    const requestOptions = {
        method: 'POST',
        body: JSON.stringify({
            name: props.name,
        })
    };
    fetch('/server.Runner/DeleteTestSuite', requestOptions)
        .then(response => response.json())
        .then(e => {
            ElMessage({
                    message: 'Deleted.',
                    type: 'success',
                })
            emit('updated')
        }).catch(e => {
            ElMessage.error('Oops, ' + e)
        });
}
</script>

<template>
    <div class="common-layout">
        <el-text class="mx-1" type="primary">{{suite.name}}</el-text>
        <el-input class="mx-1" v-model="suite.api" placeholder="API"></el-input>

        <el-button type="primary" @click="save">Save</el-button>
        <el-button type="primary" @click="del">Delete</el-button>
    </div>
</template>
