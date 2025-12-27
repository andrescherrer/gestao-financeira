import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/api/auth'
import type { LoginRequest, RegisterRequest, User } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isValidating = ref(false)
  const isValidated = ref(false)

  const isAuthenticated = computed(() => {
    // Só considerar autenticado se:
    // 1. Há token
    // 2. Token foi validado (isValidated = true)
    // 3. Token ainda está presente (não foi removido após validação)
    const hasToken = !!token.value || !!authService.getToken()
    return hasToken && isValidated.value && !!token.value
  })

  const isLoading = ref(false)

  // Inicializar do localStorage
  function init() {
    const storedToken = authService.getToken()
    if (storedToken) {
      token.value = storedToken
      // Tentar carregar dados do usuário do localStorage
      const storedUser = localStorage.getItem('auth_user')
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser)
        } catch (error) {
          console.error('Erro ao carregar dados do usuário do localStorage:', error)
          user.value = null
        }
      }
      // Resetar validação ao inicializar - será validado novamente
      isValidated.value = false
    } else {
      // Garantir que o estado está limpo se não há token
      token.value = null
      user.value = null
      isValidated.value = false
      // Limpar dados do usuário do localStorage também
      localStorage.removeItem('auth_user')
    }
  }

  // Validar token com o backend
  // IMPORTANTE: Esta função SEMPRE faz uma chamada HTTP real para validar o token.
  // Não confia em cache ou dados locais, pois o banco pode ter sido resetado.
  async function validateToken() {
    const storedToken = authService.getToken()
    if (!storedToken) {
      // Sem token, garantir que está limpo
      token.value = null
      user.value = null
      isValidated.value = true
      return false
    }

    // Evitar múltiplas validações simultâneas
    if (isValidating.value) {
      // Aguardar validação em andamento
      while (isValidating.value) {
        await new Promise(resolve => setTimeout(resolve, 50))
      }
      // Retornar baseado no estado atual após validação
      return isValidated.value && !!token.value && !!authService.getToken()
    }

    isValidating.value = true

    try {
      // CRÍTICO: SEMPRE fazer uma chamada HTTP real para validar o token.
      // Não podemos confiar em cache porque:
      // 1. O banco pode ter sido resetado (usuário não existe mais)
      // 2. O JWT pode ser válido mas o usuário foi deletado
      // 3. Dados em cache não refletem o estado real do backend
      //
      // Usamos /accounts porque é um endpoint protegido que:
      // - Verifica se o token é válido (middleware de auth)
      // - Verifica se o usuário existe no banco (busca accounts do user_id)
      // - Retorna 401/403 se o usuário não existir ou token for inválido
      const { accountService } = await import('@/api/accounts')
      
      // Fazer requisição HTTP real (sem cache)
      // Se o usuário não existir no banco, o backend retornará 401/403
      await accountService.list()
      
      // Se chegou aqui, o token é válido E o usuário existe no banco
      // Verificar se o token ainda existe (não foi removido durante a requisição)
      const currentToken = authService.getToken()
      if (!currentToken || currentToken !== storedToken) {
        // Token foi removido durante a validação (provavelmente pelo interceptor)
        logout()
        isValidated.value = true
        return false
      }
      
      // Token é válido e usuário existe no banco
      token.value = storedToken
      
      // Carregar dados do usuário do localStorage (se disponível)
      // Se não estiver no localStorage, os dados serão carregados no próximo login
      if (!user.value) {
        const storedUser = localStorage.getItem('auth_user')
        if (storedUser) {
          try {
            user.value = JSON.parse(storedUser)
          } catch (error) {
            console.error('Erro ao carregar dados do usuário do localStorage:', error)
            // Se os dados do localStorage estão corrompidos, limpar
            localStorage.removeItem('auth_user')
          }
        }
      }
      
      isValidated.value = true
      return true
    } catch (error: any) {
      // Se retornou 401 ou 403, o token é inválido ou o usuário não existe mais no banco
      if (error.response?.status === 401 || error.response?.status === 403) {
        // Limpar tudo: token, dados do usuário, localStorage
        console.warn('[Auth] Token inválido ou usuário não existe mais no banco:', error.response?.status)
        logout()
        isValidated.value = true
        return false
      }
      
      // Para outros erros (500, network, etc), por segurança também limpamos
      // Pois não podemos garantir que o token é válido sem uma resposta bem-sucedida do backend
      console.error('[Auth] Erro ao validar token:', error)
      logout()
      isValidated.value = true
      return false
    } finally {
      isValidating.value = false
    }
  }

  async function login(credentials: LoginRequest) {
    isLoading.value = true
    try {
      const response = await authService.login(credentials)
      // Salvar token no localStorage PRIMEIRO
      authService.saveToken(response.token)
      // Salvar dados do usuário no localStorage
      localStorage.setItem('auth_user', JSON.stringify(response.user))
      // Aguardar um momento para garantir que foi salvo
      await new Promise(resolve => setTimeout(resolve, 10))
      // Depois atualizar o estado reativo
      token.value = response.token
      user.value = response.user
      // Marcar como validado após login bem-sucedido
      isValidated.value = true
      return response
    } catch (error) {
      // Limpar token em caso de erro
      token.value = null
      user.value = null
      isValidated.value = false
      authService.removeToken()
      localStorage.removeItem('auth_user')
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function register(userData: RegisterRequest) {
    isLoading.value = true
    try {
      return await authService.register(userData)
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    // Remover token primeiro
    authService.removeToken()
    // Remover dados do usuário do localStorage
    localStorage.removeItem('auth_user')
    // Limpar estado reativo
    token.value = null
    user.value = null
    isValidated.value = false
    
    // Limpar dados de todas as stores para evitar dados obsoletos em cache
    // Isso é crítico quando o banco é resetado mas o frontend ainda tem cache
    Promise.all([
      import('@/stores/accounts').then(m => {
        const store = m.useAccountsStore()
        store.accounts = []
        store.currentAccount = null
        store.error = null
      }),
      import('@/stores/transactions').then(m => {
        const store = m.useTransactionsStore()
        store.transactions = []
        store.currentTransaction = null
        store.error = null
      }),
      import('@/stores/categories').then(m => {
        const store = m.useCategoriesStore()
        store.categories = []
        store.currentCategory = null
        store.error = null
      }),
    ]).catch(err => {
      // Log mas não falha se alguma store não existir
      if (import.meta.env.DEV) {
        console.warn('[Auth] Erro ao limpar stores:', err)
      }
    })
  }

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    isValidating,
    isValidated,
    init,
    validateToken,
    login,
    register,
    logout,
  }
})

