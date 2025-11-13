import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'

function removeDataTestAttrs(node) {
  if (node.type === 1 /* NodeTypes.ELEMENT */) {
    node.props = node.props.filter(prop =>
      prop.type === 6 /* NodeTypes.ATTRIBUTE */
        ? prop.name !== 'test-id'
        : true
    )
  }
}

// https://vitejs.dev/config/
export default defineConfig(({mode}) => {
    const env = loadEnv(mode, './');
    return {
  plugins: [
    vue({
      template: {
        compilerOptions: {
          nodeTransforms: process.env.NODE_ENV === 'production' ? [removeDataTestAttrs] : [],
        },
      },
    }),
    vueJsx(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: () => 'everything'
      }
    }
  },
  server: {
    proxy: {
      '/server.Runner': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/server.Mock': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/mock/server': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/browser': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/v3': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/oauth': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/api': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
      '/extensionProxy': {
          target: env.VITE_API_URL,
          changeOrigin: true,
      },
      '/data': {
        target: env.VITE_API_URL,
        changeOrigin: true,
      },
    },
  },
}});
