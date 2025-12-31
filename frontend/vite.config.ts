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
      output: {
        manualChunks: (id) => {
          // Vendor chunks
          if (id.includes('node_modules')) {
            // PrimeVue and related UI libraries
            if (id.includes('primevue') || id.includes('primeicons')) {
              return 'vendor-ui'
            }
            // ApexCharts (only used in reports)
            if (id.includes('apexcharts') || id.includes('vue3-apexcharts')) {
              return 'vendor-charts'
            }
            // PDF/Export libraries (only used when exporting)
            if (id.includes('jspdf') || id.includes('html2canvas')) {
              return 'vendor-export'
            }
            // Vue core and router
            if (id.includes('vue') || id.includes('vue-router') || id.includes('pinia')) {
              return 'vendor-vue'
            }
            // TanStack Query
            if (id.includes('@tanstack')) {
              return 'vendor-query'
            }
            // Validation libraries
            if (id.includes('vee-validate') || id.includes('zod')) {
              return 'vendor-validation'
            }
            // Other large dependencies
            if (id.includes('axios')) {
              return 'vendor-http'
            }
            // All other node_modules
            return 'vendor'
          }
        },
        chunkFileNames: 'assets/[name]-[hash].js',
        entryFileNames: 'assets/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]',
      },
    },
    // Increase chunk size warning limit (we're splitting manually)
    chunkSizeWarningLimit: 600,
    // Enable source maps for production debugging (optional, can disable for smaller bundle)
    sourcemap: false,
    // Minification - esbuild is faster and smaller by default
    minify: 'esbuild',
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
