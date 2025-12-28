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
    {
      path: '/budgets',
      name: 'budgets',
      component: () => import('@/views/BudgetsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import('@/views/ReportsView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Inicializar store se necessário
  if (!authStore.token) {
    authStore.init()
  }
  
  // Se estamos acessando rota protegida
  if (to.meta.requiresAuth) {
    const hasToken = !!localStorage.getItem('auth_token') || !!authStore.token
    
    // Se não há token, redirecionar imediatamente
    if (!hasToken) {
      authStore.logout()
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
    
    // Se há token, SEMPRE validar antes de permitir acesso
    // Isso garante que se o banco foi zerado, o token será invalidado
    try {
      const isValid = await authStore.validateToken()
      if (!isValid) {
        // Token inválido, limpar e redirecionar para login
        authStore.logout()
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }
      
      // Verificar se ainda está autenticado após validação
      if (!authStore.isAuthenticated) {
        authStore.logout()
        next({ name: 'login', query: { redirect: to.fullPath } })
        return
      }
    } catch (error) {
      // Em caso de erro na validação, considerar token inválido
      console.error('Erro ao validar token:', error)
      authStore.logout()
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
  }
  
  // Verificar autenticação usando a computed property que considera validação
  const isAuthenticated = authStore.isAuthenticated

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
