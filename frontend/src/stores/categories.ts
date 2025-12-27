import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { categoryService } from '@/api/categories'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '@/api/types'

export const useCategoriesStore = defineStore('categories', () => {
  // Estado
  const categories = ref<Category[]>([])
  const currentCategory = ref<Category | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalCategories = computed(() => categories.value.length)
  const activeCategories = computed(() =>
    categories.value.filter((category) => category.is_active)
  )
  const inactiveCategories = computed(() =>
    categories.value.filter((category) => !category.is_active)
  )

  /**
   * Lista todas as categorias do usuário
   */
  async function listCategories(isActive?: boolean) {
    isLoading.value = true
    error.value = null
    try {
      const token = localStorage.getItem('auth_token')
      if (!token) {
        throw new Error('Token de autenticação não encontrado. Faça login novamente.')
      }
      
      const response = await categoryService.list(isActive)
      // Garantir que categories é sempre um array
      categories.value = Array.isArray(response.categories) ? response.categories : []
      return { 
        categories: Array.isArray(response.categories) ? response.categories : [], 
        count: response.count || 0 
      }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error || 
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao listar categorias'
      error.value = extractErrorMessage(rawError)
      
      if (import.meta.env.DEV) {
        console.error('[Categories Store] Erro ao listar categorias:', {
          message: error.value,
          status: err.response?.status,
          statusText: err.response?.statusText,
          data: err.response?.data,
        })
      }
      
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Obtém detalhes de uma categoria específica
   */
  async function getCategory(categoryId: string) {
    isLoading.value = true
    error.value = null
    try {
      const category = await categoryService.get(categoryId)
      currentCategory.value = category

      // Atualiza a categoria na lista se já existir
      const index = categories.value.findIndex(
        (cat) => cat.category_id === categoryId
      )
      if (index !== -1) {
        categories.value[index] = category
      }

      return category
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao obter detalhes da categoria'
      error.value = extractErrorMessage(rawError)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Cria uma nova categoria
   */
  async function createCategory(data: CreateCategoryRequest) {
    isLoading.value = true
    error.value = null
    try {
      const category = await categoryService.create(data)
      categories.value.push(category)
      return category
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao criar categoria'
      error.value = extractErrorMessage(rawError)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Atualiza uma categoria existente
   */
  async function updateCategory(categoryId: string, data: UpdateCategoryRequest) {
    isLoading.value = true
    error.value = null
    try {
      const updatedCategory = await categoryService.update(categoryId, data)
      
      // Atualiza a categoria na lista
      const index = categories.value.findIndex(
        (cat) => cat.category_id === categoryId
      )
      if (index !== -1) {
        categories.value[index] = updatedCategory
      }
      
      // Atualiza currentCategory se for a mesma
      if (currentCategory.value?.category_id === categoryId) {
        currentCategory.value = updatedCategory
      }
      
      return updatedCategory
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao atualizar categoria'
      error.value = extractErrorMessage(rawError)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Deleta uma categoria
   */
  async function deleteCategory(categoryId: string) {
    isLoading.value = true
    error.value = null
    try {
      await categoryService.delete(categoryId)
      
      // Remove da lista
      categories.value = categories.value.filter(
        (cat) => cat.category_id !== categoryId
      )
      
      // Limpa currentCategory se for a mesma
      if (currentCategory.value?.category_id === categoryId) {
        currentCategory.value = null
      }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      const rawError = err.response?.data?.error ||
                       err.response?.data?.message ||
                       err.message ||
                       'Erro ao deletar categoria'
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

  function clearCurrentCategory() {
    currentCategory.value = null
  }

  return {
    // Estado
    categories,
    currentCategory,
    isLoading,
    error,
    // Computed
    totalCategories,
    activeCategories,
    inactiveCategories,
    // Ações
    listCategories,
    getCategory,
    createCategory,
    updateCategory,
    deleteCategory,
    clearError,
    clearCurrentCategory,
  }
})

