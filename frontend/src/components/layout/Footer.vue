<template>
  <footer class="border-t bg-white py-4">
    <div class="mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
        <p class="text-sm text-gray-600">
          © {{ currentYear }} Gestão Financeira. Todos os direitos reservados.
        </p>
        <div class="flex items-center gap-4 text-sm text-gray-600">
          <span>Versão {{ appVersion }}</span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { env } from '@/config/env'

const currentYear = computed(() => new Date().getFullYear())
const appVersion = computed(() => env.appVersion)
</script>

