<template>
  <el-card class="card" shadow="hover">
    <el-button type="warning" @click="ReloadMockServer({config: mockConfig})">Reload</el-button>
    <el-divider direction="vertical" />
    <el-link target="_blank" :href="link">{{ link }}</el-link>
    <el-divider direction="vertical" />
    <div class="yaml-container">
      <el-tabs v-model="tabActive">
        <el-tab-pane label="YAML" name="yaml">
          <Codemirror v-model="mockConfig" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-card>
</template>
    
<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Codemirror } from 'vue-codemirror'
import { GetMockConfig, ReloadMockServer } from '@/api/mock/mock'

const mockConfig = ref('')
const link = ref('')

GetMockConfig()
  .then((res: any) => {
    mockConfig.value = res.config
    link.value = window.location.origin + res.Prefix + '/api.json'
  })
  .catch((err: any) => {
    ElMessage({
        type: 'error',
        showClose: true,
        message: 'Oops, ' + err.message || 'Unknown error when feching test data!'
      })
  })

const tabActive = ref('yaml')
</script>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  margin-top: 1%;
  width: 100%;
  max-width: 1750px;
  height: auto;
  vertical-align: middle;
}

.yaml-container  {
  flex-grow: 1;
  overflow-y: auto !important;
}
</style>
