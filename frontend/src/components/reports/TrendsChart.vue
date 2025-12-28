<template>
  <div class="w-full">
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <Loader2 class="h-6 w-6 text-primary animate-spin" />
    </div>
    <div v-else-if="error" class="text-center py-4 text-destructive">
      Erro ao carregar dados
    </div>
    <div v-else-if="data" class="h-96">
      <!-- Chart will be implemented in FE-REP-006 -->
      <div class="flex items-center justify-center h-full text-muted-foreground">
        Gráfico de Tendências Temporais (em implementação)
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { useReports } from '@/hooks/useReports'

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
</script>

