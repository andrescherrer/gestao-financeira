import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { goalService } from '@/api/goals'
import { useAuthStore } from '@/stores/auth'
import type { Goal, CreateGoalRequest, AddContributionRequest, UpdateProgressRequest } from '@/api/types'

export const useGoalsStore = defineStore('goals', () => {
  const authStore = useAuthStore()

  // Estado
  const goals = ref<Goal[]>([])
  const currentGoal = ref<Goal | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalGoals = computed(() => goals.value.length)
  const personalGoals = computed(() =>
    goals.value.filter((goal) => goal.context === 'PERSONAL')
  )
  const businessGoals = computed(() =>
    goals.value.filter((goal) => goal.context === 'BUSINESS')
  )
  const inProgressGoals = computed(() =>
    goals.value.filter((goal) => goal.status === 'IN_PROGRESS')
  )
  const completedGoals = computed(() =>
    goals.value.filter((goal) => goal.status === 'COMPLETED')
  )
  const overdueGoals = computed(() =>
    goals.value.filter((goal) => goal.status === 'OVERDUE')
  )
  const totalTargetAmount = computed(() => {
    return goals.value.reduce((sum, goal) => {
      return sum + parseFloat(goal.target_amount || '0')
    }, 0)
  })
  const totalCurrentAmount = computed(() => {
    return goals.value.reduce((sum, goal) => {
      return sum + parseFloat(goal.current_amount || '0')
    }, 0)
  })

  /**
   * Lista todas as metas do usuário
   */
  async function listGoals(params?: {
    context?: 'PERSONAL' | 'BUSINESS'
    status?: 'IN_PROGRESS' | 'COMPLETED' | 'OVERDUE' | 'CANCELLED'
    page?: string
    limit?: string
  }) {
    isLoading.value = true
    error.value = null
    try {
      const token = localStorage.getItem('auth_token')
      if (!token) {
        throw new Error('Token de autenticação não encontrado. Faça login novamente.')
      }
      
      const response = await goalService.list(params)
      goals.value = response.goals || []
      return response
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      if (import.meta.env.DEV) {
        console.error('[Goals Store] Erro ao listar metas:', {
          message: error.value,
          status: err.response?.status,
          statusText: err.response?.statusText,
          data: err.response?.data,
        })
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Obtém detalhes de uma meta específica
   */
  async function getGoal(goalId: string) {
    isLoading.value = true
    error.value = null
    try {
      const goal = await goalService.get(goalId)
      currentGoal.value = goal

      // Atualiza a meta na lista se já existir
      const index = goals.value.findIndex(
        (g) => g.goal_id === goalId
      )
      if (index !== -1) {
        goals.value[index] = goal
      }

      return goal
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cria uma nova meta
   */
  async function createGoal(data: CreateGoalRequest) {
    isLoading.value = true
    error.value = null
    try {
      const goal = await goalService.create(data)
      goals.value.push(goal)
      return goal
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Adiciona uma contribuição à meta
   */
  async function addContribution(goalId: string, data: AddContributionRequest) {
    isLoading.value = true
    error.value = null
    try {
      const goal = await goalService.addContribution(goalId, data)
      
      // Atualiza a meta na lista
      const index = goals.value.findIndex(
        (g) => g.goal_id === goalId
      )
      if (index !== -1) {
        goals.value[index] = goal
      }
      
      // Atualiza currentGoal se for o mesmo
      if (currentGoal.value?.goal_id === goalId) {
        currentGoal.value = goal
      }
      
      return goal
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Atualiza o progresso da meta
   */
  async function updateProgress(goalId: string, data: UpdateProgressRequest) {
    isLoading.value = true
    error.value = null
    try {
      const goal = await goalService.updateProgress(goalId, data)
      
      // Atualiza a meta na lista
      const index = goals.value.findIndex(
        (g) => g.goal_id === goalId
      )
      if (index !== -1) {
        goals.value[index] = goal
      }
      
      // Atualiza currentGoal se for o mesmo
      if (currentGoal.value?.goal_id === goalId) {
        currentGoal.value = goal
      }
      
      return goal
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cancela uma meta
   */
  async function cancelGoal(goalId: string) {
    isLoading.value = true
    error.value = null
    try {
      const goal = await goalService.cancel(goalId)
      
      // Atualiza a meta na lista
      const index = goals.value.findIndex(
        (g) => g.goal_id === goalId
      )
      if (index !== -1) {
        goals.value[index] = goal
      }
      
      // Atualiza currentGoal se for o mesmo
      if (currentGoal.value?.goal_id === goalId) {
        currentGoal.value = goal
      }
      
      return goal
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Exclui uma meta
   */
  async function deleteGoal(goalId: string) {
    isLoading.value = true
    error.value = null
    try {
      await goalService.delete(goalId)
      
      // Remove da lista
      goals.value = goals.value.filter(
        (g) => g.goal_id !== goalId
      )
      
      // Limpa currentGoal se for o mesmo
      if (currentGoal.value?.goal_id === goalId) {
        currentGoal.value = null
      }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Limpa o estado
   */
  function clearError() {
    error.value = null
  }

  function clearCurrentGoal() {
    currentGoal.value = null
  }

  return {
    // Estado
    goals,
    currentGoal,
    isLoading,
    error,
    // Computed
    totalGoals,
    personalGoals,
    businessGoals,
    inProgressGoals,
    completedGoals,
    overdueGoals,
    totalTargetAmount,
    totalCurrentAmount,
    // Ações
    listGoals,
    getGoal,
    createGoal,
    addContribution,
    updateProgress,
    cancelGoal,
    deleteGoal,
    clearError,
    clearCurrentGoal,
  }
})

