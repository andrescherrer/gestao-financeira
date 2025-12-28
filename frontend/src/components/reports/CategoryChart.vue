<template>
  <div class="w-full">
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <Loader2 class="h-6 w-6 text-primary animate-spin" />
    </div>
    <div v-else-if="error" class="text-center py-4 text-destructive">
      Erro ao carregar dados: {{ (error as any)?.message || 'Erro desconhecido' }}
    </div>
    <div v-else-if="data && chartOptions" class="w-full">
      <apexchart
        type="donut"
        height="400"
        :options="chartOptions"
        :series="chartSeries"
      />
    </div>
    <div v-else class="flex items-center justify-center py-12 text-muted-foreground">
      Nenhum dado disponível
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import { useReports } from '@/hooks/useReports'
import type { ApexOptions } from 'apexcharts'

interface Props {
  year: number
  month?: number
  currency: 'BRL' | 'USD' | 'EUR'
}

const props = defineProps<Props>()

const { useCategoryReport } = useReports()

const { data, isLoading, error } = useCategoryReport({
  currency: props.currency,
})

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: props.currency || 'BRL',
  }).format(value)
}

const chartSeries = computed(() => {
  if (!data.value || !data.value.category_breakdown) return []

  // Filtrar apenas despesas para o gráfico
  const expenses = data.value.category_breakdown.filter((c) => c.type === 'EXPENSE')
  return expenses.map((c) => c.total_amount)
})

const chartLabels = computed(() => {
  if (!data.value || !data.value.category_breakdown) return []

  const expenses = data.value.category_breakdown.filter((c) => c.type === 'EXPENSE')
  return expenses.map((c) => c.category_name || 'Sem categoria')
})

const chartOptions = computed<ApexOptions>(() => {
  if (!data.value) return {} as ApexOptions

  return {
    chart: {
      type: 'donut',
      toolbar: {
        show: true,
      },
    },
    labels: chartLabels.value,
    legend: {
      position: 'bottom',
    },
    plotOptions: {
      pie: {
        donut: {
          size: '65%',
          labels: {
            show: true,
            name: {
              show: true,
              fontSize: '14px',
            },
            value: {
              show: true,
              fontSize: '16px',
              fontWeight: 600,
              formatter: (val: string) => formatCurrency(parseFloat(val)),
            },
            total: {
              show: true,
              label: 'Total',
              fontSize: '16px',
              fontWeight: 600,
              formatter: () => {
                const total = data.value?.total_expense || 0
                return formatCurrency(total)
              },
            },
          },
        },
      },
    },
    tooltip: {
      y: {
        formatter: (value: number) => formatCurrency(value),
      },
    },
    colors: [
      '#3b82f6', // blue
      '#10b981', // green
      '#f59e0b', // amber
      '#ef4444', // red
      '#8b5cf6', // purple
      '#ec4899', // pink
      '#06b6d4', // cyan
      '#f97316', // orange
    ],
  } as ApexOptions
})
</script>


