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
  async (error: AxiosError) => {
    // Log detalhado do erro em desenvolvimento
    if (import.meta.env.DEV) {
      const errorData = error.response?.data
      console.error('[API Client] Erro na requisição:', {
        url: error.config?.url,
        method: error.config?.method?.toUpperCase(),
        status: error.response?.status,
        statusText: error.response?.statusText,
        message: error.message,
        code: error.code,
        isNetworkError: !error.response, // Erro de rede (sem resposta do servidor)
        isTimeout: error.code === 'ECONNABORTED',
        errorType: (errorData as any)?.error_type,
        requestId: (errorData as any)?.request_id,
        errorDetails: (errorData as any)?.details,
      })
    }

    // Tratar erros de autenticação (401 ou 403)
    if (error.response?.status === 401 || error.response?.status === 403) {
      // Importar dinamicamente para evitar dependência circular
      const { authService } = await import('./auth')
      const { useAuthStore } = await import('@/stores/auth')
      
      // Remover token inválido
      authService.removeToken()
      
      // Limpar estado de autenticação
      const authStore = useAuthStore()
      authStore.logout()
      
      // Redirecionar para login apenas se não estiver já na página de login/register
      if (typeof window !== 'undefined') {
        const currentPath = window.location.pathname
        if (!currentPath.includes('/login') && !currentPath.includes('/register')) {
          // Usar window.location para garantir redirecionamento
          // O router guard vai tratar a navegação depois
          const redirectUrl = `/login?redirect=${encodeURIComponent(currentPath)}`
          // Usar setTimeout para garantir que o logout foi processado
          setTimeout(() => {
            window.location.href = redirectUrl
          }, 100)
        }
      }
    }
    
    // Tratar erros de rede (sem resposta do servidor)
    // Isso inclui: Failed to fetch, Network Error, timeout, etc.
    if (!error.response) {
      // Log do erro de rede
      if (import.meta.env.DEV) {
        console.error('[API Client] Erro de rede:', {
          url: error.config?.url,
          message: error.message,
          code: error.code,
          possibleCauses: [
            'Backend não está rodando',
            'Problema de conectividade',
            'CORS bloqueando requisição',
            'Timeout da requisição',
          ],
        })
      }
      
      // Se for um erro de rede e estivermos em uma rota protegida,
      // não limpar o token automaticamente (pode ser problema temporário)
      // O validateToken vai tratar isso adequadamente
    }
    
    return Promise.reject(error)
  }
)

export default apiClient

