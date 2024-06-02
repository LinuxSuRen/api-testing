<template>
  <el-aside width="collapse">
    <el-menu default-active="0" class="el-menu-vertical-demo" router>
      <template v-for="value in router.options.routes[0].children" :key="value.path">
        <el-menu-item v-show="value.meta?.isShow == true" :index="value.path">
          <el-icon>
            <component :is="value.meta?.icon"></component>
          </el-icon>
          <span>{{ t(value.meta?.title) }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </el-aside>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
</script>

<style scope>
.el-menu {
  width: 130px;
  height: calc(100vh - 59px);
  background-color: rgb(243, 244, 246);
  --el-menu-active-color: #409EFF;
  overflow-y: auto; /* 启用垂直滚动条 */
}
</style>
