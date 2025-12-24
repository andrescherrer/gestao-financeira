<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Contas', to: '/accounts' },
          { label: accountsStore.currentAccount?.name || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <i class="pi pi-spinner pi-spin text-4xl text-blue-600 mb-4"></i>
          <p class="text-gray-600">Carregando detalhes da conta...</p>
        </div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="accountsStore.error"
        class="mb-6 rounded-lg border border-red-200 bg-red-50 p-4"
      >
        <div class="flex items-center gap-2">
          <i class="pi pi-exclamation-circle text-red-600"></i>
          <p class="text-red-600">{{ accountsStore.error }}</p>
        </div>
        <div class="mt-4 flex gap-3">
          <button
            @click="handleRetry"
            class="rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-red-700"
          >
            Tentar novamente
          </button>
          <button
            @click="goBack"
            class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50"
          >
            Voltar
          </button>
        </div>
      </div>

      <!-- Account Details -->
      <div v-else-if="accountsStore.currentAccount" class="space-y-6">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold text-gray-900 mb-2">
              {{ accountsStore.currentAccount.name }}
            </h1>
            <p class="text-gray-600">
              {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <span
              class="rounded-full px-3 py-1 text-sm font-semibold"
              :class="
                accountsStore.currentAccount.is_active
                  ? 'bg-green-100 text-green-700'
                  : 'bg-gray-100 text-gray-700'
              "
            >
              {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
            </span>
            <router-link
              :to="`/accounts/${accountsStore.currentAccount.account_id}/edit`"
              class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            >
              <i class="pi pi-pencil"></i>
              Editar
            </router-link>
          </div>
        </div>

        <!-- Account Card -->
        <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <!-- Balance -->
          <div class="mb-6 rounded-lg bg-gradient-to-br from-gray-50 to-gray-100 p-6">
            <div class="text-sm font-medium text-gray-600 mb-2">Saldo Atual</div>
            <div
              class="text-4xl font-bold"
              :class="getBalanceColor(accountsStore.currentAccount.balance)"
            >
              {{ formatCurrency(accountsStore.currentAccount.balance, accountsStore.currentAccount.currency) }}
            </div>
          </div>

          <!-- Account Information Grid -->
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Tipo de Conta
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
              </div>
            </div>

            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Contexto
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ getContextLabel(accountsStore.currentAccount.context) }}
              </div>
            </div>

            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Moeda
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ accountsStore.currentAccount.currency }}
              </div>
            </div>

            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Status
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
              </div>
            </div>

            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Data de Criação
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ formatDate(accountsStore.currentAccount.created_at) }}
              </div>
            </div>

            <div class="rounded-lg border border-gray-100 bg-gray-50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-gray-500 mb-1">
                Última Atualização
              </div>
              <div class="text-lg font-semibold text-gray-900">
                {{ formatDate(accountsStore.currentAccount.updated_at) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Not Found State -->
      <div v-else class="rounded-xl border border-gray-200 bg-white p-12 text-center shadow-sm">
        <i class="pi pi-exclamation-circle text-6xl text-gray-400 mb-4"></i>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">
          Conta não encontrada
        </h3>
        <p class="text-gray-600 mb-6">
          A conta que você está procurando não existe ou foi removida.
        </p>
        <button
          @click="goBack"
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700"
        >
          <i class="pi pi-arrow-left"></i>
          Voltar para contas
        </button>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import type { Account } from '@/api/types'

const route = useRoute()
const router = useRouter()
const accountsStore = useAccountsStore()

const accountId = route.params.id as string

onMounted(async () => {
  await loadAccount()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadAccount()
    }
  }
)

async function loadAccount() {
  if (!accountId) return

  // Verifica se a conta já está na lista
  const existingAccount = accountsStore.accounts.find(
    (acc) => acc.account_id === accountId
  )

  if (existingAccount && !accountsStore.currentAccount) {
    // Se já está na lista, usa ela
    accountsStore.currentAccount = existingAccount
  } else {
    // Caso contrário, busca do servidor
    try {
      await accountsStore.getAccount(accountId)
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

function goBack() {
  router.push('/accounts')
}

function handleRetry() {
  accountsStore.clearError()
  loadAccount()
}

function getAccountTypeLabel(type: Account['type']): string {
  const labels: Record<Account['type'], string> = {
    BANK: 'Banco',
    WALLET: 'Carteira',
    INVESTMENT: 'Investimento',
    CREDIT_CARD: 'Cartão de Crédito',
  }
  return labels[type] || type
}

function getContextLabel(context: Account['context']): string {
  return context === 'PERSONAL' ? 'Pessoal' : 'Negócio'
}

function formatCurrency(amount: string, currency: Account['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function getBalanceColor(balance: string): string {
  const value = parseFloat(balance)
  if (value > 0) return 'text-green-600'
  if (value < 0) return 'text-red-600'
  return 'text-gray-900'
}

function formatDate(dateString: string): string {
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
