<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <Header></Header>
      </el-header>
      <el-container>
        <LeftSideBar />
        <el-container>
          <el-main>
            <div class="mbx">
              <el-icon><Fold /></el-icon>
              <Mbx />
            </div>
            <router-view />
          </el-main>
          <el-footer>
            <Footer />
          </el-footer>
        </el-container>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/HeaderComponent.vue'
import LeftSideBar from '@/components/LeftSideComponent.vue'
import Footer from '@/components/FooterComponent.vue'
import Mbx from '@/components/MianBaoXie.vue'
import { Fold } from '@element-plus/icons-vue'
</script>

<style>
.mbx {
  display: flex;
}

.el-icon {
  margin-right: 3px;
  cursor: pointer;
}

.common-layout {
  margin: 0 auto;
  padding: 0 auto;
  height: 100%;
}

.main-container {
  display: flex;
  background-color: rgb(243, 244, 246);
}

.el-header {
  background-color: rgb(243, 244, 246);
  text-align: center;
  line-height: normal;
}
</style>
