<template>
  <header class="sticky top-0 z-50 border-b bg-white shadow-sm">
    <div class="mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
      <!-- Logo e Título -->
      <div class="flex items-center gap-3">
        <RouterLink to="/" class="flex items-center gap-2">
          <i class="pi pi-wallet text-2xl text-blue-600"></i>
          <span class="text-xl font-bold text-gray-900">Gestão Financeira</span>
        </RouterLink>
      </div>

      <!-- Navegação Desktop -->
      <nav class="hidden items-center gap-6 md:flex">
        <RouterLink
          to="/"
          class="text-sm font-medium text-gray-700 transition-colors hover:text-blue-600"
          :class="{ 'text-blue-600': $route.name === 'home' }"
        >
          Dashboard
        </RouterLink>
        <RouterLink
          to="/accounts"
          class="text-sm font-medium text-gray-700 transition-colors hover:text-blue-600"
          :class="{ 'text-blue-600': $route.name === 'accounts' || $route.name === 'new-account' || $route.name === 'account-details' }"
        >
          Contas
        </RouterLink>
        <RouterLink
          to="/transactions"
          class="text-sm font-medium text-gray-700 transition-colors hover:text-blue-600"
          :class="{ 'text-blue-600': $route.name === 'transactions' || $route.name === 'new-transaction' || $route.name === 'transaction-details' }"
        >
          Transações
        </RouterLink>
      </nav>

      <!-- Menu do Usuário -->
      <div class="flex items-center gap-4">
        <div v-if="authStore.user" class="hidden items-center gap-2 sm:flex">
          <div class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-sm font-semibold text-blue-600">
            {{ userInitials }}
          </div>
          <span class="text-sm font-medium text-gray-700">{{ userName }}</span>
        </div>
        <button
          @click="handleLogout"
          class="flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-100 hover:text-red-600"
        >
          <i class="pi pi-sign-out"></i>
          <span class="hidden sm:inline">Sair</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const userName = computed(() => {
  if (!authStore.user) return 'Usuário'
  const firstName = authStore.user.first_name || ''
  const lastName = authStore.user.last_name || ''
  return `${firstName} ${lastName}`.trim() || authStore.user.email
})

const userInitials = computed(() => {
  if (!authStore.user) return 'U'
  const firstName = authStore.user.first_name || ''
  const lastName = authStore.user.last_name || ''
  if (firstName && lastName) {
    return `${firstName[0]}${lastName[0]}`.toUpperCase()
  }
  if (firstName) {
    return firstName[0].toUpperCase()
  }
  return authStore.user.email?.[0]?.toUpperCase() || 'U'
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

