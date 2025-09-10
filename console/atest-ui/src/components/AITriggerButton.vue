<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const isProcessing = ref(false)
const showAIDialog = ref(false)

const emits = defineEmits(['ai-trigger-clicked'])

const handleAITrigger = () => {
  if (isProcessing.value) return
  
  emits('ai-trigger-clicked')
  showAIDialog.value = true
}

const closeDialog = () => {
  showAIDialog.value = false
}

// For demo purposes - in real implementation, this would be managed by parent
const simulateProcessing = () => {
  isProcessing.value = true
  setTimeout(() => {
    isProcessing.value = false
  }, 3000)
}
</script>

<template>
  <div class="ai-trigger-container">
    <!-- Floating Action Button -->
    <el-button
      type="primary"
      circle
      size="large"
      class="ai-trigger-button"
      :class="{ 'is-processing': isProcessing }"
      :loading="isProcessing"
      @click="handleAITrigger"
      :aria-label="t ? t('ai.triggerButton') : 'AI Assistant'"
      tabindex="0"
    >
      <el-icon v-if="!isProcessing" size="24">
        <ChatLineSquare />
      </el-icon>
    </el-button>

    <!-- Simple AI Dialog Interface -->
    <el-dialog
      v-model="showAIDialog"
      title="AI Assistant"
      width="60%"
      :before-close="closeDialog"
      destroy-on-close
    >
      <div class="ai-dialog-content">
        <el-alert
          title="AI Plugin Interface"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            This is the main project's AI interface. Actual AI functionality 
            should be implemented by separate AI plugins that integrate with 
            this interface through the plugin system.
          </template>
        </el-alert>
        
        <div class="ai-placeholder">
          <p>AI plugins can be loaded here dynamically...</p>
          <el-button @click="simulateProcessing" :disabled="isProcessing">
            Simulate AI Processing
          </el-button>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">Close</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.ai-trigger-container {
  position: fixed;
  bottom: 30px;
  right: 30px;
  z-index: 1000;
}

.ai-trigger-button {
  width: 60px;
  height: 60px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  border: none;
  background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
}

.ai-trigger-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.ai-trigger-button:active {
  transform: translateY(0);
}

.ai-trigger-button.is-processing {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(var(--el-color-primary-rgb), 0.7);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(var(--el-color-primary-rgb), 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(var(--el-color-primary-rgb), 0);
  }
}

.ai-dialog-content {
  padding: 20px 0;
}

.ai-placeholder {
  margin-top: 20px;
  text-align: center;
  padding: 40px;
  border: 2px dashed var(--el-border-color);
  border-radius: 8px;
  color: var(--el-text-color-secondary);
}

/* Accessibility improvements */
.ai-trigger-button:focus {
  outline: 2px solid var(--el-color-primary);
  outline-offset: 2px;
}

@media (prefers-reduced-motion: reduce) {
  .ai-trigger-button {
    transition: none;
  }
  
  .ai-trigger-button:hover {
    transform: none;
  }
  
  .ai-trigger-button.is-processing {
    animation: none;
  }
}
</style>