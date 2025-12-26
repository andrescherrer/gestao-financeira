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
  
  // Sempre validar token se houver, mesmo em rotas públicas
  // Isso garante que o estado de autenticação está correto
  const hasToken = authStore.token || authService.getToken()
  if (hasToken) {
    const isValid = await authStore.validateToken()
    
    // Se o token é inválido e estamos em uma rota protegida, redirecionar para login
    if (!isValid) {
      const currentPath = window.location.pathname
      if (currentPath !== '/login' && currentPath !== '/register') {
        // Usar window.location para garantir redirecionamento completo
        window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        return
      }
    }
  } else {
    // Se não há token, garantir que isValidated está false
    // Isso evita problemas se o usuário tentar acessar rotas protegidas
    if (authStore.isValidated && !hasToken) {
      // Resetar validação se não há token
      authStore.logout()
    }
  }
})
</script>
