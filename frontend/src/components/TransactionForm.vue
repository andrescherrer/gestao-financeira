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
      <!-- Conta -->
      <div>
        <Label for="account_id" class="mb-1">
          Conta <span class="text-destructive">*</span>
        </Label>
        <Field
          id="account_id"
          name="account_id"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
          :class="errors.account_id ? 'border-destructive' : ''"
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
        <ErrorMessage name="account_id" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Tipo -->
      <div>
        <Label for="type" class="mb-1">
          Tipo <span class="text-destructive">*</span>
        </Label>
        <Field
          id="type"
          name="type"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
          :class="errors.type ? 'border-destructive' : ''"
          :disabled="isLoading"
        >
          <option value="">Selecione o tipo</option>
          <option value="INCOME">Receita</option>
          <option value="EXPENSE">Despesa</option>
        </Field>
        <ErrorMessage name="type" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Valor -->
      <div>
        <Label for="amount" class="mb-1">
          Valor <span class="text-destructive">*</span>
        </Label>
        <Field
          id="amount"
          name="amount"
          type="number"
          step="0.01"
          min="0.01"
          v-slot="{ field, meta }"
        >
          <Input
            id="amount"
            :name="field.name"
            type="number"
            step="0.01"
            min="0.01"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.amount || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="0.00"
            :disabled="isLoading"
          />
        </Field>
        <ErrorMessage name="amount" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Moeda -->
      <div>
        <Label for="currency" class="mb-1">
          Moeda <span class="text-destructive">*</span>
        </Label>
        <Field
          id="currency"
          name="currency"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
          :class="errors.currency ? 'border-destructive' : ''"
          :disabled="isLoading"
        >
          <option value="BRL">BRL - Real Brasileiro</option>
          <option value="USD">USD - Dólar Americano</option>
          <option value="EUR">EUR - Euro</option>
        </Field>
        <ErrorMessage name="currency" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Categoria (Opcional) -->
      <div>
        <Label for="category_id" class="mb-1">
          Categoria
        </Label>
        <Field
          id="category_id"
          name="category_id"
          v-slot="{ field, meta }"
        >
          <CategorySelect
            :model-value="field.value || ''"
            @update:model-value="(value) => { 
              const stringValue = value || ''
              field.value = stringValue
              field.onChange(stringValue)
            }"
            @blur="field.onBlur"
            placeholder="Selecione uma categoria (opcional)"
            :error="!!(errors.category_id || (meta.touched && !meta.valid))"
            :disabled="isLoading"
          />
        </Field>
        <ErrorMessage name="category_id" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Descrição -->
      <div>
        <Label for="description" class="mb-1">
          Descrição <span class="text-destructive">*</span>
        </Label>
        <Field
          id="description"
          name="description"
          v-slot="{ field, meta }"
        >
          <Textarea
            id="description"
            :name="field.name"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.description || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="Descreva a transação..."
            :disabled="isLoading"
            :rows="3"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Mínimo 3 caracteres, máximo 500 caracteres
        </p>
        <ErrorMessage name="description" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Data -->
      <div>
        <Label for="date" class="mb-1">
          Data <span class="text-destructive">*</span>
        </Label>
        <Field
          id="date"
          name="date"
          type="date"
          v-slot="{ field, meta }"
        >
          <Input
            id="date"
            :name="field.name"
            type="date"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.date || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            :disabled="isLoading"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Data da transação (formato: YYYY-MM-DD)
        </p>
        <ErrorMessage name="date" class="mt-1 text-sm text-destructive" />
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
import { ref, computed } from 'vue'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createTransactionSchema } from '@/validations/transaction'
import { useAccountsStore } from '@/stores/accounts'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { AlertCircle, Loader2, Check } from 'lucide-vue-next'
import CategorySelect from '@/components/CategorySelect.vue'
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
