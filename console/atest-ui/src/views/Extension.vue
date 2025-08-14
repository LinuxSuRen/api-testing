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
        API.GetPageOfCSS(props.name, (d) => {
            const style = document.createElement('style');
            style.textContent = d.message;
            document.head.appendChild(style);
        });

        API.GetPageOfJS(props.name, (d) => {
            const script = document.createElement('script');
            script.type = 'text/javascript';
            script.textContent = d.message;
            document.head.appendChild(script);

            const plugin = window.ATestPlugin;
            
            if (plugin && plugin.mount) {
                console.log('extension load success');
                const container = document.getElementById("plugin-container");
                if (container) {
                    container.innerHTML = ''; // Clear previous content
                    plugin.mount(container);
                }
            }
        });
    } catch (error) {
        console.log(`extension load error: ${(error as Error).message}`)
    } finally {
        console.log('extension load finally');
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
