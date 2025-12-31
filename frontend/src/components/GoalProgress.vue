<template>
  <div class="space-y-4">
    <!-- Progress Bar -->
    <div>
      <div class="mb-2 flex items-center justify-between text-sm">
        <span class="text-muted-foreground">Progresso</span>
        <span class="font-semibold">{{ progress.toFixed(0) }}%</span>
      </div>
      <div class="h-3 w-full overflow-hidden rounded-full bg-muted">
        <div
          class="h-full transition-all duration-500"
          :class="getProgressBarColor(status)"
          :style="{ width: `${Math.min(progress, 100)}%` }"
        ></div>
      </div>
    </div>

    <!-- Amounts -->
    <div class="grid grid-cols-2 gap-4">
      <div class="rounded-lg bg-muted p-4">
        <div class="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-1">
          Valor Atual
        </div>
        <div class="text-xl font-bold">
          {{ formatCurrency(currentAmount, currency) }}
        </div>
      </div>
      <div class="rounded-lg bg-muted p-4">
        <div class="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-1">
          Meta
        </div>
        <div class="text-xl font-bold">
          {{ formatCurrency(targetAmount, currency) }}
        </div>
      </div>
    </div>

    <!-- Remaining -->
    <div class="rounded-lg border p-4">
      <div class="text-sm text-muted-foreground mb-1">
        Restante para alcançar a meta
      </div>
      <div class="text-2xl font-bold">
        {{ formatCurrency(remainingAmount, currency) }}
      </div>
    </div>

    <!-- Days Remaining -->
    <div v-if="remainingDays > 0" class="rounded-lg border p-4">
      <div class="text-sm text-muted-foreground mb-1">
        Dias restantes
      </div>
      <div class="text-2xl font-bold">
        {{ remainingDays }} {{ remainingDays === 1 ? 'dia' : 'dias' }}
      </div>
    </div>
    <div v-else-if="status === 'OVERDUE'" class="rounded-lg border border-destructive bg-destructive/10 p-4">
      <div class="text-sm text-destructive font-medium">
        Meta vencida
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Goal } from '@/api/types'

interface Props {
  goal: Goal
}

const props = defineProps<Props>()

const progress = computed(() => props.goal.progress)
const status = computed(() => props.goal.status)
const currency = computed(() => props.goal.currency)
const currentAmount = computed(() => props.goal.current_amount)
const targetAmount = computed(() => props.goal.target_amount)
const remainingDays = computed(() => props.goal.remaining_days)

const remainingAmount = computed(() => {
  const target = parseFloat(props.goal.target_amount)
  const current = parseFloat(props.goal.current_amount)
  return Math.max(0, target - current).toString()
})

function formatCurrency(amount: string, currency: Goal['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function getProgressBarColor(status: Goal['status']): string {
  const colors: Record<Goal['status'], string> = {
    IN_PROGRESS: 'bg-blue-500',
    COMPLETED: 'bg-green-500',
    OVERDUE: 'bg-red-500',
    CANCELLED: 'bg-gray-500',
  }
  return colors[status] || 'bg-blue-500'
}
</script>

