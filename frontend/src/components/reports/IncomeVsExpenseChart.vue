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
        type="bar"
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
  period: 'monthly' | 'annual' | 'custom'
}

const props = defineProps<Props>()

const { useIncomeVsExpenseReport } = useReports()

const { data, isLoading, error } = useIncomeVsExpenseReport({
  currency: props.currency,
  group_by: props.period === 'monthly' ? 'month' : props.period === 'annual' ? 'month' : 'month',
})

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: props.currency || 'BRL',
  }).format(value)
}

const chartSeries = computed(() => {
  if (!data.value) return []

  if (data.value.period_breakdown && data.value.period_breakdown.length > 0) {
    // Se há breakdown por período, usar dados do breakdown
    return [
      {
        name: 'Receitas',
        data: data.value.period_breakdown.map((p) => p.total_income),
      },
      {
        name: 'Despesas',
        data: data.value.period_breakdown.map((p) => p.total_expense),
      },
    ]
  }

  // Caso contrário, usar totais gerais
  return [
    {
      name: 'Receitas',
      data: [data.value.total_income],
    },
    {
      name: 'Despesas',
      data: [data.value.total_expense],
    },
  ]
})

const chartOptions = computed<ApexOptions>(() => {
  if (!data.value) return {} as ApexOptions

  const categories = data.value.period_breakdown && data.value.period_breakdown.length > 0
    ? data.value.period_breakdown.map((p) => p.period)
    : ['Total']

  return {
    chart: {
      type: 'bar',
      toolbar: {
        show: true,
      },
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: '55%',
        borderRadius: 4,
      },
    },
    dataLabels: {
      enabled: false,
    },
    stroke: {
      show: true,
      width: 2,
      colors: ['transparent'],
    },
    xaxis: {
      categories,
    },
    yaxis: {
      labels: {
        formatter: (value: number) => formatCurrency(value),
      },
    },
    fill: {
      opacity: 1,
    },
    colors: ['#10b981', '#ef4444'], // Verde para receitas, vermelho para despesas
    legend: {
      position: 'top',
      horizontalAlign: 'right',
    },
    tooltip: {
      y: {
        formatter: (value: number) => formatCurrency(value),
      },
    },
  } as ApexOptions
})
</script>


