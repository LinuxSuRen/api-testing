<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from './net'

interface Props {
  name: string
}
const props = defineProps<Props>()
const { locale } = useI18n()

let pluginInstance: { setLocale?: (value: string) => void } | undefined

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
                const globalScope = globalThis as { ATestPlugin?: { mount?: (el: Element) => void, setLocale?: (value: string) => void } };
                const plugin = globalScope.ATestPlugin;

                if (plugin && plugin.mount) {
                    const container = document.getElementById("plugin-container");
                    if (container) {
                        container.innerHTML = ''; // Clear previous content
                        plugin.mount(container);
                        plugin.setLocale?.(locale.value);
                        pluginInstance = plugin;
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

watch(locale, (value) => {
    const normalized = value ?? 'en';
    pluginInstance?.setLocale?.(normalized);
});
</script>

<template>
    <div id="plugin-container"
        v-loading="loading"
        :element-loading-text="props.name + ' is loading...'">
    </div>
</template>
