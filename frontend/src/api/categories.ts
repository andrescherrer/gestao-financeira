import apiClient from './client'
import type {
  Category,
  CreateCategoryRequest,
  ListCategoriesResponse,
  UpdateCategoryRequest,
} from './types'

/**
 * Serviço de API para categorias
 */
export const categoryService = {
  /**
   * Lista todas as categorias do usuário autenticado
   */
  async list(isActive?: boolean): Promise<ListCategoriesResponse> {
    const params: Record<string, string> = {}
    if (isActive !== undefined) {
      params.is_active = isActive.toString()
    }

    const response = await apiClient.get<{
      message: string
      data: {
        categories: Array<{
          category_id: string
          user_id: string
          name: string
          description: string
          is_active: boolean
          created_at: string
          updated_at: string
        }>
        count: number
      }
    }>('/categories', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      categories: backendData.categories.map((cat) => ({
        category_id: cat.category_id,
        user_id: cat.user_id,
        name: cat.name,
        description: cat.description,
        is_active: cat.is_active,
        created_at: cat.created_at,
        updated_at: cat.updated_at,
      })),
      count: backendData.count,
    }
  },

  /**
   * Obtém detalhes de uma categoria específica
   */
  async get(categoryId: string): Promise<Category> {
    const response = await apiClient.get<{
      message: string
      data: {
        category_id: string
        user_id: string
        name: string
        description: string
        is_active: boolean
        created_at: string
        updated_at: string
      }
    }>(`/categories/${categoryId}`)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      category_id: backendData.category_id,
      user_id: backendData.user_id,
      name: backendData.name,
      description: backendData.description,
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria uma nova categoria
   */
  async create(data: CreateCategoryRequest): Promise<Category> {
    if (import.meta.env.DEV) {
      console.log('[CategoryService] Criando categoria com dados:', JSON.stringify(data, null, 2))
    }

    const response = await apiClient.post<{
      message: string
      data: {
        category_id: string
        user_id: string
        name: string
        description: string
        is_active: boolean
        created_at: string
      }
    }>('/categories', data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      category_id: backendData.category_id,
      user_id: backendData.user_id,
      name: backendData.name,
      description: backendData.description,
      is_active: backendData.is_active,
      created_at: backendData.created_at,
      updated_at: backendData.created_at, // Backend não retorna updated_at na criação
    }
  },

  /**
   * Atualiza uma categoria existente
   */
  async update(categoryId: string, data: UpdateCategoryRequest): Promise<Category> {
    if (import.meta.env.DEV) {
      console.log('[CategoryService] Atualizando categoria:', categoryId, JSON.stringify(data, null, 2))
    }

    const response = await apiClient.put<{
      message: string
      data: {
        category_id: string
        user_id: string
        name: string
        description: string
        is_active: boolean
        updated_at: string
      }
    }>(`/categories/${categoryId}`, data)
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      category_id: backendData.category_id,
      user_id: backendData.user_id,
      name: backendData.name,
      description: backendData.description,
      is_active: backendData.is_active,
      created_at: '', // Backend não retorna created_at na atualização
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Deleta uma categoria
   */
  async delete(categoryId: string): Promise<void> {
    if (import.meta.env.DEV) {
      console.log('[CategoryService] Deletando categoria:', categoryId)
    }

    await apiClient.delete(`/categories/${categoryId}`)
  },
}

