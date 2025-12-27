<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-lg">Progresso do Orçamento</CardTitle>
      <CardDescription v-if="categoryName">
        {{ categoryName }}
      </CardDescription>
    </CardHeader>
    <CardContent class="space-y-4">
      <!-- Loading State -->
      <div v-if="isLoading" class="flex items-center justify-center py-8">
        <Loader2 class="h-6 w-6 text-primary animate-spin" />
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-4">
        <AlertCircle class="h-6 w-6 text-destructive mx-auto mb-2" />
        <p class="text-sm text-destructive">
          {{ error.message || 'Erro ao carregar progresso' }}
        </p>
      </div>

      <!-- Progress Content -->
      <div v-else-if="progress" class="space-y-4">
        <!-- Progress Bar -->
        <div class="space-y-2">
          <div class="flex items-center justify-between text-sm">
            <span class="text-muted-foreground">Progresso</span>
            <span class="font-semibold" :class="getProgressColor(progress.percentage_used)">
              {{ progress.percentage_used.toFixed(1) }}%
            </span>
          </div>
          <div class="relative h-3 w-full overflow-hidden rounded-full bg-muted">
            <div
              class="h-full transition-all duration-500"
              :class="getProgressBarColor(progress.percentage_used, progress.is_exceeded)"
              :style="{ width: `${Math.min(progress.percentage_used, 100)}%` }"
            />
            <!-- Exceeded indicator -->
            <div
              v-if="progress.is_exceeded"
              class="absolute inset-0 h-full bg-red-500 opacity-50"
              :style="{ width: `${Math.min(progress.percentage_used, 100)}%` }"
            />
          </div>
        </div>

        <!-- Amounts -->
        <div class="grid grid-cols-3 gap-4">
          <div>
            <p class="text-xs text-muted-foreground mb-1">Orçado</p>
            <p class="text-lg font-semibold">
              {{ formatCurrency(progress.budgeted, progress.currency) }}
            </p>
          </div>
          <div>
            <p class="text-xs text-muted-foreground mb-1">Gasto</p>
            <p class="text-lg font-semibold" :class="getSpentColor(progress.is_exceeded)">
              {{ formatCurrency(progress.spent, progress.currency) }}
            </p>
          </div>
          <div>
            <p class="text-xs text-muted-foreground mb-1">Restante</p>
            <p class="text-lg font-semibold" :class="getRemainingColor(progress.remaining)">
              {{ formatCurrency(progress.remaining, progress.currency) }}
            </p>
          </div>
        </div>

        <!-- Status Badge -->
        <div class="flex items-center justify-center">
          <Badge
            :variant="progress.is_exceeded ? 'destructive' : progress.percentage_used >= 80 ? 'default' : 'secondary'"
            class="text-sm"
          >
            <AlertTriangle
              v-if="progress.is_exceeded"
              class="h-4 w-4 mr-1"
            />
            <span v-if="progress.is_exceeded">
              Orçamento Excedido
            </span>
            <span v-else-if="progress.percentage_used >= 80">
              Próximo do Limite
            </span>
            <span v-else>
              Dentro do Orçamento
            </span>
          </Badge>
        </div>

        <!-- Period Info -->
        <div class="pt-2 border-t">
          <p class="text-xs text-muted-foreground text-center">
            Período: {{ formatPeriod(progress) }}
          </p>
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useBudgets } from '@/hooks/useBudgets'
import { useCategoriesStore } from '@/stores/categories'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Loader2, AlertCircle, AlertTriangle } from 'lucide-vue-next'
import type { BudgetProgress } from '@/api/types'

interface Props {
  budgetId: string | null
  categoryId?: string
}

const props = defineProps<Props>()

const { useBudgetProgress } = useBudgets()
const categoriesStore = useCategoriesStore()

// Query de progresso
const {
  data: progress,
  isLoading,
  error,
} = useBudgetProgress(props.budgetId)

// Nome da categoria
const categoryName = computed(() => {
  if (props.categoryId) {
    const category = categoriesStore.categories.find((c) => c.category_id === props.categoryId)
    return category?.name || ''
  }
  if (progress.value?.category_id) {
    const category = categoriesStore.categories.find((c) => c.category_id === progress.value?.category_id)
    return category?.name || ''
  }
  return ''
})

// Funções auxiliares
function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  }).format(amount)
}

function formatPeriod(progress: BudgetProgress): string {
  if (progress.period_type === 'MONTHLY' && progress.month) {
    const monthNames = [
      'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
      'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
    ]
    return `${monthNames[progress.month - 1]} ${progress.year}`
  }
  return `Ano ${progress.year}`
}

function getProgressColor(percentage: number): string {
  if (percentage >= 100) return 'text-red-600'
  if (percentage >= 80) return 'text-yellow-600'
  return 'text-green-600'
}

function getProgressBarColor(percentage: number, isExceeded: boolean): string {
  if (isExceeded || percentage >= 100) return 'bg-red-500'
  if (percentage >= 80) return 'bg-yellow-500'
  return 'bg-green-500'
}

function getSpentColor(isExceeded: boolean): string {
  return isExceeded ? 'text-red-600' : 'text-foreground'
}

function getRemainingColor(remaining: number): string {
  if (remaining < 0) return 'text-red-600'
  if (remaining === 0) return 'text-yellow-600'
  return 'text-green-600'
}
</script>

