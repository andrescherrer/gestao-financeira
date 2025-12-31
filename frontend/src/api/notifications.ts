import apiClient from './client'
import type {
  Notification,
  CreateNotificationRequest,
  ListNotificationsResponse,
} from './types'

/**
 * Serviço de API para notificações
 */
export const notificationService = {
  /**
   * Lista todas as notificações do usuário autenticado
   */
  async list(params?: {
    status?: 'UNREAD' | 'READ' | 'ARCHIVED'
    type?: 'INFO' | 'WARNING' | 'SUCCESS' | 'ERROR'
    page?: number
    page_size?: number
  }): Promise<ListNotificationsResponse> {
    const response = await apiClient.get<{
      message: string
      data: {
        notifications: Array<{
          notification_id: string
          title: string
          message: string
          type: string
          status: string
          read_at?: string | null
          metadata?: Record<string, any>
          created_at: string
        }>
        total: number
        page: number
        page_size: number
        total_pages: number
      }
    }>('/notifications', {
      params,
    })
    
    // Mapear resposta do backend
    const backendData = response.data.data
    return {
      notifications: backendData.notifications.map((notif) => ({
        notification_id: notif.notification_id,
        user_id: '', // Backend não retorna user_id na listagem
        title: notif.title,
        message: notif.message,
        type: notif.type as Notification['type'],
        status: notif.status as Notification['status'],
        read_at: notif.read_at,
        metadata: notif.metadata,
        created_at: notif.created_at,
        updated_at: notif.created_at, // Backend não retorna updated_at na listagem
      })),
      total: backendData.total,
      page: backendData.page,
      page_size: backendData.page_size,
      total_pages: backendData.total_pages,
    }
  },

  /**
   * Obtém detalhes de uma notificação específica
   */
  async get(notificationId: string): Promise<Notification> {
    const response = await apiClient.get<{
      message: string
      data: {
        notification_id: string
        user_id: string
        title: string
        message: string
        type: string
        status: string
        read_at?: string | null
        metadata?: Record<string, any>
        created_at: string
        updated_at: string
      }
    }>(`/notifications/${notificationId}`)

    const backendData = response.data.data
    return {
      notification_id: backendData.notification_id,
      user_id: backendData.user_id,
      title: backendData.title,
      message: backendData.message,
      type: backendData.type as Notification['type'],
      status: backendData.status as Notification['status'],
      read_at: backendData.read_at,
      metadata: backendData.metadata,
      created_at: backendData.created_at,
      updated_at: backendData.updated_at,
    }
  },

  /**
   * Cria uma nova notificação
   */
  async create(data: CreateNotificationRequest): Promise<Notification> {
    const response = await apiClient.post<{
      message: string
      data: {
        notification_id: string
        user_id: string
        title: string
        message: string
        type: string
        status: string
        metadata?: Record<string, any>
        created_at: string
      }
    }>('/notifications', data)

    const backendData = response.data.data
    return {
      notification_id: backendData.notification_id,
      user_id: backendData.user_id,
      title: backendData.title,
      message: backendData.message,
      type: backendData.type as Notification['type'],
      status: backendData.status as Notification['status'],
      metadata: backendData.metadata,
      created_at: backendData.created_at,
      updated_at: backendData.created_at,
    }
  },

  /**
   * Marca uma notificação como lida
   */
  async markAsRead(notificationId: string): Promise<Notification> {
    const response = await apiClient.post<{
      message: string
      data: {
        notification_id: string
        status: string
        read_at: string
      }
    }>(`/notifications/${notificationId}/read`)

    const backendData = response.data.data
    // Buscar notificação completa para retornar
    return this.get(notificationId)
  },

  /**
   * Marca uma notificação como não lida
   */
  async markAsUnread(notificationId: string): Promise<Notification> {
    const response = await apiClient.post<{
      message: string
      data: {
        notification_id: string
        status: string
      }
    }>(`/notifications/${notificationId}/unread`)

    const backendData = response.data.data
    // Buscar notificação completa para retornar
    return this.get(notificationId)
  },

  /**
   * Arquiva uma notificação
   */
  async archive(notificationId: string): Promise<Notification> {
    const response = await apiClient.post<{
      message: string
      data: {
        notification_id: string
        status: string
      }
    }>(`/notifications/${notificationId}/archive`)

    const backendData = response.data.data
    // Buscar notificação completa para retornar
    return this.get(notificationId)
  },

  /**
   * Exclui uma notificação
   */
  async delete(notificationId: string): Promise<void> {
    await apiClient.delete(`/notifications/${notificationId}`)
  },
}

