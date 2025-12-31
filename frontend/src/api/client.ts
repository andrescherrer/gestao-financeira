import axios, { type AxiosInstance, type AxiosError } from 'axios'
import { env } from '@/config/env'
import { logger } from '@/utils/logger'

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
        // Log token adicionado
        logger.debug('Token adicionado ao header', {
          url: config.url,
          hasToken: !!trimmedToken,
          tokenLength: trimmedToken.length,
        })
      } else {
        logger.warn('Token encontrado mas está vazio após trim')
      }
    } else {
      // Log requisição sem token (apenas para rotas protegidas)
      if (config.url && !config.url.includes('/auth/')) {
        logger.debug('Requisição sem token', {
          url: config.url,
          hasToken: !!token,
        })
      }
    }

    // Log do body para POST/PUT/PATCH
    if (config.data && (config.method === 'post' || config.method === 'put' || config.method === 'patch')) {
      logger.debug('Request body', {
        url: config.url,
        method: config.method?.toUpperCase(),
        data: config.data,
      })
    }

    // Capturar correlation IDs dos headers de resposta anteriores
    const requestId = config.headers?.['X-Request-ID'] as string
    const traceId = config.headers?.['X-Trace-ID'] as string
    const spanId = config.headers?.['X-Span-ID'] as string
    
    if (requestId || traceId || spanId) {
      logger.setCorrelationIds(requestId, traceId, spanId)
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
    // Log detalhado do erro
    const errorData = error.response?.data as any
    const isNetworkError = !error.response
    const isTimeout = error.code === 'ECONNABORTED'

    // Capturar correlation IDs da resposta
    if (error.response) {
      const requestId = error.response.headers['x-request-id']
      const traceId = error.response.headers['x-trace-id']
      const spanId = error.response.headers['x-span-id']
      if (requestId || traceId || spanId) {
        logger.setCorrelationIds(requestId, traceId, spanId)
      }
    }

    logger.error('Erro na requisição', error, {
      url: error.config?.url,
      method: error.config?.method?.toUpperCase(),
      status: error.response?.status,
      statusText: error.response?.statusText,
      code: error.code,
      isNetworkError,
      isTimeout,
      errorType: errorData?.error_type,
      requestId: errorData?.request_id,
      errorDetails: errorData?.details,
    })

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
      logger.error('Erro de rede', error, {
        url: error.config?.url,
        code: error.code,
        possibleCauses: [
          'Backend não está rodando',
          'Problema de conectividade',
          'CORS bloqueando requisição',
          'Timeout da requisição',
        ],
      })
      
      // Se for um erro de rede e estivermos em uma rota protegida,
      // não limpar o token automaticamente (pode ser problema temporário)
      // O validateToken vai tratar isso adequadamente
    }
    
    return Promise.reject(error)
  }
)

export default apiClient

