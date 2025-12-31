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
      <!-- Nome -->
      <div>
        <Label for="name" class="mb-1">
          Nome da Meta <span class="text-destructive">*</span>
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
            placeholder="Ex: Viagem para Europa"
          />
        </Field>
        <ErrorMessage name="name" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Valor Alvo -->
      <div>
        <Label for="target_amount" class="mb-1">
          Valor Alvo <span class="text-destructive">*</span>
        </Label>
        <Field
          id="target_amount"
          name="target_amount"
          type="number"
          step="0.01"
          min="0.01"
          v-slot="{ field, meta }"
        >
          <Input
            id="target_amount"
            :name="field.name"
            type="number"
            step="0.01"
            min="0.01"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.target_amount || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="0.00"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Valor total que você deseja alcançar
        </p>
        <ErrorMessage name="target_amount" class="mt-1 text-sm text-destructive" />
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
          class="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
          :class="errors.currency ? 'border-destructive' : ''"
        >
          <option value="BRL">BRL (Real Brasileiro)</option>
          <option value="USD">USD (Dólar Americano)</option>
          <option value="EUR">EUR (Euro)</option>
        </Field>
        <ErrorMessage name="currency" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Data Limite -->
      <div>
        <Label for="deadline" class="mb-1">
          Data Limite <span class="text-destructive">*</span>
        </Label>
        <Field
          id="deadline"
          name="deadline"
          type="date"
          v-slot="{ field, meta }"
        >
          <Input
            id="deadline"
            :name="field.name"
            type="date"
            :value="field.value"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.deadline || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            :min="minDate"
          />
        </Field>
        <p class="mt-1 text-xs text-muted-foreground">
          Data em que você deseja alcançar esta meta
        </p>
        <ErrorMessage name="deadline" class="mt-1 text-sm text-destructive" />
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
      <div class="flex justify-end gap-3">
        <Button
          type="button"
          variant="outline"
          @click="$emit('cancel')"
          :disabled="isLoading"
        >
          Cancelar
        </Button>
        <Button type="submit" :disabled="isLoading">
          <Loader2 v-if="isLoading" class="mr-2 h-4 w-4 animate-spin" />
          <Check v-else class="mr-2 h-4 w-4" />
          {{ submitLabel }}
        </Button>
      </div>
    </div>
  </Form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { createGoalSchema } from '@/validations/goal'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { AlertCircle, Loader2, Check } from 'lucide-vue-next'
import type { CreateGoalFormData } from '@/validations/goal'

interface Props {
  initialValues?: Partial<CreateGoalFormData>
  submitLabel?: string
  isLoading?: boolean
}

interface Emits {
  (e: 'submit', values: CreateGoalFormData): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: false,
  submitLabel: 'Salvar',
})

const emit = defineEmits<Emits>()

const validationSchema = toTypedSchema(createGoalSchema)
const error = ref<string | null>(null)

const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

function handleSubmit(values: any) {
  error.value = null
  emit('submit', values as CreateGoalFormData)
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

