<script setup lang="ts">
import { API } from './net';

const props = defineProps({
  name: String
})

const loadPlugin = async (): Promise<void> => {
    try {
        API.GetPageOfCSS(props.name, (d) => {
            const style = document.createElement('style');
            style.type = 'text/css';
            style.textContent = d.message;
            document.head.appendChild(style);
        });

        API.GetPageOfJS(props.name, (d) => {
            const script = document.createElement('script');
            script.type = 'text/javascript';
            script.textContent = d.message;
            document.head.appendChild(script);

            // 类型安全的插件访问
            const plugin = window.ATestPlugin;
            
            if (plugin && plugin.mount) {
                console.log('插件加载成功');
                plugin.mount('#plugin-container', { 
                    message: '来自宿主的消息'
                });
            }
        });
    } catch (error) {
        console.log(`加载失败: ${(error as Error).message}`)
    } finally {
        console.log('插件加载完成');
    }
};
try {
    loadPlugin();
} catch (error) {
    console.error('插件加载失败:', error);
}
</script>

<template>
{{ props.name }}===
<div id="plugin-container"></div>
</template>
