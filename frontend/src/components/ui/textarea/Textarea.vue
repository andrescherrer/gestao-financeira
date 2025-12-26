<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { computed } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<{
  class?: HTMLAttributes["class"]
  defaultValue?: string | number
  modelValue?: string | number
  value?: string | number // Adicionado para compatibilidade com vee-validate Field
  id?: string
  name?: string
  placeholder?: string
  disabled?: boolean
  rows?: number
  onInput?: (event: Event) => void // Adicionado para compatibilidade com vee-validate Field
  onChange?: (event: Event) => void // Adicionado para compatibilidade com vee-validate Field
  onBlur?: (event: Event) => void // Adicionado para compatibilidade com vee-validate Field
}>()

const emits = defineEmits<{
  (e: "update:modelValue", payload: string | number): void
}>()

// Suporta tanto v-model quanto value do vee-validate Field
const modelValue = computed({
  get: () => {
    // Se tem value (do vee-validate Field), usa value
    if (props.value !== undefined) {
      return props.value
    }
    // Caso contrário, usa modelValue (v-model)
    return props.modelValue ?? props.defaultValue ?? ''
  },
  set: (val: string | number) => {
    emits("update:modelValue", val)
    // Se tem onChange (do vee-validate Field), chama também
    if (props.onChange) {
      const event = new Event('change', { bubbles: true })
      Object.defineProperty(event, 'target', {
        value: { value: val },
        enumerable: true
      })
      props.onChange(event as any)
    }
  }
})

function handleInput(event: Event) {
  const target = event.target as HTMLTextAreaElement
  modelValue.value = target.value
  // Chama onInput se fornecido (do vee-validate Field)
  if (props.onInput) {
    props.onInput(event)
  }
}

function handleBlur(event: Event) {
  // Chama onBlur se fornecido (do vee-validate Field)
  if (props.onBlur) {
    props.onBlur(event)
  }
}
</script>

<template>
  <textarea
    :id="id"
    :name="name"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :rows="rows"
    @input="handleInput"
    @blur="handleBlur"
    :class="cn('flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50', props.class)"
  />
</template>
