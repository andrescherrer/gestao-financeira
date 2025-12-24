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
import { ref } from 'vue'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createAccountSchema } from '@/validations/account'
import type { CreateAccountFormData } from '@/validations/account'

interface Props {
  initialValues?: Partial<CreateAccountFormData>
  submitLabel?: string
  isLoading?: boolean
}

interface Emits {
  (e: 'submit', values: CreateAccountFormData): void
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  submitLabel: 'Salvar',
})

const emit = defineEmits<Emits>()

const validationSchema = toTypedSchema(createAccountSchema)
const error = ref<string | null>(null)

function handleSubmit(values: any) {
  error.value = null
  emit('submit', values as CreateAccountFormData)
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

