import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import(/* webpackChunkName: "dashboard" */ '@/views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import(/* webpackChunkName: "auth" */ '@/views/LoginView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import(/* webpackChunkName: "auth" */ '@/views/RegisterView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/accounts',
      name: 'accounts',
      component: () => import(/* webpackChunkName: "accounts" */ '@/views/AccountsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/new',
      name: 'new-account',
      component: () => import(/* webpackChunkName: "accounts" */ '@/views/NewAccountView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/:id',
      name: 'account-details',
      component: () => import(/* webpackChunkName: "accounts" */ '@/views/AccountDetailsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/accounts/:id/edit',
      name: 'edit-account',
      component: () => import(/* webpackChunkName: "accounts" */ '@/views/EditAccountView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import(/* webpackChunkName: "transactions" */ '@/views/TransactionsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/new',
      name: 'new-transaction',
      component: () => import(/* webpackChunkName: "transactions" */ '@/views/NewTransactionView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/transactions/:id',
      name: 'transaction-details',
      component: () => import(/* webpackChunkName: "transactions" */ '@/views/TransactionDetailsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import(/* webpackChunkName: "categories" */ '@/views/CategoriesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/new',
      name: 'new-category',
      component: () => import(/* webpackChunkName: "categories" */ '@/views/NewCategoryView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/categories/:id',
      name: 'category-details',
      component: () => import(/* webpackChunkName: "categories" */ '@/views/CategoryDetailsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/budgets',
      name: 'budgets',
      component: () => import(/* webpackChunkName: "budgets" */ '@/views/BudgetsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import(/* webpackChunkName: "reports" */ '@/views/ReportsView.vue'),
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
