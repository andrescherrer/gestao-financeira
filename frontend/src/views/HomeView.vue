<template>
  <Layout>
    <div class="space-y-6">
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Dashboard' }]" />

      <!-- Header -->
      <div>
        <h1 class="text-4xl font-bold text-gray-900 mb-2">Dashboard</h1>
        <p class="text-gray-600">
          Visão geral da sua situação financeira
        </p>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <!-- Total de Contas -->
        <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 p-6 shadow-sm transition-all hover:shadow-md">
          <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-blue-200 opacity-20"></div>
          <div class="relative">
            <div class="mb-2 flex items-center gap-3">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-500 text-white">
                <i class="pi pi-wallet text-xl"></i>
              </div>
            </div>
            <div class="text-sm font-medium text-blue-700">Total de Contas</div>
            <div class="text-3xl font-bold text-blue-900">
              {{ accountsStore.totalAccounts }}
            </div>
          </div>
        </div>

        <!-- Contas Pessoais -->
        <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-purple-50 to-purple-100 p-6 shadow-sm transition-all hover:shadow-md">
          <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-purple-200 opacity-20"></div>
          <div class="relative">
            <div class="mb-2 flex items-center gap-3">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-purple-500 text-white">
                <i class="pi pi-user text-xl"></i>
              </div>
            </div>
            <div class="text-sm font-medium text-purple-700">Contas Pessoais</div>
            <div class="text-3xl font-bold text-purple-900">
              {{ accountsStore.personalAccounts.length }}
            </div>
          </div>
        </div>

        <!-- Contas de Negócio -->
        <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-indigo-50 to-indigo-100 p-6 shadow-sm transition-all hover:shadow-md">
          <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-indigo-200 opacity-20"></div>
          <div class="relative">
            <div class="mb-2 flex items-center gap-3">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-500 text-white">
                <i class="pi pi-briefcase text-xl"></i>
              </div>
            </div>
            <div class="text-sm font-medium text-indigo-700">Contas de Negócio</div>
            <div class="text-3xl font-bold text-indigo-900">
              {{ accountsStore.businessAccounts.length }}
            </div>
          </div>
        </div>

        <!-- Total de Transações -->
        <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-green-50 to-green-100 p-6 shadow-sm transition-all hover:shadow-md">
          <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-green-200 opacity-20"></div>
          <div class="relative">
            <div class="mb-2 flex items-center gap-3">
              <div class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-500 text-white">
                <i class="pi pi-list text-xl"></i>
              </div>
            </div>
            <div class="text-sm font-medium text-green-700">Transações</div>
            <div class="text-3xl font-bold text-green-900">
              {{ transactionsStore.totalTransactions }}
            </div>
          </div>
        </div>
      </div>

      <!-- Financial Overview -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Receitas vs Despesas -->
        <div class="lg:col-span-2 rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h2 class="mb-4 text-lg font-semibold text-gray-900">Resumo Financeiro</h2>
          <div class="space-y-4">
            <div class="flex items-center justify-between rounded-lg bg-green-50 p-4">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-500 text-white">
                  <i class="pi pi-arrow-down text-sm"></i>
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-600">Total de Receitas</div>
                  <div class="text-xl font-bold text-green-600">
                    {{ formatCurrency(transactionsStore.totalIncome) }}
                  </div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-xs text-gray-500">{{ transactionsStore.incomeTransactions.length }} transações</div>
              </div>
            </div>

            <div class="flex items-center justify-between rounded-lg bg-red-50 p-4">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-red-500 text-white">
                  <i class="pi pi-arrow-up text-sm"></i>
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-600">Total de Despesas</div>
                  <div class="text-xl font-bold text-red-600">
                    {{ formatCurrency(transactionsStore.totalExpense) }}
                  </div>
                </div>
              </div>
              <div class="text-right">
                <div class="text-xs text-gray-500">{{ transactionsStore.expenseTransactions.length }} transações</div>
              </div>
            </div>

            <div class="flex items-center justify-between rounded-lg bg-blue-50 p-4 border-2 border-blue-200">
              <div class="flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500 text-white">
                  <i class="pi pi-chart-line text-sm"></i>
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-600">Saldo</div>
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
        </div>

        <!-- Quick Actions -->
        <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h2 class="mb-4 text-lg font-semibold text-gray-900">Ações Rápidas</h2>
          <div class="space-y-3">
            <router-link
              to="/accounts/new"
              class="flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-4 transition-all hover:border-blue-300 hover:bg-blue-50 hover:shadow-sm"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 text-blue-600">
                <i class="pi pi-plus text-sm"></i>
              </div>
              <div class="flex-1">
                <div class="text-sm font-medium text-gray-900">Nova Conta</div>
                <div class="text-xs text-gray-500">Criar conta financeira</div>
              </div>
              <i class="pi pi-chevron-right text-gray-400"></i>
            </router-link>

            <router-link
              to="/transactions/new"
              class="flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-4 transition-all hover:border-green-300 hover:bg-green-50 hover:shadow-sm"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-green-100 text-green-600">
                <i class="pi pi-plus text-sm"></i>
              </div>
              <div class="flex-1">
                <div class="text-sm font-medium text-gray-900">Nova Transação</div>
                <div class="text-xs text-gray-500">Registrar receita ou despesa</div>
              </div>
              <i class="pi pi-chevron-right text-gray-400"></i>
            </router-link>

            <router-link
              to="/accounts"
              class="flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-4 transition-all hover:border-purple-300 hover:bg-purple-50 hover:shadow-sm"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-100 text-purple-600">
                <i class="pi pi-wallet text-sm"></i>
              </div>
              <div class="flex-1">
                <div class="text-sm font-medium text-gray-900">Ver Contas</div>
                <div class="text-xs text-gray-500">Gerenciar contas</div>
              </div>
              <i class="pi pi-chevron-right text-gray-400"></i>
            </router-link>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div v-if="transactionsStore.transactions.length > 0" class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900">Transações Recentes</h2>
          <router-link
            to="/transactions"
            class="text-sm font-medium text-blue-600 hover:text-blue-700"
          >
            Ver todas
          </router-link>
        </div>
        <div class="space-y-2">
          <div
            v-for="transaction in recentTransactions"
            :key="transaction.transaction_id"
            class="flex items-center justify-between rounded-lg border border-gray-100 p-4 transition-colors hover:bg-gray-50"
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
                <i
                  class="pi text-sm"
                  :class="transaction.type === 'INCOME' ? 'pi-arrow-down' : 'pi-arrow-up'"
                ></i>
              </div>
              <div>
                <div class="text-sm font-medium text-gray-900">{{ transaction.description }}</div>
                <div class="text-xs text-gray-500">{{ formatDate(transaction.date) }}</div>
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
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'

const accountsStore = useAccountsStore()
const transactionsStore = useTransactionsStore()

const recentTransactions = computed(() => {
  return transactionsStore.transactions
    .slice()
    .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
    .slice(0, 5)
})

onMounted(async () => {
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
  if (transactionsStore.transactions.length === 0) {
    await transactionsStore.listTransactions()
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
