import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAccountsStore } from '../accounts'
import * as accountsApi from '@/api/accounts'
import type { Account, CreateAccountRequest } from '@/api/types'

// Mock do módulo de API
vi.mock('@/api/accounts', () => ({
  accountService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
  },
}))

// Mock do auth store
vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    user: { user_id: 'user-123' },
  })),
}))

describe('Accounts Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('Estado inicial', () => {
    it('deve inicializar com estado vazio', () => {
      const store = useAccountsStore()
      
      expect(store.accounts).toEqual([])
      expect(store.currentAccount).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.totalAccounts).toBe(0)
      expect(store.activeAccounts).toEqual([])
      expect(store.personalAccounts).toEqual([])
      expect(store.businessAccounts).toEqual([])
    })
  })

  describe('listAccounts()', () => {
    it('deve listar contas com sucesso', async () => {
      const mockAccounts: Account[] = [
        {
          account_id: 'acc-1',
          user_id: 'user-123',
          name: 'Conta Corrente',
          type: 'BANK',
          balance: '1000.00',
          currency: 'BRL',
          context: 'PERSONAL',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
        {
          account_id: 'acc-2',
          user_id: 'user-123',
          name: 'Carteira',
          type: 'WALLET',
          balance: '500.00',
          currency: 'BRL',
          context: 'PERSONAL',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ]
      
      localStorage.setItem('auth_token', 'test-token')
      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: mockAccounts,
        count: 2,
      })
      
      const store = useAccountsStore()
      const result = await store.listAccounts()
      
      expect(accountsApi.accountService.list).toHaveBeenCalledWith(undefined)
      expect(store.accounts).toEqual(mockAccounts)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result.accounts).toEqual(mockAccounts)
      expect(result.count).toBe(2)
    })

    it('deve filtrar contas por contexto', async () => {
      const mockAccounts: Account[] = [
        {
          account_id: 'acc-1',
          user_id: 'user-123',
          name: 'Conta Pessoal',
          type: 'BANK',
          balance: '1000.00',
          currency: 'BRL',
          context: 'PERSONAL',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ]
      
      localStorage.setItem('auth_token', 'test-token')
      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: mockAccounts,
        count: 1,
      })
      
      const store = useAccountsStore()
      await store.listAccounts('PERSONAL')
      
      expect(accountsApi.accountService.list).toHaveBeenCalledWith('PERSONAL')
    })

    it('deve tratar erro ao listar contas', async () => {
      localStorage.setItem('auth_token', 'test-token')
      const error = new Error('Failed to fetch accounts')
      vi.mocked(accountsApi.accountService.list).mockRejectedValue(error)
      
      const store = useAccountsStore()
      
      await expect(store.listAccounts()).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })

    it('deve lançar erro se não houver token', async () => {
      const store = useAccountsStore()
      
      await expect(store.listAccounts()).rejects.toThrow('Token de autenticação não encontrado')
    })
  })

  describe('getAccount()', () => {
    it('deve obter conta específica com sucesso', async () => {
      const mockAccount: Account = {
        account_id: 'acc-1',
        user_id: 'user-123',
        name: 'Conta Corrente',
        type: 'BANK',
        balance: '1000.00',
        currency: 'BRL',
        context: 'PERSONAL',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(accountsApi.accountService.get).mockResolvedValue(mockAccount)
      
      const store = useAccountsStore()
      const result = await store.getAccount('acc-1')
      
      expect(accountsApi.accountService.get).toHaveBeenCalledWith('acc-1')
      expect(store.currentAccount).toEqual(mockAccount)
      expect(result).toEqual(mockAccount)
      expect(store.isLoading).toBe(false)
    })

    it('deve atualizar conta na lista se já existir', async () => {
      const existingAccount: Account = {
        account_id: 'acc-1',
        user_id: 'user-123',
        name: 'Conta Antiga',
        type: 'BANK',
        balance: '500.00',
        currency: 'BRL',
        context: 'PERSONAL',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      const updatedAccount: Account = {
        ...existingAccount,
        name: 'Conta Atualizada',
        balance: '1500.00',
      }
      
      const store = useAccountsStore()
      store.accounts = [existingAccount]
      
      vi.mocked(accountsApi.accountService.get).mockResolvedValue(updatedAccount)
      
      await store.getAccount('acc-1')
      
      expect(store.accounts[0]).toEqual(updatedAccount)
    })
  })

  describe('createAccount()', () => {
    it('deve criar conta com sucesso', async () => {
      const newAccountData: CreateAccountRequest = {
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 1000.00,
        currency: 'BRL',
        context: 'PERSONAL',
      }
      
      const createdAccount: Account = {
        account_id: 'acc-new',
        user_id: 'user-123',
        name: 'Nova Conta',
        type: 'BANK',
        balance: '1000.00',
        currency: 'BRL',
        context: 'PERSONAL',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(accountsApi.accountService.create).mockResolvedValue(createdAccount)
      
      const store = useAccountsStore()
      const result = await store.createAccount(newAccountData)
      
      expect(accountsApi.accountService.create).toHaveBeenCalledWith(newAccountData)
      expect(store.accounts).toContainEqual(createdAccount)
      expect(result).toEqual(createdAccount)
      expect(store.isLoading).toBe(false)
    })

    it('deve tratar erro ao criar conta', async () => {
      const newAccountData: CreateAccountRequest = {
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 1000.00,
        currency: 'BRL',
        context: 'PERSONAL',
      }
      
      const error = new Error('Failed to create account')
      vi.mocked(accountsApi.accountService.create).mockRejectedValue(error)
      
      const store = useAccountsStore()
      
      await expect(store.createAccount(newAccountData)).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })

  describe('Computed properties', () => {
    it('deve calcular totalAccounts corretamente', () => {
      const store = useAccountsStore()
      store.accounts = [
        { account_id: 'acc-1' } as Account,
        { account_id: 'acc-2' } as Account,
        { account_id: 'acc-3' } as Account,
      ]
      
      expect(store.totalAccounts).toBe(3)
    })

    it('deve filtrar activeAccounts corretamente', () => {
      const store = useAccountsStore()
      store.accounts = [
        { account_id: 'acc-1', is_active: true } as Account,
        { account_id: 'acc-2', is_active: false } as Account,
        { account_id: 'acc-3', is_active: true } as Account,
      ]
      
      expect(store.activeAccounts).toHaveLength(2)
      expect(store.activeAccounts.every(acc => acc.is_active)).toBe(true)
    })

    it('deve filtrar personalAccounts corretamente', () => {
      const store = useAccountsStore()
      store.accounts = [
        { account_id: 'acc-1', context: 'PERSONAL' } as Account,
        { account_id: 'acc-2', context: 'BUSINESS' } as Account,
        { account_id: 'acc-3', context: 'PERSONAL' } as Account,
      ]
      
      expect(store.personalAccounts).toHaveLength(2)
      expect(store.personalAccounts.every(acc => acc.context === 'PERSONAL')).toBe(true)
    })

    it('deve filtrar businessAccounts corretamente', () => {
      const store = useAccountsStore()
      store.accounts = [
        { account_id: 'acc-1', context: 'PERSONAL' } as Account,
        { account_id: 'acc-2', context: 'BUSINESS' } as Account,
        { account_id: 'acc-3', context: 'BUSINESS' } as Account,
      ]
      
      expect(store.businessAccounts).toHaveLength(2)
      expect(store.businessAccounts.every(acc => acc.context === 'BUSINESS')).toBe(true)
    })
  })

  describe('clearError()', () => {
    it('deve limpar erro', () => {
      const store = useAccountsStore()
      store.error = 'Some error'
      
      store.clearError()
      
      expect(store.error).toBeNull()
    })
  })

  describe('clearCurrentAccount()', () => {
    it('deve limpar conta atual', () => {
      const store = useAccountsStore()
      store.currentAccount = { account_id: 'acc-1' } as Account
      
      store.clearCurrentAccount()
      
      expect(store.currentAccount).toBeNull()
    })
  })
})
