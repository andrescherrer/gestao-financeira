<template>
  <aside
    class="fixed left-0 top-0 z-40 h-screen w-64 border-r border-gray-200 bg-white transition-transform"
    :class="{ '-translate-x-full': !isOpen, 'translate-x-0': isOpen }"
  >
    <div class="flex h-full flex-col">
      <!-- Logo e Nome -->
      <div class="border-b border-gray-200 px-6 py-5">
        <RouterLink to="/" class="flex items-center gap-3 transition-opacity hover:opacity-80">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-indigo-700 shadow-sm">
            <!-- Logo com ondas estilizado -->
            <svg
              class="h-6 w-6 text-white"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
          </div>
          <span class="text-lg font-bold text-gray-900">Gestão Financeira</span>
        </RouterLink>
      </div>

      <!-- Menu de Navegação -->
      <nav class="flex-1 overflow-y-auto px-4 py-6">
        <ul class="space-y-1">
          <li>
            <RouterLink
              to="/"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'home'
                  ? 'bg-gray-100 text-gray-900'
                  : 'text-gray-700 hover:bg-gray-50'
              "
            >
              <div class="flex items-center gap-3">
                <i class="pi pi-home text-base"></i>
                <span>Dashboard</span>
              </div>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/accounts"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'accounts' ||
                $route.name === 'new-account' ||
                $route.name === 'account-details' ||
                $route.name === 'edit-account'
                  ? 'bg-gray-100 text-gray-900'
                  : 'text-gray-700 hover:bg-gray-50'
              "
            >
              <div class="flex items-center gap-3">
                <i class="pi pi-wallet text-base"></i>
                <span>Contas</span>
              </div>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/transactions"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'transactions' ||
                $route.name === 'new-transaction' ||
                $route.name === 'transaction-details'
                  ? 'bg-gray-100 text-gray-900'
                  : 'text-gray-700 hover:bg-gray-50'
              "
            >
              <div class="flex items-center gap-3">
                <i class="pi pi-list text-base"></i>
                <span>Transações</span>
              </div>
            </RouterLink>
          </li>
        </ul>
      </nav>
    </div>
  </aside>
</template>

<script setup lang="ts">
defineProps<{
  isOpen?: boolean
}>()
</script>
