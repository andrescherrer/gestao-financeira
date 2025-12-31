import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useInvestmentsStore } from '../investments'
import { useAuthStore } from '../auth'
import * as investmentService from '@/api/investments'
import type { Investment, CreateInvestmentRequest, UpdateInvestmentRequest } from '@/api/types'

// Mock do investmentService
vi.mock('@/api/investments', () => ({
  investmentService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}))

// Mock do errorTranslations
vi.mock('@/utils/errorTranslations', () => ({
  extractErrorMessage: (err: any) => err.message || 'Erro desconhecido',
}))

describe('Investments Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  const mockInvestment: Investment = {
    investment_id: '123e4567-e89b-12d3-a456-426614174000',
    user_id: 'user-123',
    account_id: 'account-123',
    type: 'STOCK',
    name: 'Petrobras',
    ticker: 'PETR4',
    purchase_date: '2024-01-15',
    purchase_amount: '1000.00',
    current_value: '1200.00',
    currency: 'BRL',
    quantity: '100',
    context: 'PERSONAL',
    return_absolute: '200.00',
    return_percentage: 20.0,
    created_at: '2024-01-15T10:00:00Z',
    updated_at: '2024-01-15T10:00:00Z',
  }

  describe('listInvestments', () => {
    it('should list investments successfully', async () => {
      localStorage.setItem('auth_token', 'test-token')
      const store = useInvestmentsStore()
      const mockResponse = {
        investments: [mockInvestment],
        count: 1,
      }

      vi.mocked(investmentService.investmentService.list).mockResolvedValue(mockResponse)

      const result = await store.listInvestments()

      expect(store.investments).toHaveLength(1)
      expect(store.investments[0]).toEqual(mockInvestment)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result).toEqual(mockResponse)
    })

    it('should handle errors when listing investments', async () => {
      localStorage.setItem('auth_token', 'test-token')
      const store = useInvestmentsStore()
      const error = new Error('Failed to list investments')

      vi.mocked(investmentService.investmentService.list).mockRejectedValue(error)

      await expect(store.listInvestments()).rejects.toThrow()
      expect(store.error).toBe('Failed to list investments')
      expect(store.isLoading).toBe(false)
    })
  })

  describe('getInvestment', () => {
    it('should get investment successfully', async () => {
      const store = useInvestmentsStore()

      vi.mocked(investmentService.investmentService.get).mockResolvedValue(mockInvestment)

      const result = await store.getInvestment(mockInvestment.investment_id)

      expect(store.currentInvestment).toEqual(mockInvestment)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result).toEqual(mockInvestment)
    })

    it('should update investment in list if exists', async () => {
      const store = useInvestmentsStore()
      store.investments = [mockInvestment]

      const updatedInvestment = { ...mockInvestment, current_value: '1300.00' }
      vi.mocked(investmentService.investmentService.get).mockResolvedValue(updatedInvestment)

      await store.getInvestment(mockInvestment.investment_id)

      expect(store.investments[0].current_value).toBe('1300.00')
    })
  })

  describe('createInvestment', () => {
    it('should create investment successfully', async () => {
      const store = useInvestmentsStore()
      const createData: CreateInvestmentRequest = {
        account_id: 'account-123',
        type: 'STOCK',
        name: 'Petrobras',
        ticker: 'PETR4',
        purchase_date: '2024-01-15',
        purchase_amount: 1000.0,
        currency: 'BRL',
        quantity: 100,
        context: 'PERSONAL',
      }

      vi.mocked(investmentService.investmentService.create).mockResolvedValue(mockInvestment)

      const result = await store.createInvestment(createData)

      expect(store.investments).toHaveLength(1)
      expect(store.investments[0]).toEqual(mockInvestment)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result).toEqual(mockInvestment)
    })
  })

  describe('updateInvestment', () => {
    it('should update investment successfully', async () => {
      const store = useInvestmentsStore()
      store.investments = [mockInvestment]
      store.currentInvestment = mockInvestment

      const updateData: UpdateInvestmentRequest = {
        current_value: 1300.0,
      }

      const updatedInvestment = {
        ...mockInvestment,
        current_value: '1300.00',
        return_absolute: '300.00',
        return_percentage: 30.0,
      }

      vi.mocked(investmentService.investmentService.update).mockResolvedValue(updatedInvestment)

      const result = await store.updateInvestment(mockInvestment.investment_id, updateData)

      expect(store.investments[0].current_value).toBe('1300.00')
      expect(store.currentInvestment?.current_value).toBe('1300.00')
      expect(result).toEqual(updatedInvestment)
    })
  })

  describe('deleteInvestment', () => {
    it('should delete investment successfully', async () => {
      const store = useInvestmentsStore()
      store.investments = [mockInvestment]
      store.currentInvestment = mockInvestment

      vi.mocked(investmentService.investmentService.delete).mockResolvedValue()

      await store.deleteInvestment(mockInvestment.investment_id)

      expect(store.investments).toHaveLength(0)
      expect(store.currentInvestment).toBeNull()
      expect(store.isLoading).toBe(false)
    })
  })

  describe('computed properties', () => {
    it('should calculate totalInvestments correctly', () => {
      const store = useInvestmentsStore()
      store.investments = [mockInvestment, { ...mockInvestment, investment_id: 'another-id' }]

      expect(store.totalInvestments).toBe(2)
    })

    it('should filter personal investments correctly', () => {
      const store = useInvestmentsStore()
      store.investments = [
        mockInvestment,
        { ...mockInvestment, investment_id: 'another-id', context: 'BUSINESS' },
      ]

      expect(store.personalInvestments).toHaveLength(1)
      expect(store.personalInvestments[0].context).toBe('PERSONAL')
    })

    it('should calculate totalValue correctly', () => {
      const store = useInvestmentsStore()
      store.investments = [
        mockInvestment,
        { ...mockInvestment, investment_id: 'another-id', current_value: '500.00' },
      ]

      expect(store.totalValue).toBe(1700.0)
    })

    it('should calculate totalReturn correctly', () => {
      const store = useInvestmentsStore()
      store.investments = [
        mockInvestment,
        { ...mockInvestment, investment_id: 'another-id', return_absolute: '-50.00' },
      ]

      expect(store.totalReturn).toBe(150.0)
    })
  })

  describe('clearError', () => {
    it('should clear error', () => {
      const store = useInvestmentsStore()
      store.error = 'Some error'

      store.clearError()

      expect(store.error).toBeNull()
    })
  })

  describe('clearCurrentInvestment', () => {
    it('should clear current investment', () => {
      const store = useInvestmentsStore()
      store.currentInvestment = mockInvestment

      store.clearCurrentInvestment()

      expect(store.currentInvestment).toBeNull()
    })
  })
})

