<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from './net'
import { Cache } from './cache'

interface Props {
  name: string
}
const props = defineProps<Props>()
const loading = ref(true)

// Prepare i18n for plugin context
const { t, locale } = useI18n()
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

            // Implement retry mechanism with exponential backoff and context handover
            const checkPluginLoad = (retries = 0, maxRetries = 10) => {
                const globalScope = globalThis as {
                    ATestPlugin?: { mount?: (el: Element, context?: unknown) => void }
                };
                const plugin = globalScope.ATestPlugin;

                if (plugin && typeof plugin.mount === 'function') {
                    console.log('extension load success');
                    const container = document.getElementById('plugin-container');
                    if (container) {
                        container.innerHTML = '';

                        const context = {
                            i18n: { t, locale },
                            API,
                            Cache
                        };

                        try {
                            plugin.mount(container, context);
                        } catch (error) {
                            console.error('extension mount error:', error);
                        }
                    }
                    loading.value = false;
                } else if (retries < maxRetries) {
                    const delay = 50 + retries * 50;
                    setTimeout(() => checkPluginLoad(retries + 1, maxRetries), delay);
                } else {
                    loading.value = false;
                }
            };

            checkPluginLoad();
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
