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
          <span>Voltar para contas</span>
        </button>
        <h1 class="text-4xl font-bold mb-2">Detalhes da Conta</h1>
      </div>

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
        class="rounded-md bg-red-50 border border-red-200 p-4 mb-6"
      >
        <div class="flex items-center gap-2">
          <i class="pi pi-exclamation-circle text-red-600"></i>
          <p class="text-red-600">{{ accountsStore.error }}</p>
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

      <!-- Account Details -->
      <div v-else-if="accountsStore.currentAccount" class="space-y-6">
        <!-- Account Card -->
        <div class="rounded-lg border border-gray-200 bg-white p-6">
          <div class="mb-6 flex items-start justify-between">
            <div>
              <h2 class="text-2xl font-bold text-gray-900 mb-2">
                {{ accountsStore.currentAccount.name }}
              </h2>
              <p class="text-gray-600">
                {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
              </p>
            </div>
            <span
              class="rounded-full px-3 py-1 text-sm font-medium"
              :class="
                accountsStore.currentAccount.is_active
                  ? 'bg-green-100 text-green-700'
                  : 'bg-gray-100 text-gray-700'
              "
            >
              {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
            </span>
          </div>

          <!-- Balance -->
          <div class="mb-6 rounded-lg bg-gray-50 p-6">
            <div class="text-sm text-gray-600 mb-2">Saldo Atual</div>
            <div
              class="text-4xl font-bold"
              :class="getBalanceColor(accountsStore.currentAccount.balance)"
            >
              {{ formatCurrency(accountsStore.currentAccount.balance, accountsStore.currentAccount.currency) }}
            </div>
          </div>

          <!-- Account Information Grid -->
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <div class="text-sm text-gray-600 mb-1">Tipo de Conta</div>
              <div class="text-lg font-medium text-gray-900">
                {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Contexto</div>
              <div class="text-lg font-medium text-gray-900">
                {{ getContextLabel(accountsStore.currentAccount.context) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Moeda</div>
              <div class="text-lg font-medium text-gray-900">
                {{ accountsStore.currentAccount.currency }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Status</div>
              <div class="text-lg font-medium text-gray-900">
                {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Data de Criação</div>
              <div class="text-lg font-medium text-gray-900">
                {{ formatDate(accountsStore.currentAccount.created_at) }}
              </div>
            </div>

            <div>
              <div class="text-sm text-gray-600 mb-1">Última Atualização</div>
              <div class="text-lg font-medium text-gray-900">
                {{ formatDate(accountsStore.currentAccount.updated_at) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex gap-3">
          <router-link
            :to="`/accounts/${accountsStore.currentAccount.account_id}/edit`"
            class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 transition-colors"
          >
            <i class="pi pi-pencil"></i>
            Editar Conta
          </router-link>
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
          Conta não encontrada
        </h3>
        <p class="text-gray-600 mb-6">
          A conta que você está procurando não existe ou foi removida.
        </p>
        <button
          @click="goBack"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 transition-colors"
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

