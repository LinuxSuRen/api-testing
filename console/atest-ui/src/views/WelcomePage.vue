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
    <div class="container">
        <el-card class="welcome-card">
            <template #header>
                <div class="card-header">
                    <span>Welcome to atest</span>
                </div>
            </template>
            
            <div class="welcome-content">
                <p>Use atest to improve your code quality!</p>
                <p class="guide-text">Please read the following guide if this is your first time using atest:</p>
                
                <el-steps direction="vertical" :active="3" class="guide-steps">
                    <el-step title="Create a store" description="Create a store for saving the data" />
                    <el-step title="Create test suite" description="Create a test suite on the left panel" />
                    <el-step title="Create test case" description="Select a suite, then create the test case" />
                </el-steps>
                
                <div class="document-link">
                    <el-link type="primary" href="https://linuxsuren.github.io/api-testing/" target="_blank" :icon="Document" class="doc-link">
                        View official documentation
                    </el-link>
                </div>
            </div>
        </el-card>

        <el-divider border-style="dashed" />

        <div class="dependency-section">
            <el-row :gutter=20>
                <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
                    <el-card class="dependency-card">
                        <template #header>
                            <div class="card-header">
                                <span>Golang Dependencies</span>
                                <el-tag type="info" size="small">{{ Object.keys(sbomItems.go || {}).length }} packages</el-tag>
                            </div>
                        </template>
                        <el-scrollbar :height="scrollbarHeight" always>
                            <ul class="dependency-list">
                                <li v-for="(k, v) in sbomItems.go" :key="v" class="dependency-item">
                                    <span class="package-name">{{ v }}</span>
                                    <el-tag size="small" class="version-tag">{{ k }}</el-tag>
                                </li>
                            </ul>
                        </el-scrollbar>
                    </el-card>
                </el-col>
                
                <el-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12">
                    <el-card class="dependency-card">
                        <template #header>
                            <div class="card-header">
                                <span>JavaScript Dependencies</span>
                                <el-tag type="info" size="small">
                                    {{ Object.keys(sbomItems.js?.dependencies || {}).length + Object.keys(sbomItems.js?.devDependencies || {}).length }} packages
                                </el-tag>
                            </div>
                        </template>
                        <el-tabs type="border-card">
                            <el-tab-pane label="Dependencies">
                                <el-scrollbar :height="scrollbarHeight - 50" always>
                                    <ul class="dependency-list">
                                        <li v-for="(k, v) in sbomItems.js?.dependencies" :key="v" class="dependency-item">
                                            <span class="package-name">{{ v }}</span>
                                            <el-tag size="small" class="version-tag">{{ k }}</el-tag>
                                        </li>
                                    </ul>
                                </el-scrollbar>
                            </el-tab-pane>
                            <el-tab-pane label="Dev Dependencies">
                                <el-scrollbar :height="scrollbarHeight - 50" always>
                                    <ul class="dependency-list">
                                        <li v-for="(k, v) in sbomItems.js?.devDependencies" :key="v" class="dependency-item">
                                            <span class="package-name">{{ v }}</span>
                                            <el-tag size="small" class="version-tag">{{ k }}</el-tag>
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

<style scoped>
.container {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    box-sizing: border-box;
}

.welcome-card {
    margin-bottom: 20px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
}

.welcome-content {
    padding: 10px;
}

.guide-text {
    margin: 15px 0;
    color: #606266;
}

.guide-steps {
    margin: 20px 0;
}

.document-link {
    margin-top: 20px;
    text-align: center;
}

.doc-link {
    font-size: 16px;
}

.dependency-section {
    margin-top: 30px;
}

.dependency-card {
    height: 100%;
    margin-bottom: 20px;
}

.dependency-list {
    list-style: none;
    padding: 0;
    margin: 0;
}

.dependency-item {
    padding: 8px 12px;
    border-bottom: 1px solid #ebeef5;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: 5px;
}

.dependency-item:last-child {
    border-bottom: none;
}

.package-name {
    font-family: monospace;
    color: #409eff;
    word-break: break-word;
}

.version-tag {
    margin-left: 10px;
    font-family: monospace;
}

/* Responsive adjustments */
@media (max-width: 768px) {
    .container {
        padding: 15px;
    }
    
    .card-header {
        flex-direction: column;
        align-items: flex-start;
    }
    
    .guide-steps {
        margin: 15px 0;
    }
    
    .dependency-item {
        flex-direction: column;
        align-items: flex-start;
        padding: 10px 0;
    }
    
    .version-tag {
        margin-left: 0;
        margin-top: 5px;
    }
}

@media (max-width: 480px) {
    .container {
        padding: 10px;
    }
    
    .welcome-content {
        padding: 5px;
    }
    
    .dependency-card {
        margin-bottom: 15px;
    }
}
</style>

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
