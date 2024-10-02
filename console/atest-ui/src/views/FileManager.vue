<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const fileLists = ref([])

onMounted(() => {
  const savedFile = localStorage.getItem("fileLists")
  if (savedFile) {
    fileLists.value = JSON.parse(savedFile)
  }
})

const handleFileUpload = (file) => {
  const reader = new FileReader()

  reader.onload = (e: any) => {
    const fileData = {
      name: file.name,
      type: file.type,
      size: file.size,
      content: e.target.result // 文件内容作为 Base64 字符串
    }

    fileLists.value.push(fileData) // 将文件信息推入fileList

    localStorage.setItem('fileLists', JSON.stringify(fileLists.value)) // 存入localStorage

    ElMessage({
      message: '文件上传成功',
      type: 'success'
    })
  }

  reader.readAsDataURL(file.raw) // 读取文件内容为Base64编码
  return false
}

const deleteFile = (file) => {
  const index = fileLists.value.findIndex(item => item.name === file.name)
  if (index !== -1) {
    fileLists.value.splice(index, 1)
    localStorage.setItem('fileLists', JSON.stringify(fileLists.value))
  }
}

</script>

<template>
  <div>
    <h2>文件管理器</h2>
    
    <el-upload
      action="#"
      :auto-upload="false"
      :on-change="handleFileUpload"
      :on-remove="deleteFile"
      :file-list="fileLists"
      multiple
    >
      <el-button type="primary">点击上传</el-button>
      <template #tip>
        <div class="el-upload__tip">
          可以上传任意类型的文件
        </div>
      </template>
    </el-upload>
    
  </div>
</template>
