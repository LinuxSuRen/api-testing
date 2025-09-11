/*
Copyright 2023-2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { config } from '@vue/test-utils'
import ElementPlus from 'element-plus'
import { createI18n } from 'vue-i18n'

// Create a mock i18n instance for testing
const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      ai: {
        status: {
          healthy: 'Healthy',
          unhealthy: 'Unhealthy',
          unknown: 'Unknown'
        },
        trigger: {
          button: 'AI Assistant',
          dialog: {
            title: 'AI Assistant',
            placeholder: 'Ask me anything about your API tests...',
            send: 'Send',
            cancel: 'Cancel'
          }
        },
        triggerButton: 'AI Assistant'
      },
      'AI Assistant': 'AI Assistant',
      'AI Plugin Status': 'AI Plugin Status'
    },
    zh: {
      ai: {
        status: {
          healthy: '健康',
          unhealthy: '不健康',
          unknown: '未知'
        },
        trigger: {
          button: 'AI助手',
          dialog: {
            title: 'AI助手',
            placeholder: '询问关于API测试的任何问题...',
            send: '发送',
            cancel: '取消'
          }
        },
        triggerButton: 'AI助手'
      },
      'AI Assistant': 'AI助手',
      'AI Plugin Status': 'AI插件状态'
    }
  }
})

// Global plugins configuration for tests
config.global.plugins = [ElementPlus, i18n]

// Global stubs for Element Plus components
config.global.stubs = {
  'el-icon': {
    template: '<div class="el-icon"><slot /></div>'
  },
  'el-tooltip': {
    template: '<div class="el-tooltip"><slot /></div>'
  },
  'el-badge': {
    template: '<div class="el-badge"><slot /></div>'
  },
  'el-button': {
    template: '<button class="el-button" @click="$emit(\'click\')"><slot /></button>',
    emits: ['click']
  },
  'el-dialog': {
    template: '<div class="el-dialog" v-if="modelValue"><slot /></div>',
    props: ['modelValue'],
    emits: ['update:modelValue']
  },
  'el-input': {
    template: '<input class="el-input" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue'],
    emits: ['update:modelValue']
  },
  'el-alert': {
    template: '<div class="el-alert"><slot /></div>'
  },
  'Loading': {
    template: '<div class="loading">Loading...</div>'
  },
  'ChatLineSquare': {
    template: '<div class="chat-line-square-icon"></div>'
  }
}

// Mock global properties that might be used in components
config.global.mocks = {
  $t: (key: string) => key
}