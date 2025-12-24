<template>
  <Form
    @submit="handleSubmit"
    :validation-schema="validationSchema"
    :initial-values="initialValues"
    v-slot="{ errors }"
  >
    <div v-if="error" class="mb-4 rounded-md bg-red-50 border border-red-200 p-3">
      <div class="flex items-center gap-2">
        <i class="pi pi-exclamation-circle text-red-600"></i>
        <p class="text-sm text-red-600">{{ error }}</p>
      </div>
    </div>

    <div class="space-y-6">
      <!-- Conta -->
      <div>
        <label for="account_id" class="block text-sm font-medium text-gray-700 mb-1">
          Conta <span class="text-red-500">*</span>
        </label>
        <Field
          id="account_id"
          name="account_id"
          as="select"
          class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
          :class="
            errors.account_id
              ? 'border-red-300 focus:border-red-500'
              : 'border-gray-300 focus:border-blue-500'
          "
          :disabled="isLoading"
        >
          <option value="">Selecione uma conta</option>
          <option
            v-for="account in accounts"
            :key="account.account_id"
            :value="account.account_id"
          >
            {{ account.name }} ({{ account.currency }})
          </option>
        </Field>
        <ErrorMessage name="account_id" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Tipo -->
      <div>
        <label for="type" class="block text-sm font-medium text-gray-700 mb-1">
          Tipo <span class="text-red-500">*</span>
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
          :disabled="isLoading"
        >
          <option value="">Selecione o tipo</option>
          <option value="INCOME">Receita</option>
          <option value="EXPENSE">Despesa</option>
        </Field>
        <ErrorMessage name="type" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Valor -->
      <div>
        <label for="amount" class="block text-sm font-medium text-gray-700 mb-1">
          Valor <span class="text-red-500">*</span>
        </label>
        <Field
          id="amount"
          name="amount"
          type="number"
          step="0.01"
          min="0.01"
          class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
          :class="
            errors.amount
              ? 'border-red-300 focus:border-red-500'
              : 'border-gray-300 focus:border-blue-500'
          "
          placeholder="0.00"
          :disabled="isLoading"
        />
        <ErrorMessage name="amount" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Moeda -->
      <div>
        <label for="currency" class="block text-sm font-medium text-gray-700 mb-1">
          Moeda <span class="text-red-500">*</span>
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
          :disabled="isLoading"
        >
          <option value="BRL">BRL - Real Brasileiro</option>
          <option value="USD">USD - Dólar Americano</option>
          <option value="EUR">EUR - Euro</option>
        </Field>
        <ErrorMessage name="currency" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Descrição -->
      <div>
        <label for="description" class="block text-sm font-medium text-gray-700 mb-1">
          Descrição <span class="text-red-500">*</span>
        </label>
        <Field
          id="description"
          name="description"
          as="textarea"
          rows="3"
          class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
          :class="
            errors.description
              ? 'border-red-300 focus:border-red-500'
              : 'border-gray-300 focus:border-blue-500'
          "
          placeholder="Descreva a transação..."
          :disabled="isLoading"
        />
        <p class="mt-1 text-xs text-gray-500">
          Mínimo 3 caracteres, máximo 500 caracteres
        </p>
        <ErrorMessage name="description" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Data -->
      <div>
        <label for="date" class="block text-sm font-medium text-gray-700 mb-1">
          Data <span class="text-red-500">*</span>
        </label>
        <Field
          id="date"
          name="date"
          type="date"
          class="block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
          :class="
            errors.date
              ? 'border-red-300 focus:border-red-500'
              : 'border-gray-300 focus:border-blue-500'
          "
          :disabled="isLoading"
        />
        <p class="mt-1 text-xs text-gray-500">
          Data da transação (formato: YYYY-MM-DD)
        </p>
        <ErrorMessage name="date" class="mt-1 text-sm text-red-600" />
      </div>

      <!-- Actions Slot -->
      <slot name="actions" :isLoading="isLoading" :errors="errors">
        <div class="flex gap-3 pt-4">
          <button
            type="submit"
            :disabled="isLoading"
            class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            <i v-if="isLoading" class="pi pi-spinner pi-spin"></i>
            <i v-else class="pi pi-check"></i>
            {{ submitLabel || 'Salvar' }}
          </button>
        </div>
      </slot>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createTransactionSchema } from '@/validations/transaction'
import { useAccountsStore } from '@/stores/accounts'
import type { CreateTransactionFormData } from '@/validations/transaction'

interface Props {
  initialValues?: Partial<CreateTransactionFormData>
  submitLabel?: string
  isLoading?: boolean
}

interface Emits {
  (e: 'submit', values: CreateTransactionFormData): void
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  submitLabel: 'Salvar',
})

const emit = defineEmits<Emits>()

const accountsStore = useAccountsStore()

const validationSchema = toTypedSchema(createTransactionSchema)
const error = ref<string | null>(null)

const accounts = computed(() => accountsStore.accounts)

function handleSubmit(values: any) {
  error.value = null
  emit('submit', values as CreateTransactionFormData)
}

// Expor função para definir erro externamente
function setError(message: string) {
  error.value = message
}

function clearError() {
  error.value = null
}

defineExpose({
  setError,
  clearError,
})
</script>

