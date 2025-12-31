import { ref, onMounted, onUnmounted } from 'vue'
import { useNotificationsStore } from '@/stores/notifications'
import { useAuthStore } from '@/stores/auth'
import { env } from '@/config/env'
import { toast } from '@/components/ui/toast'
import type { Notification } from '@/api/types'

interface WebSocketMessage {
  type: 'notification' | 'notification_update' | 'pong'
  data?: any
  notification_id?: string
  update_type?: string
}

/**
 * Composable para gerenciar conexão WebSocket de notificações
 */
export function useWebSocketNotifications() {
  const notificationsStore = useNotificationsStore()
  const authStore = useAuthStore()

  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = ref(3000) // 3 segundos inicial
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let pingInterval: ReturnType<typeof setInterval> | null = null

  /**
   * Conecta ao WebSocket
   */
  function connect() {
    if (ws.value?.readyState === WebSocket.OPEN) {
      return // Já conectado
    }

    const token = authStore.token
    if (!token) {
      console.warn('[WebSocket] Token não encontrado, não é possível conectar')
      return
    }

    try {
      // Construir URL do WebSocket
      const wsUrl = env.apiUrl.replace('http://', 'ws://').replace('https://', 'wss://')
      const url = `${wsUrl}/ws/notifications?token=${encodeURIComponent(token)}`

      ws.value = new WebSocket(url)

      ws.value.onopen = () => {
        console.log('[WebSocket] Conectado')
        isConnected.value = true
        reconnectAttempts.value = 0
        reconnectDelay.value = 3000

        // Iniciar ping
        startPing()
      }

      ws.value.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          handleMessage(message)
        } catch (error) {
          console.error('[WebSocket] Erro ao processar mensagem:', error)
        }
      }

      ws.value.onerror = (error) => {
        console.error('[WebSocket] Erro:', error)
        isConnected.value = false
      }

      ws.value.onclose = (event) => {
        console.log('[WebSocket] Desconectado', event.code, event.reason)
        isConnected.value = false
        stopPing()

        // Tentar reconectar se não foi fechado intencionalmente
        if (event.code !== 1000 && reconnectAttempts.value < maxReconnectAttempts) {
          scheduleReconnect()
        }
      }
    } catch (error) {
      console.error('[WebSocket] Erro ao conectar:', error)
      scheduleReconnect()
    }
  }

  /**
   * Desconecta do WebSocket
   */
  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    stopPing()

    if (ws.value) {
      ws.value.close(1000, 'Desconexão intencional')
      ws.value = null
    }
    isConnected.value = false
  }

  /**
   * Agenda reconexão
   */
  function scheduleReconnect() {
    if (reconnectTimer) {
      return // Já agendado
    }

    reconnectAttempts.value++
    if (reconnectAttempts.value > maxReconnectAttempts) {
      console.warn('[WebSocket] Máximo de tentativas de reconexão atingido')
      return
    }

    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      console.log(`[WebSocket] Tentando reconectar (${reconnectAttempts.value}/${maxReconnectAttempts})...`)
      connect()
    }, reconnectDelay.value)

    // Aumentar delay exponencialmente
    reconnectDelay.value = Math.min(reconnectDelay.value * 2, 30000) // Max 30s
  }

  /**
   * Inicia ping periódico
   */
  function startPing() {
    if (pingInterval) {
      return
    }

    pingInterval = setInterval(() => {
      if (ws.value?.readyState === WebSocket.OPEN) {
        ws.value.send(JSON.stringify({ type: 'ping' }))
      }
    }, 50000) // Ping a cada 50 segundos
  }

  /**
   * Para ping
   */
  function stopPing() {
    if (pingInterval) {
      clearInterval(pingInterval)
      pingInterval = null
    }
  }

  /**
   * Processa mensagens recebidas
   */
  function handleMessage(message: WebSocketMessage) {
    switch (message.type) {
      case 'notification':
        handleNewNotification(message.data)
        break
      case 'notification_update':
        handleNotificationUpdate(message)
        break
      case 'pong':
        // Resposta ao ping, não precisa fazer nada
        break
      default:
        console.warn('[WebSocket] Tipo de mensagem desconhecido:', message.type)
    }
  }

  /**
   * Processa nova notificação
   */
  function handleNewNotification(data: any) {
    if (!data) return

    const notification: Notification = {
      notification_id: data.notification_id,
      user_id: data.user_id || '',
      title: data.title,
      message: data.message || '',
      type: data.type,
      status: data.status || 'UNREAD',
      read_at: data.read_at || null,
      metadata: data.metadata || {},
      created_at: data.created_at,
      updated_at: data.created_at,
    }

    // Adicionar ao store
    notificationsStore.addNotification(notification)

    // Exibir toast
    const typeLabels: Record<Notification['type'], string> = {
      INFO: 'Informação',
      WARNING: 'Aviso',
      SUCCESS: 'Sucesso',
      ERROR: 'Erro',
    }

    toast.info(notification.title, {
      description: notification.message,
      duration: 5000,
    })
  }

  /**
   * Processa atualização de notificação
   */
  function handleNotificationUpdate(message: WebSocketMessage) {
    // Atualizar notificação no store se existir
    if (message.notification_id) {
      // Buscar notificação atualizada
      notificationsStore.getNotification(message.notification_id).catch(() => {
        // Se não conseguir buscar, apenas logar
        console.log('[WebSocket] Atualização de notificação recebida:', message)
      })
    }
  }

  // Conectar ao montar
  onMounted(() => {
    if (authStore.isAuthenticated) {
      connect()
    }
  })

  // Desconectar ao desmontar
  onUnmounted(() => {
    disconnect()
  })

  // Reconectar quando autenticação mudar
  const unwatchAuth = authStore.$subscribe(() => {
    if (authStore.isAuthenticated && !isConnected.value) {
      connect()
    } else if (!authStore.isAuthenticated && isConnected.value) {
      disconnect()
    }
  })

  // Cleanup do watcher
  onUnmounted(() => {
    unwatchAuth()
  })

  return {
    isConnected,
    connect,
    disconnect,
  }
}

