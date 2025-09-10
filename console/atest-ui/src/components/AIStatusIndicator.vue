<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { API, AIPluginHealth } from '../views/net'

const aiPluginHealth = ref<Record<string, AIPluginHealth>>({})
const loading = ref(true)
const hasAIPlugins = ref(false)
let healthCheckInterval: number | null = null

const loadAIPluginHealth = () => {
  API.GetAllAIPluginHealth(
    (healthData) => {
      aiPluginHealth.value = healthData
      hasAIPlugins.value = Object.keys(healthData).length > 0
      loading.value = false
    },
    () => {
      loading.value = false
    }
  )
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'online':
      return 'success'
    case 'offline':
      return 'danger'
    case 'error':
      return 'danger'
    case 'processing':
      return 'warning'
    default:
      return 'info'
  }
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'online':
      return 'CircleCheck'
    case 'offline':
      return 'CircleClose'
    case 'error':
      return 'Warning'
    case 'processing':
      return 'Loading'
    default:
      return 'QuestionFilled'
  }
}

onMounted(() => {
  loadAIPluginHealth()
  // Poll for health status every 30 seconds
  healthCheckInterval = window.setInterval(loadAIPluginHealth, 30000)
})

onUnmounted(() => {
  if (healthCheckInterval) {
    clearInterval(healthCheckInterval)
  }
})
</script>

<template>
  <div class="ai-status-indicator" v-if="hasAIPlugins">
    <el-tooltip content="AI Plugin Status" placement="top">
      <div class="ai-status-container">
        <el-icon v-if="loading" class="is-loading">
          <Loading />
        </el-icon>
        <template v-else>
          <el-badge 
            v-for="(health, name) in aiPluginHealth" 
            :key="name"
            :type="getStatusColor(health.status)"
            :dot="true"
            class="ai-plugin-badge"
          >
            <el-tooltip 
              :content="`${health.name}: ${health.status}${health.errorMessage ? ' - ' + health.errorMessage : ''}`"
              placement="top"
            >
              <el-icon :class="`status-${health.status}`">
                <component :is="getStatusIcon(health.status)" />
              </el-icon>
            </el-tooltip>
          </el-badge>
        </template>
      </div>
    </el-tooltip>
  </div>
</template>

<style scoped>
.ai-status-indicator {
  display: flex;
  align-items: center;
  margin-left: 10px;
}

.ai-status-container {
  display: flex;
  align-items: center;
  gap: 4px;
}

.ai-plugin-badge {
  display: inline-flex;
  align-items: center;
}

.status-online {
  color: var(--el-color-success);
}

.status-offline {
  color: var(--el-color-danger);
}

.status-error {
  color: var(--el-color-danger);
}

.status-processing {
  color: var(--el-color-warning);
  animation: spin 2s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>