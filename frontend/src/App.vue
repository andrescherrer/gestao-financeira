<template>
  <router-view />
  <Toaster />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { Toaster } from '@/components/ui/toast'

const authStore = useAuthStore()

onMounted(() => {
  // Inicializar apenas uma vez ao montar
  authStore.init()
})
</script>
