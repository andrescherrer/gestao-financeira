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
        <TransactionForm
          ref="formRef"
          :isLoading="transactionsStore.isLoading"
          submitLabel="Criar Transação"
          :initial-values="initialValues"
          @submit="handleSubmit"
        >
          <template #actions="{ isLoading: formLoading }">
            <div class="flex gap-3 pt-4">
              <button
                type="submit"
                :disabled="formLoading"
                class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                <i v-if="formLoading" class="pi pi-spinner pi-spin"></i>
                <i v-else class="pi pi-check"></i>
                {{ formLoading ? 'Criando...' : 'Criar Transação' }}
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
          </template>
        </TransactionForm>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import TransactionForm from '@/components/TransactionForm.vue'
import type { CreateTransactionFormData } from '@/validations/transaction'
import type { CreateTransactionRequest } from '@/api/types'

const router = useRouter()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const formRef = ref<InstanceType<typeof TransactionForm> | null>(null)

// Inicializar data com hoje
const today = new Date().toISOString().split('T')[0]

const initialValues = computed(() => ({
  date: today,
  currency: 'BRL' as const,
}))

onMounted(async () => {
  // Carregar contas se ainda não foram carregadas
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
})

async function handleSubmit(values: CreateTransactionFormData) {
  if (formRef.value) {
    formRef.value.clearError()
  }

  try {
    // Preparar dados para API
    // Backend espera amount como number (float64) e currency obrigatório
    let amount = 0.0
    if (values.amount) {
      const parsed = parseFloat(values.amount)
      if (isNaN(parsed) || !isFinite(parsed) || parsed <= 0) {
        throw new Error('Valor deve ser um número maior que zero')
      }
      amount = parsed
    }

    // Validar descrição (backend exige mínimo 3 caracteres)
    const trimmedDescription = values.description.trim()
    if (trimmedDescription.length < 3) {
      throw new Error('Descrição deve ter no mínimo 3 caracteres')
    }

    const transactionData: CreateTransactionRequest = {
      account_id: values.account_id,
      type: values.type,
      amount: amount, // Sempre enviar como number
      currency: values.currency || 'BRL', // Sempre enviar currency (obrigatório)
      description: trimmedDescription,
      date: values.date, // ISO 8601 format: YYYY-MM-DD
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

    const errorMessage =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao criar transação. Tente novamente.'
    
    if (formRef.value) {
      formRef.value.setError(errorMessage)
    }
  }
}

function goBack() {
  router.push('/transactions')
}
</script>
