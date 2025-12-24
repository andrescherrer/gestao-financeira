import apiClient from './client'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from './types'

/**
 * Serviço de autenticação
 */
export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials)
    return response.data
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
    localStorage.removeItem('auth_token')
    document.cookie = 'auth_token=; path=/; max-age=0'
  },

  getToken(): string | null {
    return localStorage.getItem('auth_token')
  },

  isAuthenticated(): boolean {
    return this.getToken() !== null
  },
}

