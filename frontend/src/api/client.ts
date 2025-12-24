import axios, { type AxiosInstance, type AxiosError } from 'axios'
import { env } from '@/config/env'

/**
 * Cliente HTTP configurado
 */
export const apiClient: AxiosInstance = axios.create({
  baseURL: env.apiUrl,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Interceptor para adicionar token JWT
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token && config.headers) {
      // Garantir que o token não está vazio
      const trimmedToken = token.trim()
      if (trimmedToken) {
        config.headers.Authorization = `Bearer ${trimmedToken}`
      }
    }
    return config
  },
  (error: AxiosError) => {
    return Promise.reject(error)
  }
)

// Interceptor para tratar respostas e erros
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    // Apenas rejeitar o erro - não fazer redirecionamento automático
    // O redirecionamento será tratado pelos componentes que capturam o erro
    return Promise.reject(error)
  }
)

export default apiClient

