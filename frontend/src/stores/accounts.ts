import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { accountService } from '@/api/accounts'
import { useAuthStore } from '@/stores/auth'
import type { Account, CreateAccountRequest } from '@/api/types'

export const useAccountsStore = defineStore('accounts', () => {
  const authStore = useAuthStore()

  // Estado
  const accounts = ref<Account[]>([])
  const currentAccount = ref<Account | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalAccounts = computed(() => accounts.value.length)
  const activeAccounts = computed(() =>
    accounts.value.filter((account) => account.is_active)
  )
  const personalAccounts = computed(() =>
    accounts.value.filter((account) => account.context === 'PERSONAL')
  )
  const businessAccounts = computed(() =>
    accounts.value.filter((account) => account.context === 'BUSINESS')
  )

  /**
   * Lista todas as contas do usuário
   */
  async function listAccounts(context?: 'PERSONAL' | 'BUSINESS') {
    isLoading.value = true
    error.value = null
    try {
      // Verificar se há token antes de fazer a requisição
      const token = localStorage.getItem('auth_token')
      if (!token) {
        throw new Error('Token de autenticação não encontrado. Faça login novamente.')
      }
      
      const response = await accountService.list(context)
      accounts.value = response.accounts || []
      return { accounts: response.accounts, count: response.count }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      // Log detalhado do erro em desenvolvimento
      if (import.meta.env.DEV) {
        console.error('[Accounts Store] Erro ao listar contas:', {
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
   * Obtém detalhes de uma conta específica
   */
  async function getAccount(accountId: string) {
    isLoading.value = true
    error.value = null
    try {
      const account = await accountService.get(accountId)
      currentAccount.value = account

      // Atualiza a conta na lista se já existir
      const index = accounts.value.findIndex(
        (acc) => acc.account_id === accountId
      )
      if (index !== -1) {
        accounts.value[index] = account
      }

      return account
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cria uma nova conta
   */
  async function createAccount(data: CreateAccountRequest) {
    isLoading.value = true
    error.value = null
    try {
      const account = await accountService.create(data)
      accounts.value.push(account)
      return account
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Atualiza uma conta específica (útil após operações de transação)
   */
  async function refreshAccount(accountId: string) {
    try {
      const account = await accountService.get(accountId)
      
      // Atualiza a conta na lista se já existir
      const index = accounts.value.findIndex(
        (acc) => acc.account_id === accountId
      )
      if (index !== -1) {
        accounts.value[index] = account
      }
      
      // Atualiza currentAccount se for a mesma
      if (currentAccount.value?.account_id === accountId) {
        currentAccount.value = account
      }
      
      return account
    } catch (err: any) {
      // Log erro mas não falha silenciosamente
      if (import.meta.env.DEV) {
        console.warn('[Accounts Store] Erro ao atualizar conta:', err)
      }
      // Não propagar erro para não quebrar o fluxo principal
      return null
    }
  }

  /**
   * Limpa o estado
   */
  function clearError() {
    error.value = null
  }

  function clearCurrentAccount() {
    currentAccount.value = null
  }

  return {
    // Estado
    accounts,
    currentAccount,
    isLoading,
    error,
    // Computed
    totalAccounts,
    activeAccounts,
    personalAccounts,
    businessAccounts,
    // Ações
    listAccounts,
    getAccount,
    createAccount,
    refreshAccount,
    clearError,
    clearCurrentAccount,
  }
})

