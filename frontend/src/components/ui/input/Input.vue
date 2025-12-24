<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { useVModel } from "@vueuse/core"
import { computed } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<{
  defaultValue?: string | number
  modelValue?: string | number
  value?: string | number
  class?: HTMLAttributes["class"]
  type?: string
  id?: string
  name?: string
  placeholder?: string
  disabled?: boolean
  onInput?: (event: Event) => void
  onChange?: (event: Event) => void
  onBlur?: (event: Event) => void
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
  const target = event.target as HTMLInputElement
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
  <input
    :id="id"
    :name="name"
    :type="type || 'text'"
    :value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    @input="handleInput"
    @blur="handleBlur"
    :class="cn('flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50', props.class)"
  />
</template>
