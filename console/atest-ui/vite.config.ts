import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
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
export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          nodeTransforms: true ? [removeDataTestAttrs] : [],
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
  server: {
    proxy: {
      '/server.Runner': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/server.Mock': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/mock/server': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/browser': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/v3': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/oauth': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
      '/api': {
        target: 'http://127.0.0.1:9090',
        changeOrigin: true,
      },
    },
  },
})
