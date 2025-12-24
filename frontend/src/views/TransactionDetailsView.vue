<template>
  <Layout>
    <div>
      <!-- Header com botão voltar -->
      <div class="mb-6">
        <button
          @click="goBack"
          class="mb-4 flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
        >
          <i class="pi pi-arrow-left"></i>
          <span>Voltar para transações</span>
        </button>
        <h1 class="text-4xl font-bold mb-2">Detalhes da Transação</h1>
      </div>

      <!-- Loading State -->
      <div v-if="transactionsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <i class="pi pi-spinner pi-spin text-4xl text-blue-600 mb-4"></i>
          <p class="text-gray-600">Carregando detalhes da transação...</p>
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
        <div class="mt-4 flex gap-3">
          <button
            @click="handleRetry"
            class="rounded-md bg-red-600 px-4 py-2 text-white hover:bg-red-700 transition-colors"
          >
            Tentar novamente
          </button>
          <button
            @click="goBack"
            class="rounded-md border border-gray-300 px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Voltar
          </button>
        </div>
      </div>

      <!-- Transaction Details -->
      <div v-else-if="transactionsStore.currentTransaction" class="space-y-6">
        <!-- Transaction Card -->
        <div class="rounded-lg border border-gray-200 bg-white p-6">
          <div class="mb-6 flex items-start justify-between">
            <div>
              <h2 class="text-2xl font-bold text-gray-900 mb-2">
                {{ transactionsStore.currentTransaction.description }}
              </h2>
              <p class="text-gray-600">
                {{ getTransactionTypeLabel(transactionsStore.currentTransaction.type) }}
              </p>
            </div>
            <span
              class="rounded-full px-3 py-1 text-sm font-medium"
              :class="
                transactionsStore.currentTransaction.type === 'INCOME'
                  ? 'bg-green-100 text-green-700'
                  : 'bg-red-100 text-red-700'
              "
            >
              {{ transactionsStore.currentTransaction.type === 'INCOME' ? 'Receita' : 'Despesa' }}
            </span>
          </div>

          <!-- Amount -->
          <div class="mb-6 rounded-lg bg-gray-50 p-6">
            <div class="text-sm text-gray-600 mb-2">Valor</div>
            <div
              class="text-4xl font-bold"
              :class="
                transactionsStore.currentTransaction.type === 'INCOME'
                  ? 'text-green-600'
                  : 'text-red-600'
              "
            >
              {{ formatCurrency(
                parseFloat(transactionsStore.currentTransaction.amount),
                transactionsStore.currentTransaction.currency
              ) }}
            </div>
          </div>

          <!-- Transaction Information Grid -->
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <div class="text-sm text-gray-600 mb-1">Tipo</div>
              <div class="text-lg font-medium text-gray-900">
                {{ getTransactionTypeLabel(transactionsStore.currentTransaction.type) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Conta</div>
              <div class="text-lg font-medium text-gray-900">
                {{ getAccountName(transactionsStore.currentTransaction.account_id) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Data</div>
              <div class="text-lg font-medium text-gray-900">
                {{ formatDate(transactionsStore.currentTransaction.date) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Moeda</div>
              <div class="text-lg font-medium text-gray-900">
                {{ transactionsStore.currentTransaction.currency }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Data de Criação</div>
              <div class="text-lg font-medium text-gray-900">
                {{ formatDateTime(transactionsStore.currentTransaction.created_at) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Última Atualização</div>
              <div class="text-lg font-medium text-gray-900">
                {{ formatDateTime(transactionsStore.currentTransaction.updated_at) }}
              </div>
            </div>
          </div>

          <!-- Description -->
          <div class="mt-6 pt-6 border-t border-gray-200">
            <div class="text-sm text-gray-600 mb-2">Descrição</div>
            <div class="text-lg text-gray-900">
              {{ transactionsStore.currentTransaction.description }}
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <button
            @click="goBack"
            class="inline-flex items-center gap-2 rounded-md border border-gray-300 px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors"
          >
            <i class="pi pi-arrow-left"></i>
            Voltar
          </button>
        </div>
      </div>

      <!-- Not Found State -->
      <div v-else class="rounded-lg border border-gray-200 bg-white p-12 text-center">
        <i class="pi pi-exclamation-circle text-6xl text-gray-400 mb-4"></i>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">
          Transação não encontrada
        </h3>
        <p class="text-gray-600 mb-6">
          A transação que você está procurando não existe ou foi removida.
        </p>
        <button
          @click="goBack"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 transition-colors"
        >
          <i class="pi pi-arrow-left"></i>
          Voltar para transações
        </button>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import type { Transaction } from '@/api/types'

const route = useRoute()
const router = useRouter()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()

const transactionId = route.params.id as string

onMounted(async () => {
  // Carregar contas se ainda não foram carregadas
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
  await loadTransaction()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadTransaction()
    }
  }
)

async function loadTransaction() {
  if (!transactionId) return

  // Verifica se a transação já está na lista
  const existingTransaction = transactionsStore.transactions.find(
    (tx) => tx.transaction_id === transactionId
  )

  if (existingTransaction && !transactionsStore.currentTransaction) {
    // Se já está na lista, usa ela
    transactionsStore.currentTransaction = existingTransaction
  } else {
    // Caso contrário, busca do servidor
    try {
      await transactionsStore.getTransaction(transactionId)
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

function goBack() {
  router.push('/transactions')
}

function handleRetry() {
  transactionsStore.clearError()
  loadTransaction()
}

function getTransactionTypeLabel(type: Transaction['type']): string {
  return type === 'INCOME' ? 'Receita' : 'Despesa'
}

function getAccountName(accountId: string): string {
  const account = accountsStore.accounts.find((acc) => acc.account_id === accountId)
  return account?.name || 'Conta não encontrada'
}

function formatCurrency(value: number, currency: Transaction['currency']): string {
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

function formatDateTime(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
</script>
