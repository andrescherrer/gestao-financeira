<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Transações', to: '/transactions' },
          { label: 'Nova Transação' },
        ]"
      />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Nova Transação</h1>
        <p class="text-muted-foreground">
          Preencha os dados para criar uma nova transação financeira
        </p>
      </div>

      <!-- Form -->
      <Card>
        <CardContent class="p-6">
          <TransactionForm
            ref="formRef"
            :isLoading="transactionsStore.isLoading"
            submitLabel="Criar Transação"
            :initial-values="initialValues"
            @submit="handleSubmit"
          >
            <template #actions="{ isLoading: formLoading }">
              <div class="flex gap-3 pt-4">
                <Button
                  type="submit"
                  :disabled="formLoading"
                >
                  <Loader2 v-if="formLoading" class="h-4 w-4 animate-spin" />
                  <Check v-else class="h-4 w-4" />
                  {{ formLoading ? 'Criando...' : 'Criar Transação' }}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  @click="goBack"
                >
                  <X class="h-4 w-4" />
                  Cancelar
                </Button>
              </div>
            </template>
          </TransactionForm>
        </CardContent>
      </Card>
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
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Loader2, Check, X } from 'lucide-vue-next'
import type { CreateTransactionFormData } from '@/validations/transaction'
import type { CreateTransactionRequest } from '@/api/types'

const router = useRouter()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()
const formRef = ref<InstanceType<typeof TransactionForm> | null>(null)

const today = new Date().toISOString().split('T')[0]

const initialValues = computed(() => ({
  date: today,
  currency: 'BRL' as const,
}))

onMounted(async () => {
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
})

async function handleSubmit(values: CreateTransactionFormData) {
  if (formRef.value) {
    formRef.value.clearError()
  }

  try {
    let amount = 0.0
    if (values.amount) {
      const parsed = parseFloat(values.amount)
      if (isNaN(parsed) || !isFinite(parsed) || parsed <= 0) {
        throw new Error('Valor deve ser um número maior que zero')
      }
      amount = parsed
    }

    const trimmedDescription = values.description.trim()
    if (trimmedDescription.length < 3) {
      throw new Error('Descrição deve ter no mínimo 3 caracteres')
    }

    const transactionData: CreateTransactionRequest = {
      account_id: values.account_id,
      type: values.type,
      amount: amount,
      currency: values.currency || 'BRL',
      description: trimmedDescription,
      date: values.date,
    }

    if (import.meta.env.DEV) {
      console.log('[NewTransactionView] Dados sendo enviados:', transactionData)
    }

    const transaction = await transactionsStore.createTransaction(transactionData)
    
    // Atualizar saldo da conta em tempo real
    // O backend já atualizou o saldo via event handler, mas precisamos atualizar no frontend
    await accountsStore.refreshAccount(transaction.account_id)
    
    // Mostrar toast de sucesso
    const { toast } = await import('@/components/ui/toast')
    toast.success('Transação criada com sucesso!', {
      description: `Saldo da conta atualizado automaticamente.`,
    })
    
    router.push(`/transactions/${transaction.transaction_id}`)
  } catch (err: any) {
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
    
    // Mostrar toast de erro
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao criar transação', {
      description: errorMessage,
    })
    
    if (formRef.value) {
      formRef.value.setError(errorMessage)
    }
  }
}

function goBack() {
  router.push('/transactions')
}
</script>
