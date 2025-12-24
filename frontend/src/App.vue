<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Inicializar imediatamente ao montar o app
authStore.init()

onMounted(() => {
  // Garantir que está inicializado após montagem
  authStore.init()
})
</script>
