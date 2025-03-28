<template>
  <el-autocomplete
    v-model="input"
    clearable
    :fetch-suggestions="querySearch"
    @select="handleSelect"
    @keyup.enter="handleEnter"
    :placeholder="props.placeholder"
  >
    <template #default="{ item }">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span>{{ item.value }}</span>
        <el-icon @click.stop="deleteHistoryItem(item)">
          <delete />
        </el-icon>
      </div>
    </template>
  </el-autocomplete>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue'
import { ElAutocomplete, ElIcon } from 'element-plus'
import { Delete } from '@element-plus/icons-vue'

const props = defineProps({
  maxItems: {
    type: Number,
    default: 10
  },
  group: {
    type: String,
    default: 'history'
  },
  storage: {
    type: String,
    default: 'localStorage'
  },
  callback: {
    type: Function,
    default: () => true
  },
  placeholder: {
    type: String,
    default: 'Type something'
  }
})

const input = ref('')
const suggestions = ref([])
interface HistoryItem {
  value: string
  count: number
  timestamp: number
}

const querySearch = (queryString: string, cb: any) => {
  const results = suggestions.value.filter((item: HistoryItem) => item.value.includes(queryString))
  cb(results)
}

const handleSelect = (item: HistoryItem) => {
  input.value = item.value
}

const handleEnter = async () => {
  if (props.callback) {
    const result = await props.callback()
    if (!result) {
      return
    }
  }
  if (input.value === '') {
    return;
  }

  const history = JSON.parse(getStorage().getItem(props.group) || '[]')
  const existingItem = history.find((item: HistoryItem) => item.value === input.value)

  if (existingItem) {
    existingItem.count++
    existingItem.timestamp = Date.now()
  } else {
    history.push({ value: input.value, count: 1, timestamp: Date.now() })
  }

  if (history.length > props.maxItems) {
    history.sort((a: HistoryItem, b: HistoryItem) => a.count - b.count || a.timestamp - b.timestamp)
    history.shift()
  }

  getStorage().setItem(props.group, JSON.stringify(history))
  suggestions.value = history
}

const loadHistory = () => {
  suggestions.value = JSON.parse(getStorage().getItem(props.group) || '[]')
}

const deleteHistoryItem = (item: HistoryItem) => {
  const history = JSON.parse(getStorage().getItem(props.group) || '[]')
  const updatedHistory = history.filter((historyItem: HistoryItem) => historyItem.value !== item.value)
  getStorage().setItem(props.group, JSON.stringify(updatedHistory))
  suggestions.value = updatedHistory
}

const getStorage = () => {
    switch (props.storage) {
      case 'localStorage':
        return localStorage
      case 'sessionStorage':
        return sessionStorage
      default:
        return localStorage
    }
}

loadHistory()
</script>
