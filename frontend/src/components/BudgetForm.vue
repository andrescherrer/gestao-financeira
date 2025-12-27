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
      <!-- Categoria -->
      <div>
        <Label for="category_id" class="mb-1">
          Categoria <span class="text-destructive">*</span>
        </Label>
        <Field
          id="category_id"
          name="category_id"
          v-slot="{ field, meta }"
        >
          <CategorySelect
            :model-value="field.value"
            @update:model-value="field.onChange"
            @blur="field.onBlur"
            :error="!!(errors.category_id || (meta.touched && !meta.valid))"
            placeholder="Selecione uma categoria"
          />
        </Field>
        <ErrorMessage name="category_id" class="mt-1 text-sm text-destructive" />
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
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            type="number"
            step="0.01"
            min="0.01"
            :class="errors.amount || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="0.00"
          />
        </Field>
        <ErrorMessage name="amount" class="mt-1 text-sm text-destructive" />
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

      <!-- Tipo de Período -->
      <div>
        <Label for="period_type" class="mb-1">
          Tipo de Período <span class="text-destructive">*</span>
        </Label>
        <Field
          id="period_type"
          name="period_type"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.period_type ? 'border-destructive' : ''"
          @change="handlePeriodTypeChange"
        >
          <option value="">Selecione o tipo</option>
          <option value="MONTHLY">Mensal</option>
          <option value="YEARLY">Anual</option>
        </Field>
        <ErrorMessage name="period_type" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Ano -->
      <div>
        <Label for="year" class="mb-1">
          Ano <span class="text-destructive">*</span>
        </Label>
        <Field
          id="year"
          name="year"
          type="number"
          min="1900"
          max="3000"
          v-slot="{ field, meta }"
        >
          <Input
            id="year"
            :name="field.name"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            type="number"
            min="1900"
            max="3000"
            :class="errors.year || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            :placeholder="currentYear.toString()"
          />
        </Field>
        <ErrorMessage name="year" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Mês (apenas para MONTHLY) -->
      <div v-if="showMonthField">
        <Label for="month" class="mb-1">
          Mês <span class="text-destructive">*</span>
        </Label>
        <Field
          id="month"
          name="month"
          as="select"
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.month ? 'border-destructive' : ''"
        >
          <option value="">Selecione o mês</option>
          <option :value="1">Janeiro</option>
          <option :value="2">Fevereiro</option>
          <option :value="3">Março</option>
          <option :value="4">Abril</option>
          <option :value="5">Maio</option>
          <option :value="6">Junho</option>
          <option :value="7">Julho</option>
          <option :value="8">Agosto</option>
          <option :value="9">Setembro</option>
          <option :value="10">Outubro</option>
          <option :value="11">Novembro</option>
          <option :value="12">Dezembro</option>
        </Field>
        <ErrorMessage name="month" class="mt-1 text-sm text-destructive" />
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

      <!-- Botões -->
      <div class="flex gap-3 justify-end">
        <Button
          type="button"
          variant="outline"
          @click="$emit('cancel')"
          :disabled="isLoading"
        >
          Cancelar
        </Button>
        <Button
          type="submit"
          :disabled="isLoading"
        >
          <Loader2 v-if="isLoading" class="h-4 w-4 mr-2 animate-spin" />
          {{ submitLabel }}
        </Button>
      </div>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createBudgetSchema, updateBudgetSchema, type CreateBudgetFormData, type UpdateBudgetFormData } from '@/validations/budget'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import CategorySelect from '@/components/CategorySelect.vue'
import { Loader2, AlertCircle } from 'lucide-vue-next'
import type { Budget } from '@/api/types'

interface Props {
  budget?: Budget | null
  isLoading?: boolean
  error?: string | null
  submitLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  budget: null,
  isLoading: false,
  error: null,
  submitLabel: 'Salvar',
})

const emit = defineEmits<{
  'submit': [data: CreateBudgetFormData | UpdateBudgetFormData]
  'cancel': []
}>()

// Schema de validação (criação ou atualização)
const validationSchema = computed(() => {
  return props.budget ? toTypedSchema(updateBudgetSchema) : toTypedSchema(createBudgetSchema)
})

// Valores iniciais
const initialValues = computed(() => {
  if (props.budget) {
    return {
      category_id: props.budget.category_id,
      amount: props.budget.amount,
      currency: props.budget.currency,
      period_type: props.budget.period_type,
      year: props.budget.year,
      month: props.budget.month,
      context: props.budget.context,
    }
  }
  const currentDate = new Date()
  return {
    category_id: '',
    amount: '',
    currency: 'BRL',
    period_type: '',
    year: currentDate.getFullYear(),
    month: currentDate.getMonth() + 1,
    context: '',
  }
})

// Ano atual
const currentYear = new Date().getFullYear()

// Mostrar campo de mês apenas para MONTHLY
const periodType = ref<string>(initialValues.value.period_type || '')
const showMonthField = computed(() => periodType.value === 'MONTHLY')

function handlePeriodTypeChange(event: Event) {
  const target = event.target as HTMLSelectElement
  periodType.value = target.value
}

function handleSubmit(values: any) {
  // Limpar mês se for YEARLY
  if (values.period_type === 'YEARLY') {
    values.month = undefined
  }
  emit('submit', values)
}
</script>

