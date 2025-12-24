<template>
  <Form
    @submit="handleSubmit"
    :validation-schema="validationSchema"
    :initial-values="initialValues"
    v-slot="{ errors }"
  >
    <Card v-if="error" class="mb-4 border-destructive">
      <CardContent class="p-3">
        <div class="flex items-center gap-2">
          <AlertCircle class="h-4 w-4 text-destructive" />
          <p class="text-sm text-destructive">{{ error }}</p>
        </div>
      </CardContent>
    </Card>

    <div class="space-y-6">
      <!-- Nome -->
      <div>
        <Label for="name" class="mb-1">
          Nome da Conta <span class="text-destructive">*</span>
        </Label>
        <Field
          id="name"
          name="name"
          type="text"
          as-child
        >
          <Input
            :class="errors.name ? 'border-destructive' : ''"
            placeholder="Ex: Conta Corrente Banco do Brasil"
          />
        </Field>
        <ErrorMessage name="name" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Tipo -->
      <div>
        <Label for="type" class="mb-1">
          Tipo de Conta <span class="text-destructive">*</span>
        </Label>
        <Field
          id="type"
          name="type"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.type ? 'border-destructive' : ''"
        >
          <option value="">Selecione o tipo</option>
          <option value="BANK">Banco</option>
          <option value="WALLET">Carteira</option>
          <option value="INVESTMENT">Investimento</option>
          <option value="CREDIT_CARD">Cartão de Crédito</option>
        </Field>
        <ErrorMessage name="type" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Contexto -->
      <div>
        <Label for="context" class="mb-1">
          Contexto <span class="text-destructive">*</span>
        </Label>
        <Field
          id="context"
          name="context"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.context ? 'border-destructive' : ''"
        >
          <option value="">Selecione o contexto</option>
          <option value="PERSONAL">Pessoal</option>
          <option value="BUSINESS">Negócio</option>
        </Field>
        <ErrorMessage name="context" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Moeda -->
      <div>
        <Label for="currency" class="mb-1">
          Moeda
        </Label>
        <Field
          id="currency"
          name="currency"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.currency ? 'border-destructive' : ''"
        >
          <option value="BRL">BRL - Real Brasileiro</option>
          <option value="USD">USD - Dólar Americano</option>
          <option value="EUR">EUR - Euro</option>
        </Field>
        <ErrorMessage name="currency" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Saldo Inicial -->
      <div>
        <Label for="initial_balance" class="mb-1">
          Saldo Inicial
        </Label>
        <Field
          id="initial_balance"
          name="initial_balance"
          type="number"
          step="0.01"
          as-child
        >
          <Input
            :class="errors.initial_balance ? 'border-destructive' : ''"
            placeholder="0.00"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Deixe em branco para começar com saldo zero
        </p>
        <ErrorMessage name="initial_balance" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Actions Slot -->
      <slot name="actions" :isLoading="isLoading" :errors="errors">
        <div class="flex gap-3 pt-4">
          <Button
            type="submit"
            :disabled="isLoading"
          >
            <Loader2 v-if="isLoading" class="h-4 w-4 animate-spin" />
            <Check v-else class="h-4 w-4" />
            {{ submitLabel || 'Salvar' }}
          </Button>
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
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { AlertCircle, Loader2, Check } from 'lucide-vue-next'
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
