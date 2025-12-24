<template>
  <header class="sticky top-0 z-50 h-16 border-b border-border bg-background">
    <div class="flex h-full items-center justify-between px-6">
      <!-- Left: Breadcrumbs ou título da página -->
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <span class="font-medium text-foreground">{{ pageTitle }}</span>
      </div>

      <!-- Right: Notificações e Usuário -->
      <div class="flex items-center gap-4">
        <!-- Ícone de Notificação -->
        <Button
          variant="ghost"
          size="icon"
          class="relative"
          title="Notificações"
        >
          <Bell class="h-5 w-5" />
          <span
            v-if="hasNotifications"
            class="absolute right-2 top-2 h-2 w-2 rounded-full bg-destructive"
          ></span>
        </Button>

        <!-- Avatar do Usuário -->
        <div v-if="authStore.user" class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 text-sm font-semibold text-white shadow-sm">
            {{ userInitials }}
          </div>
          <div class="hidden flex-col sm:flex">
            <span class="text-sm font-semibold text-foreground">{{ userName }}</span>
            <span class="text-xs text-muted-foreground">{{ authStore.user.email }}</span>
          </div>
          <!-- Menu Dropdown do Usuário -->
          <div class="relative" ref="userMenuRef">
            <Button
              variant="ghost"
              size="icon"
              @click="showUserMenu = !showUserMenu"
            >
              <ChevronDown class="h-4 w-4" />
            </Button>
            <div
              v-if="showUserMenu"
              class="absolute right-0 mt-2 w-48 rounded-lg border border-border bg-popover p-1 shadow-lg"
            >
              <Button
                variant="ghost"
                class="flex w-full items-center gap-2 justify-start"
                @click="handleLogout"
              >
                <LogOut class="h-4 w-4" />
                <span>Sair</span>
              </Button>
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
import { Button } from '@/components/ui/button'
import { Bell, ChevronDown, LogOut } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const showUserMenu = ref(false)
const hasNotifications = ref(false)
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
