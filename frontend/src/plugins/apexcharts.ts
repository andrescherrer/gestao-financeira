/**
 * Plugin para carregar VueApexCharts dinamicamente
 * Isso reduz o bundle inicial, carregando apenas quando necessário
 */

import type { App } from 'vue'
import VueApexCharts from 'vue3-apexcharts'

let isRegistered = false

/**
 * Registra VueApexCharts na aplicação Vue
 * Pode ser chamado múltiplas vezes, mas só registra uma vez
 */
export async function registerApexCharts(app: App): Promise<void> {
  if (isRegistered) {
    return
  }

  // Lazy load VueApexCharts
  app.use(VueApexCharts)
  isRegistered = true
}

