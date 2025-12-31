import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { reactMockPlugin } from './vite.config.plugins'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    reactMockPlugin(), // Deve vir antes dos outros plugins
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  assetsInclude: ['**/*.woff', '**/*.woff2', '**/*.ttf', '**/*.eot', '**/*.svg'],
  optimizeDeps: {
    exclude: ['cypress', 'cypress-axe'],
  },
  build: {
    rollupOptions: {
      external: (id) => {
        // Exclude Cypress from build
        if (id.includes('cypress')) {
          return true
        }
        return false
      },
    },
  },
  // Exclude Cypress files from build
  publicDir: 'public',
  preview: {
    host: '0.0.0.0',
    port: 4173,
    strictPort: false,
    // Allow access from any host (needed for Docker containers)
    cors: true,
  },
})
