<script setup lang="ts">
import { NewSuggestedAPIsQuery } from './types'
import type { TestCase } from './types'

const props = defineProps({
  suite: String,
  address: String,
  method: String,
  operationId: String
})
// const emit = defineEmits(['update:address', 'update:method', 'update:operationId'])
const querySuggestedAPIs = NewSuggestedAPIsQuery('', props.suite!)

const handleSelect = (item: TestCase) => {
  // emit('update:address', item.request.api)
  // emit('update:method', item.request.method)
  // emit('update:operationId', item.name)
}
</script>

<template>
  <el-autocomplete
    :value="address"
    @input="$emit('update:address', $event)"
    :fetch-suggestions="querySuggestedAPIs"
    @select="handleSelect"
    placeholder="API Address"
    style="width: 70%; margin-left: 5px; margin-right: 5px"
  >
    <template #default="{ item }">
      <div class="value">{{ item.request.method }}</div>
      <span class="link">{{ item.request.api }}</span>
    </template>
  </el-autocomplete>
</template>
../types/types../types/types