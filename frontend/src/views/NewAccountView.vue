<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Contas', to: '/accounts' },
          { label: 'Nova Conta' },
        ]"
      />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Nova Conta</h1>
        <p class="text-muted-foreground">
          Preencha os dados para criar uma nova conta financeira
        </p>
      </div>

      <!-- Form -->
      <Card>
        <CardContent class="p-6">
          <AccountForm
            ref="formRef"
            :isLoading="accountsStore.isLoading"
            submitLabel="Criar Conta"
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
                  {{ formLoading ? 'Criando...' : 'Criar Conta' }}
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import AccountForm from '@/components/AccountForm.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Loader2, Check, X } from 'lucide-vue-next'
import type { CreateAccountFormData } from '@/validations/account'
import type { CreateAccountRequest } from '@/api/types'

const router = useRouter()
const accountsStore = useAccountsStore()
const formRef = ref<InstanceType<typeof AccountForm> | null>(null)

async function handleSubmit(values: CreateAccountFormData) {
  if (formRef.value) {
    formRef.value.clearError()
  }

  try {
    let initialBalance = 0.0
    if (values.initial_balance) {
      const parsed = parseFloat(values.initial_balance)
      if (isNaN(parsed) || !isFinite(parsed)) {
        throw new Error('Saldo inicial deve ser um número válido')
      }
      initialBalance = parsed
    }

    const trimmedName = values.name.trim()
    if (trimmedName.length < 3) {
      throw new Error('Nome da conta deve ter no mínimo 3 caracteres')
    }

    const accountData: CreateAccountRequest = {
      name: trimmedName,
      type: values.type,
      context: values.context,
      currency: values.currency || 'BRL',
      initial_balance: initialBalance,
    }

    if (import.meta.env.DEV) {
      console.log('[NewAccountView] Dados sendo enviados:', accountData)
    }

    const account = await accountsStore.createAccount(accountData)
    router.push(`/accounts/${account.account_id}`)
  } catch (err: any) {
    if (import.meta.env.DEV) {
      console.error('[NewAccountView] Erro ao criar conta:', {
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
