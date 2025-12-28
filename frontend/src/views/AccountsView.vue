<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Contas' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl sm:text-4xl font-bold mb-2">Contas</h1>
          <p class="text-sm sm:text-base text-muted-foreground">
            Gerencie suas contas bancárias e financeiras
          </p>
        </div>
        <Button as-child class="w-full sm:w-auto">
          <router-link to="/accounts/new">
            <Plus class="h-4 w-4 mr-2" />
            Nova Conta
          </router-link>
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando contas...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="accountsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ accountsStore.error }}</p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="link"
              @click="handleRetry"
              class="text-destructive"
            >
              Tentar novamente
            </Button>
            <Button
              v-if="isAuthError"
              variant="outline"
              @click="handleLogin"
              class="text-destructive border-destructive"
            >
              Fazer login novamente
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <EmptyState
        v-else-if="accountsStore.accounts.length === 0"
        :icon="Wallet"
        title="Nenhuma conta encontrada"
        description="Comece criando sua primeira conta financeira"
        action-label="Criar Primeira Conta"
        action-to="/accounts/new"
        :action-icon="Plus"
      />

      <!-- Accounts List -->
      <div v-else class="space-y-4">
        <!-- Stats -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-blue-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500 text-white">
                    <Wallet class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-blue-700">Total de Contas</div>
                <div class="text-3xl font-bold text-blue-900">
                  {{ accountsStore.totalAccounts }}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-purple-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-500 text-white">
                    <User class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-purple-700">Contas Pessoais</div>
                <div class="text-3xl font-bold text-purple-900">
                  {{ accountsStore.personalAccounts.length }}
                </div>
              </div>
            </CardContent>
          </Card>
          <Card class="group relative overflow-hidden">
            <CardContent class="p-6">
              <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-indigo-200 opacity-20"></div>
              <div class="relative">
                <div class="mb-2 flex items-center gap-3">
                  <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-500 text-white">
                    <Briefcase class="h-5 w-5" />
                  </div>
                </div>
                <div class="text-sm font-medium text-indigo-700">Contas de Negócio</div>
                <div class="text-3xl font-bold text-indigo-900">
                  {{ accountsStore.businessAccounts.length }}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Accounts Grid -->
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          <AccountCard
            v-for="account in accountsStore.accounts"
            :key="account.account_id"
            :account="account"
          />
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import AccountCard from '@/components/AccountCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Wallet, User, Briefcase, Plus, Loader2, AlertCircle } from 'lucide-vue-next'

const router = useRouter()
const accountsStore = useAccountsStore()
const authStore = useAuthStore()

const isAuthError = computed(() => {
  const error = accountsStore.error?.toLowerCase() || ''
  return error.includes('token') || 
         error.includes('autenticação') || 
         error.includes('unauthorized') ||
         error.includes('invalid') ||
         error.includes('expired')
})

onMounted(async () => {
  // Garantir que o token está inicializado
  if (!authStore.token) {
    authStore.init()
  }
  
  // Aguardar um pouco para garantir que validateToken terminou (se ainda estiver rodando)
  if (authStore.isValidating) {
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  
  // Só carregar se não tiver dados E não estiver carregando
  // validateToken pode já ter carregado os dados através da store
  if (accountsStore.accounts.length === 0 && !accountsStore.isLoading) {
    await accountsStore.listAccounts()
  }
})

function handleRetry() {
  accountsStore.clearError()
  accountsStore.listAccounts()
}

function handleLogin() {
  authStore.logout()
  router.push({ name: 'login', query: { redirect: '/accounts' } })
}
</script>
