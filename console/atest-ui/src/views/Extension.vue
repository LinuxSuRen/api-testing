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
                const globalScope = globalThis as { ATestPlugin?: { mount?: (el: Element) => void } };
                const plugin = globalScope.ATestPlugin;

                if (plugin && plugin.mount) {
                    const container = document.getElementById("plugin-container");
                    if (container) {
                        container.innerHTML = ''; // Clear previous content
                        plugin.mount(container);
                        loading.value = false;
                    } else {
                        loading.value = false;
                    }
                } else if (retries < maxRetries) {
                    // Incremental retry mechanism: 50ms, 100ms, 150ms...
                    const delay = 50 + retries * 50;
                    setTimeout(() => checkPluginLoad(retries + 1, maxRetries), delay);
                } else {
                    loading.value = false;
                }
            };

            // Start the retry mechanism
            checkPluginLoad();
        });
    } catch (error) {
        loading.value = false; // Set loading to false on error
        console.error('Failed to load extension assets', error);
    }
};

loadPlugin().catch((error) => {
    loading.value = false;
    console.error('Failed to initialize extension plugin', error);
});
</script>

<template>
    <div id="plugin-container"
        v-loading="loading"
        :element-loading-text="props.name + ' is loading...'">
    </div>
</template>
