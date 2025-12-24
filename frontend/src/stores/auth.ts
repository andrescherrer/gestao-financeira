import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/api/auth'
import type { LoginRequest, RegisterRequest, User } from '@/api/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  const isAuthenticated = computed(() => {
    return !!token.value || authService.isAuthenticated()
  })

  const isLoading = ref(false)

  // Inicializar do localStorage
  function init() {
    const storedToken = authService.getToken()
    if (storedToken) {
      token.value = storedToken
    }
  }

  async function login(credentials: LoginRequest) {
    isLoading.value = true
    try {
      const response = await authService.login(credentials)
      // Salvar token primeiro
      authService.saveToken(response.token)
      // Depois atualizar o estado
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
    token.value = null
    user.value = null
    authService.removeToken()
  }

  return {
    user,
    token,
    isAuthenticated,
    isLoading,
    init,
    login,
    register,
    logout,
  }
})

