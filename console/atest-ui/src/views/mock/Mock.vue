<script setup lang="ts">
import { ref } from 'vue';
import { Codemirror } from 'vue-codemirror';

const mockConfig = ref('');
const link = ref('')
API.GetMockConfig((d) => {
  mockConfig.value = d.Config
  link.value = window.location.origin + d.Prefix + "/api.json"
})
const tabActive = ref('yaml')
</script>
    
<template>
  <div>
    <el-button type="warning" @click="API.ReloadMockServer(mockConfig)">Reload</el-button>
    <el-divider direction="vertical" />
    <el-link target="_blank" :href="link">{{ link }}</el-link> <!-- Noncompliant -->
  </div>
  <div>
    <el-tabs v-model="tabActive">
      <el-tab-pane label="YAML" name="yaml">
        <Codemirror v-model="mockConfig" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
    