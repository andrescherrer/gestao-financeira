import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: () => import('@/views/AccountsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/new',
      name: 'new-account',
      component: () => import('@/views/NewAccountView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/:id',
      name: 'account-details',
      component: () => import('@/views/AccountDetailsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/:id/edit',
      name: 'edit-account',
      component: () => import('@/views/EditAccountView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('@/views/TransactionsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/new',
      name: 'new-transaction',
      component: () => import('@/views/NewTransactionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/:id',
      name: 'transaction-details',
      component: () => import('@/views/TransactionDetailsView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Verificar autenticação baseado no token no localStorage e na store
  const hasTokenInStorage = !!localStorage.getItem('auth_token')
  const hasTokenInStore = !!authStore.token
  
  // Se não há token na store mas há no storage, inicializar
  if (!hasTokenInStore && hasTokenInStorage) {
    authStore.init()
  }
  
  // Verificar autenticação (usar storage como fonte da verdade)
  const isAuthenticated = hasTokenInStorage || hasTokenInStore

  if (to.meta.requiresAuth && !isAuthenticated) {
    // Limpar qualquer token residual
    authStore.logout()
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresGuest && isAuthenticated) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
