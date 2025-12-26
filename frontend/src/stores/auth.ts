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
    // Sempre verificar o localStorage como fonte da verdade
    return !!token.value || !!authService.getToken()
  })

  const isLoading = ref(false)

  // Inicializar do localStorage
  function init() {
    const storedToken = authService.getToken()
    if (storedToken) {
      token.value = storedToken
    } else {
      // Garantir que o estado está limpo se não há token
      token.value = null
      user.value = null
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
      return isValidated.value && !!token.value
    }

    isValidating.value = true

    try {
      // Tentar fazer uma requisição protegida para validar o token
      // Usamos /accounts porque é um endpoint simples e protegido
      const { accountService } = await import('@/api/accounts')
      await accountService.list()
      
      // Se chegou aqui, o token é válido
      token.value = storedToken
      isValidated.value = true
      // Nota: não temos endpoint /me, então não podemos carregar user aqui
      // Mas o token é válido, então mantemos o estado
      return true
    } catch (error: any) {
      // Se retornou 401, o token é inválido
      if (error.response?.status === 401) {
        // Limpar token inválido
        logout()
        isValidated.value = true
        return false
      }
      // Outros erros podem ser temporários, mas por segurança, limpamos
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
      // Aguardar um momento para garantir que foi salvo
      await new Promise(resolve => setTimeout(resolve, 10))
      // Depois atualizar o estado reativo
      token.value = response.token
      user.value = response.user
      return response
    } catch (error) {
      // Limpar token em caso de erro
      token.value = null
      user.value = null
      authService.removeToken()
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

