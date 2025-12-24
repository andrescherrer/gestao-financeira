import apiClient from './client'
import type {
  Account,
  CreateAccountRequest,
  ListAccountsResponse,
} from './types'

/**
 * Serviço de API para contas
 */
export const accountService = {
  /**
   * Lista todas as contas do usuário autenticado
   */
  async list(context?: 'PERSONAL' | 'BUSINESS'): Promise<ListAccountsResponse> {
    const params = context ? { context } : {}
    const response = await apiClient.get<ListAccountsResponse>('/accounts', {
      params,
    })
    return response.data
  },

  /**
   * Obtém detalhes de uma conta específica
   */
  async get(accountId: string): Promise<Account> {
    const response = await apiClient.get<Account>(`/accounts/${accountId}`)
    return response.data
  },

  /**
   * Cria uma nova conta
   */
  async create(data: CreateAccountRequest): Promise<Account> {
    const response = await apiClient.post<Account>('/accounts', data)
    return response.data
  },
}

