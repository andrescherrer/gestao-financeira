import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'
import * as authApi from '@/api/auth'
import type { LoginRequest, RegisterRequest, User } from '@/api/types'

// Mock do módulo de API
vi.mock('@/api/auth', () => ({
  authService: {
    login: vi.fn(),
    register: vi.fn(),
    saveToken: vi.fn(),
    removeToken: vi.fn(),
    getToken: vi.fn(),
  },
}))

// Mock do módulo de accounts store (usado em validateToken)
vi.mock('@/stores/accounts', () => ({
  useAccountsStore: vi.fn(() => ({
    accounts: [],
    isLoading: false,
    listAccounts: vi.fn().mockResolvedValue({ accounts: [], count: 0 }),
  })),
}))

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    sessionStorage.clear()
    vi.clearAllMocks()
  })

  describe('Estado inicial', () => {
    it('deve inicializar com estado vazio', () => {
      const store = useAuthStore()
      
      expect(store.user).toBeNull()
      expect(store.token).toBeNull()
      expect(store.isAuthenticated).toBe(false)
      expect(store.isLoading).toBe(false)
      expect(store.isValidating).toBe(false)
      expect(store.isValidated).toBe(false)
    })
  })

  describe('init()', () => {
    it('deve carregar token e usuário do localStorage', () => {
      const mockToken = 'test-token'
      const mockUser: User = {
        user_id: 'user-123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        full_name: 'Test User',
      }
      
      localStorage.setItem('auth_token', mockToken)
      localStorage.setItem('auth_user', JSON.stringify(mockUser))
      
      vi.mocked(authApi.authService.getToken).mockReturnValue(mockToken)
      
      const store = useAuthStore()
      store.init()
      
      expect(store.token).toBe(mockToken)
      expect(store.user).toEqual(mockUser)
    })

    it('deve limpar estado se não houver token', () => {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      vi.mocked(authApi.authService.getToken).mockReturnValue(null)
      
      const store = useAuthStore()
      store.init()
      
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isValidated).toBe(false)
    })

    it('deve tratar erro ao parsear usuário do localStorage', () => {
      localStorage.setItem('auth_token', 'test-token')
      localStorage.setItem('auth_user', 'invalid-json')
      
      vi.mocked(authApi.authService.getToken).mockReturnValue('test-token')
      
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
      
      const store = useAuthStore()
      store.init()
      
      expect(store.user).toBeNull()
      expect(consoleSpy).toHaveBeenCalled()
      
      consoleSpy.mockRestore()
    })
  })

  describe('login()', () => {
    it('deve fazer login com sucesso', async () => {
      const mockToken = 'test-token'
      const mockUser: User = {
        user_id: 'user-123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        full_name: 'Test User',
      }
      
      const credentials: LoginRequest = {
        email: 'test@example.com',
        password: 'password123',
      }
      
      vi.mocked(authApi.authService.login).mockResolvedValue({
        token: mockToken,
        user: mockUser,
      })
      
      const store = useAuthStore()
      const result = await store.login(credentials)
      
      expect(authApi.authService.login).toHaveBeenCalledWith(credentials)
      expect(authApi.authService.saveToken).toHaveBeenCalledWith(mockToken)
      expect(store.token).toBe(mockToken)
      expect(store.user).toEqual(mockUser)
      expect(store.isValidated).toBe(true)
      expect(store.isLoading).toBe(false)
      expect(result.token).toBe(mockToken)
      expect(result.user).toEqual(mockUser)
    })

    it('deve tratar erro no login', async () => {
      const credentials: LoginRequest = {
        email: 'test@example.com',
        password: 'wrong-password',
      }
      
      const error = new Error('Invalid credentials')
      vi.mocked(authApi.authService.login).mockRejectedValue(error)
      
      const store = useAuthStore()
      
      await expect(store.login(credentials)).rejects.toThrow('Invalid credentials')
      
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isValidated).toBe(false)
      expect(authApi.authService.removeToken).toHaveBeenCalled()
    })

    it('deve definir isLoading durante o login', async () => {
      const mockToken = 'test-token'
      const mockUser: User = {
        user_id: 'user-123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        full_name: 'Test User',
      }
      
      vi.mocked(authApi.authService.login).mockImplementation(
        () => new Promise(resolve => setTimeout(() => resolve({
          token: mockToken,
          user: mockUser,
        }), 100))
      )
      
      const store = useAuthStore()
      const loginPromise = store.login({
        email: 'test@example.com',
        password: 'password123',
      })
      
      expect(store.isLoading).toBe(true)
      
      await loginPromise
      
      expect(store.isLoading).toBe(false)
    })
  })

  describe('register()', () => {
    it('deve registrar novo usuário com sucesso', async () => {
      const userData: RegisterRequest = {
        email: 'newuser@example.com',
        password: 'password123',
        first_name: 'New',
        last_name: 'User',
      }
      
      const mockResponse = {
        message: 'User registered successfully',
        data: {
          user_id: 'user-456',
          email: 'newuser@example.com',
          first_name: 'New',
          last_name: 'User',
          full_name: 'New User',
        },
      }
      
      vi.mocked(authApi.authService.register).mockResolvedValue(mockResponse)
      
      const store = useAuthStore()
      const result = await store.register(userData)
      
      expect(authApi.authService.register).toHaveBeenCalledWith(userData)
      expect(result).toEqual(mockResponse)
      expect(store.isLoading).toBe(false)
    })

    it('deve tratar erro no registro', async () => {
      const userData: RegisterRequest = {
        email: 'existing@example.com',
        password: 'password123',
        first_name: 'Existing',
        last_name: 'User',
      }
      
      const error = new Error('Email already exists')
      vi.mocked(authApi.authService.register).mockRejectedValue(error)
      
      const store = useAuthStore()
      
      await expect(store.register(userData)).rejects.toThrow('Email already exists')
      expect(store.isLoading).toBe(false)
    })
  })

  describe('logout()', () => {
    it('deve fazer logout e limpar estado', () => {
      const store = useAuthStore()
      
      // Simular usuário logado
      store.token = 'test-token'
      store.user = {
        user_id: 'user-123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        full_name: 'Test User',
      }
      store.isValidated = true
      
      localStorage.setItem('auth_token', 'test-token')
      localStorage.setItem('auth_user', JSON.stringify(store.user))
      
      store.logout()
      
      expect(authApi.authService.removeToken).toHaveBeenCalled()
      expect(store.token).toBeNull()
      expect(store.user).toBeNull()
      expect(store.isValidated).toBe(false)
      // Verificar que removeToken foi chamado (que limpa o localStorage)
      expect(authApi.authService.removeToken).toHaveBeenCalled()
    })
  })

  describe('isAuthenticated', () => {
    it('deve retornar false quando não há token', () => {
      const store = useAuthStore()
      expect(store.isAuthenticated).toBe(false)
    })

    it('deve retornar false quando há token mas não foi validado', () => {
      const store = useAuthStore()
      store.token = 'test-token'
      store.isValidated = false
      
      expect(store.isAuthenticated).toBe(false)
    })

    it('deve retornar true quando há token e foi validado', () => {
      const store = useAuthStore()
      store.token = 'test-token'
      store.isValidated = true
      
      expect(store.isAuthenticated).toBe(true)
    })
  })
})
