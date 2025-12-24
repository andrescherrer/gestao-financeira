<template>
  <aside
    class="fixed left-0 top-16 z-40 h-[calc(100vh-4rem)] w-64 border-r bg-white transition-transform md:translate-x-0"
    :class="{ '-translate-x-full': !isOpen, 'translate-x-0': isOpen }"
  >
    <nav class="flex h-full flex-col p-4">
      <ul class="space-y-2">
        <li>
          <RouterLink
            to="/"
            class="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors"
            :class="
              $route.name === 'home'
                ? 'bg-blue-50 text-blue-600'
                : 'text-gray-700 hover:bg-gray-100'
            "
          >
            <i class="pi pi-home"></i>
            <span>Dashboard</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink
            to="/accounts"
            class="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors"
            :class="
              $route.name === 'accounts' ||
              $route.name === 'new-account' ||
              $route.name === 'account-details'
                ? 'bg-blue-50 text-blue-600'
                : 'text-gray-700 hover:bg-gray-100'
            "
          >
            <i class="pi pi-wallet"></i>
            <span>Contas</span>
          </RouterLink>
        </li>
        <li>
          <RouterLink
            to="/transactions"
            class="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors"
            :class="
              $route.name === 'transactions' ||
              $route.name === 'new-transaction' ||
              $route.name === 'transaction-details'
                ? 'bg-blue-50 text-blue-600'
                : 'text-gray-700 hover:bg-gray-100'
            "
          >
            <i class="pi pi-list"></i>
            <span>Transações</span>
          </RouterLink>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup lang="ts">
defineProps<{
  isOpen?: boolean
}>()
</script>

