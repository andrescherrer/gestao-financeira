import apiClient from './client'
import type {
  Investment,
  CreateInvestmentRequest,
  UpdateInvestmentRequest,
  ListInvestmentsResponse,
} from './types'

/**
 * Serviço de API para investimentos
 */
export const investmentService = {
  /**
   * Lista todos os investimentos do usuário autenticado
   */
  async list(params?: {
    context?: 'PERSONAL' | 'BUSINESS'
    type?: 'STOCK' | 'FUND' | 'CDB' | 'TREASURY' | 'CRYPTO' | 'OTHER'
    page?: string
    limit?: string
  }): Promise<ListInvestmentsResponse> {
    const response = await apiClient.get<{
      message: string
      data: {
        investments: Array<{
          investment_id: string
          user_id: string
          account_id: string
          type: string
          name: string
          ticker?: string
          purchase_date: string
          purchase_amount: number
          current_value: number
          currency: string
          quantity?: number
          context: string
          return_absolute: number
          return_percentage: number
          created_at: string
          updated_at: string
        }>
        count: number
        pagination?: {
          page: number
          limit: number
          total: number
          total_pages: number
          has_next: boolean
          has_prev: boolean
        }
      }
    }>('/investments', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      investments: backendData.investments.map((inv) => ({
        investment_id: inv.investment_id,
        user_id: inv.user_id,
        account_id: inv.account_id,
        type: inv.type as Investment['type'],
        name: inv.name,
        ticker: inv.ticker,
        purchase_date: inv.purchase_date,
        purchase_amount: inv.purchase_amount.toString(),
        current_value: inv.current_value.toString(),
        currency: inv.currency as Investment['currency'],
        quantity: inv.quantity?.toString(),
        context: inv.context as Investment['context'],
        return_absolute: inv.return_absolute.toString(),
        return_percentage: inv.return_percentage,
        created_at: inv.created_at,
        updated_at: inv.updated_at,
      })),
      count: backendData.count,
      pagination: backendData.pagination,
    }
  },

  /**
   * Obtém detalhes de um investimento específico
   */
  async get(investmentId: string): Promise<Investment> {
    const response = await apiClient.get<{
      message: string
      data: {
        investment_id: string
        user_id: string
        account_id: string
        type: string
        name: string
        ticker?: string
        purchase_date: string
        purchase_amount: number
        current_value: number
        currency: string
        quantity?: number
        context: string
        return_absolute: number
        return_percentage: number
        created_at: string
        updated_at: string
      }
    }>(`/investments/${investmentId}`)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      investment_id: backendData.investment_id,
      user_id: backendData.user_id,
      account_id: backendData.account_id,
      type: backendData.type as Investment['type'],
      name: backendData.name,
      ticker: backendData.ticker,
      purchase_date: backendData.purchase_date,
      purchase_amount: backendData.purchase_amount.toString(),
      current_value: backendData.current_value.toString(),
      currency: backendData.currency as Investment['currency'],
      quantity: backendData.quantity?.toString(),
      context: backendData.context as Investment['context'],
      return_absolute: backendData.return_absolute.toString(),
      return_percentage: backendData.return_percentage,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria um novo investimento
   */
  async create(data: CreateInvestmentRequest): Promise<Investment> {
    if (import.meta.env.DEV) {
      console.log('[InvestmentService] Criando investimento com dados:', JSON.stringify(data, null, 2))
    }

    const response = await apiClient.post<{
      message: string
      data: {
        investment_id: string
        user_id: string
        account_id: string
        type: string
        name: string
        ticker?: string
        purchase_date: string
        purchase_amount: number
        current_value: number
        currency: string
        quantity?: number
        context: string
        created_at: string
      }
    }>('/investments', data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      investment_id: backendData.investment_id,
      user_id: backendData.user_id,
      account_id: backendData.account_id,
      type: backendData.type as Investment['type'],
      name: backendData.name,
      ticker: backendData.ticker,
      purchase_date: backendData.purchase_date,
      purchase_amount: backendData.purchase_amount.toString(),
      current_value: backendData.current_value.toString(),
      currency: backendData.currency as Investment['currency'],
      quantity: backendData.quantity?.toString(),
      context: backendData.context as Investment['context'],
      return_absolute: '0',
      return_percentage: 0,
      created_at: backendData.created_at,
      updated_at: backendData.created_at,
    }
  },

  /**
   * Atualiza um investimento
   */
  async update(investmentId: string, data: UpdateInvestmentRequest): Promise<Investment> {
    const response = await apiClient.put<{
      message: string
      data: {
        investment_id: string
        user_id: string
        account_id: string
        type: string
        name: string
        ticker?: string
        purchase_date: string
        purchase_amount: number
        current_value: number
        currency: string
        quantity?: number
        context: string
        return_absolute: number
        return_percentage: number
        updated_at: string
      }
    }>(`/investments/${investmentId}`, data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      investment_id: backendData.investment_id,
      user_id: backendData.user_id,
      account_id: backendData.account_id,
      type: backendData.type as Investment['type'],
      name: backendData.name,
      ticker: backendData.ticker,
      purchase_date: backendData.purchase_date,
      purchase_amount: backendData.purchase_amount.toString(),
      current_value: backendData.current_value.toString(),
      currency: backendData.currency as Investment['currency'],
      quantity: backendData.quantity?.toString(),
      context: backendData.context as Investment['context'],
      return_absolute: backendData.return_absolute.toString(),
      return_percentage: backendData.return_percentage,
      created_at: '', // Backend não retorna na atualização
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Exclui um investimento
   */
  async delete(investmentId: string): Promise<void> {
    await apiClient.delete<{
      message: string
      data: {
        success: boolean
      }
    }>(`/investments/${investmentId}`)
  },
}

