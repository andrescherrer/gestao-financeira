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
    const response = await apiClient.get<{
      message: string
      data: {
        accounts: Array<{
          account_id: string
          user_id: string
          name: string
          type: string
          balance: number
          currency: string
          context: string
          is_active: boolean
          created_at: string
          updated_at: string
        }>
        count: number
      }
    }>('/accounts', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      accounts: backendData.accounts.map((acc) => ({
        account_id: acc.account_id,
        user_id: acc.user_id,
        name: acc.name,
        type: acc.type as Account['type'],
        balance: acc.balance.toString(),
        currency: acc.currency as Account['currency'],
        context: acc.context as Account['context'],
        is_active: acc.is_active,
        created_at: acc.created_at,
        updated_at: acc.updated_at,
      })),
      count: backendData.count,
    }
  },

  /**
   * Obtém detalhes de uma conta específica
   */
  async get(accountId: string): Promise<Account> {
    const response = await apiClient.get<{
      message: string
      data: {
        account_id: string
        user_id: string
        name: string
        type: string
        balance: number
        currency: string
        context: string
        is_active: boolean
        created_at: string
        updated_at: string
      }
    }>(`/accounts/${accountId}`)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      account_id: backendData.account_id,
      user_id: backendData.user_id,
      name: backendData.name,
      type: backendData.type as Account['type'],
      balance: backendData.balance.toString(),
      currency: backendData.currency as Account['currency'],
      context: backendData.context as Account['context'],
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria uma nova conta
   */
  async create(data: CreateAccountRequest): Promise<Account> {
    const response = await apiClient.post<{
      message: string
      data: {
        account_id: string
        user_id: string
        name: string
        type: string
        balance: number
        currency: string
        context: string
        is_active: boolean
        created_at: string
      }
    }>('/accounts', data)
    
    // Mapear resposta do backend para o formato esperado pelo frontend
    const backendData = response.data.data
    return {
      account_id: backendData.account_id,
      user_id: backendData.user_id,
      name: backendData.name,
      type: backendData.type as Account['type'],
      balance: backendData.balance.toString(),
      currency: backendData.currency as Account['currency'],
      context: backendData.context as Account['context'],
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.created_at, // Backend não retorna updated_at na criação
    }
  },
}

