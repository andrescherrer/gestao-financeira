<template>
  <Card
    class="group relative cursor-pointer overflow-hidden transition-all hover:border-primary hover:shadow-lg"
    @click="handleClick"
    role="button"
    :aria-label="`Meta ${goal.name}, progresso ${goal.progress.toFixed(0)}%`"
    tabindex="0"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <!-- Background gradient effect -->
    <div
      class="absolute inset-0 opacity-0 transition-opacity group-hover:opacity-5"
      :class="getGoalGradient(goal.status)"
    ></div>

    <CardContent class="relative p-6">
      <div class="mb-4 flex items-start justify-between">
        <div class="flex items-center gap-3">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-lg"
            :class="getGoalIconBg(goal.status)"
          >
            <component :is="getGoalIcon(goal.status)" class="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-foreground">
              {{ goal.name }}
            </h3>
            <p class="text-sm text-muted-foreground">
              {{ getStatusLabel(goal.status) }}
            </p>
          </div>
        </div>
        <Badge
          :variant="getStatusBadgeVariant(goal.status)"
        >
          {{ getStatusLabel(goal.status) }}
        </Badge>
      </div>

      <!-- Progress Bar -->
      <div class="mb-4">
        <div class="mb-2 flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Progresso</span>
          <span class="font-semibold">{{ goal.progress.toFixed(0) }}%</span>
        </div>
        <div class="h-2 w-full overflow-hidden rounded-full bg-muted">
          <div
            class="h-full transition-all duration-300"
            :class="getProgressBarColor(goal.status)"
            :style="{ width: `${Math.min(goal.progress, 100)}%` }"
          ></div>
        </div>
      </div>

      <div class="mb-4 space-y-2">
        <div class="rounded-lg bg-muted p-4">
          <div class="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-1">
            Valor Atual
          </div>
          <div class="text-2xl font-bold">
            {{ formatCurrency(goal.current_amount, goal.currency) }}
          </div>
        </div>
        <div class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Meta:</span>
          <span class="font-medium">
            {{ formatCurrency(goal.target_amount, goal.currency) }}
          </span>
        </div>
        <div class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Restante:</span>
          <span class="font-medium">
            {{ formatCurrency(getRemainingAmount(), goal.currency) }}
          </span>
        </div>
        <div v-if="goal.remaining_days > 0" class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Dias restantes:</span>
          <span class="font-medium">{{ goal.remaining_days }} dias</span>
        </div>
        <div v-else-if="goal.status === 'OVERDUE'" class="flex items-center justify-between text-sm text-destructive">
          <span>Vencida</span>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <Badge
          :variant="goal.context === 'PERSONAL' ? 'secondary' : 'outline'"
          class="bg-purple-100 text-purple-700"
        >
          {{ getContextLabel(goal.context) }}
        </Badge>
        <ChevronRight
          class="h-5 w-5 text-muted-foreground transition-all group-hover:translate-x-1 group-hover:text-primary"
        />
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Target,
  CheckCircle2,
  AlertCircle,
  XCircle,
  ChevronRight,
} from 'lucide-vue-next'
import type { Goal } from '@/api/types'

interface Props {
  goal: Goal
}

const props = defineProps<Props>()
const router = useRouter()

function handleClick() {
  router.push(`/goals/${props.goal.goal_id}`)
}

function getStatusLabel(status: Goal['status']): string {
  const labels: Record<Goal['status'], string> = {
    IN_PROGRESS: 'Em Progresso',
    COMPLETED: 'Concluída',
    OVERDUE: 'Vencida',
    CANCELLED: 'Cancelada',
  }
  return labels[status] || status
}

function getContextLabel(context: Goal['context']): string {
  return context === 'PERSONAL' ? 'Pessoal' : 'Negócio'
}

function formatCurrency(amount: string, currency: Goal['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function getRemainingAmount(): string {
  const target = parseFloat(props.goal.target_amount)
  const current = parseFloat(props.goal.current_amount)
  const remaining = Math.max(0, target - current)
  return remaining.toString()
}

function getGoalIcon(status: Goal['status']) {
  const icons: Record<Goal['status'], any> = {
    IN_PROGRESS: Target,
    COMPLETED: CheckCircle2,
    OVERDUE: AlertCircle,
    CANCELLED: XCircle,
  }
  return icons[status] || Target
}

function getGoalIconBg(status: Goal['status']): string {
  const colors: Record<Goal['status'], string> = {
    IN_PROGRESS: 'bg-blue-500',
    COMPLETED: 'bg-green-500',
    OVERDUE: 'bg-red-500',
    CANCELLED: 'bg-gray-500',
  }
  return colors[status] || 'bg-blue-500'
}

function getGoalGradient(status: Goal['status']): string {
  const gradients: Record<Goal['status'], string> = {
    IN_PROGRESS: 'bg-gradient-to-br from-blue-500 to-blue-600',
    COMPLETED: 'bg-gradient-to-br from-green-500 to-green-600',
    OVERDUE: 'bg-gradient-to-br from-red-500 to-red-600',
    CANCELLED: 'bg-gradient-to-br from-gray-500 to-gray-600',
  }
  return gradients[status] || 'bg-gradient-to-br from-blue-500 to-blue-600'
}

function getStatusBadgeVariant(status: Goal['status']): 'default' | 'secondary' | 'destructive' | 'outline' {
  const variants: Record<Goal['status'], 'default' | 'secondary' | 'destructive' | 'outline'> = {
    IN_PROGRESS: 'default',
    COMPLETED: 'default',
    OVERDUE: 'destructive',
    CANCELLED: 'outline',
  }
  return variants[status] || 'default'
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

