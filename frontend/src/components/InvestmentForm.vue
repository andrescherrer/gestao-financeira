<template>
  <Form
    @submit="handleSubmit"
    :validation-schema="validationSchema"
    :initial-values="initialValues"
    v-slot="{ errors, values }"
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
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.account_id ? 'border-destructive' : ''"
        >
          <option value="">Selecione a conta</option>
          <option
            v-for="account in accounts"
            :key="account.account_id"
            :value="account.account_id"
          >
            {{ account.name }} ({{ formatCurrency(account.balance, account.currency) }})
          </option>
        </Field>
        <ErrorMessage name="account_id" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Tipo -->
      <div>
        <Label for="type" class="mb-1">
          Tipo de Investimento <span class="text-destructive">*</span>
        </Label>
        <Field
          id="type"
          name="type"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.type ? 'border-destructive' : ''"
        >
          <option value="">Selecione o tipo</option>
          <option value="STOCK">Ação</option>
          <option value="FUND">Fundo</option>
          <option value="CDB">CDB</option>
          <option value="TREASURY">Tesouro Direto</option>
          <option value="CRYPTO">Criptomoeda</option>
          <option value="OTHER">Outro</option>
        </Field>
        <ErrorMessage name="type" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Nome -->
      <div>
        <Label for="name" class="mb-1">
          Nome do Investimento <span class="text-destructive">*</span>
        </Label>
        <Field
          id="name"
          name="name"
          type="text"
          v-slot="{ field, meta }"
        >
          <Input
            id="name"
            :name="field.name"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.name || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="Ex: Petrobras"
          />
        </Field>
        <ErrorMessage name="name" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Ticker -->
      <div>
        <Label for="ticker" class="mb-1">
          Ticker (Opcional)
        </Label>
        <Field
          id="ticker"
          name="ticker"
          type="text"
          v-slot="{ field, meta }"
        >
          <Input
            id="ticker"
            :name="field.name"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.ticker || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="Ex: PETR4"
            maxlength="20"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Código de negociação (ex: PETR4, AAPL)
        </p>
        <ErrorMessage name="ticker" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Data de Compra -->
      <div>
        <Label for="purchase_date" class="mb-1">
          Data de Compra <span class="text-destructive">*</span>
        </Label>
        <Field
          id="purchase_date"
          name="purchase_date"
          type="date"
          v-slot="{ field, meta }"
        >
          <Input
            id="purchase_date"
            :name="field.name"
            type="date"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.purchase_date || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            :max="maxDate"
          />
        </Field>
        <ErrorMessage name="purchase_date" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Valor de Compra -->
      <div>
        <Label for="purchase_amount" class="mb-1">
          Valor de Compra <span class="text-destructive">*</span>
        </Label>
        <Field
          id="purchase_amount"
          name="purchase_amount"
          type="number"
          step="0.01"
          v-slot="{ field, meta }"
        >
          <Input
            id="purchase_amount"
            :name="field.name"
            type="number"
            step="0.01"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.purchase_amount || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="0.00"
            min="0.01"
          />
        </Field>
        <ErrorMessage name="purchase_amount" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Quantidade -->
      <div v-if="requiresQuantity(values?.type)">
        <Label for="quantity" class="mb-1">
          Quantidade <span class="text-destructive">*</span>
        </Label>
        <Field
          id="quantity"
          name="quantity"
          type="number"
          step="0.0001"
          v-slot="{ field, meta }"
        >
          <Input
            id="quantity"
            :name="field.name"
            type="number"
            step="0.0001"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.quantity || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="0.0000"
            min="0.0001"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Quantidade de ações, cotas ou moedas
        </p>
        <ErrorMessage name="quantity" class="mt-1 text-sm text-destructive" />
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
import { createInvestmentSchema } from '@/validations/investment'
import { useAccountsStore } from '@/stores/accounts'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { AlertCircle, Loader2, Check } from 'lucide-vue-next'
import type { CreateInvestmentFormData } from '@/validations/investment'
import type { Account } from '@/api/types'

interface Props {
  initialValues?: Partial<CreateInvestmentFormData>
  submitLabel?: string
  isLoading?: boolean
}

interface Emits {
  (e: 'submit', values: CreateInvestmentFormData): void
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  submitLabel: 'Salvar',
})

const emit = defineEmits<Emits>()

const validationSchema = toTypedSchema(createInvestmentSchema)
const error = ref<string | null>(null)
const accountsStore = useAccountsStore()

const accounts = computed(() => accountsStore.accounts)
const maxDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

function requiresQuantity(type?: string): boolean {
  return type === 'STOCK' || type === 'FUND' || type === 'CRYPTO'
}

function formatCurrency(amount: string, currency: Account['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function handleSubmit(values: any) {
  error.value = null
  emit('submit', values as CreateInvestmentFormData)
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

