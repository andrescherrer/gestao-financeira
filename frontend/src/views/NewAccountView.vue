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
        <Form
          @submit="handleSubmit"
          :validation-schema="validationSchema"
          v-slot="{ errors }"
        >
          <div v-if="error" class="mb-4 rounded-md bg-red-50 border border-red-200 p-3">
            <div class="flex items-center gap-2">
              <i class="pi pi-exclamation-circle text-red-600"></i>
              <p class="text-sm text-red-600">{{ error }}</p>
            </div>
          </div>

          <div class="space-y-6">
            <!-- Nome -->
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 mb-1">
                Nome da Conta <span class="text-red-500">*</span>
              </label>
              <Field
                id="name"
                name="name"
                type="text"
                class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
                :class="
                  errors.name
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                "
                placeholder="Ex: Conta Corrente Banco do Brasil"
              />
              <ErrorMessage name="name" class="mt-1 text-sm text-red-600" />
            </div>

            <!-- Tipo -->
            <div>
              <label for="type" class="block text-sm font-medium text-gray-700 mb-1">
                Tipo de Conta <span class="text-red-500">*</span>
              </label>
              <Field
                id="type"
                name="type"
                as="select"
                class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
                :class="
                  errors.type
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                "
              >
                <option value="">Selecione o tipo</option>
                <option value="BANK">Banco</option>
                <option value="WALLET">Carteira</option>
                <option value="INVESTMENT">Investimento</option>
                <option value="CREDIT_CARD">Cartão de Crédito</option>
              </Field>
              <ErrorMessage name="type" class="mt-1 text-sm text-red-600" />
            </div>

            <!-- Contexto -->
            <div>
              <label for="context" class="block text-sm font-medium text-gray-700 mb-1">
                Contexto <span class="text-red-500">*</span>
              </label>
              <Field
                id="context"
                name="context"
                as="select"
                class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
                :class="
                  errors.context
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                "
              >
                <option value="">Selecione o contexto</option>
                <option value="PERSONAL">Pessoal</option>
                <option value="BUSINESS">Negócio</option>
              </Field>
              <ErrorMessage name="context" class="mt-1 text-sm text-red-600" />
            </div>

            <!-- Moeda -->
            <div>
              <label for="currency" class="block text-sm font-medium text-gray-700 mb-1">
                Moeda
              </label>
              <Field
                id="currency"
                name="currency"
                as="select"
                class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
                :class="
                  errors.currency
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                "
              >
                <option value="BRL">BRL - Real Brasileiro</option>
                <option value="USD">USD - Dólar Americano</option>
                <option value="EUR">EUR - Euro</option>
              </Field>
              <ErrorMessage name="currency" class="mt-1 text-sm text-red-600" />
            </div>

            <!-- Saldo Inicial -->
            <div>
              <label for="initial_balance" class="block text-sm font-medium text-gray-700 mb-1">
                Saldo Inicial
              </label>
              <Field
                id="initial_balance"
                name="initial_balance"
                type="number"
                step="0.01"
                class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
                :class="
                  errors.initial_balance
                    ? 'border-red-300 focus:border-red-500'
                    : 'border-gray-300 focus:border-blue-500'
                "
                placeholder="0.00"
              />
              <p class="mt-1 text-xs text-gray-500">
                Deixe em branco para começar com saldo zero
              </p>
              <ErrorMessage name="initial_balance" class="mt-1 text-sm text-red-600" />
            </div>

            <!-- Actions -->
            <div class="flex gap-3 pt-4">
              <button
                type="submit"
                :disabled="accountsStore.isLoading"
                class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                <i v-if="accountsStore.isLoading" class="pi pi-spinner pi-spin"></i>
                <i v-else class="pi pi-check"></i>
                {{ accountsStore.isLoading ? 'Criando...' : 'Criar Conta' }}
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
          </div>
        </Form>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import { createAccountSchema } from '@/validations/account'

const validationSchema = toTypedSchema(createAccountSchema)

const router = useRouter()
const accountsStore = useAccountsStore()

const error = ref<string | null>(null)

async function handleSubmit(values: any) {
  error.value = null
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
    error.value =
      err.response?.data?.message ||
      err.response?.data?.error ||
      err.message ||
      'Erro ao criar conta. Tente novamente.'
  }
}

function goBack() {
  router.push('/accounts')
}
</script>

