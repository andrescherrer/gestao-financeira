<template>
  <router-view />
  <Toaster />
  <PWAUpdatePrompt />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/api/auth'
import { Toaster } from '@/components/ui/toast'
import { useTheme } from '@/composables/useTheme'
import PWAUpdatePrompt from '@/components/PWAUpdatePrompt.vue'

const router = useRouter()
const authStore = useAuthStore()

// Inicializar tema
useTheme()

onMounted(async () => {
  // Inicializar token do localStorage
  authStore.init()
  
  // A validação do token já é feita no router guard (beforeEach)
  // Não precisamos validar novamente aqui para evitar chamadas duplicadas
  // O router guard já cuida de redirecionar se necessário
})
</script>
