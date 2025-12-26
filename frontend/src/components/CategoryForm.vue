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
          Nome da Categoria <span class="text-destructive">*</span>
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
            placeholder="Ex: Alimentação"
          />
        </Field>
        <ErrorMessage name="name" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Descrição -->
      <div>
        <Label for="description" class="mb-1">
          Descrição
        </Label>
        <Field
          id="description"
          name="description"
          v-slot="{ field, meta }"
        >
          <Textarea
            id="description"
            :name="field.name"
            :value="field.value || ''"
            @input="field.onInput"
            @change="field.onChange"
            @blur="field.onBlur"
            :class="errors.description || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
            placeholder="Descreva a categoria (opcional)"
            :rows="4"
          />
        </Field>
        <ErrorMessage name="description" class="mt-1 text-sm text-destructive" />
      </div>

      <!-- Actions -->
      <div class="flex gap-3 justify-end">
        <Button
          type="button"
          variant="outline"
          @click="$emit('cancel')"
        >
          Cancelar
        </Button>
        <Button
          type="submit"
          :disabled="isSubmitting"
        >
          <Loader2 v-if="isSubmitting" class="h-4 w-4 mr-2 animate-spin" />
          <Check v-else-if="submitted" class="h-4 w-4 mr-2" />
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
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Loader2, Check, AlertCircle } from 'lucide-vue-next'
import type { CreateCategoryFormData, UpdateCategoryFormData } from '@/validations/category'
import { createCategorySchema, updateCategorySchema } from '@/validations/category'

interface Props {
  initialValues?: Partial<CreateCategoryFormData | UpdateCategoryFormData>
  validationSchema?: ReturnType<typeof toTypedSchema<typeof createCategorySchema>>
  submitLabel?: string
  isSubmitting?: boolean
  error?: string | null
}

const props = withDefaults(defineProps<Props>(), {
  submitLabel: 'Salvar',
  isSubmitting: false,
  error: null,
})

const defaultValidationSchema = toTypedSchema(createCategorySchema)
const validationSchema = computed(() => props.validationSchema || defaultValidationSchema)

const emit = defineEmits<{
  submit: [values: CreateCategoryFormData | UpdateCategoryFormData]
  cancel: []
}>()

const submitted = ref(false)

function handleSubmit(values: CreateCategoryFormData | UpdateCategoryFormData) {
  submitted.value = true
  emit('submit', values)
  // Reset submitted after a delay to show checkmark
  setTimeout(() => {
    submitted.value = false
  }, 2000)
}

// Expose methods for parent component
defineExpose({
  clearError: () => {
    // Error is managed by parent
  },
  setError: (message: string) => {
    // Error is managed by parent
  },
})
</script>

