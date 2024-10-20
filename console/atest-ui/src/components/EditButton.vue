<script setup lang="ts">
import { nextTick, ref } from 'vue'
import { ElInput } from 'element-plus'
import type { InputInstance } from 'element-plus'

const props = defineProps({
    value: String,
})

const emit = defineEmits(['changed'])
const inputVisible = ref(false)
const inputValue = ref('')
const InputRef = ref<InputInstance>()

const showInput = () => {
  inputVisible.value = true
  inputValue.value = props.value ?? ''
  nextTick(() => {
    InputRef.value!.input!.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value && props.value !== inputValue.value) {
    emit('changed', inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}
</script>

<template>
    <span class="flex gap-2">
      <el-input
        v-if="inputVisible"
        ref="InputRef"
        v-model="inputValue"
        class="w-20"
        style="width: 200px"
        @keyup.enter="handleInputConfirm"
        @blur="handleInputConfirm"
      />
      <el-button v-else class="button-new-tag" size="small" @click="showInput">
        {{ value }}
      </el-button>
    </span>
</template>
