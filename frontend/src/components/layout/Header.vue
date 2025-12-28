<template>
  <header class="sticky top-0 z-50 h-16 border-b border-border bg-background">
    <div class="flex h-full items-center justify-between px-4 sm:px-6">
      <!-- Left: Menu Mobile e Título -->
      <div class="flex items-center gap-2 sm:gap-4">
        <!-- Botão Menu Mobile -->
        <Button
          variant="ghost"
          size="icon"
          class="md:hidden"
          @click="$emit('toggle-sidebar')"
          :aria-label="sidebarOpen ? 'Fechar menu' : 'Abrir menu'"
          :aria-expanded="sidebarOpen"
          aria-controls="sidebar"
        >
          <Menu v-if="!sidebarOpen" class="h-5 w-5" aria-hidden="true" />
          <X v-else class="h-5 w-5" aria-hidden="true" />
        </Button>
        <span class="font-medium text-foreground text-sm sm:text-base">{{ pageTitle }}</span>
      </div>

      <!-- Right: Notificações e Usuário -->
      <div class="flex items-center gap-2 sm:gap-4">
        <!-- Toggle de Tema -->
        <Button
          variant="ghost"
          size="icon"
          @click="toggleTheme"
          :title="isDark ? 'Alternar para modo claro' : 'Alternar para modo escuro'"
          :aria-label="isDark ? 'Alternar para modo claro' : 'Alternar para modo escuro'"
          aria-pressed="false"
        >
          <Sun v-if="isDark" class="h-5 w-5" aria-hidden="true" />
          <Moon v-else class="h-5 w-5" aria-hidden="true" />
        </Button>

        <!-- Ícone de Notificação -->
        <Button
          variant="ghost"
          size="icon"
          class="relative"
          title="Notificações"
          aria-label="Notificações"
          :aria-describedby="hasNotifications ? 'notification-badge' : undefined"
        >
          <Bell class="h-5 w-5" aria-hidden="true" />
          <span
            v-if="hasNotifications"
            id="notification-badge"
            class="absolute right-2 top-2 h-2 w-2 rounded-full bg-destructive"
            aria-label="Você tem notificações não lidas"
          ></span>
        </Button>

        <!-- Avatar do Usuário com Dropdown Menu -->
        <!-- Só mostrar se realmente estiver autenticado E validado -->
        <DropdownMenu v-if="authStore.isAuthenticated && authStore.isValidated && !authStore.isValidating && authStore.token">
          <DropdownMenuTrigger as-child>
            <Button
              variant="ghost"
              class="flex items-center gap-2 h-auto p-1 hover:bg-accent"
            >
              <div class="flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 text-sm font-semibold text-white shadow-sm">
                {{ userInitials }}
              </div>
              <div class="hidden flex-col items-start sm:flex">
                <span class="text-sm font-semibold text-foreground">{{ userName }}</span>
                <span class="text-xs text-muted-foreground">
                  {{ authStore.user?.email || 'Usuário autenticado' }}
                </span>
              </div>
              <ChevronDown class="h-4 w-4 text-muted-foreground" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-56">
            <DropdownMenuLabel>
              <div class="flex flex-col space-y-1">
                <p class="text-sm font-medium leading-none">{{ userName }}</p>
                <p class="text-xs leading-none text-muted-foreground">
                  {{ authStore.user?.email || 'Usuário autenticado' }}
                </p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleProfile">
              <User class="h-4 w-4" />
              <span>Perfil</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="handleSettings">
              <Settings class="h-4 w-4" />
              <span>Configurações</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleLogout" class="text-destructive focus:text-destructive">
              <LogOut class="h-4 w-4" />
              <span>Sair</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Bell, ChevronDown, LogOut, User, Settings, Sun, Moon, Menu, X } from 'lucide-vue-next'
import { useTheme } from '@/composables/useTheme'

interface Props {
  sidebarOpen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  sidebarOpen: false,
})

defineEmits<{
  'toggle-sidebar': []
}>()

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { isDark, toggleTheme } = useTheme()

const hasNotifications = ref(false)

const userName = computed(() => {
  if (!authStore.user) {
    // Se não há user mas há token, mostrar placeholder
    return authStore.token ? 'Usuário' : 'Usuário'
  }
  const firstName = authStore.user.first_name || ''
  const lastName = authStore.user.last_name || ''
  return `${firstName} ${lastName}`.trim() || authStore.user.email
})

const userInitials = computed(() => {
  if (!authStore.user) {
    // Se não há user mas há token, mostrar 'U'
    return 'U'
  }
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

async function handleLogout() {
  try {
    // Limpar estado da store primeiro
    authStore.logout()
    // Aguardar um momento para garantir que o localStorage foi limpo
    await new Promise(resolve => setTimeout(resolve, 100))
    // Redirecionar para login
    await router.push('/login')
    // Forçar reload para garantir que tudo está limpo
    window.location.reload()
  } catch (error) {
    console.error('Erro ao fazer logout:', error)
    // Mesmo em caso de erro, tentar redirecionar
    window.location.href = '/login'
  }
}

function handleProfile() {
  // TODO: Implementar página de perfil
  console.log('Perfil')
}

function handleSettings() {
  // TODO: Implementar página de configurações
  console.log('Configurações')
}
</script>
