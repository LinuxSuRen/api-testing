<script setup lang="ts">
import { ref } from 'vue'
import { API } from '../views/net'
import { Document } from '@element-plus/icons-vue'

interface SBOM {
    go: {}
    js: {
        dependencies: {}
        devDependencies: {}
    }
}
const sbomItems = ref({} as SBOM)
API.SBOM((d) => {
    sbomItems.value = d
})
</script>

<template>
    <div class="w-full max-w-[1200px] mx-auto p-[20px] box-border max-[768px]:p-[15px] max-[480px]:p-[10px]">
        <el-card class="mb-[20px]">
            <template #header>
                <div class="flex justify-between items-center flex-wrap gap-[10px] max-[768px]:flex-col max-[768px]:items-start">
                    <span>Welcome to atest</span>
                </div>
            </template>
            
            <div class="p-[10px] max-[480px]:p-[5px]">
                <p>Use atest to improve your code quality!</p>
                <p class="my-[15px] text-[#606266]">Please read the following guide if this is your first time using atest:</p>
                
                <el-steps direction="vertical" :active="3" class="my-[20px] max-[768px]:my-[15px]">
                    <el-step title="Create a store" description="Create a store for saving the data" />
                    <el-step title="Create test suite" description="Create a test suite on the left panel" />
                    <el-step title="Create test case" description="Select a suite, then create the test case" />
                </el-steps>
                
                <div class="mt-[20px] text-center">
                    <el-link type="primary" href="https://linuxsuren.github.io/api-testing/" target="_blank" :icon="Document" class="text-[16px]">
                        View official documentation
                    </el-link>
                </div>
            </div>
        </el-card>

        <el-divider border-style="dashed" />

        <div class="mt-[30px]">
            <el-row :gutter="20">
                <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
                    <el-card class="h-full mb-[20px] max-[480px]:mb-[15px]">
                        <template #header>
                            <div class="flex justify-between items-center flex-wrap gap-[10px] max-[768px]:flex-col max-[768px]:items-start">
                                <span>Golang Dependencies</span>
                                <el-tag type="info" size="small">{{ Object.keys(sbomItems.go || {}).length }} packages</el-tag>
                            </div>
                        </template>
                        <el-scrollbar :height="scrollbarHeight" always>
                            <ul class="list-none p-0 m-0">
                                <li v-for="(k, v) in sbomItems.go" :key="v" class="py-[8px] px-[12px] max-[768px]:py-[10px] max-[768px]:px-0 border-b border-[#ebeef5] last:border-b-0 flex justify-between items-center flex-wrap gap-[5px] max-[768px]:flex-col max-[768px]:items-start">
                                    <span class="font-mono text-[#409eff] [word-break:break-word]">{{ v }}</span>
                                    <el-tag size="small" class="ml-[10px] max-[768px]:ml-0 max-[768px]:mt-[5px] font-mono">{{ k }}</el-tag>
                                </li>
                            </ul>
                        </el-scrollbar>
                    </el-card>
                </el-col>
                
                <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
                    <el-card class="h-full mb-[20px] max-[480px]:mb-[15px]">
                        <template #header>
                            <div class="flex justify-between items-center flex-wrap gap-[10px] max-[768px]:flex-col max-[768px]:items-start">
                                <span>JavaScript Dependencies</span>
                                <el-tag type="info" size="small">
                                    {{ Object.keys(sbomItems.js?.dependencies || {}).length + Object.keys(sbomItems.js?.devDependencies || {}).length }} packages
                                </el-tag>
                            </div>
                        </template>
                        <el-tabs type="border-card">
                            <el-tab-pane label="Dependencies">
                                <el-scrollbar :height="scrollbarHeight - 50" always>
                                    <ul class="list-none p-0 m-0">
                                        <li v-for="(k, v) in sbomItems.js?.dependencies" :key="v" class="py-[8px] px-[12px] max-[768px]:py-[10px] max-[768px]:px-0 border-b border-[#ebeef5] last:border-b-0 flex justify-between items-center flex-wrap gap-[5px] max-[768px]:flex-col max-[768px]:items-start">
                                            <span class="font-mono text-[#409eff] [word-break:break-word]">{{ v }}</span>
                                            <el-tag size="small" class="ml-[10px] max-[768px]:ml-0 max-[768px]:mt-[5px] font-mono">{{ k }}</el-tag>
                                        </li>
                                    </ul>
                                </el-scrollbar>
                            </el-tab-pane>
                            <el-tab-pane label="Dev Dependencies">
                                <el-scrollbar :height="scrollbarHeight - 50" always>
                                    <ul class="list-none p-0 m-0">
                                        <li v-for="(k, v) in sbomItems.js?.devDependencies" :key="v" class="py-[8px] px-[12px] max-[768px]:py-[10px] max-[768px]:px-0 border-b border-[#ebeef5] last:border-b-0 flex justify-between items-center flex-wrap gap-[5px] max-[768px]:flex-col max-[768px]:items-start">
                                            <span class="font-mono text-[#409eff] [word-break:break-word]">{{ v }}</span>
                                            <el-tag size="small" class="ml-[10px] max-[768px]:ml-0 max-[768px]:mt-[5px] font-mono">{{ k }}</el-tag>
                                        </li>
                                    </ul>
                                </el-scrollbar>
                            </el-tab-pane>
                        </el-tabs>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </div>
</template>



<script lang="ts">
export default {
    computed: {
        scrollbarHeight() {
            // Adjust scrollbar height based on screen size
            if (window.innerWidth < 768) {
                return 200;
            }
            return 300;
        }
    }
}
</script>
