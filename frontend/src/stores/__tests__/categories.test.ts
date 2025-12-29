import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCategoriesStore } from '../categories'
import * as categoriesApi from '@/api/categories'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/api/types'

// Mock do módulo de API
vi.mock('@/api/categories', () => ({
  categoryService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('Categories Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('Estado inicial', () => {
    it('deve inicializar com estado vazio', () => {
      const store = useCategoriesStore()
      
      expect(store.categories).toEqual([])
      expect(store.currentCategory).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(store.totalCategories).toBe(0)
      expect(store.activeCategories).toEqual([])
      expect(store.inactiveCategories).toEqual([])
    })
  })

  describe('listCategories()', () => {
    it('deve listar categorias com sucesso', async () => {
      const mockCategories: Category[] = [
        {
          category_id: 'cat-1',
          user_id: 'user-123',
          name: 'Alimentação',
          slug: 'alimentacao',
          description: 'Gastos com alimentação',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
        {
          category_id: 'cat-2',
          user_id: 'user-123',
          name: 'Transporte',
          slug: 'transporte',
          description: 'Gastos com transporte',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ]
      
      localStorage.setItem('auth_token', 'test-token')
      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: mockCategories,
        count: 2,
      })
      
      const store = useCategoriesStore()
      const result = await store.listCategories()
      
      expect(categoriesApi.categoryService.list).toHaveBeenCalledWith(undefined)
      expect(store.categories).toEqual(mockCategories)
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
      expect(result.categories).toEqual(mockCategories)
      expect(result.count).toBe(2)
    })

    it('deve filtrar categorias ativas', async () => {
      const mockCategories: Category[] = [
        {
          category_id: 'cat-1',
          user_id: 'user-123',
          name: 'Alimentação',
          slug: 'alimentacao',
          description: 'Gastos com alimentação',
          is_active: true,
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ]
      
      localStorage.setItem('auth_token', 'test-token')
      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: mockCategories,
        count: 1,
      })
      
      const store = useCategoriesStore()
      await store.listCategories(true)
      
      expect(categoriesApi.categoryService.list).toHaveBeenCalledWith(true)
    })

    it('deve tratar erro ao listar categorias', async () => {
      localStorage.setItem('auth_token', 'test-token')
      const error = new Error('Failed to fetch categories')
      vi.mocked(categoriesApi.categoryService.list).mockRejectedValue(error)
      
      const store = useCategoriesStore()
      
      await expect(store.listCategories()).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })

    it('deve garantir que categories é sempre um array', async () => {
      localStorage.setItem('auth_token', 'test-token')
      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: null as any,
        count: 0,
      })
      
      const store = useCategoriesStore()
      await store.listCategories()
      
      expect(Array.isArray(store.categories)).toBe(true)
      expect(store.categories).toEqual([])
    })
  })

  describe('getCategory()', () => {
    it('deve obter categoria específica com sucesso', async () => {
      const mockCategory: Category = {
        category_id: 'cat-1',
        user_id: 'user-123',
        name: 'Alimentação',
        slug: 'alimentacao',
        description: 'Gastos com alimentação',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(categoriesApi.categoryService.get).mockResolvedValue(mockCategory)
      
      const store = useCategoriesStore()
      const result = await store.getCategory('cat-1')
      
      expect(categoriesApi.categoryService.get).toHaveBeenCalledWith('cat-1')
      expect(store.currentCategory).toEqual(mockCategory)
      expect(result).toEqual(mockCategory)
      expect(store.isLoading).toBe(false)
    })

    it('deve atualizar categoria na lista se já existir', async () => {
      const existingCategory: Category = {
        category_id: 'cat-1',
        user_id: 'user-123',
        name: 'Alimentação Antiga',
        slug: 'alimentacao-antiga',
        description: 'Descrição antiga',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      const updatedCategory: Category = {
        ...existingCategory,
        name: 'Alimentação Atualizada',
        description: 'Nova descrição',
      }
      
      const store = useCategoriesStore()
      store.categories = [existingCategory]
      
      vi.mocked(categoriesApi.categoryService.get).mockResolvedValue(updatedCategory)
      
      await store.getCategory('cat-1')
      
      expect(store.categories[0]).toEqual(updatedCategory)
    })
  })

  describe('createCategory()', () => {
    it('deve criar categoria com sucesso', async () => {
      const newCategoryData: CreateCategoryRequest = {
        name: 'Nova Categoria',
        description: 'Descrição da nova categoria',
      }
      
      const createdCategory: Category = {
        category_id: 'cat-new',
        user_id: 'user-123',
        name: 'Nova Categoria',
        slug: 'nova-categoria',
        description: 'Descrição da nova categoria',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      vi.mocked(categoriesApi.categoryService.create).mockResolvedValue(createdCategory)
      
      const store = useCategoriesStore()
      const result = await store.createCategory(newCategoryData)
      
      expect(categoriesApi.categoryService.create).toHaveBeenCalledWith(newCategoryData)
      expect(store.categories).toContainEqual(createdCategory)
      expect(result).toEqual(createdCategory)
      expect(store.isLoading).toBe(false)
    })

    it('deve tratar erro ao criar categoria', async () => {
      const newCategoryData: CreateCategoryRequest = {
        name: 'Nova Categoria',
        description: 'Descrição',
      }
      
      const error = new Error('Failed to create category')
      vi.mocked(categoriesApi.categoryService.create).mockRejectedValue(error)
      
      const store = useCategoriesStore()
      
      await expect(store.createCategory(newCategoryData)).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })

  describe('updateCategory()', () => {
    it('deve atualizar categoria com sucesso', async () => {
      const existingCategory: Category = {
        category_id: 'cat-1',
        user_id: 'user-123',
        name: 'Alimentação',
        slug: 'alimentacao',
        description: 'Descrição antiga',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      const updateData: UpdateCategoryRequest = {
        name: 'Alimentação Atualizada',
        description: 'Nova descrição',
      }
      
      const updatedCategory: Category = {
        ...existingCategory,
        name: 'Alimentação Atualizada',
        description: 'Nova descrição',
      }
      
      const store = useCategoriesStore()
      store.categories = [existingCategory]
      store.currentCategory = existingCategory
      
      vi.mocked(categoriesApi.categoryService.update).mockResolvedValue(updatedCategory)
      
      const result = await store.updateCategory('cat-1', updateData)
      
      expect(categoriesApi.categoryService.update).toHaveBeenCalledWith('cat-1', updateData)
      expect(store.categories[0]).toEqual(updatedCategory)
      expect(store.currentCategory).toEqual(updatedCategory)
      expect(result).toEqual(updatedCategory)
      expect(store.isLoading).toBe(false)
    })
  })

  describe('deleteCategory()', () => {
    it('deve deletar categoria com sucesso', async () => {
      const categoryToDelete: Category = {
        category_id: 'cat-1',
        user_id: 'user-123',
        name: 'Alimentação',
        slug: 'alimentacao',
        description: 'Descrição',
        is_active: true,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }
      
      const store = useCategoriesStore()
      store.categories = [categoryToDelete]
      store.currentCategory = categoryToDelete
      
      vi.mocked(categoriesApi.categoryService.delete).mockResolvedValue(undefined)
      
      await store.deleteCategory('cat-1')
      
      expect(categoriesApi.categoryService.delete).toHaveBeenCalledWith('cat-1')
      expect(store.categories).not.toContainEqual(categoryToDelete)
      expect(store.currentCategory).toBeNull()
      expect(store.isLoading).toBe(false)
    })

    it('deve tratar erro ao deletar categoria', async () => {
      const error = new Error('Failed to delete category')
      vi.mocked(categoriesApi.categoryService.delete).mockRejectedValue(error)
      
      const store = useCategoriesStore()
      
      await expect(store.deleteCategory('cat-1')).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })

  describe('Computed properties', () => {
    it('deve calcular totalCategories corretamente', () => {
      const store = useCategoriesStore()
      store.categories = [
        { category_id: 'cat-1' } as Category,
        { category_id: 'cat-2' } as Category,
      ]
      
      expect(store.totalCategories).toBe(2)
    })

    it('deve filtrar activeCategories corretamente', () => {
      const store = useCategoriesStore()
      store.categories = [
        { category_id: 'cat-1', is_active: true } as Category,
        { category_id: 'cat-2', is_active: false } as Category,
        { category_id: 'cat-3', is_active: true } as Category,
      ]
      
      expect(store.activeCategories).toHaveLength(2)
      expect(store.activeCategories.every(cat => cat.is_active)).toBe(true)
    })

    it('deve filtrar inactiveCategories corretamente', () => {
      const store = useCategoriesStore()
      store.categories = [
        { category_id: 'cat-1', is_active: true } as Category,
        { category_id: 'cat-2', is_active: false } as Category,
        { category_id: 'cat-3', is_active: false } as Category,
      ]
      
      expect(store.inactiveCategories).toHaveLength(2)
      expect(store.inactiveCategories.every(cat => !cat.is_active)).toBe(true)
    })
  })
})
