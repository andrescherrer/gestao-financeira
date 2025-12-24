<script setup lang="ts">
import type { HTMLAttributes } from "vue"
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
}>()

const emits = defineEmits<{
  (e: "update:modelValue", payload: string | number): void
  (e: "input", event: Event): void
  (e: "change", event: Event): void
  (e: "blur", event: Event): void
}>()

// Suporta tanto v-model quanto value do vee-validate Field
const inputValue = computed({
  get: () => {
    // Se tem value (do vee-validate Field), usa value
    if (props.value !== undefined && props.value !== null) {
      return props.value
    }
    // Caso contrário, usa modelValue (v-model)
    return props.modelValue ?? props.defaultValue ?? ''
  },
  set: (val: string | number) => {
    emits("update:modelValue", val)
  }
})

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  inputValue.value = target.value
  emits("input", event)
}

function handleChange(event: Event) {
  emits("change", event)
}

function handleBlur(event: Event) {
  emits("blur", event)
}
</script>

<template>
  <input
    :id="id"
    :name="name"
    :type="type || 'text'"
    :value="inputValue"
    :placeholder="placeholder"
    :disabled="disabled"
    @input="handleInput"
    @change="handleChange"
    @blur="handleBlur"
    :class="cn('flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50', props.class)"
  />
</template>
