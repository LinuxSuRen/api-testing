<script setup lang="ts">
const props = defineProps({
    name: String,
})
</script>

<template>
    <div class="common-layout">
        <el-text class="mx-1" type="primary">{{props.name}}</el-text>
    </div>
</template>
