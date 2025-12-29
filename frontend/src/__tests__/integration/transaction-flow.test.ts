import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useAccountsStore } from '@/stores/accounts'
import { useTransactionsStore } from '@/stores/transactions'
import * as authApi from '@/api/auth'
import * as accountsApi from '@/api/accounts'
import * as transactionsApi from '@/api/transactions'
import { setupIntegrationTests, createMockUser, createMockAccount, createMockTransaction } from './setup'

// Mock das APIs
vi.mock('@/api/auth', () => ({
  authService: {
    getToken: vi.fn(),
  },
}))

vi.mock('@/api/accounts', () => ({
  accountService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
  },
}))

vi.mock('@/api/transactions', () => ({
  transactionService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
  },
}))

describe('Fluxo de Transações - Integração', () => {
  setupIntegrationTests()

  beforeEach(async () => {
    // Simular usuário autenticado
    const mockUser = createMockUser()
    const mockToken = 'test-token-123'
    const mockAccount = createMockAccount()

    vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)
    localStorage.setItem('auth_token', mockToken)
    localStorage.setItem('auth_user', JSON.stringify(mockUser))

    const authStore = useAuthStore()
    authStore.token = mockToken
    authStore.user = mockUser
    authStore.isValidated = true

    // Mock: carregar contas
    vi.mocked(accountsApi.accountService.list).mockResolvedValue({
      accounts: [mockAccount] as any,
      count: 1,
    })

    const accountsStore = useAccountsStore()
    await accountsStore.listAccounts()
  })

  describe('Criação de Transação', () => {
    it('deve criar transação e atualizar a lista automaticamente', async () => {
      const existingTransactions: any[] = []
      const newTransaction = createMockTransaction()

      // Mock: listar transações existentes
      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: existingTransactions,
        count: 0,
      })

      // Mock: criar nova transação
      vi.mocked(transactionsApi.transactionService.create).mockResolvedValue(newTransaction as any)

      const transactionsStore = useTransactionsStore()

      // Carregar transações iniciais
      await transactionsStore.listTransactions()
      expect(transactionsStore.transactions).toHaveLength(0)

      // Criar nova transação
      const createdTransaction = await transactionsStore.createTransaction({
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
      })

      // Verificar que a transação foi criada
      expect(createdTransaction).toEqual(newTransaction)
      expect(transactionsApi.transactionService.create).toHaveBeenCalledWith({
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
      })

      // Verificar que a lista foi atualizada automaticamente
      expect(transactionsStore.transactions).toHaveLength(1)
      expect(transactionsStore.transactions).toContainEqual(newTransaction)
      expect(transactionsStore.isLoading).toBe(false)
    })

    it('deve atualizar propriedades computadas após criar transação', async () => {
      const incomeTransaction = createMockTransaction()
      const expenseTransaction = {
        ...createMockTransaction(),
        transaction_id: 'tx-2',
        type: 'EXPENSE' as const,
        amount: '500.00',
      }

      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: [],
        count: 0,
      })

      vi.mocked(transactionsApi.transactionService.create)
        .mockResolvedValueOnce(incomeTransaction as any)
        .mockResolvedValueOnce(expenseTransaction as any)

      const transactionsStore = useTransactionsStore()

      // Criar receita
      await transactionsStore.createTransaction({
        account_id: 'acc-1',
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Salário',
        date: '2024-01-01',
      })

      // Criar despesa
      await transactionsStore.createTransaction({
        account_id: 'acc-1',
        type: 'EXPENSE',
        amount: 500.00,
        currency: 'BRL',
        description: 'Compras',
        date: '2024-01-02',
      })

      // Verificar propriedades computadas
      expect(transactionsStore.totalTransactions).toBe(2)
      expect(transactionsStore.incomeTransactions).toHaveLength(1)
      expect(transactionsStore.expenseTransactions).toHaveLength(1)
      expect(transactionsStore.totalIncome).toBe(1000.00)
      expect(transactionsStore.totalExpense).toBe(500.00)
      expect(transactionsStore.balance).toBe(500.00)
    })

    it('deve tratar erro ao criar transação e manter estado consistente', async () => {
      const existingTransactions = [createMockTransaction()]

      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: existingTransactions as any,
        count: 1,
      })

      const error = new Error('Failed to create transaction')
      vi.mocked(transactionsApi.transactionService.create).mockRejectedValue(error)

      const transactionsStore = useTransactionsStore()
      await transactionsStore.listTransactions()

      const initialCount = transactionsStore.transactions.length

      // Tentar criar transação
      await expect(
        transactionsStore.createTransaction({
          account_id: 'acc-1',
          type: 'INCOME',
          amount: 1000.00,
          currency: 'BRL',
          description: 'Salário',
          date: '2024-01-01',
        })
      ).rejects.toThrow()

      // Verificar que a lista não foi alterada
      expect(transactionsStore.transactions).toHaveLength(initialCount)
      expect(transactionsStore.error).toBeTruthy()
      expect(transactionsStore.isLoading).toBe(false)
    })
  })

  describe('Listagem de Transações', () => {
    it('deve carregar transações e atualizar propriedades computadas', async () => {
      const transactions = [
        createMockTransaction(),
        {
          ...createMockTransaction(),
          transaction_id: 'tx-2',
          type: 'EXPENSE' as const,
          amount: '500.00',
        },
        {
          ...createMockTransaction(),
          transaction_id: 'tx-3',
          type: 'INCOME' as const,
          amount: '2000.00',
        },
      ]

      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: transactions as any,
        count: 3,
      })

      const transactionsStore = useTransactionsStore()

      await transactionsStore.listTransactions()

      // Verificar lista
      expect(transactionsStore.transactions).toHaveLength(3)
      expect(transactionsStore.isLoading).toBe(false)
      expect(transactionsStore.error).toBeNull()

      // Verificar propriedades computadas
      expect(transactionsStore.totalTransactions).toBe(3)
      expect(transactionsStore.incomeTransactions).toHaveLength(2)
      expect(transactionsStore.expenseTransactions).toHaveLength(1)
      expect(transactionsStore.totalIncome).toBe(3000.00)
      expect(transactionsStore.totalExpense).toBe(500.00)
      expect(transactionsStore.balance).toBe(2500.00)
    })

    it('deve filtrar transações por accountId', async () => {
      const transactions = [createMockTransaction()]

      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: transactions as any,
        count: 1,
      })

      const transactionsStore = useTransactionsStore()

      await transactionsStore.listTransactions('acc-1')

      expect(transactionsApi.transactionService.list).toHaveBeenCalledWith('acc-1', undefined)
    })

    it('deve filtrar transações por tipo', async () => {
      const transactions = [createMockTransaction()]

      vi.mocked(transactionsApi.transactionService.list).mockResolvedValue({
        transactions: transactions as any,
        count: 1,
      })

      const transactionsStore = useTransactionsStore()

      await transactionsStore.listTransactions(undefined, 'INCOME')

      expect(transactionsApi.transactionService.list).toHaveBeenCalledWith(undefined, 'INCOME')
    })
  })

  describe('Fluxo Completo: Conta -> Transação', () => {
    it('deve criar conta e depois criar transação nessa conta', async () => {
      const newAccount = {
        ...createMockAccount(),
        account_id: 'acc-new',
        name: 'Nova Conta',
      }

      const newTransaction = {
        ...createMockTransaction(),
        account_id: 'acc-new',
        transaction_id: 'tx-new',
      }

      // Mock: criar conta
      vi.mocked(accountsApi.accountService.create).mockResolvedValue(newAccount as any)

      // Mock: criar transação
      vi.mocked(transactionsApi.transactionService.create).mockResolvedValue(newTransaction as any)

      const accountsStore = useAccountsStore()
      const transactionsStore = useTransactionsStore()

      // Criar conta
      const createdAccount = await accountsStore.createAccount({
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 0,
        currency: 'BRL',
        context: 'PERSONAL',
      })

      expect(createdAccount.account_id).toBe('acc-new')

      // Criar transação na conta criada
      const createdTransaction = await transactionsStore.createTransaction({
        account_id: createdAccount.account_id,
        type: 'INCOME',
        amount: 1000.00,
        currency: 'BRL',
        description: 'Primeira transação',
        date: '2024-01-01',
      })

      expect(createdTransaction.account_id).toBe(createdAccount.account_id)
      expect(transactionsStore.transactions).toContainEqual(createdTransaction)
    })
  })
})
