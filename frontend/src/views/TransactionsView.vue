<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Transações' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Transações</h1>
          <p class="text-gray-600">
            Gerencie suas receitas e despesas
          </p>
        </div>
        <router-link
          to="/transactions/new"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
        >
          <i class="pi pi-plus"></i>
          Nova Transação
        </router-link>
      </div>

      <!-- Loading State -->
      <div v-if="transactionsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <i class="pi pi-spinner pi-spin text-4xl text-blue-600 mb-4"></i>
          <p class="text-gray-600">Carregando transações...</p>
        </div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="transactionsStore.error"
        class="rounded-md bg-red-50 border border-red-200 p-4 mb-6"
      >
        <div class="flex items-center gap-2">
          <i class="pi pi-exclamation-circle text-red-600"></i>
          <p class="text-red-600">{{ transactionsStore.error }}</p>
        </div>
        <button
          @click="handleRetry"
          class="mt-3 text-sm text-red-600 hover:text-red-700 underline"
        >
          Tentar novamente
        </button>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="transactionsStore.transactions.length === 0"
        class="rounded-lg border border-gray-200 bg-white p-12 text-center"
      >
        <i class="pi pi-credit-card text-6xl text-gray-400 mb-4"></i>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">
          Nenhuma transação encontrada
        </h3>
        <p class="text-gray-600 mb-6">
          Comece registrando sua primeira transação financeira
        </p>
        <router-link
          to="/transactions/new"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
        >
          <i class="pi pi-plus"></i>
          Criar Primeira Transação
        </router-link>
      </div>

      <!-- Transactions List -->
      <div v-else class="space-y-4">
        <!-- Stats -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-50 to-gray-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-gray-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-500 text-white">
                  <i class="pi pi-list"></i>
                </div>
              </div>
              <div class="text-sm font-medium text-gray-700">Total de Transações</div>
              <div class="text-3xl font-bold text-gray-900">
                {{ transactionsStore.totalTransactions }}
              </div>
            </div>
          </div>
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-green-50 to-green-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-green-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500 text-white">
                  <i class="pi pi-arrow-down"></i>
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
          </div>
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-red-50 to-red-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-red-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-red-500 text-white">
                  <i class="pi pi-arrow-up"></i>
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
          </div>
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 p-6 shadow-sm transition-all hover:shadow-md border-2"
            :class="transactionsStore.balance >= 0 ? 'border-green-200' : 'border-red-200'">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full opacity-20"
              :class="transactionsStore.balance >= 0 ? 'bg-green-200' : 'bg-red-200'"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg text-white"
                  :class="transactionsStore.balance >= 0 ? 'bg-green-500' : 'bg-red-500'">
                  <i class="pi pi-chart-line"></i>
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
          </div>
        </div>

        <!-- Filters -->
        <div class="flex flex-wrap gap-3 items-center">
          <label class="text-sm font-medium text-gray-700">Filtros:</label>
          <select
            v-model="selectedType"
            @change="handleFilterChange"
            class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Todos os tipos</option>
            <option value="INCOME">Receitas</option>
            <option value="EXPENSE">Despesas</option>
          </select>
          <select
            v-model="selectedAccountId"
            @change="handleFilterChange"
            class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
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
          <button
            v-if="selectedType || selectedAccountId"
            @click="clearFilters"
            class="text-sm text-gray-600 hover:text-gray-900 underline"
          >
            Limpar filtros
          </button>
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

const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()

const selectedType = ref<'INCOME' | 'EXPENSE' | ''>('')
const selectedAccountId = ref<string>('')

onMounted(async () => {
  // Carregar contas se ainda não foram carregadas
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
  // Carregar transações se ainda não foram carregadas
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
