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
  
  // Verificar se estamos em uma rota protegida
  const currentPath = window.location.pathname
  const isProtectedRoute = currentPath !== '/login' && currentPath !== '/register'
  
  // Se estamos em rota protegida, validar token
  if (isProtectedRoute) {
    const hasToken = authStore.token || authService.getToken()
    
    if (!hasToken) {
      // Não há token, redirecionar imediatamente
      authStore.logout()
      window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
      return
    }
    
    // Validar token
    try {
      const isValid = await authStore.validateToken()
      if (!isValid) {
        // Token inválido, limpar e redirecionar
        authStore.logout()
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        return
      }
      
      // Verificar se ainda está autenticado após validação
      if (!authStore.isAuthenticated) {
        authStore.logout()
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        return
      }
    } catch (error) {
      // Em caso de erro, considerar token inválido
      console.error('Erro ao validar token no App.vue:', error)
      authStore.logout()
      window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
      return
    }
  }
})
</script>
