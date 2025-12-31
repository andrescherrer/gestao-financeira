import apiClient from './client'
import type {
  Goal,
  CreateGoalRequest,
  AddContributionRequest,
  UpdateProgressRequest,
  ListGoalsResponse,
} from './types'

/**
 * Serviço de API para metas
 */
export const goalService = {
  /**
   * Lista todas as metas do usuário autenticado
   */
  async list(params?: {
    context?: 'PERSONAL' | 'BUSINESS'
    status?: 'IN_PROGRESS' | 'COMPLETED' | 'OVERDUE' | 'CANCELLED'
    page?: string
    limit?: string
  }): Promise<ListGoalsResponse> {
    const response = await apiClient.get<{
      message: string
      data: {
        goals: Array<{
          goal_id: string
          user_id: string
          name: string
          target_amount: number
          current_amount: number
          currency: string
          deadline: string
          context: string
          status: string
          progress: number
          remaining_days: number
          created_at: string
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
    }>('/goals', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      goals: backendData.goals.map((goal) => ({
        goal_id: goal.goal_id,
        user_id: goal.user_id,
        name: goal.name,
        target_amount: goal.target_amount.toString(),
        current_amount: goal.current_amount.toString(),
        currency: goal.currency as Goal['currency'],
        deadline: goal.deadline,
        context: goal.context as Goal['context'],
        status: goal.status as Goal['status'],
        progress: goal.progress,
        remaining_days: goal.remaining_days,
        created_at: goal.created_at,
        updated_at: goal.created_at, // Backend não retorna updated_at na listagem
      })),
      count: backendData.count,
      pagination: backendData.pagination,
    }
  },

  /**
   * Obtém detalhes de uma meta específica
   */
  async get(goalId: string): Promise<Goal> {
    const response = await apiClient.get<{
      message: string
      data: {
        goal_id: string
        user_id: string
        name: string
        target_amount: number
        current_amount: number
        currency: string
        deadline: string
        context: string
        status: string
        progress: number
        remaining_days: number
        created_at: string
        updated_at: string
      }
    }>(`/goals/${goalId}`)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      goal_id: backendData.goal_id,
      user_id: backendData.user_id,
      name: backendData.name,
      target_amount: backendData.target_amount.toString(),
      current_amount: backendData.current_amount.toString(),
      currency: backendData.currency as Goal['currency'],
      deadline: backendData.deadline,
      context: backendData.context as Goal['context'],
      status: backendData.status as Goal['status'],
      progress: backendData.progress,
      remaining_days: backendData.remaining_days,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria uma nova meta
   */
  async create(data: CreateGoalRequest): Promise<Goal> {
    if (import.meta.env.DEV) {
      console.log('[GoalService] Criando meta com dados:', JSON.stringify(data, null, 2))
    }

    const response = await apiClient.post<{
      message: string
      data: {
        goal_id: string
        user_id: string
        name: string
        target_amount: number
        current_amount: number
        currency: string
        deadline: string
        context: string
        status: string
        progress: number
        remaining_days: number
        created_at: string
      }
    }>('/goals', data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      goal_id: backendData.goal_id,
      user_id: backendData.user_id,
      name: backendData.name,
      target_amount: backendData.target_amount.toString(),
      current_amount: backendData.current_amount.toString(),
      currency: backendData.currency as Goal['currency'],
      deadline: backendData.deadline,
      context: backendData.context as Goal['context'],
      status: backendData.status as Goal['status'],
      progress: backendData.progress,
      remaining_days: backendData.remaining_days,
      created_at: backendData.created_at,
      updated_at: backendData.created_at,
    }
  },

  /**
   * Adiciona uma contribuição à meta
   */
  async addContribution(goalId: string, data: AddContributionRequest): Promise<Goal> {
    const response = await apiClient.post<{
      message: string
      data: {
        goal_id: string
        current_amount: number
        target_amount: number
        currency: string
        progress: number
        status: string
        remaining_days: number
        updated_at: string
      }
    }>(`/goals/${goalId}/contribute`, data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    // Precisamos buscar a meta completa para ter todos os dados
    return this.get(goalId)
  },

  /**
   * Atualiza o progresso da meta
   */
  async updateProgress(goalId: string, data: UpdateProgressRequest): Promise<Goal> {
    const response = await apiClient.put<{
      message: string
      data: {
        goal_id: string
        current_amount: number
        target_amount: number
        currency: string
        progress: number
        status: string
        remaining_days: number
        updated_at: string
      }
    }>(`/goals/${goalId}/progress`, data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    // Precisamos buscar a meta completa para ter todos os dados
    return this.get(goalId)
  },

  /**
   * Cancela uma meta
   */
  async cancel(goalId: string): Promise<Goal> {
    await apiClient.post<{
      message: string
      data: {
        goal_id: string
        status: string
        updated_at: string
      }
    }>(`/goals/${goalId}/cancel`)
    
    // Buscar meta atualizada
    return this.get(goalId)
  },

  /**
   * Exclui uma meta
   */
  async delete(goalId: string): Promise<void> {
    await apiClient.delete<{
      message: string
      data: {
        success: boolean
      }
    }>(`/goals/${goalId}`)
  },
}

