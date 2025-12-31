<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Investimentos' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl sm:text-4xl font-bold mb-2">Investimentos</h1>
          <p class="text-sm sm:text-base text-muted-foreground">
            Gerencie seus investimentos e acompanhe a rentabilidade
          </p>
        </div>
        <Button as-child class="w-full sm:w-auto">
          <router-link to="/investments/new">
            <Plus class="h-4 w-4 mr-2" />
            Novo Investimento
          </router-link>
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="investmentsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando investimentos...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="investmentsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ investmentsStore.error }}</p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="link"
              @click="handleRetry"
              class="text-destructive"
            >
              Tentar novamente
            </Button>
            <Button
              v-if="isAuthError"
              variant="outline"
              @click="handleLogin"
              class="text-destructive border-destructive"
            >
              Fazer login novamente
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <EmptyState
        v-else-if="investmentsStore.investments.length === 0"
        :icon="TrendingUp"
        title="Nenhum investimento encontrado"
        description="Comece adicionando seu primeiro investimento para acompanhar sua carteira."
      >
        <Button as-child>
          <router-link to="/investments/new">
            <Plus class="h-4 w-4 mr-2" />
            Novo Investimento
          </router-link>
        </Button>
      </EmptyState>

      <!-- Investments Grid -->
      <div v-else class="space-y-6">
        <!-- Summary Cards -->
        <div class="grid gap-4 md:grid-cols-3">
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Total Investido
              </div>
              <div class="text-2xl font-bold">
                {{ formatCurrency(investmentsStore.totalValue.toString(), 'BRL') }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Retorno Total
              </div>
              <div
                class="text-2xl font-bold"
                :class="investmentsStore.totalReturn >= 0 ? 'text-green-600' : 'text-red-600'"
              >
                {{ formatCurrency(investmentsStore.totalReturn.toString(), 'BRL') }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Total de Investimentos
              </div>
              <div class="text-2xl font-bold">
                {{ investmentsStore.totalInvestments }}
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Investments List -->
        <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <InvestmentCard
            v-for="investment in investmentsStore.investments"
            :key="investment.investment_id"
            :investment="investment"
          />
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useInvestmentsStore } from '@/stores/investments'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import InvestmentCard from '@/components/InvestmentCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus, Loader2, AlertCircle, TrendingUp } from 'lucide-vue-next'

const router = useRouter()
const investmentsStore = useInvestmentsStore()
const authStore = useAuthStore()

const isAuthError = computed(() => {
  const error = investmentsStore.error?.toLowerCase() || ''
  return error.includes('token') || 
         error.includes('autenticação') || 
         error.includes('unauthorized') ||
         error.includes('invalid') ||
         error.includes('expired')
})

function formatCurrency(amount: string, currency: string): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

onMounted(async () => {
  if (!authStore.token) {
    authStore.init()
  }
  
  if (authStore.isValidating) {
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  
  if (investmentsStore.investments.length === 0 && !investmentsStore.isLoading) {
    await investmentsStore.listInvestments()
  }
})

function handleRetry() {
  investmentsStore.clearError()
  investmentsStore.listInvestments()
}

function handleLogin() {
  authStore.logout()
  router.push({ name: 'login', query: { redirect: '/investments' } })
}
</script>

