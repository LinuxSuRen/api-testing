<template>
  <div class="container">
    <div class="logo_container">
      <a href="#">
        <img class="logo_container_img" src="../assets/logo/logo.png" alt="logo" />
        &nbsp;&nbsp;
        <!-- todo -->
        <!-- <p class="pro_name">API Testing</p> -->
      </a>
    </div>
    <div>
      <ui>
        <li>
          <LanguageSwitch />
        </li>
        <li>
          <a href="https://linuxsuren.github.io/api-testing/" target="_blank">
            <p>Document</p>
          </a>
        </li>
        <li>
          <a href="https://github.com/linuxsuren/api-testing.git" target="_blank"
            ><img class="github" src="../assets/github/github-mark.png" alt="github"
          /></a>
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
import { ref } from 'vue'
// import { Sunny, Moon } from '@element-plus/icons-vue'
import { setAsDarkTheme } from '../theme/theme'
import LanguageSwitch from '@/components/LanguageSwitch.vue'

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
  line-height: 40px;
  justify-content: space-between;
  flex-direction: row;
  align-items: center;
}

.logo_container {
  display: flex;
  align-items: center;
}

.logo_container_img {
  /* 保持图片等比例缩放 */
  height: 50px;
  vertical-align: middle;
}

.pro_name {
  display: inline-block;
  vertical-align: middle;
}

p,
a {
  font-weight: 700;
  color: black;
}

ui {
  display: flex;
  list-style: none;
}

li {
  display: flex;
  padding-right: 1em;
  vertical-align: middle;
  align-items: center;
}

.el-switch {
  --el-switch-on-color: black;
  --el-switch-off-color: white;
}

.github {
  vertical-align: middle;
  width: 30px;
  height: 30px;
}
</style>
