import apiClient from './client'
import type {
  Budget,
  CreateBudgetRequest,
  UpdateBudgetRequest,
  ListBudgetsResponse,
  BudgetProgress,
} from './types'

/**
 * Serviço de API para orçamentos
 */
export const budgetService = {
  /**
   * Lista todos os orçamentos do usuário autenticado
   */
  async list(params?: {
    category_id?: string
    period_type?: 'MONTHLY' | 'YEARLY'
    year?: number
    month?: number
    context?: 'PERSONAL' | 'BUSINESS'
    is_active?: boolean
  }): Promise<ListBudgetsResponse> {
    const response = await apiClient.get<{
      message: string
      data: {
        budgets: Array<{
          budget_id: string
          user_id: string
          category_id: string
          amount: number
          currency: string
          period_type: string
          year: number
          month?: number
          context: string
          is_active: boolean
          created_at: string
          updated_at: string
        }>
        total: number
      }
    }>('/budgets', {
      params,
    })

    const backendData = response.data.data
    return {
      budgets: backendData.budgets.map((budget) => ({
        budget_id: budget.budget_id,
        user_id: budget.user_id,
        category_id: budget.category_id,
        amount: budget.amount,
        currency: budget.currency as Budget['currency'],
        period_type: budget.period_type as Budget['period_type'],
        year: budget.year,
        month: budget.month,
        context: budget.context as Budget['context'],
        is_active: budget.is_active,
        created_at: budget.created_at,
        updated_at: budget.updated_at,
      })),
      total: backendData.total,
    }
  },

  /**
   * Obtém detalhes de um orçamento específico
   */
  async get(budgetId: string): Promise<Budget> {
    const response = await apiClient.get<{
      message: string
      data: {
        budget_id: string
        user_id: string
        category_id: string
        amount: number
        currency: string
        period_type: string
        year: number
        month?: number
        context: string
        is_active: boolean
        created_at: string
        updated_at: string
      }
    }>(`/budgets/${budgetId}`)

    const backendData = response.data.data
    return {
      budget_id: backendData.budget_id,
      user_id: backendData.user_id,
      category_id: backendData.category_id,
      amount: backendData.amount,
      currency: backendData.currency as Budget['currency'],
      period_type: backendData.period_type as Budget['period_type'],
      year: backendData.year,
      month: backendData.month,
      context: backendData.context as Budget['context'],
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria um novo orçamento
   */
  async create(data: CreateBudgetRequest): Promise<Budget> {
    if (import.meta.env.DEV) {
      console.log('[BudgetService] Criando orçamento com dados:', JSON.stringify(data, null, 2))
    }

    const response = await apiClient.post<{
      message: string
      data: {
        budget_id: string
        user_id: string
        category_id: string
        amount: number
        currency: string
        period_type: string
        year: number
        month?: number
        context: string
        is_active: boolean
        created_at: string
      }
    }>('/budgets', data)

    const backendData = response.data.data
    return {
      budget_id: backendData.budget_id,
      user_id: backendData.user_id,
      category_id: backendData.category_id,
      amount: backendData.amount,
      currency: backendData.currency as Budget['currency'],
      period_type: backendData.period_type as Budget['period_type'],
      year: backendData.year,
      month: backendData.month,
      context: backendData.context as Budget['context'],
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.created_at,
    }
  },

  /**
   * Atualiza um orçamento existente
   */
  async update(budgetId: string, data: UpdateBudgetRequest): Promise<Budget> {
    const response = await apiClient.put<{
      message: string
      data: {
        budget_id: string
        user_id: string
        category_id: string
        amount: number
        currency: string
        period_type: string
        year: number
        month?: number
        context: string
        is_active: boolean
        created_at: string
        updated_at: string
      }
    }>(`/budgets/${budgetId}`, data)

    const backendData = response.data.data
    return {
      budget_id: backendData.budget_id,
      user_id: backendData.user_id,
      category_id: backendData.category_id,
      amount: backendData.amount,
      currency: backendData.currency as Budget['currency'],
      period_type: backendData.period_type as Budget['period_type'],
      year: backendData.year,
      month: backendData.month,
      context: backendData.context as Budget['context'],
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Deleta um orçamento
   */
  async delete(budgetId: string): Promise<void> {
    await apiClient.delete(`/budgets/${budgetId}`)
  },

  /**
   * Obtém o progresso de um orçamento
   */
  async getProgress(budgetId: string): Promise<BudgetProgress> {
    const response = await apiClient.get<{
      message: string
      data: {
        budget_id: string
        category_id: string
        budgeted: number
        spent: number
        remaining: number
        percentage_used: number
        currency: string
        is_exceeded: boolean
        period_type: string
        year: number
        month?: number
      }
    }>(`/budgets/${budgetId}/progress`)

    const backendData = response.data.data
    return {
      budget_id: backendData.budget_id,
      category_id: backendData.category_id,
      budgeted: backendData.budgeted,
      spent: backendData.spent,
      remaining: backendData.remaining,
      percentage_used: backendData.percentage_used,
      currency: backendData.currency as BudgetProgress['currency'],
      is_exceeded: backendData.is_exceeded,
      period_type: backendData.period_type as BudgetProgress['period_type'],
      year: backendData.year,
      month: backendData.month,
    }
  },
}

