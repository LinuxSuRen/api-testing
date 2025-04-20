<script setup lang="ts">
import { ref } from 'vue'

const simpleButton = ref(localStorage.getItem('button.style') === 'simple')
</script>

<template>
    <el-button v-bind="$attrs">
        <slot v-if="!simpleButton || !($attrs.icon)"></slot>
    </el-button>
</template>
