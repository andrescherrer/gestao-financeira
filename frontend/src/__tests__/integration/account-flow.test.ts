import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useAccountsStore } from '@/stores/accounts'
import * as authApi from '@/api/auth'
import * as accountsApi from '@/api/accounts'
import { setupIntegrationTests, createMockUser, createMockAccount } from './setup'

// Mock das APIs
vi.mock('@/api/auth', () => ({
  authService: {
    login: vi.fn(),
    saveToken: vi.fn(),
    removeToken: vi.fn(),
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

describe('Fluxo de Contas - Integração', () => {
  setupIntegrationTests()

  beforeEach(async () => {
    // Simular usuário autenticado
    const mockUser = createMockUser()
    const mockToken = 'test-token-123'

    vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)
    localStorage.setItem('auth_token', mockToken)
    localStorage.setItem('auth_user', JSON.stringify(mockUser))

    const authStore = useAuthStore()
    authStore.token = mockToken
    authStore.user = mockUser
    authStore.isValidated = true
  })

  describe('Criação de Conta', () => {
    it('deve criar conta e atualizar a lista automaticamente', async () => {
      const existingAccounts = [createMockAccount()]
      const newAccount = {
        ...createMockAccount(),
        account_id: 'acc-new',
        name: 'Nova Conta',
      }

      // Mock: listar contas existentes
      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: existingAccounts as any,
        count: 1,
      })

      // Mock: criar nova conta
      vi.mocked(accountsApi.accountService.create).mockResolvedValue(newAccount as any)

      const accountsStore = useAccountsStore()

      // Carregar contas iniciais
      await accountsStore.listAccounts()
      expect(accountsStore.accounts).toHaveLength(1)

      // Criar nova conta
      const createdAccount = await accountsStore.createAccount({
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 1000.00,
        currency: 'BRL',
        context: 'PERSONAL',
      })

      // Verificar que a conta foi criada
      expect(createdAccount).toEqual(newAccount)
      expect(accountsApi.accountService.create).toHaveBeenCalledWith({
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 1000.00,
        currency: 'BRL',
        context: 'PERSONAL',
      })

      // Verificar que a lista foi atualizada automaticamente
      expect(accountsStore.accounts).toHaveLength(2)
      expect(accountsStore.accounts).toContainEqual(newAccount)
      expect(accountsStore.isLoading).toBe(false)
    })

    it('deve atualizar propriedades computadas após criar conta', async () => {
      const newAccount = {
        ...createMockAccount(),
        account_id: 'acc-new',
        name: 'Nova Conta',
        is_active: true,
        context: 'PERSONAL' as const,
      }

      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: [],
        count: 0,
      })

      vi.mocked(accountsApi.accountService.create).mockResolvedValue(newAccount as any)

      const accountsStore = useAccountsStore()

      // Criar conta
      await accountsStore.createAccount({
        name: 'Nova Conta',
        type: 'BANK',
        initial_balance: 1000.00,
        currency: 'BRL',
        context: 'PERSONAL',
      })

      // Verificar propriedades computadas
      expect(accountsStore.totalAccounts).toBe(1)
      expect(accountsStore.activeAccounts).toHaveLength(1)
      expect(accountsStore.personalAccounts).toHaveLength(1)
      expect(accountsStore.businessAccounts).toHaveLength(0)
    })

    it('deve tratar erro ao criar conta e manter estado consistente', async () => {
      const existingAccounts = [createMockAccount()]

      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: existingAccounts as any,
        count: 1,
      })

      const error = new Error('Failed to create account')
      vi.mocked(accountsApi.accountService.create).mockRejectedValue(error)

      const accountsStore = useAccountsStore()
      await accountsStore.listAccounts()

      const initialCount = accountsStore.accounts.length

      // Tentar criar conta
      await expect(
        accountsStore.createAccount({
          name: 'Nova Conta',
          type: 'BANK',
          initial_balance: 1000.00,
          currency: 'BRL',
          context: 'PERSONAL',
        })
      ).rejects.toThrow()

      // Verificar que a lista não foi alterada
      expect(accountsStore.accounts).toHaveLength(initialCount)
      expect(accountsStore.error).toBeTruthy()
      expect(accountsStore.isLoading).toBe(false)
    })
  })

  describe('Listagem de Contas', () => {
    it('deve carregar contas e atualizar propriedades computadas', async () => {
      const accounts = [
        createMockAccount(),
        {
          ...createMockAccount(),
          account_id: 'acc-2',
          is_active: false,
          context: 'BUSINESS' as const,
        },
        {
          ...createMockAccount(),
          account_id: 'acc-3',
          is_active: true,
          context: 'BUSINESS' as const,
        },
      ]

      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: accounts as any,
        count: 3,
      })

      const accountsStore = useAccountsStore()

      await accountsStore.listAccounts()

      // Verificar lista
      expect(accountsStore.accounts).toHaveLength(3)
      expect(accountsStore.isLoading).toBe(false)
      expect(accountsStore.error).toBeNull()

      // Verificar propriedades computadas
      expect(accountsStore.totalAccounts).toBe(3)
      expect(accountsStore.activeAccounts).toHaveLength(2)
      expect(accountsStore.personalAccounts).toHaveLength(1)
      expect(accountsStore.businessAccounts).toHaveLength(2)
    })

    it('deve filtrar contas por contexto', async () => {
      const personalAccounts = [
        createMockAccount(),
        {
          ...createMockAccount(),
          account_id: 'acc-2',
          context: 'PERSONAL' as const,
        },
      ]

      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: personalAccounts as any,
        count: 2,
      })

      const accountsStore = useAccountsStore()

      await accountsStore.listAccounts('PERSONAL')

      expect(accountsApi.accountService.list).toHaveBeenCalledWith('PERSONAL')
      expect(accountsStore.personalAccounts).toHaveLength(2)
    })
  })

  describe('Obter Conta Específica', () => {
    it('deve obter conta e atualizar lista se já existir', async () => {
      const account = createMockAccount()
      const updatedAccount = {
        ...account,
        name: 'Conta Atualizada',
        balance: '2000.00',
      }

      vi.mocked(accountsApi.accountService.list).mockResolvedValue({
        accounts: [account] as any,
        count: 1,
      })

      vi.mocked(accountsApi.accountService.get).mockResolvedValue(updatedAccount as any)

      const accountsStore = useAccountsStore()
      await accountsStore.listAccounts()

      // Obter conta específica
      const result = await accountsStore.getAccount(account.account_id)

      expect(result).toEqual(updatedAccount)
      expect(accountsStore.currentAccount).toEqual(updatedAccount)
      expect(accountsStore.accounts[0]).toEqual(updatedAccount)
    })
  })
})
