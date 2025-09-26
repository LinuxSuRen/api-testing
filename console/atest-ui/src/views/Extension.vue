<script setup lang="ts">
import { ref } from 'vue'
import { API } from './net';

interface Props {
  name: string
}
const props = defineProps<Props>()
const loading = ref(true)
const loadPlugin = async (): Promise<void> => {
    try {
        // First load CSS
        API.GetPageOfCSS(props.name, (d) => {
            const style = document.createElement('style');
            style.textContent = d.message;
            document.head.appendChild(style);
        });

        // Then load JS and mount plugin
        API.GetPageOfJS(props.name, (d) => {
            const script = document.createElement('script');
            script.type = 'text/javascript';
            script.textContent = d.message;
            document.head.appendChild(script);

            // Implement retry mechanism with exponential backoff
            const checkPluginLoad = (retries = 0, maxRetries = 10) => {
                const plugin = (window as any).ATestPlugin;

                console.log(`Plugin load attempt ${retries + 1}/${maxRetries + 1}`);

                if (plugin && plugin.mount) {
                    console.log('extension load success');
                    const container = document.getElementById("plugin-container");
                    if (container) {
                        container.innerHTML = ''; // Clear previous content
                        plugin.mount(container);
                        loading.value = false;
                    } else {
                        console.error('Plugin container not found');
                        loading.value = false;
                    }
                } else if (retries < maxRetries) {
                    // Incremental retry mechanism: 50ms, 100ms, 150ms...
                    const delay = 50 + retries * 50;
                    console.log(`ATestPlugin not ready, retrying in ${delay}ms (attempt ${retries + 1}/${maxRetries + 1})`);
                    setTimeout(() => checkPluginLoad(retries + 1, maxRetries), delay);
                } else {
                    console.error('ATestPlugin not found or missing mount method after max retries');
                    console.error('Window.ATestPlugin value:', (window as any).ATestPlugin);
                    loading.value = false;
                }
            };

            // Start the retry mechanism
            checkPluginLoad();
        });
    } catch (error) {
        console.log(`extension load error: ${(error as Error).message}`);
        loading.value = false; // Set loading to false on error
    }
};
try {
    loadPlugin();
} catch (error) {
    console.error('extension load error:', error);
}
</script>

<template>
    <div id="plugin-container"
        v-loading="loading"
        :element-loading-text="props.name + ' is loading...'">
    </div>
</template>
