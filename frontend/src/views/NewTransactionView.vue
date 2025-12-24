<template>
  <Layout>
    <div>
      <!-- Header -->
      <div class="mb-6">
        <button
          @click="goBack"
          class="mb-4 flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
        >
          <i class="pi pi-arrow-left"></i>
          <span>Voltar para transações</span>
        </button>
        <h1 class="text-4xl font-bold mb-2">Nova Transação</h1>
        <p class="text-gray-600">
          Preencha os dados para criar uma nova transação financeira
        </p>
      </div>

      <!-- Form -->
      <div class="rounded-lg border border-gray-200 bg-white p-6">
        <form @submit.prevent="handleSubmit">
          <!-- Error Message -->
          <div
            v-if="error"
            class="mb-6 rounded-md bg-red-50 border border-red-200 p-4"
          >
            <div class="flex items-center gap-2">
              <i class="pi pi-exclamation-circle text-red-600"></i>
              <p class="text-red-600">{{ error }}</p>
            </div>
          </div>

          <!-- Account -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Conta *
            </label>
            <select
              v-model="formData.account_id"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              :disabled="transactionsStore.isLoading"
            >
              <option value="">Selecione uma conta</option>
              <option
                v-for="account in accountsStore.accounts"
                :key="account.account_id"
                :value="account.account_id"
              >
                {{ account.name }} ({{ account.currency }})
              </option>
            </select>
          </div>

          <!-- Type -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tipo *
            </label>
            <select
              v-model="formData.type"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              :disabled="transactionsStore.isLoading"
            >
              <option value="">Selecione o tipo</option>
              <option value="INCOME">Receita</option>
              <option value="EXPENSE">Despesa</option>
            </select>
          </div>

          <!-- Amount -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Valor *
            </label>
            <input
              v-model="formData.amount"
              type="number"
              step="0.01"
              min="0.01"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="0.00"
              :disabled="transactionsStore.isLoading"
            />
          </div>

          <!-- Currency -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Moeda *
            </label>
            <select
              v-model="formData.currency"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              :disabled="transactionsStore.isLoading"
            >
              <option value="BRL">BRL - Real Brasileiro</option>
              <option value="USD">USD - Dólar Americano</option>
              <option value="EUR">EUR - Euro</option>
            </select>
          </div>

          <!-- Description -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Descrição *
            </label>
            <textarea
              v-model="formData.description"
              rows="3"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Descreva a transação..."
              :disabled="transactionsStore.isLoading"
            ></textarea>
            <p class="mt-1 text-xs text-gray-500">
              Mínimo 3 caracteres, máximo 500 caracteres
            </p>
          </div>

          <!-- Date -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Data *
            </label>
            <input
              v-model="formData.date"
              type="date"
              class="w-full rounded-md border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              :disabled="transactionsStore.isLoading"
            />
            <p class="mt-1 text-xs text-gray-500">
              Data da transação (formato: YYYY-MM-DD)
            </p>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-4">
            <button
              type="submit"
              :disabled="transactionsStore.isLoading"
              class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
            >
              <i v-if="transactionsStore.isLoading" class="pi pi-spinner pi-spin"></i>
              <i v-else class="pi pi-check"></i>
              {{ transactionsStore.isLoading ? 'Criando...' : 'Criar Transação' }}
            </button>
            <button
              type="button"
              @click="goBack"
              class="inline-flex items-center gap-2 rounded-md border border-gray-300 px-4 py-2 text-gray-700 hover:bg-gray-50 transition-colors"
            >
              <i class="pi pi-times"></i>
              Cancelar
            </button>
          </div>
        </form>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import type { CreateTransactionRequest } from '@/api/types'

const router = useRouter()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()

const error = ref<string | null>(null)

// Inicializar data com hoje
const today = new Date().toISOString().split('T')[0]

const formData = ref({
  account_id: '',
  type: '' as 'INCOME' | 'EXPENSE' | '',
  amount: '',
  currency: 'BRL' as 'BRL' | 'USD' | 'EUR',
  description: '',
  date: today,
})

onMounted(async () => {
  // Carregar contas se ainda não foram carregadas
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
})

async function handleSubmit() {
  error.value = null

  try {
    // Validar dados básicos
    if (!formData.value.account_id) {
      throw new Error('Selecione uma conta')
    }
    if (!formData.value.type) {
      throw new Error('Selecione o tipo de transação')
    }
    if (!formData.value.amount || parseFloat(formData.value.amount) <= 0) {
      throw new Error('Valor deve ser maior que zero')
    }
    if (!formData.value.description || formData.value.description.trim().length < 3) {
      throw new Error('Descrição deve ter no mínimo 3 caracteres')
    }
    if (!formData.value.date) {
      throw new Error('Data é obrigatória')
    }

    // Preparar dados para API
    // Backend espera amount como number (float64) e currency obrigatório
    const transactionData: CreateTransactionRequest = {
      account_id: formData.value.account_id,
      type: formData.value.type as 'INCOME' | 'EXPENSE',
      amount: parseFloat(formData.value.amount), // Converter para number
      currency: formData.value.currency, // Sempre presente (default 'BRL')
      description: formData.value.description.trim(),
      date: formData.value.date, // ISO 8601 format: YYYY-MM-DD
    }

    // Log em desenvolvimento para debug
    if (import.meta.env.DEV) {
      console.log('[NewTransactionView] Dados sendo enviados:', transactionData)
    }

    const transaction = await transactionsStore.createTransaction(transactionData)
    
    // Redirecionar para detalhes da transação criada
    router.push(`/transactions/${transaction.transaction_id}`)
  } catch (err: any) {
    // Log detalhado do erro em desenvolvimento
    if (import.meta.env.DEV) {
      console.error('[NewTransactionView] Erro ao criar transação:', {
        message: err.message,
        response: err.response?.data,
        status: err.response?.status,
        statusText: err.response?.statusText,
      })
    }

    error.value =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao criar transação. Tente novamente.'
  }
}

function goBack() {
  router.push('/transactions')
}
</script>
