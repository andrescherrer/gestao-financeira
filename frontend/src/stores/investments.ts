import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { investmentService } from '@/api/investments'
import { useAuthStore } from '@/stores/auth'
import type { Investment, CreateInvestmentRequest, UpdateInvestmentRequest } from '@/api/types'

export const useInvestmentsStore = defineStore('investments', () => {
  const authStore = useAuthStore()

  // Estado
  const investments = ref<Investment[]>([])
  const currentInvestment = ref<Investment | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalInvestments = computed(() => investments.value.length)
  const personalInvestments = computed(() =>
    investments.value.filter((investment) => investment.context === 'PERSONAL')
  )
  const businessInvestments = computed(() =>
    investments.value.filter((investment) => investment.context === 'BUSINESS')
  )
  const totalValue = computed(() => {
    return investments.value.reduce((sum, inv) => {
      return sum + parseFloat(inv.current_value || '0')
    }, 0)
  })
  const totalReturn = computed(() => {
    return investments.value.reduce((sum, inv) => {
      return sum + parseFloat(inv.return_absolute || '0')
    }, 0)
  })

  /**
   * Lista todos os investimentos do usuário
   */
  async function listInvestments(params?: {
    context?: 'PERSONAL' | 'BUSINESS'
    type?: 'STOCK' | 'FUND' | 'CDB' | 'TREASURY' | 'CRYPTO' | 'OTHER'
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
      
      const response = await investmentService.list(params)
      investments.value = response.investments || []
      return response
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      if (import.meta.env.DEV) {
        console.error('[Investments Store] Erro ao listar investimentos:', {
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
   * Obtém detalhes de um investimento específico
   */
  async function getInvestment(investmentId: string) {
    isLoading.value = true
    error.value = null
    try {
      const investment = await investmentService.get(investmentId)
      currentInvestment.value = investment

      // Atualiza o investimento na lista se já existir
      const index = investments.value.findIndex(
        (inv) => inv.investment_id === investmentId
      )
      if (index !== -1) {
        investments.value[index] = investment
      }

      return investment
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cria um novo investimento
   */
  async function createInvestment(data: CreateInvestmentRequest) {
    isLoading.value = true
    error.value = null
    try {
      const investment = await investmentService.create(data)
      investments.value.push(investment)
      return investment
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Atualiza um investimento
   */
  async function updateInvestment(investmentId: string, data: UpdateInvestmentRequest) {
    isLoading.value = true
    error.value = null
    try {
      const investment = await investmentService.update(investmentId, data)
      
      // Atualiza o investimento na lista
      const index = investments.value.findIndex(
        (inv) => inv.investment_id === investmentId
      )
      if (index !== -1) {
        investments.value[index] = investment
      }
      
      // Atualiza currentInvestment se for o mesmo
      if (currentInvestment.value?.investment_id === investmentId) {
        currentInvestment.value = investment
      }
      
      return investment
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Exclui um investimento
   */
  async function deleteInvestment(investmentId: string) {
    isLoading.value = true
    error.value = null
    try {
      await investmentService.delete(investmentId)
      
      // Remove da lista
      investments.value = investments.value.filter(
        (inv) => inv.investment_id !== investmentId
      )
      
      // Limpa currentInvestment se for o mesmo
      if (currentInvestment.value?.investment_id === investmentId) {
        currentInvestment.value = null
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

  function clearCurrentInvestment() {
    currentInvestment.value = null
  }

  return {
    // Estado
    investments,
    currentInvestment,
    isLoading,
    error,
    // Computed
    totalInvestments,
    personalInvestments,
    businessInvestments,
    totalValue,
    totalReturn,
    // Ações
    listInvestments,
    getInvestment,
    createInvestment,
    updateInvestment,
    deleteInvestment,
    clearError,
    clearCurrentInvestment,
  }
})

