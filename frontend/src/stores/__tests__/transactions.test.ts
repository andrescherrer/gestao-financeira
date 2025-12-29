import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTransactionsStore } from '../transactions'
import * as transactionsApi from '@/api/transactions'
import type { Transaction, CreateTransactionRequest } from '@/api/types'

// Mock do módulo de API
vi.mock('@/api/transactions', () => ({
  transactionService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
  },
}))

describe('Transactions Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('Estado inicial', () => {
    it('deve inicializar com estado vazio', () => {
      const store = useTransactionsStore()
      
      expect(store.transactions).toEqual([])
      expect(store.currentTransaction).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.totalTransactions).toBe(0)
      expect(store.incomeTransactions).toEqual([])
      expect(store.expenseTransactions).toEqual([])
      expect(store.totalIncome).toBe(0)
      expect(store.totalExpense).toBe(0)
      expect(store.balance).toBe(0)
    })
  })

  describe('listTransactions()', () => {
    it('deve listar transações com sucesso', async () => {
      const mockTransactions: Transaction[] = [
        {
          transaction_id: 'tx-1',
          account_id: 'acc-1',
          user_id: 'user-123',
          type: 'INCOME',
          amount: '1000.00',
          currency: 'BRL',
          description: 'Salário',
          date: '2024-01-01',
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
        {
          transaction_id: 'tx-2',
          account_id: 'acc-1',
          user_id: 'user-123',
          type: 'EXPENSE',
          amount: '500.00',
          currency: 'BRL',
          description: 'Compras',
          date: '2024-01-02',
          created_at: '2024-01-02T00:00:00Z',
          updated_at: '2024-01-02T00:00:00Z',
        },
      ]
      
      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: mockTransactions,
        count: 2,
      })
      
      const store = useTransactionsStore()
      const result = await store.listTransactions()
      
      expect(transactionsApi.transactionService.list).toHaveBeenCalledWith(undefined, undefined)
      expect(store.transactions).toEqual(mockTransactions)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result.transactions).toEqual(mockTransactions)
      expect(result.count).toBe(2)
    })

    it('deve filtrar transações por accountId', async () => {
      const mockTransactions: Transaction[] = []
      
      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: mockTransactions,
        count: 0,
      })
      
      const store = useTransactionsStore()
      await store.listTransactions('acc-1')
      
      expect(transactionsApi.transactionService.list).toHaveBeenCalledWith('acc-1', undefined)
    })

    it('deve filtrar transações por tipo', async () => {
      const mockTransactions: Transaction[] = []
      
      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: mockTransactions,
        count: 0,
      })
      
      const store = useTransactionsStore()
      await store.listTransactions(undefined, 'INCOME')
      
      expect(transactionsApi.transactionService.list).toHaveBeenCalledWith(undefined, 'INCOME')
    })

    it('deve tratar erro ao listar transações', async () => {
      const error = new Error('Failed to fetch transactions')
      vi.mocked(transactionsApi.transactionService.list).mockRejectedValue(error)
      
      const store = useTransactionsStore()
      
      await expect(store.listTransactions()).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })

  describe('getTransaction()', () => {
    it('deve obter transação específica com sucesso', async () => {
      const mockTransaction: Transaction = {
        transaction_id: 'tx-1',
        account_id: 'acc-1',
        user_id: 'user-123',
        type: 'INCOME',
        amount: '1000.00',
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(transactionsApi.transactionService.get).mockResolvedValue(mockTransaction)
      
      const store = useTransactionsStore()
      const result = await store.getTransaction('tx-1')
      
      expect(transactionsApi.transactionService.get).toHaveBeenCalledWith('tx-1')
      expect(store.currentTransaction).toEqual(mockTransaction)
      expect(result).toEqual(mockTransaction)
      expect(store.isLoading).toBe(false)
    })

    it('deve atualizar transação na lista se já existir', async () => {
      const existingTransaction: Transaction = {
        transaction_id: 'tx-1',
        account_id: 'acc-1',
        user_id: 'user-123',
        type: 'INCOME',
        amount: '1000.00',
        currency: 'BRL',
        description: 'Salário Antigo',
        date: '2024-01-01',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      const updatedTransaction: Transaction = {
        ...existingTransaction,
        description: 'Salário Atualizado',
        amount: '1500.00',
      }
      
      const store = useTransactionsStore()
      store.transactions = [existingTransaction]
      
      vi.mocked(transactionsApi.transactionService.get).mockResolvedValue(updatedTransaction)
      
      await store.getTransaction('tx-1')
      
      expect(store.transactions[0]).toEqual(updatedTransaction)
    })
  })

  describe('createTransaction()', () => {
    it('deve criar transação com sucesso', async () => {
      const newTransactionData: CreateTransactionRequest = {
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
      }
      
      const createdTransaction: Transaction = {
        transaction_id: 'tx-new',
        account_id: 'acc-1',
        user_id: 'user-123',
        type: 'INCOME',
        amount: '1000.00',
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(transactionsApi.transactionService.create).mockResolvedValue(createdTransaction)
      
      const store = useTransactionsStore()
      const result = await store.createTransaction(newTransactionData)
      
      expect(transactionsApi.transactionService.create).toHaveBeenCalledWith(newTransactionData)
      expect(store.transactions).toContainEqual(createdTransaction)
      expect(result).toEqual(createdTransaction)
      expect(store.isLoading).toBe(false)
    })

    it('deve tratar erro ao criar transação', async () => {
      const newTransactionData: CreateTransactionRequest = {
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
      }
      
      const error = new Error('Failed to create transaction')
      vi.mocked(transactionsApi.transactionService.create).mockRejectedValue(error)
      
      const store = useTransactionsStore()
      
      await expect(store.createTransaction(newTransactionData)).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })

  describe('Computed properties', () => {
    it('deve calcular totalTransactions corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1' } as Transaction,
        { transaction_id: 'tx-2' } as Transaction,
        { transaction_id: 'tx-3' } as Transaction,
      ]
      
      expect(store.totalTransactions).toBe(3)
    })

    it('deve filtrar incomeTransactions corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1', type: 'INCOME' } as Transaction,
        { transaction_id: 'tx-2', type: 'EXPENSE' } as Transaction,
        { transaction_id: 'tx-3', type: 'INCOME' } as Transaction,
      ]
      
      expect(store.incomeTransactions).toHaveLength(2)
      expect(store.incomeTransactions.every(tx => tx.type === 'INCOME')).toBe(true)
    })

    it('deve filtrar expenseTransactions corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1', type: 'INCOME' } as Transaction,
        { transaction_id: 'tx-2', type: 'EXPENSE' } as Transaction,
        { transaction_id: 'tx-3', type: 'EXPENSE' } as Transaction,
      ]
      
      expect(store.expenseTransactions).toHaveLength(2)
      expect(store.expenseTransactions.every(tx => tx.type === 'EXPENSE')).toBe(true)
    })

    it('deve calcular totalIncome corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1', type: 'INCOME', amount: '1000.00' } as Transaction,
        { transaction_id: 'tx-2', type: 'INCOME', amount: '500.00' } as Transaction,
        { transaction_id: 'tx-3', type: 'EXPENSE', amount: '200.00' } as Transaction,
      ]
      
      expect(store.totalIncome).toBe(1500.00)
    })

    it('deve calcular totalExpense corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1', type: 'INCOME', amount: '1000.00' } as Transaction,
        { transaction_id: 'tx-2', type: 'EXPENSE', amount: '500.00' } as Transaction,
        { transaction_id: 'tx-3', type: 'EXPENSE', amount: '200.00' } as Transaction,
      ]
      
      expect(store.totalExpense).toBe(700.00)
    })

    it('deve calcular balance corretamente', () => {
      const store = useTransactionsStore()
      store.transactions = [
        { transaction_id: 'tx-1', type: 'INCOME', amount: '1000.00' } as Transaction,
        { transaction_id: 'tx-2', type: 'EXPENSE', amount: '500.00' } as Transaction,
        { transaction_id: 'tx-3', type: 'EXPENSE', amount: '200.00' } as Transaction,
      ]
      
      expect(store.balance).toBe(300.00)
    })
  })

  describe('clearError()', () => {
    it('deve limpar erro', () => {
      const store = useTransactionsStore()
      store.error = 'Some error'
      
      store.clearError()
      
      expect(store.error).toBeNull()
    })
  })

  describe('clearCurrentTransaction()', () => {
    it('deve limpar transação atual', () => {
      const store = useTransactionsStore()
      store.currentTransaction = { transaction_id: 'tx-1' } as Transaction
      
      store.clearCurrentTransaction()
      
      expect(store.currentTransaction).toBeNull()
    })
  })
})
