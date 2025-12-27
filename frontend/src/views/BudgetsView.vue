<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Orçamentos' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Orçamentos</h1>
          <p class="text-muted-foreground">
            Gerencie seus orçamentos mensais e anuais
          </p>
        </div>
        <Button as-child>
          <router-link to="/budgets/new">
            <Plus class="h-4 w-4 mr-2" />
            Novo Orçamento
          </router-link>
        </Button>
      </div>

      <!-- Filters -->
      <Card class="mb-6">
        <CardContent class="p-4">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
            <div>
              <label class="text-sm font-medium mb-2 block">Período</label>
              <Select v-model="filters.period_type" @update:model-value="handleFilterChange">
                <SelectTrigger>
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todos</SelectItem>
                  <SelectItem value="MONTHLY">Mensal</SelectItem>
                  <SelectItem value="YEARLY">Anual</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <label class="text-sm font-medium mb-2 block">Ano</label>
              <Input
                v-model.number="filters.year"
                type="number"
                :min="2020"
                :max="2100"
                placeholder="Ano"
                @input="handleFilterChange"
              />
            </div>
            <div>
              <label class="text-sm font-medium mb-2 block">Mês</label>
              <Select :model-value="filters.month?.toString() || ''" @update:model-value="handleMonthChange">
                <SelectTrigger>
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todos</SelectItem>
                  <SelectItem :value="1">Janeiro</SelectItem>
                  <SelectItem :value="2">Fevereiro</SelectItem>
                  <SelectItem :value="3">Março</SelectItem>
                  <SelectItem :value="4">Abril</SelectItem>
                  <SelectItem :value="5">Maio</SelectItem>
                  <SelectItem :value="6">Junho</SelectItem>
                  <SelectItem :value="7">Julho</SelectItem>
                  <SelectItem :value="8">Agosto</SelectItem>
                  <SelectItem :value="9">Setembro</SelectItem>
                  <SelectItem :value="10">Outubro</SelectItem>
                  <SelectItem :value="11">Novembro</SelectItem>
                  <SelectItem :value="12">Dezembro</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <label class="text-sm font-medium mb-2 block">Contexto</label>
              <Select v-model="filters.context" @update:model-value="handleFilterChange">
                <SelectTrigger>
                  <SelectValue placeholder="Todos" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Todos</SelectItem>
                  <SelectItem value="PERSONAL">Pessoal</SelectItem>
                  <SelectItem value="BUSINESS">Negócio</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Loading State -->
      <div v-if="isLoadingBudgets" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando orçamentos...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="budgetsError"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">
              {{ (budgetsError as any)?.message || 'Erro ao carregar orçamentos' }}
            </p>
          </div>
          <Button
            variant="link"
            @click="refetchBudgets"
            class="text-destructive"
          >
            Tentar novamente
          </Button>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <EmptyState
        v-else-if="budgets.length === 0"
        :icon="Target"
        title="Nenhum orçamento encontrado"
        description="Comece criando seu primeiro orçamento"
        action-label="Criar Primeiro Orçamento"
        action-to="/budgets/new"
        :action-icon="Plus"
      />

      <!-- Budgets List -->
      <div v-else class="space-y-4">
        <!-- Stats -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-green-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500 text-white">
                    <Target class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-green-700">Total de Orçamentos</div>
                <div class="text-3xl font-bold text-green-900">
                  {{ total }}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-blue-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500 text-white">
                    <Calendar class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-blue-700">Orçamentos Mensais</div>
                <div class="text-3xl font-bold text-blue-900">
                  {{ monthlyBudgets.length }}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-purple-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-500 text-white">
                    <Calendar class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-purple-700">Orçamentos Anuais</div>
                <div class="text-3xl font-bold text-purple-900">
                  {{ yearlyBudgets.length }}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Budgets Grid -->
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          <Card
            v-for="budget in budgets"
            :key="budget.budget_id"
            class="hover:shadow-md transition-shadow cursor-pointer"
            @click="goToBudgetDetails(budget.budget_id)"
          >
            <CardHeader>
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <CardTitle class="text-lg">
                    {{ getCategoryName(budget.category_id) }}
                  </CardTitle>
                  <CardDescription class="mt-1">
                    {{ formatPeriod(budget) }}
                  </CardDescription>
                </div>
                <Badge :variant="budget.is_active ? 'default' : 'secondary'">
                  {{ budget.is_active ? 'Ativo' : 'Inativo' }}
                </Badge>
              </div>
            </CardHeader>
            <CardContent>
              <div class="space-y-2">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-muted-foreground">Valor:</span>
                  <span class="text-lg font-semibold">
                    {{ formatCurrency(budget.amount, budget.currency) }}
                  </span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-muted-foreground">Contexto:</span>
                  <Badge variant="outline">
                    {{ budget.context === 'PERSONAL' ? 'Pessoal' : 'Negócio' }}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBudgets } from '@/hooks/useBudgets'
import { useCategoriesStore } from '@/stores/categories'
import { useBudgetAlerts } from '@/composables/useBudgetAlerts'
import Layout from '@/components/layout/Layout.vue'
import EmptyState from '@/components/EmptyState.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Target, Plus, Loader2, AlertCircle, Calendar } from 'lucide-vue-next'
import type { Budget } from '@/api/types'
import type { AcceptableValue } from 'reka-ui'

const router = useRouter()
const categoriesStore = useCategoriesStore()

// Filtros
const filters = ref<{
  period_type?: 'MONTHLY' | 'YEARLY' | ''
  year?: number
  month?: number
  context?: 'PERSONAL' | 'BUSINESS' | ''
}>({
  period_type: '',
  year: new Date().getFullYear(),
  month: undefined,
  context: '',
})

// Query params para o hook
const queryParams = computed(() => {
  const params: any = {}
  if (filters.value.period_type) {
    params.period_type = filters.value.period_type
  }
  if (filters.value.year) {
    params.year = filters.value.year
  }
  if (filters.value.month) {
    params.month = filters.value.month
  }
  if (filters.value.context) {
    params.context = filters.value.context
  }
  return Object.keys(params).length > 0 ? params : undefined
})

// Hook de orçamentos
const {
  budgets,
  total,
  isLoadingBudgets,
  budgetsError,
  refetchBudgets,
} = useBudgets(queryParams.value)

// Computed
const monthlyBudgets = computed(() => {
  const budgetsList = Array.isArray(budgets) ? budgets : []
  return budgetsList.filter((b: Budget) => b.period_type === 'MONTHLY')
})

const yearlyBudgets = computed(() => {
  const budgetsList = Array.isArray(budgets) ? budgets : []
  return budgetsList.filter((b: Budget) => b.period_type === 'YEARLY')
})

// Carregar categorias se necessário
if (categoriesStore.categories.length === 0) {
  categoriesStore.listCategories(true)
}

// Funções auxiliares
function formatCurrency(amount: number, currency: string): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  }).format(amount)
}

function formatPeriod(budget: Budget): string {
  if (budget.period_type === 'MONTHLY' && budget.month) {
    const monthNames = [
      'Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
      'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro'
    ]
    return `${monthNames[budget.month - 1]} ${budget.year}`
  }
  return `Ano ${budget.year}`
}

function getCategoryName(categoryId: string): string {
  const category = categoriesStore.categories.find((c) => c.category_id === categoryId)
  return category?.name || 'Categoria não encontrada'
}

function goToBudgetDetails(budgetId: string) {
  router.push(`/budgets/${budgetId}`)
}

function handleFilterChange() {
  // O hook vai refetch automaticamente quando queryParams mudar
}

function handleMonthChange(value: AcceptableValue) {
  filters.value.month = value ? parseInt(String(value), 10) : undefined
  handleFilterChange()
}

// Watch para refetch quando filtros mudarem
watch(queryParams, () => {
  refetchBudgets()
}, { deep: true })

// Monitorar alertas de orçamento
const { monitorAllBudgets, clearAlerts } = useBudgetAlerts()
let cleanupAlerts: (() => void) | null = null

onMounted(() => {
  // Iniciar monitoramento de alertas após um pequeno delay
  // para garantir que os dados foram carregados
  setTimeout(() => {
    cleanupAlerts = monitorAllBudgets()
  }, 1000)
})

onUnmounted(() => {
  if (cleanupAlerts) {
    cleanupAlerts()
  }
  clearAlerts()
})
</script>

