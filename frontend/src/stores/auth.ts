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
      // Tentar reutilizar a chamada de accounts se já estiver em andamento ou concluída
      // Isso evita múltiplas chamadas HTTP desnecessárias
      const { useAccountsStore } = await import('@/stores/accounts')
      const accountsStore = useAccountsStore()
      
      // Se accounts já está carregando, aguardar e usar o resultado
      if (accountsStore.isLoading) {
        // Aguardar até que a requisição de accounts termine
        while (accountsStore.isLoading) {
          await new Promise(resolve => setTimeout(resolve, 50))
        }
        // Se não houve erro, o token é válido
        if (!accountsStore.error) {
          const currentToken = authService.getToken()
          if (currentToken && currentToken === storedToken) {
            token.value = storedToken
            if (!user.value) {
              const storedUser = localStorage.getItem('auth_user')
              if (storedUser) {
                try {
                  user.value = JSON.parse(storedUser)
                } catch (error) {
                  console.error('Erro ao carregar dados do usuário do localStorage:', error)
                }
              }
            }
            isValidated.value = true
            isValidating.value = false
            return true
          }
        }
      }
      
      // Se accounts já tem dados, o token é válido (não precisa fazer nova requisição)
      if (accountsStore.accounts.length > 0) {
        const currentToken = authService.getToken()
        if (currentToken && currentToken === storedToken) {
          token.value = storedToken
          if (!user.value) {
            const storedUser = localStorage.getItem('auth_user')
            if (storedUser) {
              try {
                user.value = JSON.parse(storedUser)
              } catch (error) {
                console.error('Erro ao carregar dados do usuário do localStorage:', error)
              }
            }
          }
          isValidated.value = true
          isValidating.value = false
          return true
        }
      }
      
      // Se não podemos reutilizar, fazer uma requisição leve para validar
      // Usamos /accounts porque é um endpoint simples e protegido
      const { accountService } = await import('@/api/accounts')
      await accountService.list()
      
      // Se chegou aqui, o token é válido
      // Verificar se o token ainda existe (não foi removido durante a requisição)
      const currentToken = authService.getToken()
      if (!currentToken || currentToken !== storedToken) {
        // Token foi removido durante a validação
        logout()
        isValidated.value = true
        return false
      }
      
      token.value = storedToken
      // Garantir que os dados do usuário estão carregados
      // Se não estão na store, tentar carregar do localStorage
      if (!user.value) {
        const storedUser = localStorage.getItem('auth_user')
        if (storedUser) {
          try {
            user.value = JSON.parse(storedUser)
          } catch (error) {
            console.error('Erro ao carregar dados do usuário do localStorage:', error)
          }
        }
      }
      isValidated.value = true
      return true
    } catch (error: any) {
      // Se retornou 401 ou 403, o token é inválido ou expirado
      if (error.response?.status === 401 || error.response?.status === 403) {
        // Limpar token inválido
        logout()
        isValidated.value = true
        return false
      }
      // Para outros erros (500, network, etc), por segurança também limpamos
      // Pois não podemos garantir que o token é válido
      console.error('Erro ao validar token:', error)
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

