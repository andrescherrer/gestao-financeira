<template>
  <aside
    class="fixed left-0 top-0 z-40 h-screen w-64 border-r border-border bg-background transition-transform"
    :class="{ '-translate-x-full': !isOpen, 'translate-x-0': isOpen }"
  >
    <div class="flex h-full flex-col">
      <!-- Logo e Nome -->
      <div class="border-b border-border px-6 py-5">
        <RouterLink to="/" class="flex items-center gap-3 transition-opacity hover:opacity-80">
          <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-indigo-700 shadow-sm">
            <Wallet class="h-6 w-6 text-white" />
          </div>
          <span class="text-lg font-bold text-foreground">Gestão Financeira</span>
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
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <Home class="h-4 w-4" />
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
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <Wallet class="h-4 w-4" />
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
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <List class="h-4 w-4" />
                <span>Transações</span>
              </div>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/categories"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'categories' ||
                $route.name === 'new-category' ||
                $route.name === 'category-details'
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <Tag class="h-4 w-4" />
                <span>Categorias</span>
              </div>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/budgets"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'budgets' ||
                $route.name === 'new-budget' ||
                $route.name === 'budget-details'
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <Target class="h-4 w-4" />
                <span>Orçamentos</span>
              </div>
            </RouterLink>
          </li>
          <li>
            <RouterLink
              to="/reports"
              class="group flex items-center justify-between rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="
                $route.name === 'reports'
                  ? 'bg-accent text-accent-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              "
            >
              <div class="flex items-center gap-3">
                <BarChart class="h-4 w-4" />
                <span>Relatórios</span>
              </div>
            </RouterLink>
          </li>
        </ul>
      </nav>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { Home, Wallet, List, Tag, Target, BarChart } from 'lucide-vue-next'

defineProps<{
  isOpen?: boolean
}>()
</script>
