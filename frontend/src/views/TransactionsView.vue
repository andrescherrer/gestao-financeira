<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Transações' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Transações</h1>
          <p class="text-muted-foreground">
            Gerencie suas receitas e despesas
          </p>
        </div>
        <Button as-child>
          <router-link to="/transactions/new">
            <Plus class="h-4 w-4 mr-2" />
            Nova Transação
          </router-link>
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="transactionsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando transações...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="transactionsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ transactionsStore.error }}</p>
          </div>
          <Button
            variant="link"
            @click="handleRetry"
            class="text-destructive"
          >
            Tentar novamente
          </Button>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <Card
        v-else-if="transactionsStore.transactions.length === 0"
      >
        <CardContent class="p-12 text-center">
          <CreditCard class="mx-auto h-16 w-16 text-muted-foreground mb-4" />
          <h3 class="text-xl font-semibold text-foreground mb-2">
            Nenhuma transação encontrada
          </h3>
          <p class="text-muted-foreground mb-6">
            Comece registrando sua primeira transação financeira
          </p>
          <Button as-child>
            <router-link to="/transactions/new">
              <Plus class="h-4 w-4 mr-2" />
              Criar Primeira Transação
            </router-link>
          </Button>
        </CardContent>
      </Card>

      <!-- Transactions List -->
      <div v-else class="space-y-4">
        <!-- Stats -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-gray-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-500 text-white">
                    <List class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-gray-700">Total de Transações</div>
                <div class="text-3xl font-bold text-gray-900">
                  {{ transactionsStore.totalTransactions }}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-green-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500 text-white">
                    <ArrowDown class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-green-700">Receitas</div>
                <div class="text-3xl font-bold text-green-600">
                  {{ formatCurrency(transactionsStore.totalIncome) }}
                </div>
                <div class="text-xs text-green-600 mt-1 font-medium">
                  {{ transactionsStore.incomeTransactions.length }} transações
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-red-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-red-500 text-white">
                    <ArrowUp class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-red-700">Despesas</div>
                <div class="text-3xl font-bold text-red-600">
                  {{ formatCurrency(transactionsStore.totalExpense) }}
                </div>
                <div class="text-xs text-red-600 mt-1 font-medium">
                  {{ transactionsStore.expenseTransactions.length }} transações
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden border-2"
            :class="transactionsStore.balance >= 0 ? 'border-green-200' : 'border-red-200'">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full opacity-20"
                :class="transactionsStore.balance >= 0 ? 'bg-green-200' : 'bg-red-200'"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg text-white"
                    :class="transactionsStore.balance >= 0 ? 'bg-green-500' : 'bg-red-500'">
                    <TrendingUp class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium"
                  :class="transactionsStore.balance >= 0 ? 'text-green-700' : 'text-red-700'">Saldo</div>
                <div
                  class="text-3xl font-bold"
                  :class="
                    transactionsStore.balance >= 0 ? 'text-green-600' : 'text-red-600'
                  "
                >
                  {{ formatCurrency(transactionsStore.balance) }}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Filters -->
        <div class="flex flex-wrap gap-3 items-center">
          <label class="text-sm font-medium text-foreground">Filtros:</label>
          <select
            v-model="selectedType"
            @change="handleFilterChange"
            class="flex h-9 w-full max-w-xs rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          >
            <option value="">Todos os tipos</option>
            <option value="INCOME">Receitas</option>
            <option value="EXPENSE">Despesas</option>
          </select>
          <select
            v-model="selectedAccountId"
            @change="handleFilterChange"
            class="flex h-9 w-full max-w-xs rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          >
            <option value="">Todas as contas</option>
            <option
              v-for="account in accountsStore.accounts"
              :key="account.account_id"
              :value="account.account_id"
            >
              {{ account.name }}
            </option>
          </select>
          <Button
            v-if="selectedType || selectedAccountId"
            variant="ghost"
            size="sm"
            @click="clearFilters"
          >
            Limpar filtros
          </Button>
        </div>

        <!-- Transactions Table -->
        <TransactionTable :transactions="transactionsStore.transactions" />
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import TransactionTable from '@/components/TransactionTable.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { List, ArrowDown, ArrowUp, TrendingUp, CreditCard, Plus, Loader2, AlertCircle } from 'lucide-vue-next'

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()

const selectedType = ref<'INCOME' | 'EXPENSE' | ''>('')
const selectedAccountId = ref<string>('')

onMounted(async () => {
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
  if (transactionsStore.transactions.length === 0) {
    await transactionsStore.listTransactions()
  }
})

async function handleFilterChange() {
  await transactionsStore.listTransactions(
    selectedAccountId.value || undefined,
    selectedType.value || undefined
  )
}

function clearFilters() {
  selectedType.value = ''
  selectedAccountId.value = ''
  transactionsStore.listTransactions()
}

function handleRetry() {
  transactionsStore.clearError()
  transactionsStore.listTransactions(
    selectedAccountId.value || undefined,
    selectedType.value || undefined
  )
}

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value)
}
</script>
