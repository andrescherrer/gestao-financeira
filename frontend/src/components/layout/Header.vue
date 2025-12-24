<template>
  <header class="sticky top-0 z-50 border-b bg-gradient-to-r from-white to-gray-50 shadow-sm backdrop-blur-sm">
    <div class="mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
      <!-- Logo e Título -->
      <div class="flex items-center gap-3">
        <RouterLink to="/" class="flex items-center gap-2 transition-transform hover:scale-105">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 text-white shadow-md">
            <i class="pi pi-wallet text-xl"></i>
          </div>
          <span class="text-xl font-bold bg-gradient-to-r from-blue-600 to-blue-800 bg-clip-text text-transparent">Gestão Financeira</span>
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
        <div v-if="authStore.user" class="hidden items-center gap-3 sm:flex">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-blue-600 text-sm font-semibold text-white shadow-md">
            {{ userInitials }}
          </div>
          <div class="flex flex-col">
            <span class="text-sm font-semibold text-gray-900">{{ userName }}</span>
            <span class="text-xs text-gray-500">{{ authStore.user.email }}</span>
          </div>
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
  if (firstName && firstName[0]) {
    return firstName[0].toUpperCase()
  }
  return authStore.user.email?.[0]?.toUpperCase() || 'U'
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

