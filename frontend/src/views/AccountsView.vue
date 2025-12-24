<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Contas' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Contas</h1>
          <p class="text-muted-foreground">
            Gerencie suas contas bancárias e financeiras
          </p>
        </div>
        <Button as-child>
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
          <Button
            variant="link"
            @click="handleRetry"
            class="text-destructive"
          >
            Tentar novamente
          </Button>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <Card
        v-else-if="accountsStore.accounts.length === 0"
      >
        <CardContent class="p-12 text-center">
          <Wallet class="mx-auto h-16 w-16 text-muted-foreground mb-4" />
          <h3 class="text-xl font-semibold text-foreground mb-2">
            Nenhuma conta encontrada
          </h3>
          <p class="text-muted-foreground mb-6">
            Comece criando sua primeira conta financeira
          </p>
          <Button as-child>
            <router-link to="/accounts/new">
              <Plus class="h-4 w-4 mr-2" />
              Criar Primeira Conta
            </router-link>
          </Button>
        </CardContent>
      </Card>

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
import { onMounted } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import AccountCard from '@/components/AccountCard.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Wallet, User, Briefcase, Plus, Loader2, AlertCircle } from 'lucide-vue-next'

const accountsStore = useAccountsStore()

onMounted(async () => {
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
})

function handleRetry() {
  accountsStore.clearError()
  accountsStore.listAccounts()
}
</script>
