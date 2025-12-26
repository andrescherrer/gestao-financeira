<template>
  <Select
    :model-value="modelValue"
    @update:model-value="handleChange"
    :disabled="isLoading || disabled"
  >
    <SelectTrigger :class="errorClass" @blur="$emit('blur')">
      <SelectValue :placeholder="placeholder" />
    </SelectTrigger>
    <SelectContent>
      <SelectGroup>
        <SelectLabel v-if="showLabel">{{ label }}</SelectLabel>
        <SelectItem value="">
          <span class="text-muted-foreground">{{ placeholder }}</span>
        </SelectItem>
        <SelectItem
          v-for="category in activeCategories"
          :key="category.category_id"
          :value="category.category_id"
        >
          {{ category.name }}
        </SelectItem>
      </SelectGroup>
    </SelectContent>
  </Select>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import type { AcceptableValue } from 'reka-ui'
import { useCategoriesStore } from '@/stores/categories'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

interface Props {
  modelValue?: string
  placeholder?: string
  label?: string
  showLabel?: boolean
  disabled?: boolean
  error?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Selecione uma categoria',
  label: 'Categorias',
  showLabel: false,
  disabled: false,
  error: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'blur': []
}>()

const categoriesStore = useCategoriesStore()

const activeCategories = computed(() => categoriesStore.activeCategories)
const isLoading = computed(() => categoriesStore.isLoading)

const errorClass = computed(() => {
  return props.error ? 'border-destructive' : ''
})

onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.listCategories(true) // Apenas categorias ativas
  }
})

function handleChange(value: AcceptableValue) {
  // Garantir que sempre emite uma string
  if (value === null || value === undefined || value === '') {
    emit('update:modelValue', '')
  } else if (typeof value === 'object') {
    // Se for um objeto (não deveria acontecer, mas por segurança)
    console.warn('CategorySelect recebeu um objeto em vez de string:', value)
    emit('update:modelValue', '')
  } else {
    emit('update:modelValue', String(value))
  }
}
</script>

