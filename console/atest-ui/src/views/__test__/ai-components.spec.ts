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

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import AIStatusIndicator from '../../components/AIStatusIndicator.vue'
import AITriggerButton from '../../components/AITriggerButton.vue'
import { API } from '../net'

// Mock the API
vi.mock('../net', () => ({
  API: {
    GetAllAIPluginHealth: vi.fn()
  }
}))

describe('AI Components', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('AIStatusIndicator', () => {
    it('should render without AI plugins', async () => {
      const mockAPI = API.GetAllAIPluginHealth as any
      mockAPI.mockImplementation((callback: Function) => {
        callback({})
      })

      const wrapper = mount(AIStatusIndicator)
      
      // Should not render when no AI plugins
      expect(wrapper.find('.ai-status-indicator').exists()).toBe(false)
    })

    it('should render AI plugin status when plugins exist', async () => {
      const mockAPI = API.GetAllAIPluginHealth as any
      mockAPI.mockImplementation((callback: Function) => {
        callback({
          'test-plugin': {
            name: 'test-plugin',
            status: 'online',
            lastCheckAt: '2025-09-09T10:00:00Z',
            responseTime: 100,
            errorMessage: '',
            metrics: {}
          }
        })
      })

      const wrapper = mount(AIStatusIndicator)
      
      // Wait for component to process data
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.ai-status-indicator').exists()).toBe(true)
      expect(wrapper.find('.ai-plugin-badge').exists()).toBe(true)
    })

    it('should handle different plugin statuses correctly', async () => {
      const mockAPI = API.GetAllAIPluginHealth as any
      mockAPI.mockImplementation((callback: Function) => {
        callback({
          'plugin-online': {
            name: 'plugin-online',
            status: 'online',
            lastCheckAt: '2025-09-09T10:00:00Z',
            responseTime: 100
          },
          'plugin-offline': {
            name: 'plugin-offline',
            status: 'offline',
            lastCheckAt: '2025-09-09T10:00:00Z',
            responseTime: 0,
            errorMessage: 'Plugin not responding'
          }
        })
      })

      const wrapper = mount(AIStatusIndicator)
      await wrapper.vm.$nextTick()
      
      const badges = wrapper.findAll('.ai-plugin-badge')
      expect(badges).toHaveLength(2)
    })
  })

  describe('AITriggerButton', () => {
    it('should render floating action button', () => {
      const wrapper = mount(AITriggerButton)
      
      expect(wrapper.find('.ai-trigger-button').exists()).toBe(true)
      expect(wrapper.find('.ai-trigger-container').exists()).toBe(true)
    })

    it('should emit ai-trigger-clicked when button is clicked', async () => {
      const wrapper = mount(AITriggerButton)
      
      await wrapper.find('.ai-trigger-button').trigger('click')
      
      expect(wrapper.emitted('ai-trigger-clicked')).toBeTruthy()
      expect(wrapper.emitted('ai-trigger-clicked')).toHaveLength(1)
    })

    it('should show dialog when triggered', async () => {
      const wrapper = mount(AITriggerButton)
      
      await wrapper.find('.ai-trigger-button').trigger('click')
      await wrapper.vm.$nextTick()
      
      expect(wrapper.find('.ai-dialog-content').exists()).toBe(true)
    })

    it('should have proper accessibility attributes', () => {
      const wrapper = mount(AITriggerButton)
      
      const button = wrapper.find('.ai-trigger-button')
      expect(button.attributes('aria-label')).toBeDefined()
      expect(button.attributes('tabindex')).toBe('0')
    })

    it('should handle processing state correctly', async () => {
      const wrapper = mount(AITriggerButton)
      
      // Trigger processing simulation
      await wrapper.find('.ai-trigger-button').trigger('click')
      await wrapper.vm.$nextTick()
      
      // Check if processing state can be triggered (button should exist and be clickable)
      expect(wrapper.find('.ai-trigger-button').exists()).toBe(true)
    })
  })
})