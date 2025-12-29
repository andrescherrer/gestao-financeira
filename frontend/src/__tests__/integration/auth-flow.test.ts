import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useAccountsStore } from '@/stores/accounts'
import * as authApi from '@/api/auth'
import { setupIntegrationTests, createMockUser } from './setup'

// Mock das APIs
vi.mock('@/api/auth', () => ({
  authService: {
    login: vi.fn(),
    register: vi.fn(),
    saveToken: vi.fn(),
    removeToken: vi.fn(),
    getToken: vi.fn(),
  },
}))

vi.mock('@/api/accounts', () => ({
  accountService: {
    list: vi.fn(),
  },
}))

describe('Fluxo de Autenticação - Integração', () => {
  setupIntegrationTests()

  describe('Login Flow', () => {
    it('deve fazer login completo e atualizar todas as stores', async () => {
      const mockUser = createMockUser()
      const mockToken = 'test-token-123'
      const mockAccounts = [
        {
          account_id: 'acc-1',
          name: 'Conta Teste',
          currency: 'BRL',
        },
      ]

      // Mock do login
      vi.mocked(authApi.authService.login).mockResolvedValue({
        token: mockToken,
        user: mockUser,
      })

      // Mock do getToken
      vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)

      // Mock do listAccounts (usado na validação)
      const { accountService } = await import('@/api/accounts')
      vi.mocked(accountService.list).mockResolvedValue({
        accounts: mockAccounts as any,
        count: 1,
      })

      const authStore = useAuthStore()
      const accountsStore = useAccountsStore()

      // Estado inicial
      expect(authStore.isAuthenticated).toBe(false)
      expect(authStore.token).toBeNull()
      expect(authStore.user).toBeNull()

      // Fazer login
      const loginResult = await authStore.login({
        email: 'test@example.com',
        password: 'password123',
      })

      // Verificar resultado do login
      expect(loginResult.token).toBe(mockToken)
      expect(loginResult.user).toEqual(mockUser)
      expect(authApi.authService.login).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      })

      // Verificar que o token foi salvo
      expect(authApi.authService.saveToken).toHaveBeenCalledWith(mockToken)

      // Verificar estado da auth store
      expect(authStore.token).toBe(mockToken)
      expect(authStore.user).toEqual(mockUser)
      expect(authStore.isValidated).toBe(true)
      expect(authStore.isAuthenticated).toBe(true)
      expect(authStore.isLoading).toBe(false)

      // Verificar que saveToken foi chamado (que salva no localStorage)
      // O saveToken é mockado, então verificamos que foi chamado
      expect(authApi.authService.saveToken).toHaveBeenCalledWith(mockToken)
    })

    it('deve permitir carregar dados após login bem-sucedido', async () => {
      const mockUser = createMockUser()
      const mockToken = 'test-token-123'
      const mockAccounts = [
        {
          account_id: 'acc-1',
          name: 'Conta Teste',
          currency: 'BRL',
        },
      ]

      // Mock do login
      vi.mocked(authApi.authService.login).mockResolvedValue({
        token: mockToken,
        user: mockUser,
      })

      vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)

      // Mock do listAccounts
      const { accountService } = await import('@/api/accounts')
      vi.mocked(accountService.list).mockResolvedValue({
        accounts: mockAccounts as any,
        count: 1,
      })

      const authStore = useAuthStore()
      const accountsStore = useAccountsStore()

      // Fazer login
      await authStore.login({
        email: 'test@example.com',
        password: 'password123',
      })

      // Verificar que após login, o token está validado
      expect(authStore.isValidated).toBe(true)
      expect(authStore.token).toBe(mockToken)
      expect(authStore.isAuthenticated).toBe(true)

      // Garantir que o token está no localStorage (saveToken faz isso, mas está mockado)
      // Vamos simular que foi salvo
      localStorage.setItem('auth_token', mockToken)

      // Verificar que podemos carregar dados após autenticação
      // Isso demonstra que a integração entre auth e accounts funciona
      await accountsStore.listAccounts()
      
      // Verificar que conseguimos carregar contas após autenticação
      expect(accountsStore.accounts.length).toBeGreaterThan(0)
      expect(accountsStore.accounts[0].account_id).toBe('acc-1')
    })

    it('deve tratar erro de login e limpar estado', async () => {
      const error = new Error('Invalid credentials')
      vi.mocked(authApi.authService.login).mockRejectedValue(error)

      const authStore = useAuthStore()

      // Tentar fazer login
      await expect(
        authStore.login({
          email: 'wrong@example.com',
          password: 'wrong-password',
        })
      ).rejects.toThrow('Invalid credentials')

      // Verificar que o estado foi limpo
      expect(authStore.token).toBeNull()
      expect(authStore.user).toBeNull()
      expect(authStore.isValidated).toBe(false)
      expect(authStore.isAuthenticated).toBe(false)
      expect(authApi.authService.removeToken).toHaveBeenCalled()
    })
  })

  describe('Logout Flow', () => {
    it('deve fazer logout completo e limpar todas as stores', async () => {
      const mockUser = createMockUser()
      const mockToken = 'test-token-123'

      // Simular usuário logado
      vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)
      localStorage.setItem('auth_token', mockToken)
      localStorage.setItem('auth_user', JSON.stringify(mockUser))

      const authStore = useAuthStore()
      const accountsStore = useAccountsStore()
      const { useTransactionsStore } = await import('@/stores/transactions')
      const { useCategoriesStore } = await import('@/stores/categories')

      const transactionsStore = useTransactionsStore()
      const categoriesStore = useCategoriesStore()

      // Simular dados nas stores
      authStore.token = mockToken
      authStore.user = mockUser
      authStore.isValidated = true
      accountsStore.accounts = [{ account_id: 'acc-1' } as any]
      transactionsStore.transactions = [{ transaction_id: 'tx-1' } as any]
      categoriesStore.categories = [{ category_id: 'cat-1' } as any]

      // Fazer logout
      authStore.logout()

      // Aguardar um pouco para que as promises de limpeza sejam executadas
      await new Promise(resolve => setTimeout(resolve, 100))

      // Verificar que removeToken foi chamado
      expect(authApi.authService.removeToken).toHaveBeenCalled()

      // Verificar que a auth store foi limpa
      expect(authStore.token).toBeNull()
      expect(authStore.user).toBeNull()
      expect(authStore.isValidated).toBe(false)
      expect(authStore.isAuthenticated).toBe(false)

      // Verificar que outras stores foram limpas (o logout limpa através de promises)
      // Pode levar um momento para as promises serem resolvidas
      expect(accountsStore.accounts).toEqual([])
      expect(transactionsStore.transactions).toEqual([])
      expect(categoriesStore.categories).toEqual([])

      // Verificar localStorage (removeToken limpa isso)
      expect(authApi.authService.removeToken).toHaveBeenCalled()
    })
  })

  describe('Register Flow', () => {
    it('deve registrar novo usuário com sucesso', async () => {
      const mockResponse = {
        message: 'User registered successfully',
        data: {
          user_id: 'user-new',
          email: 'newuser@example.com',
          first_name: 'New',
          last_name: 'User',
          full_name: 'New User',
        },
      }

      vi.mocked(authApi.authService.register).mockResolvedValue(mockResponse)

      const authStore = useAuthStore()

      const result = await authStore.register({
        email: 'newuser@example.com',
        password: 'password123',
        first_name: 'New',
        last_name: 'User',
      })

      expect(authApi.authService.register).toHaveBeenCalledWith({
        email: 'newuser@example.com',
        password: 'password123',
        first_name: 'New',
        last_name: 'User',
      })

      expect(result).toEqual(mockResponse)
      expect(authStore.isLoading).toBe(false)
    })
  })
})
