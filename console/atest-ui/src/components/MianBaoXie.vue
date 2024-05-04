<template>
  <div class="pb-6">
    <el-breadcrumb separator="/">
      <el-breadcrumb-item v-for="(item, index) in $route.matched" :key="index">{{ item.meta.title }}</el-breadcrumb-item>
    </el-breadcrumb>
  </div>
</template>
