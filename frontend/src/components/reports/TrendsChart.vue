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
        type="line"
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
  currency: 'BRL' | 'USD' | 'EUR'
  period: 'monthly' | 'annual' | 'custom'
}

const props = defineProps<Props>()

const { useAnnualReport } = useReports()

const { data, isLoading, error } = useAnnualReport({
  year: props.year,
  currency: props.currency,
})

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: props.currency || 'BRL',
  }).format(value)
}

const monthNames = [
  'Jan', 'Fev', 'Mar', 'Abr', 'Mai', 'Jun',
  'Jul', 'Ago', 'Set', 'Out', 'Nov', 'Dez'
]

const chartSeries = computed(() => {
  if (!data.value || !data.value.monthly_breakdown) return []

  return [
    {
      name: 'Receitas',
      data: data.value.monthly_breakdown.map((m) => m.total_income),
    },
    {
      name: 'Despesas',
      data: data.value.monthly_breakdown.map((m) => m.total_expense),
    },
    {
      name: 'Saldo',
      data: data.value.monthly_breakdown.map((m) => m.balance),
    },
  ]
})

const chartOptions = computed<ApexOptions>(() => {
  if (!data.value) return {} as ApexOptions

  const categories = data.value.monthly_breakdown
    ? data.value.monthly_breakdown.map((m) => monthNames[m.month - 1])
    : []

  return {
    chart: {
      type: 'line',
      toolbar: {
        show: true,
      },
      zoom: {
        enabled: true,
      },
    },
    stroke: {
      curve: 'smooth',
      width: 3,
    },
    markers: {
      size: 5,
    },
    xaxis: {
      categories,
    },
    yaxis: {
      labels: {
        formatter: (value: number) => formatCurrency(value),
      },
    },
    colors: ['#10b981', '#ef4444', '#3b82f6'], // Verde (receitas), Vermelho (despesas), Azul (saldo)
    legend: {
      position: 'top',
      horizontalAlign: 'right',
    },
    tooltip: {
      y: {
        formatter: (value: number) => formatCurrency(value),
      },
    },
    grid: {
      borderColor: '#e5e7eb',
    },
  } as ApexOptions
})
</script>


