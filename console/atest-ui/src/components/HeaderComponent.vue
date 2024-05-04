<template>
  <div class="container">
    <div>
      <!-- todo: add logo -->
      <a href="#">API Testing</a>
    </div>
    <div>
      <ui>
        <li>
          <a href="https://linuxsuren.github.io/api-testing/" target="_blank">
            <p>Document</p>
          </a>
        </li>
        <li>
          <a href="https://github.com/linuxsuren/api-testing.git" target="_blank"><img
              src="../assets/github/github-mark.png" alt="github" /></a>
        </li>
        <!-- todo：主题样式显示有误。 -->
        <!-- <li>
          <el-tooltip content="点击以切换主题" placement="bottom" effect="light">
            <el-switch
            v-model="darkOpt"
            inline-prompt
            :active-action-icon="Sunny"
            :inactive-action-icon="Moon"
            @click="setTheme()"
          />
          </el-tooltip>
        </li> -->
      </ui>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, inject } from 'vue'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { setAsDarkTheme } from '../theme/theme'

const darkOpt = ref(false)

const setTheme = () => {
  console.log(darkOpt.value)
  setAsDarkTheme(darkOpt.value)
  localStorage.setItem('darkOpt', darkOpt.value)
}
</script>

<style scoped>
.container {
  display: flex;
  align-items: center;
  line-height: 40px;
  justify-content: space-around;
}

ui {
  display: flex;
  list-style: none;

  li {
    line-height: 50px;
    margin-left: 20%;
  }
}

p,
a {
  font-weight: 700;
  color: black;
}

.el-switch {
  --el-switch-on-color: black;
  --el-switch-off-color: white;
}

img {
  vertical-align: middle;
  width: 30px;
  height: 30px;
}
</style>
