<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Contas', to: '/accounts' },
          { label: accountsStore.currentAccount?.name || 'Editar', to: `/accounts/${accountId}` },
          { label: 'Editar' },
        ]"
      />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Editar Conta</h1>
        <p class="text-muted-foreground">
          Atualize os dados da conta financeira
        </p>
      </div>

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading && !initialValues" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando dados da conta...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="accountsStore.error && !initialValues"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ accountsStore.error }}</p>
          </div>
          <div class="flex gap-3">
            <Button
              @click="handleRetry"
              variant="destructive"
            >
              Tentar novamente
            </Button>
            <Button
              variant="outline"
              @click="goBack"
            >
              Voltar
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Form -->
      <Card v-else-if="initialValues">
        <CardContent class="p-6">
          <AccountForm
            ref="formRef"
            :initial-values="initialValues"
            :isLoading="accountsStore.isLoading"
            submitLabel="Salvar Alterações"
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
                  {{ formLoading ? 'Salvando...' : 'Salvar Alterações' }}
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
          </AccountForm>
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import AccountForm from '@/components/AccountForm.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Loader2, Check, X, AlertCircle } from 'lucide-vue-next'
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
