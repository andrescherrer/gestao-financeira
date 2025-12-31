import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { notificationService } from '@/api/notifications'
import { useAuthStore } from '@/stores/auth'
import type { Notification, CreateNotificationRequest } from '@/api/types'

export const useNotificationsStore = defineStore('notifications', () => {
  const authStore = useAuthStore()

  // Estado
  const notifications = ref<Notification[]>([])
  const currentNotification = ref<Notification | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalNotifications = computed(() => notifications.value.length)
  const unreadNotifications = computed(() =>
    notifications.value.filter((notif) => notif.status === 'UNREAD')
  )
  const readNotifications = computed(() =>
    notifications.value.filter((notif) => notif.status === 'READ')
  )
  const archivedNotifications = computed(() =>
    notifications.value.filter((notif) => notif.status === 'ARCHIVED')
  )
  const unreadCount = computed(() => unreadNotifications.value.length)
  const infoNotifications = computed(() =>
    notifications.value.filter((notif) => notif.type === 'INFO')
  )
  const warningNotifications = computed(() =>
    notifications.value.filter((notif) => notif.type === 'WARNING')
  )
  const successNotifications = computed(() =>
    notifications.value.filter((notif) => notif.type === 'SUCCESS')
  )
  const errorNotifications = computed(() =>
    notifications.value.filter((notif) => notif.type === 'ERROR')
  )

  /**
   * Lista todas as notificações do usuário
   */
  async function listNotifications(params?: {
    status?: 'UNREAD' | 'READ' | 'ARCHIVED'
    type?: 'INFO' | 'WARNING' | 'SUCCESS' | 'ERROR'
    page?: number
    page_size?: number
  }) {
    isLoading.value = true
    error.value = null
    try {
      const token = localStorage.getItem('auth_token')
      if (!token) {
        throw new Error('Token de autenticação não encontrado. Faça login novamente.')
      }
      
      const response = await notificationService.list(params)
      notifications.value = response.notifications || []
      return response
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      if (import.meta.env.DEV) {
        console.error('[Notifications Store] Erro ao listar notificações:', {
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
  },

  /**
   * Obtém uma notificação específica
   */
  async function getNotification(notificationId: string) {
    isLoading.value = true
    error.value = null
    try {
      const notification = await notificationService.get(notificationId)
      currentNotification.value = notification
      
      // Atualizar na lista se já existir
      const index = notifications.value.findIndex(
        (n) => n.notification_id === notificationId
      )
      if (index !== -1) {
        notifications.value[index] = notification
      }
      
      return notification
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      if (import.meta.env.DEV) {
        console.error('[Notifications Store] Erro ao obter notificação:', {
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
  },

  /**
   * Cria uma nova notificação
   */
  async function createNotification(data: CreateNotificationRequest) {
    isLoading.value = true
    error.value = null
    try {
      const notification = await notificationService.create(data)
      notifications.value.unshift(notification) // Adicionar no início
      return notification
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      
      if (import.meta.env.DEV) {
        console.error('[Notifications Store] Erro ao criar notificação:', {
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
  },

  /**
   * Marca uma notificação como lida
   */
  async function markAsRead(notificationId: string) {
    error.value = null
    try {
      const notification = await notificationService.markAsRead(notificationId)
      
      // Atualizar na lista
      const index = notifications.value.findIndex(
        (n) => n.notification_id === notificationId
      )
      if (index !== -1) {
        notifications.value[index] = notification
      }
      
      // Atualizar currentNotification se for a mesma
      if (currentNotification.value?.notification_id === notificationId) {
        currentNotification.value = notification
      }
      
      return notification
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    }
  },

  /**
   * Marca uma notificação como não lida
   */
  async function markAsUnread(notificationId: string) {
    error.value = null
    try {
      const notification = await notificationService.markAsUnread(notificationId)
      
      // Atualizar na lista
      const index = notifications.value.findIndex(
        (n) => n.notification_id === notificationId
      )
      if (index !== -1) {
        notifications.value[index] = notification
      }
      
      // Atualizar currentNotification se for a mesma
      if (currentNotification.value?.notification_id === notificationId) {
        currentNotification.value = notification
      }
      
      return notification
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    }
  },

  /**
   * Arquiva uma notificação
   */
  async function archiveNotification(notificationId: string) {
    error.value = null
    try {
      const notification = await notificationService.archive(notificationId)
      
      // Atualizar na lista
      const index = notifications.value.findIndex(
        (n) => n.notification_id === notificationId
      )
      if (index !== -1) {
        notifications.value[index] = notification
      }
      
      // Atualizar currentNotification se for a mesma
      if (currentNotification.value?.notification_id === notificationId) {
        currentNotification.value = notification
      }
      
      return notification
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    }
  },

  /**
   * Exclui uma notificação
   */
  async function deleteNotification(notificationId: string) {
    error.value = null
    try {
      await notificationService.delete(notificationId)
      
      // Remover da lista
      notifications.value = notifications.value.filter(
        (n) => n.notification_id !== notificationId
      )
      
      // Limpar currentNotification se for a mesma
      if (currentNotification.value?.notification_id === notificationId) {
        currentNotification.value = null
      }
    } catch (err: any) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      error.value = extractErrorMessage(err)
      throw err
    }
  },

  /**
   * Adiciona uma notificação recebida via WebSocket
   */
  function addNotification(notification: Notification) {
    // Verificar se já existe
    const index = notifications.value.findIndex(
      (n) => n.notification_id === notification.notification_id
    )
    if (index !== -1) {
      notifications.value[index] = notification
    } else {
      notifications.value.unshift(notification) // Adicionar no início
    }
  },

  /**
   * Limpa o erro
   */
  function clearError() {
    error.value = null
  },

  /**
   * Limpa a notificação atual
   */
  function clearCurrentNotification() {
    currentNotification.value = null
  },

  return {
    // Estado
    notifications,
    currentNotification,
    isLoading,
    error,
    // Computed
    totalNotifications,
    unreadNotifications,
    readNotifications,
    archivedNotifications,
    unreadCount,
    infoNotifications,
    warningNotifications,
    successNotifications,
    errorNotifications,
    // Actions
    listNotifications,
    getNotification,
    createNotification,
    markAsRead,
    markAsUnread,
    archiveNotification,
    deleteNotification,
    addNotification,
    clearError,
    clearCurrentNotification,
  }
})

