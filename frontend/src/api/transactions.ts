import apiClient from './client'
import type {
  Transaction,
  CreateTransactionRequest,
  ListTransactionsResponse,
} from './types'

/**
 * Serviço de API para transações
 */
export const transactionService = {
  /**
   * Lista todas as transações do usuário autenticado
   */
  async list(accountId?: string, type?: 'INCOME' | 'EXPENSE'): Promise<ListTransactionsResponse> {
    const params: Record<string, string> = {}
    if (accountId) params.account_id = accountId
    if (type) params.type = type

    const response = await apiClient.get<{
      message: string
      data: {
        transactions: Array<{
          transaction_id: string
          user_id: string
          account_id: string
          type: string
          amount: number
          currency: string
          description: string
          date: string
          created_at: string
          updated_at: string
        }>
        count: number
      }
    }>('/transactions', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      transactions: backendData.transactions.map((tx) => ({
        transaction_id: tx.transaction_id,
        user_id: tx.user_id,
        account_id: tx.account_id,
        type: tx.type as Transaction['type'],
        amount: tx.amount.toString(),
        currency: tx.currency as Transaction['currency'],
        description: tx.description,
        date: tx.date,
        created_at: tx.created_at,
        updated_at: tx.updated_at,
      })),
      count: backendData.count,
    }
  },

  /**
   * Obtém detalhes de uma transação específica
   */
  async get(transactionId: string): Promise<Transaction> {
    const response = await apiClient.get<{
      message: string
      data: {
        transaction_id: string
        user_id: string
        account_id: string
        type: string
        amount: number
        currency: string
        description: string
        date: string
        created_at: string
        updated_at: string
      }
    }>(`/transactions/${transactionId}`)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      transaction_id: backendData.transaction_id,
      user_id: backendData.user_id,
      account_id: backendData.account_id,
      type: backendData.type as Transaction['type'],
      amount: backendData.amount.toString(),
      currency: backendData.currency as Transaction['currency'],
      description: backendData.description,
      date: backendData.date,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria uma nova transação
   */
  async create(data: CreateTransactionRequest): Promise<Transaction> {
    // Log em desenvolvimento para debug
    if (import.meta.env.DEV) {
      console.log('[TransactionService] Criando transação com dados:', JSON.stringify(data, null, 2))
    }

    const response = await apiClient.post<{
      message: string
      data: {
        transaction_id: string
        user_id: string
        account_id: string
        type: string
        amount: number
        currency: string
        description: string
        date: string
        created_at: string
        updated_at: string
      }
    }>('/transactions', data)
    
    // Mapear resposta do backend para o formato esperado pelo frontend
    const backendData = response.data.data
    return {
      transaction_id: backendData.transaction_id,
      user_id: backendData.user_id,
      account_id: backendData.account_id,
      type: backendData.type as Transaction['type'],
      amount: backendData.amount.toString(),
      currency: backendData.currency as Transaction['currency'],
      description: backendData.description,
      date: backendData.date,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },
}

