<template>
  <header class="sticky top-0 z-50 h-16 border-b border-gray-200 bg-white">
    <div class="flex h-full items-center justify-between px-6">
      <!-- Left: Breadcrumbs ou título da página -->
      <div class="flex items-center gap-2 text-sm text-gray-600">
        <span class="font-medium text-gray-900">{{ pageTitle }}</span>
      </div>

      <!-- Right: Notificações e Usuário -->
      <div class="flex items-center gap-4">
        <!-- Ícone de Notificação -->
        <button
          class="relative flex h-10 w-10 items-center justify-center rounded-lg text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
          title="Notificações"
        >
          <i class="pi pi-bell text-lg"></i>
          <span
            v-if="hasNotifications"
            class="absolute right-2 top-2 h-2 w-2 rounded-full bg-red-500"
          ></span>
        </button>

        <!-- Avatar do Usuário -->
        <div v-if="authStore.user" class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 text-sm font-semibold text-white shadow-sm">
            {{ userInitials }}
          </div>
          <div class="hidden flex-col sm:flex">
            <span class="text-sm font-semibold text-gray-900">{{ userName }}</span>
            <span class="text-xs text-gray-500">{{ authStore.user.email }}</span>
          </div>
          <!-- Menu Dropdown do Usuário -->
          <div class="relative" ref="userMenuRef">
            <button
              @click="showUserMenu = !showUserMenu"
              class="flex items-center text-gray-600 hover:text-gray-900"
            >
              <i class="pi pi-chevron-down text-sm"></i>
            </button>
            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 rounded-lg border border-gray-200 bg-white py-1 shadow-lg"
            >
              <button
                @click="handleLogout"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 transition-colors hover:bg-gray-50 hover:text-red-600"
              >
                <i class="pi pi-sign-out"></i>
                <span>Sair</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const showUserMenu = ref(false)
const hasNotifications = ref(false) // Pode ser conectado a um store de notificações
const userMenuRef = ref<HTMLElement | null>(null)

function handleClickOutside(event: MouseEvent) {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target as Node)) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

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

const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    home: 'Dashboard',
    accounts: 'Contas',
    'new-account': 'Nova Conta',
    'account-details': 'Detalhes da Conta',
    'edit-account': 'Editar Conta',
    transactions: 'Transações',
    'new-transaction': 'Nova Transação',
    'transaction-details': 'Detalhes da Transação',
  }
  return titles[route.name as string] || 'Gestão Financeira'
})

function handleLogout() {
  showUserMenu.value = false
  authStore.logout()
  router.push('/login')
}
</script>
