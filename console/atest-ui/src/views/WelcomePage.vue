<script setup lang="ts">
import { ref } from 'vue'
import { API } from '../views/net'

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
    <div>Welcome to use atest to improve your code quality!</div>
    <div>Please read the following guide if this is your first time to use atest.</div>
    <li>Create a store for saving the data</li>
    <li>Create a test suite on the left panel</li>
    <li>Select a suite, then create the test case</li>

    <div>
        Please get more details from the <a href="https://linuxsuren.github.io/api-testing/" target="_blank" rel="noopener">official document</a>.
    </div>

    <el-divider/>

    <div>
        Golang dependencies:
        <div>
            <el-scrollbar height="200px" always>
                <li v-for="k, v in sbomItems.go">
                    {{ v }}@{{ k }}
                </li>
            </el-scrollbar>
        </div>
    </div>

    <div>
        JavaScript dependencies:
        <div>
            <el-scrollbar height="200px" always>
                <li v-for="k, v in sbomItems.js.dependencies">
                    {{ v }}@{{ k }}
                </li>
                <li v-for="k, v in sbomItems.js.devDependencies">
                    {{ v }}@{{ k }}
                </li>
            </el-scrollbar>
        </div>
    </div>
</template>
