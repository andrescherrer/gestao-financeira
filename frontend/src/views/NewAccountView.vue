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
          <span>Voltar para contas</span>
        </button>
        <h1 class="text-4xl font-bold mb-2">Nova Conta</h1>
        <p class="text-gray-600">
          Preencha os dados para criar uma nova conta financeira
        </p>
      </div>

      <!-- Form -->
      <div class="rounded-lg border border-gray-200 bg-white p-6">
        <AccountForm
          ref="formRef"
          :isLoading="accountsStore.isLoading"
          submitLabel="Criar Conta"
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
                {{ formLoading ? 'Criando...' : 'Criar Conta' }}
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
        </AccountForm>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import AccountForm from '@/components/AccountForm.vue'
import type { CreateAccountFormData } from '@/validations/account'

const router = useRouter()
const accountsStore = useAccountsStore()
const formRef = ref<InstanceType<typeof AccountForm> | null>(null)

async function handleSubmit(values: CreateAccountFormData) {
  if (formRef.value) {
    formRef.value.clearError()
  }

  try {
    // Preparar dados para API
    const accountData = {
      name: values.name,
      type: values.type,
      context: values.context,
      currency: values.currency || 'BRL',
      initial_balance: values.initial_balance
        ? parseFloat(values.initial_balance).toFixed(2)
        : undefined,
    }

    const account = await accountsStore.createAccount(accountData)
    
    // Redirecionar para detalhes da conta criada
    router.push(`/accounts/${account.account_id}`)
  } catch (err: any) {
    const errorMessage =
      err.response?.data?.message ||
      err.response?.data?.error ||
      err.message ||
      'Erro ao criar conta. Tente novamente.'
    
    if (formRef.value) {
      formRef.value.setError(errorMessage)
    }
  }
}

function goBack() {
  router.push('/accounts')
}
</script>

