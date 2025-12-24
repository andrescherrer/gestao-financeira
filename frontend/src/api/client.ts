import axios, { type AxiosInstance, type AxiosError } from 'axios'
import { env } from '@/config/env'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'

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
      config.headers.Authorization = `Bearer ${token}`
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
    if (error.response) {
      const status = error.response.status

      // 401 Unauthorized - token inválido ou expirado
      if (status === 401) {
        // Ignorar 401 em rotas de autenticação (login/register)
        const currentPath = window.location.pathname
        if (currentPath === '/login' || currentPath === '/register') {
          return Promise.reject(error)
        }

        const authStore = useAuthStore()
        authStore.logout()
        
        // Usar setTimeout para evitar problemas de importação circular
        setTimeout(() => {
          if (router.currentRoute.value.path !== '/login') {
            router.push({ name: 'login', query: { redirect: router.currentRoute.value.fullPath } })
          }
        }, 0)
      }
    }

    return Promise.reject(error)
  }
)

export default apiClient

