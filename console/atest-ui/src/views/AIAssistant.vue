<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { ChatDotRound, Share, Refresh } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { API } from './net'

const { t } = useI18n()

interface ChatMessage {
  id: string
  type: 'user' | 'assistant'
  content: string
  timestamp: Date
}

const messages = ref<ChatMessage[]>([])
const inputMessage = ref('')
const isLoading = ref(false)
const chatContainer = ref<HTMLElement>()

// Add welcome message
messages.value.push({
  id: '1',
  type: 'assistant',
  content: 'Hello! I\'m your AI assistant. I can help you with:\n\n• Converting natural language to SQL queries\n• Generating test cases\n• Optimizing database queries\n\nHow can I assist you today?',
  timestamp: new Date()
})

const sendMessage = async () => {
  if (!inputMessage.value.trim() || isLoading.value) return

  const userMessage: ChatMessage = {
    id: Date.now().toString(),
    type: 'user',
    content: inputMessage.value,
    timestamp: new Date()
  }

  messages.value.push(userMessage)
  const userInput = inputMessage.value
  inputMessage.value = ''
  isLoading.value = true

  try {
    // Call AI service through the backend API
    // For now, we'll simulate the AI response
    await simulateAIResponse(userInput)
  } catch (error) {
    ElMessage.error('Failed to get AI response: ' + error)
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

const simulateAIResponse = async (userInput: string) => {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 1000))

  let response = ''
  const lowerInput = userInput.toLowerCase()

  if (lowerInput.includes('sql') || lowerInput.includes('query') || lowerInput.includes('database')) {
    response = `Here's a SQL query based on your request:\n\n\`\`\`sql\nSELECT * FROM users WHERE status = 'active'\nORDER BY created_at DESC\nLIMIT 10;\n\`\`\`\n\nThis query will retrieve the 10 most recently created active users.`
  } else if (lowerInput.includes('test') || lowerInput.includes('case')) {
    response = `Here's a test case structure for your API:\n\n\`\`\`json\n{\n  "name": "Test User Creation",\n  "request": {\n    "method": "POST",\n    "url": "/api/users",\n    "body": {\n      "name": "John Doe",\n      "email": "john@example.com"\n    }\n  },\n  "response": {\n    "statusCode": 201,\n    "body": {\n      "id": "{{generated_id}}",\n      "name": "John Doe",\n      "email": "john@example.com"\n    }\n  }\n}\n\`\`\`\n\nThis test case validates user creation functionality.`
  } else {
    response = `I understand you're asking about: "${userInput}"\n\nI'm currently in MVP mode and can help with:\n• SQL query generation\n• Test case creation\n• Query optimization\n\nCould you please specify what type of assistance you need?`
  }

  const assistantMessage: ChatMessage = {
    id: Date.now().toString(),
    type: 'assistant',
    content: response,
    timestamp: new Date()
  }

  messages.value.push(assistantMessage)
}

const clearChat = () => {
  messages.value = [{
    id: '1',
    type: 'assistant',
    content: 'Hello! I\'m your AI assistant. I can help you with:\n\n• Converting natural language to SQL queries\n• Generating test cases\n• Optimizing database queries\n\nHow can I assist you today?',
    timestamp: new Date()
  }]
}

const scrollToBottom = () => {
  setTimeout(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  }, 100)
}

const handleKeyPress = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const formatMessage = (content: string) => {
  // Simple formatting for code blocks and line breaks
  return content
    .replace(/\n/g, '<br>')
    .replace(/```(\w+)?\n([\s\S]*?)```/g, '<pre><code class="language-$1">$2</code></pre>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
}
</script>

<template>
  <div class="ai-assistant-container">
    <div class="page-header">
      <span class="page-title">
        <el-icon><ChatDotRound /></el-icon>
        AI Assistant
      </span>
      <el-button type="primary" @click="clearChat" :icon="Refresh">Clear Chat</el-button>
    </div>

    <div class="chat-container" ref="chatContainer">
      <div 
        v-for="message in messages" 
        :key="message.id" 
        :class="['message', message.type]"
      >
        <div class="message-content">
          <div class="message-header">
            <span class="message-sender">
              {{ message.type === 'user' ? 'You' : 'AI Assistant' }}
            </span>
            <span class="message-time">
              {{ message.timestamp.toLocaleTimeString() }}
            </span>
          </div>
          <div class="message-text" v-html="formatMessage(message.content)"></div>
        </div>
      </div>
      
      <div v-if="isLoading" class="message assistant">
        <div class="message-content">
          <div class="message-header">
            <span class="message-sender">AI Assistant</span>
          </div>
          <div class="message-text">
            <el-icon class="is-loading"><ChatDotRound /></el-icon>
            Thinking...
          </div>
        </div>
      </div>
    </div>

    <div class="input-container">
      <el-input
        v-model="inputMessage"
        type="textarea"
        :rows="3"
        placeholder="Ask me anything about SQL queries, test cases, or database optimization..."
        :disabled="isLoading"
        @keypress="handleKeyPress"
        class="message-input"
      />
      <el-button 
        type="primary" 
        @click="sendMessage" 
        :disabled="!inputMessage.trim() || isLoading"
        :icon="Share"
        class="send-button"
      >
        Send
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.ai-assistant-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 60px);
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--el-border-color);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chat-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background: var(--el-bg-color-page);
  border-radius: 8px;
  margin-bottom: 20px;
}

.message {
  margin-bottom: 20px;
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 70%;
  padding: 12px 16px;
  border-radius: 12px;
  background: var(--el-bg-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.message.user .message-content {
  background: var(--el-color-primary);
  color: white;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
  opacity: 0.7;
}

.message-sender {
  font-weight: 600;
}

.message-time {
  font-size: 11px;
}

.message-text {
  line-height: 1.5;
  word-wrap: break-word;
}

.message-text :deep(pre) {
  background: var(--el-fill-color-light);
  padding: 12px;
  border-radius: 6px;
  margin: 8px 0;
  overflow-x: auto;
}

.message-text :deep(code) {
  background: var(--el-fill-color-light);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.message-input {
  flex: 1;
}

.send-button {
  height: 40px;
}

.is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
