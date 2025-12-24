<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Contas' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Contas</h1>
          <p class="text-gray-600">
            Gerencie suas contas bancárias e financeiras
          </p>
        </div>
        <router-link
          to="/accounts/new"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
        >
          <i class="pi pi-plus"></i>
          Nova Conta
        </router-link>
      </div>

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <i class="pi pi-spinner pi-spin text-4xl text-blue-600 mb-4"></i>
          <p class="text-gray-600">Carregando contas...</p>
        </div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="accountsStore.error"
        class="rounded-md bg-red-50 border border-red-200 p-4 mb-6"
      >
        <div class="flex items-center gap-2">
          <i class="pi pi-exclamation-circle text-red-600"></i>
          <p class="text-red-600">{{ accountsStore.error }}</p>
        </div>
        <button
          @click="handleRetry"
          class="mt-3 text-sm text-red-600 hover:text-red-700 underline"
        >
          Tentar novamente
        </button>
      </div>

      <!-- Empty State -->
      <div
        v-else-if="accountsStore.accounts.length === 0"
        class="rounded-lg border border-gray-200 bg-white p-12 text-center"
      >
        <i class="pi pi-wallet text-6xl text-gray-400 mb-4"></i>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">
          Nenhuma conta encontrada
        </h3>
        <p class="text-gray-600 mb-6">
          Comece criando sua primeira conta financeira
        </p>
        <router-link
          to="/accounts/new"
          class="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
        >
          <i class="pi pi-plus"></i>
          Criar Primeira Conta
        </router-link>
      </div>

      <!-- Accounts List -->
      <div v-else class="space-y-4">
        <!-- Stats -->
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-blue-50 to-blue-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-blue-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500 text-white">
                  <i class="pi pi-wallet"></i>
                </div>
              </div>
              <div class="text-sm font-medium text-blue-700">Total de Contas</div>
              <div class="text-3xl font-bold text-blue-900">
                {{ accountsStore.totalAccounts }}
              </div>
            </div>
          </div>
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-purple-50 to-purple-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-purple-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-500 text-white">
                  <i class="pi pi-user"></i>
                </div>
              </div>
              <div class="text-sm font-medium text-purple-700">Contas Pessoais</div>
              <div class="text-3xl font-bold text-purple-900">
                {{ accountsStore.personalAccounts.length }}
              </div>
            </div>
          </div>
          <div class="group relative overflow-hidden rounded-xl bg-gradient-to-br from-indigo-50 to-indigo-100 p-6 shadow-sm transition-all hover:shadow-md">
            <div class="absolute right-0 top-0 -mr-4 -mt-4 h-24 w-24 rounded-full bg-indigo-200 opacity-20"></div>
            <div class="relative">
              <div class="mb-2 flex items-center gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-500 text-white">
                  <i class="pi pi-briefcase"></i>
                </div>
              </div>
              <div class="text-sm font-medium text-indigo-700">Contas de Negócio</div>
              <div class="text-3xl font-bold text-indigo-900">
                {{ accountsStore.businessAccounts.length }}
              </div>
            </div>
          </div>
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

