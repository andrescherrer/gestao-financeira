import apiClient from './client'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from './types'

/**
 * Serviço de autenticação
 */
export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<{
      message: string
      data: {
        token: string
        user_id: string
        email: string
        first_name: string
        last_name: string
        full_name: string
        expires_in: number
      }
    }>('/auth/login', credentials)
    
    // Mapear resposta do backend para o formato esperado pelo frontend
    const backendData = response.data.data
    return {
      token: backendData.token,
      user: {
        user_id: backendData.user_id,
        email: backendData.email,
        first_name: backendData.first_name,
        last_name: backendData.last_name,
        full_name: backendData.full_name,
      },
    }
  },

  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    const response = await apiClient.post<RegisterResponse>('/auth/register', userData)
    return response.data
  },

  saveToken(token: string): void {
    localStorage.setItem('auth_token', token)
    const maxAge = 60 * 60 * 24 * 7 // 7 dias
    document.cookie = `auth_token=${token}; path=/; max-age=${maxAge}; SameSite=Lax; Secure=${window.location.protocol === 'https:'}`
  },

  removeToken(): void {
    // Remover do localStorage
    localStorage.removeItem('auth_token')
    // Remover cookie (múltiplas tentativas para garantir)
    document.cookie = 'auth_token=; path=/; max-age=0; expires=Thu, 01 Jan 1970 00:00:00 GMT'
    document.cookie = 'auth_token=; path=/; max-age=0'
    // Limpar sessionStorage também (caso esteja sendo usado)
    sessionStorage.removeItem('auth_token')
  },

  getToken(): string | null {
    return localStorage.getItem('auth_token')
  },

  isAuthenticated(): boolean {
    return this.getToken() !== null
  },
}

