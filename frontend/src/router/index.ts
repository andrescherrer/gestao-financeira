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
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/views/CategoriesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/new',
      name: 'new-category',
      component: () => import('@/views/NewCategoryView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/:id',
      name: 'category-details',
      component: () => import('@/views/CategoryDetailsView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Verificar autenticação baseado no token no localStorage e na store
  const hasTokenInStorage = !!localStorage.getItem('auth_token')
  const hasTokenInStore = !!authStore.token
  
  // Se não há token na store mas há no storage, inicializar
  if (!hasTokenInStore && hasTokenInStorage) {
    authStore.init()
  }
  
  // Se há token e estamos acessando rota protegida, validar token
  if (to.meta.requiresAuth && (hasTokenInStorage || hasTokenInStore)) {
    // Se ainda não foi validado, validar agora
    if (!authStore.isValidated) {
      const isValid = await authStore.validateToken()
      if (!isValid) {
        // Token inválido, redirecionar para login
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }
    }
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
