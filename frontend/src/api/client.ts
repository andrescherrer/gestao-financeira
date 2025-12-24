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
        // Debug: log apenas em desenvolvimento
        if (import.meta.env.DEV) {
          console.log('[API Client] Token adicionado ao header:', {
            url: config.url,
            hasToken: !!trimmedToken,
            tokenLength: trimmedToken.length,
          })
        }
      } else {
        console.warn('[API Client] Token encontrado mas está vazio após trim')
      }
    } else {
      // Debug: log apenas em desenvolvimento
      if (import.meta.env.DEV && config.url && !config.url.includes('/auth/')) {
        console.warn('[API Client] Requisição sem token:', {
          url: config.url,
          hasToken: !!token,
        })
      }
    }

    // Log do body em desenvolvimento para POST/PUT/PATCH
    if (import.meta.env.DEV && config.data && (config.method === 'post' || config.method === 'put' || config.method === 'patch')) {
      console.log('[API Client] Request body:', {
        url: config.url,
        method: config.method?.toUpperCase(),
        data: config.data,
        dataString: JSON.stringify(config.data, null, 2),
      })
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

