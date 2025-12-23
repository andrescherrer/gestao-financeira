import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { env } from '@/lib/config/env';

// Base URL da API - obtida das variáveis de ambiente
const API_BASE_URL = env.apiUrl;

// Criar instância do Axios
export const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // 30 segundos
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para adicionar token JWT nas requisições
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Buscar token do localStorage (ou outro storage)
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Interceptor para tratar respostas e erros
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error: AxiosError) => {
    // Tratar erros HTTP
    if (error.response) {
      const status = error.response.status;
      
      // 401 Unauthorized - token inválido ou expirado
      if (status === 401) {
        // Limpar token e redirecionar para login
        if (typeof window !== 'undefined') {
          localStorage.removeItem('auth_token');
          // Redirecionar apenas se não estiver na página de login
          if (window.location.pathname !== '/login') {
            window.location.href = '/login';
          }
        }
      }
      
      // 403 Forbidden - sem permissão
      if (status === 403) {
        // Tratar erro de permissão
        console.error('Acesso negado');
      }
      
      // 404 Not Found
      if (status === 404) {
        console.error('Recurso não encontrado');
      }
      
      // 500 Internal Server Error
      if (status >= 500) {
        console.error('Erro no servidor');
      }
    } else if (error.request) {
      // Requisição foi feita mas não houve resposta
      console.error('Erro de conexão com o servidor');
    } else {
      // Erro ao configurar a requisição
      console.error('Erro ao fazer requisição:', error.message);
    }
    
    return Promise.reject(error);
  }
);

export default apiClient;

