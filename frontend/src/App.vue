<template>
  <router-view />
  <Toaster />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/api/auth'
import { Toaster } from '@/components/ui/toast'

const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  // Inicializar token do localStorage
  authStore.init()
  
  // Se há token, validar com o backend
  if (authStore.token || authService.getToken()) {
    const isValid = await authStore.validateToken()
    
    // Se o token é inválido e estamos em uma rota protegida, redirecionar para login
    if (!isValid) {
      const currentPath = window.location.pathname
      if (currentPath !== '/login' && currentPath !== '/register') {
        await router.push({ name: 'login', query: { redirect: currentPath } })
      }
    }
  }
})
</script>
