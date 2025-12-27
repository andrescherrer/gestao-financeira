import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { transactionService } from '@/api/transactions'
import type { Transaction, CreateTransactionRequest } from '@/api/types'

export const useTransactionsStore = defineStore('transactions', () => {
  // Estado
  const transactions = ref<Transaction[]>([])
  const currentTransaction = ref<Transaction | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalTransactions = computed(() => transactions.value.length)
  const incomeTransactions = computed(() =>
    transactions.value.filter((tx) => tx.type === 'INCOME')
  )
  const expenseTransactions = computed(() =>
    transactions.value.filter((tx) => tx.type === 'EXPENSE')
  )
  const totalIncome = computed(() => {
    return incomeTransactions.value.reduce((sum, tx) => {
      return sum + parseFloat(tx.amount)
    }, 0)
  })
  const totalExpense = computed(() => {
    return expenseTransactions.value.reduce((sum, tx) => {
      return sum + parseFloat(tx.amount)
    }, 0)
  })
  const balance = computed(() => {
    return totalIncome.value - totalExpense.value
  })

  /**
   * Lista todas as transações do usuário
   */
  async function listTransactions(accountId?: string, type?: 'INCOME' | 'EXPENSE') {
    isLoading.value = true
    error.value = null
    try {
      const response = await transactionService.list(accountId, type)
      transactions.value = response.transactions || []
      return { transactions: response.transactions, count: response.count }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao listar transações'
      error.value = extractErrorMessage(rawError)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Obtém detalhes de uma transação específica
   */
  async function getTransaction(transactionId: string) {
    isLoading.value = true
    error.value = null
    try {
      const transaction = await transactionService.get(transactionId)
      currentTransaction.value = transaction

      // Atualiza a transação na lista se já existir
      const index = transactions.value.findIndex(
        (tx) => tx.transaction_id === transactionId
      )
      if (index !== -1) {
        transactions.value[index] = transaction
      }

      return transaction
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao obter detalhes da transação'
      error.value = extractErrorMessage(rawError)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cria uma nova transação
   */
  async function createTransaction(data: CreateTransactionRequest) {
    isLoading.value = true
    error.value = null
    try {
      const transaction = await transactionService.create(data)
      transactions.value.push(transaction)
      return transaction
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao criar transação'
      error.value = extractErrorMessage(rawError)
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

  function clearCurrentTransaction() {
    currentTransaction.value = null
  }

  return {
    // Estado
    transactions,
    currentTransaction,
    isLoading,
    error,
    // Computed
    totalTransactions,
    incomeTransactions,
    expenseTransactions,
    totalIncome,
    totalExpense,
    balance,
    // Ações
    listTransactions,
    getTransaction,
    createTransaction,
    clearError,
    clearCurrentTransaction,
  }
})

