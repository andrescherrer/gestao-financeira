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
        <h1 class="text-4xl font-bold mb-2">Editar Conta</h1>
        <p class="text-gray-600">
          Atualize os dados da conta financeira
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading && !initialValues" class="flex items-center justify-center py-12">
        <div class="text-center">
          <i class="pi pi-spinner pi-spin text-4xl text-blue-600 mb-4"></i>
          <p class="text-gray-600">Carregando dados da conta...</p>
        </div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="accountsStore.error && !initialValues"
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

      <!-- Form -->
      <div v-else-if="initialValues" class="rounded-lg border border-gray-200 bg-white p-6">
        <AccountForm
          ref="formRef"
          :initial-values="initialValues"
          :isLoading="accountsStore.isLoading"
          submitLabel="Salvar Alterações"
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
                {{ formLoading ? 'Salvando...' : 'Salvar Alterações' }}
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import AccountForm from '@/components/AccountForm.vue'
import type { CreateAccountFormData } from '@/validations/account'

const route = useRoute()
const router = useRouter()
const accountsStore = useAccountsStore()
const formRef = ref<InstanceType<typeof AccountForm> | null>(null)

const accountId = route.params.id as string

const initialValues = computed<Partial<CreateAccountFormData> | null>(() => {
  const account = accountsStore.currentAccount
  if (!account) return null

  return {
    name: account.name,
    type: account.type,
    context: account.context,
    currency: account.currency,
    initial_balance: account.balance !== '0.00' ? account.balance : undefined,
  }
})

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

  if (existingAccount) {
    accountsStore.currentAccount = existingAccount
  } else {
    try {
      await accountsStore.getAccount(accountId)
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

async function handleSubmit(values: CreateAccountFormData) {
  if (formRef.value) {
    formRef.value.clearError()
  }

  try {
    // Por enquanto, apenas redirecionar de volta
    // TODO: Implementar atualização de conta quando o backend suportar
    router.push(`/accounts/${accountId}`)
  } catch (err: any) {
    const errorMessage =
      err.response?.data?.message ||
      err.response?.data?.error ||
      err.message ||
      'Erro ao atualizar conta. Tente novamente.'
    
    if (formRef.value) {
      formRef.value.setError(errorMessage)
    }
  }
}

function goBack() {
  router.push(`/accounts/${accountId}`)
}

function handleRetry() {
  accountsStore.clearError()
  loadAccount()
}
</script>

