import apiClient from './client';
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from './types';

/**
 * Serviço de autenticação
 */
export const authService = {
  /**
   * Realiza login do usuário
   */
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);
    return response.data;
  },

  /**
   * Registra um novo usuário
   */
  async register(userData: RegisterRequest): Promise<RegisterResponse> {
    const response = await apiClient.post<RegisterResponse>('/auth/register', userData);
    return response.data;
  },

  /**
   * Salva o token de autenticação
   */
  saveToken(token: string): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
    }
  },

  /**
   * Remove o token de autenticação
   */
  removeToken(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
    }
  },

  /**
   * Obtém o token de autenticação
   */
  getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  },

  /**
   * Verifica se o usuário está autenticado
   */
  isAuthenticated(): boolean {
    return this.getToken() !== null;
  },
};

