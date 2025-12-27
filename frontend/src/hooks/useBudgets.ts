import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { budgetService } from '@/api/budgets'
import type {
  Budget,
  CreateBudgetRequest,
  UpdateBudgetRequest,
  ListBudgetsResponse,
  BudgetProgress,
} from '@/api/types'

/**
 * Hook para gerenciar orçamentos usando TanStack Query
 */
export function useBudgets(params?: {
  category_id?: string
  period_type?: 'MONTHLY' | 'YEARLY'
  year?: number
  month?: number
  context?: 'PERSONAL' | 'BUSINESS'
  is_active?: boolean
}) {
  const queryClient = useQueryClient()

  // Query key base para orçamentos
  const budgetsQueryKey = ['budgets', params]

  /**
   * Lista todos os orçamentos
   */
  const {
    data: budgetsData,
    isLoading: isLoadingBudgets,
    error: budgetsError,
    refetch: refetchBudgets,
  } = useQuery<ListBudgetsResponse>({
    queryKey: budgetsQueryKey,
    queryFn: () => budgetService.list(params),
  })

  /**
   * Obtém um orçamento específico
   */
  const useBudget = (budgetId: string | null) => {
    return useQuery<Budget>({
      queryKey: ['budget', budgetId],
      queryFn: () => budgetService.get(budgetId!),
      enabled: !!budgetId,
    })
  }

  /**
   * Obtém o progresso de um orçamento
   */
  const useBudgetProgress = (budgetId: string | null) => {
    return useQuery<BudgetProgress>({
      queryKey: ['budget-progress', budgetId],
      queryFn: () => budgetService.getProgress(budgetId!),
      enabled: !!budgetId,
      refetchInterval: 60000, // Refetch a cada minuto para atualizar progresso
    })
  }

  /**
   * Cria um novo orçamento
   */
  const createBudgetMutation = useMutation({
    mutationFn: (data: CreateBudgetRequest) => budgetService.create(data),
    onSuccess: () => {
      // Invalidar queries de orçamentos para refetch
      queryClient.invalidateQueries({ queryKey: ['budgets'] })
    },
  })

  /**
   * Atualiza um orçamento
   */
  const updateBudgetMutation = useMutation({
    mutationFn: ({ budgetId, data }: { budgetId: string; data: UpdateBudgetRequest }) =>
      budgetService.update(budgetId, data),
    onSuccess: (_, variables) => {
      // Invalidar queries relacionadas
      queryClient.invalidateQueries({ queryKey: ['budgets'] })
      queryClient.invalidateQueries({ queryKey: ['budget', variables.budgetId] })
      queryClient.invalidateQueries({ queryKey: ['budget-progress', variables.budgetId] })
    },
  })

  /**
   * Deleta um orçamento
   */
  const deleteBudgetMutation = useMutation({
    mutationFn: (budgetId: string) => budgetService.delete(budgetId),
    onSuccess: () => {
      // Invalidar queries de orçamentos
      queryClient.invalidateQueries({ queryKey: ['budgets'] })
    },
  })

  return {
    // Data
    budgets: budgetsData.value?.budgets ?? [],
    total: budgetsData.value?.total ?? 0,
    isLoadingBudgets,
    budgetsError,

    // Actions
    refetchBudgets,

    // Hooks
    useBudget,
    useBudgetProgress,

    // Mutations
    createBudget: createBudgetMutation.mutate,
    createBudgetAsync: createBudgetMutation.mutateAsync,
    isCreatingBudget: createBudgetMutation.isPending,
    createBudgetError: createBudgetMutation.error,

    updateBudget: updateBudgetMutation.mutate,
    updateBudgetAsync: updateBudgetMutation.mutateAsync,
    isUpdatingBudget: updateBudgetMutation.isPending,
    updateBudgetError: updateBudgetMutation.error,

    deleteBudget: deleteBudgetMutation.mutate,
    deleteBudgetAsync: deleteBudgetMutation.mutateAsync,
    isDeletingBudget: deleteBudgetMutation.isPending,
    deleteBudgetError: deleteBudgetMutation.error,
  }
}

