<template>
  <Layout>
    <div class="space-y-6">
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Dashboard' }]" />

      <!-- Header -->
      <div>
        <h1 class="text-4xl font-bold text-foreground mb-2">Dashboard</h1>
        <p class="text-muted-foreground">
          Visão geral da sua situação financeira
        </p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <!-- Total de Contas -->
        <Card class="group relative overflow-hidden">
          <CardContent class="p-6">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-blue-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-500 text-white">
                  <Wallet class="h-6 w-6" />
                </div>
              </div>
              <div class="text-sm font-medium text-blue-700">Total de Contas</div>
              <div class="text-3xl font-bold text-blue-900">
                {{ accountsStore.totalAccounts }}
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Contas Pessoais -->
        <Card class="group relative overflow-hidden">
          <CardContent class="p-6">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-purple-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-purple-500 text-white">
                  <User class="h-6 w-6" />
                </div>
              </div>
              <div class="text-sm font-medium text-purple-700">Contas Pessoais</div>
              <div class="text-3xl font-bold text-purple-900">
                {{ accountsStore.personalAccounts.length }}
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Contas de Negócio -->
        <Card class="group relative overflow-hidden">
          <CardContent class="p-6">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-indigo-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-500 text-white">
                  <Briefcase class="h-6 w-6" />
                </div>
              </div>
              <div class="text-sm font-medium text-indigo-700">Contas de Negócio</div>
              <div class="text-3xl font-bold text-indigo-900">
                {{ accountsStore.businessAccounts.length }}
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Total de Transações -->
        <Card class="group relative overflow-hidden">
          <CardContent class="p-6">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-green-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-500 text-white">
                  <List class="h-6 w-6" />
                </div>
              </div>
              <div class="text-sm font-medium text-green-700">Transações</div>
              <div class="text-3xl font-bold text-green-900">
                {{ transactionsStore.totalTransactions }}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Financial Overview -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Receitas vs Despesas -->
        <Card class="lg:col-span-2">
          <CardHeader>
            <CardTitle>Resumo Financeiro</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-4">
              <div class="flex items-center justify-between rounded-lg bg-green-50 p-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500 text-white">
                    <ArrowDown class="h-5 w-5" />
                  </div>
                  <div>
                    <div class="text-sm font-medium text-muted-foreground">Total de Receitas</div>
                    <div class="text-xl font-bold text-green-600">
                      {{ formatCurrency(transactionsStore.totalIncome) }}
                    </div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs text-muted-foreground">{{ transactionsStore.incomeTransactions.length }} transações</div>
                </div>
              </div>

              <div class="flex items-center justify-between rounded-lg bg-red-50 p-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-red-500 text-white">
                    <ArrowUp class="h-5 w-5" />
                  </div>
                  <div>
                    <div class="text-sm font-medium text-muted-foreground">Total de Despesas</div>
                    <div class="text-xl font-bold text-red-600">
                      {{ formatCurrency(transactionsStore.totalExpense) }}
                    </div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xs text-muted-foreground">{{ transactionsStore.expenseTransactions.length }} transações</div>
                </div>
              </div>

              <div class="flex items-center justify-between rounded-lg bg-blue-50 p-4 border-2 border-blue-200">
                <div class="flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500 text-white">
                    <TrendingUp class="h-5 w-5" />
                  </div>
                  <div>
                    <div class="text-sm font-medium text-muted-foreground">Saldo</div>
                    <div
                      class="text-xl font-bold"
                      :class="transactionsStore.balance >= 0 ? 'text-green-600' : 'text-red-600'"
                    >
                      {{ formatCurrency(transactionsStore.balance) }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Quick Actions -->
        <Card>
          <CardHeader>
            <CardTitle>Ações Rápidas</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="space-y-3">
              <Button
                variant="outline"
                class="w-full justify-start h-auto p-4"
                as-child
              >
                <router-link to="/accounts/new">
                  <div class="flex items-center gap-3 w-full">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                      <Plus class="h-5 w-5" />
                    </div>
                    <div class="flex-1 text-left">
                      <div class="text-sm font-medium text-foreground">Nova Conta</div>
                      <div class="text-xs text-muted-foreground">Criar conta financeira</div>
                    </div>
                    <ChevronRight class="h-4 w-4 text-muted-foreground" />
                  </div>
                </router-link>
              </Button>

              <Button
                variant="outline"
                class="w-full justify-start h-auto p-4"
                as-child
              >
                <router-link to="/transactions/new">
                  <div class="flex items-center gap-3 w-full">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 text-green-600">
                      <Plus class="h-5 w-5" />
                    </div>
                    <div class="flex-1 text-left">
                      <div class="text-sm font-medium text-foreground">Nova Transação</div>
                      <div class="text-xs text-muted-foreground">Registrar receita ou despesa</div>
                    </div>
                    <ChevronRight class="h-4 w-4 text-muted-foreground" />
                  </div>
                </router-link>
              </Button>

              <Button
                variant="outline"
                class="w-full justify-start h-auto p-4"
                as-child
              >
                <router-link to="/accounts">
                  <div class="flex items-center gap-3 w-full">
                    <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100 text-purple-600">
                      <Wallet class="h-5 w-5" />
                    </div>
                    <div class="flex-1 text-left">
                      <div class="text-sm font-medium text-foreground">Ver Contas</div>
                      <div class="text-xs text-muted-foreground">Gerenciar contas</div>
                    </div>
                    <ChevronRight class="h-4 w-4 text-muted-foreground" />
                  </div>
                </router-link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Recent Transactions -->
      <Card v-if="transactionsStore.transactions.length > 0">
        <CardHeader>
          <div class="flex items-center justify-between">
            <CardTitle>Transações Recentes</CardTitle>
            <Button variant="link" as-child>
              <router-link to="/transactions">
                Ver todas
              </router-link>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div class="space-y-2">
            <div
              v-for="transaction in recentTransactions"
              :key="transaction.transaction_id"
              class="flex items-center justify-between rounded-lg border border-border p-4 transition-colors hover:bg-muted/50"
            >
              <div class="flex items-center gap-3">
                <div
                  class="flex h-10 w-10 items-center justify-center rounded-lg"
                  :class="
                    transaction.type === 'INCOME'
                      ? 'bg-green-100 text-green-600'
                      : 'bg-red-100 text-red-600'
                  "
                >
                  <ArrowDown v-if="transaction.type === 'INCOME'" class="h-5 w-5" />
                  <ArrowUp v-else class="h-5 w-5" />
                </div>
                <div>
                  <div class="text-sm font-medium text-foreground">{{ transaction.description }}</div>
                  <div class="text-xs text-muted-foreground">{{ formatDate(transaction.date) }}</div>
                </div>
              </div>
              <div
                class="text-sm font-semibold"
                :class="transaction.type === 'INCOME' ? 'text-green-600' : 'text-red-600'"
              >
                {{ formatCurrency(parseFloat(transaction.amount)) }}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Wallet, User, Briefcase, List, ArrowDown, ArrowUp, TrendingUp, Plus, ChevronRight } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()

const recentTransactions = computed(() => {
  return transactionsStore.transactions
    .slice()
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
    .slice(0, 5)
})

onMounted(async () => {
  // Validar token antes de carregar dados
  // O validateToken já foi chamado no router guard, então só verificamos se está autenticado
  if (!authStore.isAuthenticated) {
    authStore.logout()
    router.push({ name: 'login', query: { redirect: '/' } })
    return
  }

  // Carregar dados apenas se autenticado e validado
  // Se validateToken já carregou accounts (através da store), não precisa chamar novamente
  try {
    // Aguardar um pouco para garantir que validateToken terminou (se ainda estiver rodando)
    if (authStore.isValidating) {
      await new Promise(resolve => setTimeout(resolve, 100))
    }
    
    // Só carregar se não tiver dados E não estiver carregando
    if (accountsStore.accounts.length === 0 && !accountsStore.isLoading) {
      await accountsStore.listAccounts()
    }
    if (transactionsStore.transactions.length === 0 && !transactionsStore.isLoading) {
      await transactionsStore.listTransactions()
    }
  } catch (error: any) {
    // Se houver erro 401 ou 403, token pode ter sido invalidado
    if (error.response?.status === 401 || error.response?.status === 403) {
      authStore.logout()
      router.push({ name: 'login', query: { redirect: '/' } })
      return
    }
    // Outros erros podem ser ignorados ou tratados conforme necessário
    console.error('Erro ao carregar dados do dashboard:', error)
  }
})

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value)
}

function formatDate(dateString: string): string {
  if (!dateString) return 'Data inválida'
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return 'Data inválida'
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}
</script>
