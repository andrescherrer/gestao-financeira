<template>
  <div class="relative overflow-hidden" :class="containerClass">
    <!-- Placeholder enquanto carrega -->
    <div
      v-if="!loaded && showPlaceholder"
      class="absolute inset-0 bg-muted animate-pulse"
      :aria-hidden="true"
    ></div>

    <!-- Imagem otimizada -->
    <img
      :src="src"
      :srcset="srcset"
      :sizes="sizes"
      :alt="alt"
      :loading="lazy ? 'lazy' : 'eager'"
      :decoding="decoding"
      :class="imageClass"
      @load="handleLoad"
      @error="handleError"
      :style="imageStyle"
    />

    <!-- Erro de carregamento -->
    <div
      v-if="hasError"
      class="absolute inset-0 flex items-center justify-center bg-muted"
      :aria-label="`Erro ao carregar imagem: ${alt}`"
    >
      <ImageOff class="h-8 w-8 text-muted-foreground" aria-hidden="true" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ImageOff } from 'lucide-vue-next'

interface Props {
  src: string
  alt: string
  srcset?: string
  sizes?: string
  lazy?: boolean
  decoding?: 'async' | 'auto' | 'sync'
  containerClass?: string
  imageClass?: string
  imageStyle?: string | Record<string, string>
  showPlaceholder?: boolean
  aspectRatio?: string // Ex: "16/9", "1/1"
  objectFit?: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down'
}

const props = withDefaults(defineProps<Props>(), {
  lazy: true,
  decoding: 'async',
  showPlaceholder: true,
  objectFit: 'cover',
})

const loaded = ref(false)
const hasError = ref(false)

function handleLoad() {
  loaded.value = true
  hasError.value = false
}

function handleError() {
  hasError.value = true
  loaded.value = false
}

const imageStyle = computed(() => {
  const style: Record<string, string> = {
    objectFit: props.objectFit,
  }

  if (props.aspectRatio) {
    style.aspectRatio = props.aspectRatio
  }

  if (props.imageStyle && typeof props.imageStyle === 'object') {
    Object.assign(style, props.imageStyle)
  }

  return style
})
</script>

