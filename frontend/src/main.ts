import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import VueApexCharts from 'vue3-apexcharts'
import App from './App.vue'
import router from './router'
import './config/vee-validate'
import { logger } from './utils/logger'

// Configure TanStack Query
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // 5 minutes
    },
  },
})

// Configurar captura de erros globais
if (typeof window !== 'undefined') {
  // Capturar erros não tratados
  window.addEventListener('error', (event) => {
    logger.error('Unhandled error', event.error || event.message, {
      filename: event.filename,
      lineno: event.lineno,
      colno: event.colno,
    })
  })

  // Capturar promessas rejeitadas
  window.addEventListener('unhandledrejection', (event) => {
    logger.error('Unhandled promise rejection', event.reason, {
      type: 'unhandledrejection',
    })
  })

  // Log de inicialização
  logger.info('Application initialized', {
    userAgent: navigator.userAgent,
    url: window.location.href,
  })
}

const app = createApp(App)

// Configurar error handler global do Vue
app.config.errorHandler = (err, instance, info) => {
  logger.error('Vue error', err, {
    component: instance?.$options.name || 'Unknown',
    info,
  })
}

app.use(createPinia())
app.use(router)
app.use(VueQueryPlugin, { queryClient })
app.use(VueApexCharts)

app.mount('#app')
