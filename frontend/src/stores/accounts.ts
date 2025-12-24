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
      const response = await accountService.list(context)
      accounts.value = response.accounts || []
      return { accounts: response.accounts, count: response.count }
    } catch (err: any) {
      error.value =
        err.response?.data?.message ||
        err.message ||
        'Erro ao listar contas'
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
      error.value =
        err.response?.data?.message ||
        err.message ||
        'Erro ao obter detalhes da conta'
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
      error.value =
        err.response?.data?.message ||
        err.message ||
        'Erro ao criar conta'
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
    clearError,
    clearCurrentAccount,
  }
})

