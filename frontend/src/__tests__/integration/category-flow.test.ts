import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useCategoriesStore } from '@/stores/categories'
import * as authApi from '@/api/auth'
import * as categoriesApi from '@/api/categories'
import { setupIntegrationTests, createMockUser, createMockCategory } from './setup'

// Mock das APIs
vi.mock('@/api/auth', () => ({
  authService: {
    getToken: vi.fn(),
  },
}))

vi.mock('@/api/categories', () => ({
  categoryService: {
    list: vi.fn(),
    get: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('Fluxo de Categorias - Integração', () => {
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

  describe('Criação de Categoria', () => {
    it('deve criar categoria e atualizar a lista automaticamente', async () => {
      const existingCategories: any[] = []
      const newCategory = createMockCategory()

      // Mock: listar categorias existentes
      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: existingCategories,
        count: 0,
      })

      // Mock: criar nova categoria
      vi.mocked(categoriesApi.categoryService.create).mockResolvedValue(newCategory as any)

      const categoriesStore = useCategoriesStore()

      // Carregar categorias iniciais
      await categoriesStore.listCategories()
      expect(categoriesStore.categories).toHaveLength(0)

      // Criar nova categoria
      const createdCategory = await categoriesStore.createCategory({
        name: 'Alimentação',
        description: 'Gastos com alimentação',
      })

      // Verificar que a categoria foi criada
      expect(createdCategory).toEqual(newCategory)
      expect(categoriesApi.categoryService.create).toHaveBeenCalledWith({
        name: 'Alimentação',
        description: 'Gastos com alimentação',
      })

      // Verificar que a lista foi atualizada automaticamente
      expect(categoriesStore.categories).toHaveLength(1)
      expect(categoriesStore.categories).toContainEqual(newCategory)
      expect(categoriesStore.isLoading).toBe(false)
    })

    it('deve atualizar propriedades computadas após criar categoria', async () => {
      const activeCategory = createMockCategory()
      const inactiveCategory = {
        ...createMockCategory(),
        category_id: 'cat-2',
        name: 'Categoria Inativa',
        is_active: false,
      }

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: [],
        count: 0,
      })

      vi.mocked(categoriesApi.categoryService.create)
        .mockResolvedValueOnce(activeCategory as any)
        .mockResolvedValueOnce(inactiveCategory as any)

      const categoriesStore = useCategoriesStore()

      // Criar categoria ativa
      await categoriesStore.createCategory({
        name: 'Alimentação',
        description: 'Gastos com alimentação',
      })

      // Criar categoria inativa
      await categoriesStore.createCategory({
        name: 'Categoria Inativa',
        description: 'Descrição',
      })

      // Verificar propriedades computadas
      expect(categoriesStore.totalCategories).toBe(2)
      expect(categoriesStore.activeCategories).toHaveLength(1)
      expect(categoriesStore.inactiveCategories).toHaveLength(1)
    })

    it('deve tratar erro ao criar categoria e manter estado consistente', async () => {
      const existingCategories = [createMockCategory()]

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: existingCategories as any,
        count: 1,
      })

      const error = new Error('Failed to create category')
      vi.mocked(categoriesApi.categoryService.create).mockRejectedValue(error)

      const categoriesStore = useCategoriesStore()
      await categoriesStore.listCategories()

      const initialCount = categoriesStore.categories.length

      // Tentar criar categoria
      await expect(
        categoriesStore.createCategory({
          name: 'Nova Categoria',
          description: 'Descrição',
        })
      ).rejects.toThrow()

      // Verificar que a lista não foi alterada
      expect(categoriesStore.categories).toHaveLength(initialCount)
      expect(categoriesStore.error).toBeTruthy()
      expect(categoriesStore.isLoading).toBe(false)
    })
  })

  describe('Atualização de Categoria', () => {
    it('deve atualizar categoria e atualizar lista e currentCategory', async () => {
      const existingCategory = createMockCategory()
      const updatedCategory = {
        ...existingCategory,
        name: 'Alimentação Atualizada',
        description: 'Nova descrição',
      }

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: [existingCategory] as any,
        count: 1,
      })

      vi.mocked(categoriesApi.categoryService.update).mockResolvedValue(updatedCategory as any)

      const categoriesStore = useCategoriesStore()
      await categoriesStore.listCategories()
      categoriesStore.currentCategory = existingCategory

      // Atualizar categoria
      const result = await categoriesStore.updateCategory(existingCategory.category_id, {
        name: 'Alimentação Atualizada',
        description: 'Nova descrição',
      })

      // Verificar resultado
      expect(result).toEqual(updatedCategory)
      expect(categoriesApi.categoryService.update).toHaveBeenCalledWith(
        existingCategory.category_id,
        {
          name: 'Alimentação Atualizada',
          description: 'Nova descrição',
        }
      )

      // Verificar que lista e currentCategory foram atualizados
      expect(categoriesStore.categories[0]).toEqual(updatedCategory)
      expect(categoriesStore.currentCategory).toEqual(updatedCategory)
    })
  })

  describe('Deleção de Categoria', () => {
    it('deve deletar categoria e remover da lista', async () => {
      const categoryToDelete = createMockCategory()
      const otherCategory = {
        ...createMockCategory(),
        category_id: 'cat-2',
        name: 'Outra Categoria',
      }

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: [categoryToDelete, otherCategory] as any,
        count: 2,
      })

      vi.mocked(categoriesApi.categoryService.delete).mockResolvedValue(undefined)

      const categoriesStore = useCategoriesStore()
      await categoriesStore.listCategories()
      categoriesStore.currentCategory = categoryToDelete

      // Deletar categoria
      await categoriesStore.deleteCategory(categoryToDelete.category_id)

      // Verificar que foi removida da lista
      expect(categoriesStore.categories).toHaveLength(1)
      expect(categoriesStore.categories).not.toContainEqual(categoryToDelete)
      expect(categoriesStore.categories).toContainEqual(otherCategory)

      // Verificar que currentCategory foi limpo
      expect(categoriesStore.currentCategory).toBeNull()
    })
  })

  describe('Listagem de Categorias', () => {
    it('deve carregar categorias e atualizar propriedades computadas', async () => {
      const categories = [
        createMockCategory(),
        {
          ...createMockCategory(),
          category_id: 'cat-2',
          is_active: false,
        },
        {
          ...createMockCategory(),
          category_id: 'cat-3',
          is_active: true,
        },
      ]

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: categories as any,
        count: 3,
      })

      const categoriesStore = useCategoriesStore()

      await categoriesStore.listCategories()

      // Verificar lista
      expect(categoriesStore.categories).toHaveLength(3)
      expect(categoriesStore.isLoading).toBe(false)
      expect(categoriesStore.error).toBeNull()

      // Verificar propriedades computadas
      expect(categoriesStore.totalCategories).toBe(3)
      expect(categoriesStore.activeCategories).toHaveLength(2)
      expect(categoriesStore.inactiveCategories).toHaveLength(1)
    })

    it('deve filtrar categorias ativas', async () => {
      const activeCategories = [createMockCategory()]

      vi.mocked(categoriesApi.categoryService.list).mockResolvedValue({
        categories: activeCategories as any,
        count: 1,
      })

      const categoriesStore = useCategoriesStore()

      await categoriesStore.listCategories(true)

      expect(categoriesApi.categoryService.list).toHaveBeenCalledWith(true)
    })
  })
})
