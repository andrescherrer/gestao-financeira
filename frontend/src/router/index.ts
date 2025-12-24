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
  
  // Sempre inicializar para garantir que o token está sincronizado
  authStore.init()

  // Verificar autenticação baseado no token no localStorage (fonte da verdade)
  const hasToken = !!authStore.token || !!localStorage.getItem('auth_token')

  if (to.meta.requiresAuth && !hasToken) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresGuest && hasToken) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
