import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { VitePWA } from 'vite-plugin-pwa'
import { reactMockPlugin } from './vite.config.plugins'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    reactMockPlugin(), // Deve vir antes dos outros plugins
    vue(),
    vueDevTools(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'mask-icon.svg'],
      manifest: {
        name: 'Gestão Financeira',
        short_name: 'Gestão Fin',
        description: 'Sistema de gestão financeira pessoal',
        theme_color: '#ffffff',
        background_color: '#ffffff',
        display: 'standalone',
        orientation: 'portrait',
        scope: '/',
        start_url: '/',
        icons: [
          {
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable',
          },
        ],
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff,woff2}'],
        runtimeCaching: [
          {
            urlPattern: /^https:\/\/api\./,
            handler: 'NetworkFirst',
            options: {
              cacheName: 'api-cache',
              expiration: {
                maxEntries: 50,
                maxAgeSeconds: 5 * 60, // 5 minutes
              },
              cacheableResponse: {
                statuses: [0, 200],
              },
            },
          },
          {
            urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp)$/,
            handler: 'CacheFirst',
            options: {
              cacheName: 'images-cache',
              expiration: {
                maxEntries: 100,
                maxAgeSeconds: 30 * 24 * 60 * 60, // 30 days
              },
            },
          },
        ],
      },
      devOptions: {
        enabled: false, // Desabilitar em desenvolvimento por padrão
      },
    }),
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
