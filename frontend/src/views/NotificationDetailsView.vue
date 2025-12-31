<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Notificações', to: '/notifications' },
          { label: notificationStore.currentNotification?.title || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="notificationStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando detalhes da notificação...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="notificationStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ notificationStore.error }}</p>
          </div>
          <div class="flex gap-3">
            <Button
              @click="handleRetry"
              variant="destructive"
            >
              Tentar novamente
            </Button>
            <Button
              variant="outline"
              @click="goBack"
            >
              Voltar
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Notification Details -->
      <div v-else-if="notificationStore.currentNotification" class="space-y-6">
        <!-- Header -->
        <Card>
          <CardContent class="p-6">
            <div class="flex items-start gap-4 mb-4">
              <div
                class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg"
                :class="getTypeIconBg(notificationStore.currentNotification.type)"
              >
                <component
                  :is="getTypeIcon(notificationStore.currentNotification.type)"
                  class="h-6 w-6 text-white"
                />
              </div>
              <div class="flex-1">
                <div class="mb-2 flex items-start justify-between gap-4">
                  <h1 class="text-2xl font-bold text-foreground">
                    {{ notificationStore.currentNotification.title }}
                  </h1>
                  <Badge
                    :variant="getTypeBadgeVariant(notificationStore.currentNotification.type)"
                  >
                    {{ getTypeLabel(notificationStore.currentNotification.type) }}
                  </Badge>
                </div>
                <p class="text-muted-foreground mb-4">
                  {{ notificationStore.currentNotification.message }}
                </p>
                <div class="flex items-center gap-4 text-sm text-muted-foreground">
                  <span>
                    Criada em: {{ formatDate(notificationStore.currentNotification.created_at) }}
                  </span>
                  <span v-if="notificationStore.currentNotification.read_at">
                    Lida em: {{ formatDate(notificationStore.currentNotification.read_at) }}
                  </span>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Actions -->
        <Card>
          <CardContent class="p-6">
            <h2 class="text-xl font-semibold mb-4">Ações</h2>
            <div class="flex flex-wrap gap-3">
              <Button
                v-if="notificationStore.currentNotification.status === 'UNREAD'"
                @click="handleMarkAsRead"
                :disabled="isMarkingAsRead"
              >
                <CheckCircle2 v-if="!isMarkingAsRead" class="h-4 w-4 mr-2" />
                <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
                Marcar como lida
              </Button>
              <Button
                v-else-if="notificationStore.currentNotification.status === 'READ'"
                variant="outline"
                @click="handleMarkAsUnread"
                :disabled="isMarkingAsUnread"
              >
                <Circle v-if="!isMarkingAsUnread" class="h-4 w-4 mr-2" />
                <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
                Marcar como não lida
              </Button>
              <Button
                v-if="notificationStore.currentNotification.status !== 'ARCHIVED'"
                variant="outline"
                @click="handleArchive"
                :disabled="isArchiving"
              >
                <Archive v-if="!isArchiving" class="h-4 w-4 mr-2" />
                <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
                Arquivar
              </Button>
              <Button
                variant="destructive"
                @click="handleDelete"
                :disabled="isDeleting"
              >
                <Trash2 v-if="!isDeleting" class="h-4 w-4 mr-2" />
                <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
                Excluir
              </Button>
            </div>
          </CardContent>
        </Card>

        <!-- Metadata -->
        <Card v-if="notificationStore.currentNotification.metadata && Object.keys(notificationStore.currentNotification.metadata).length > 0">
          <CardContent class="p-6">
            <h2 class="text-xl font-semibold mb-4">Informações Adicionais</h2>
            <div class="space-y-2">
              <div
                v-for="(value, key) in notificationStore.currentNotification.metadata"
                :key="key"
                class="flex items-center justify-between py-2 border-b last:border-0"
              >
                <span class="font-medium text-muted-foreground">{{ key }}:</span>
                <span class="text-foreground">{{ formatMetadataValue(value) }}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Delete Confirmation Dialog -->
      <ConfirmDialog
        v-model:open="showDeleteDialog"
        title="Excluir Notificação"
        description="Tem certeza que deseja excluir esta notificação? Esta ação não pode ser desfeita."
        confirm-text="Excluir"
        cancel-text="Cancelar"
        variant="destructive"
        @confirm="confirmDelete"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotificationsStore } from '@/stores/notifications'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  Info,
  AlertTriangle,
  CheckCircle2,
  XCircle,
  Bell,
  Loader2,
  AlertCircle,
  Circle,
  Archive,
  Trash2,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const notificationStore = useNotificationsStore()

const isMarkingAsRead = ref(false)
const isMarkingAsUnread = ref(false)
const isArchiving = ref(false)
const isDeleting = ref(false)
const showDeleteDialog = ref(false)

const notificationId = route.params.id as string

const loadNotification = async () => {
  try {
    await notificationStore.getNotification(notificationId)
    // Marcar como lida automaticamente ao abrir
    if (notificationStore.currentNotification?.status === 'UNREAD') {
      await handleMarkAsRead()
    }
  } catch (error) {
    console.error('Erro ao carregar notificação:', error)
  }
}

const handleRetry = () => {
  notificationStore.clearError()
  loadNotification()
}

const goBack = () => {
  router.push('/notifications')
}

const handleMarkAsRead = async () => {
  if (!notificationStore.currentNotification) return
  isMarkingAsRead.value = true
  try {
    await notificationStore.markAsRead(notificationStore.currentNotification.notification_id)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Notificação marcada como lida')
  } catch (error) {
    console.error('Erro ao marcar como lida:', error)
  } finally {
    isMarkingAsRead.value = false
  }
}

const handleMarkAsUnread = async () => {
  if (!notificationStore.currentNotification) return
  isMarkingAsUnread.value = true
  try {
    await notificationStore.markAsUnread(notificationStore.currentNotification.notification_id)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Notificação marcada como não lida')
  } catch (error) {
    console.error('Erro ao marcar como não lida:', error)
  } finally {
    isMarkingAsUnread.value = false
  }
}

const handleArchive = async () => {
  if (!notificationStore.currentNotification) return
  isArchiving.value = true
  try {
    await notificationStore.archiveNotification(notificationStore.currentNotification.notification_id)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Notificação arquivada')
    router.push('/notifications')
  } catch (error) {
    console.error('Erro ao arquivar:', error)
  } finally {
    isArchiving.value = false
  }
}

const handleDelete = async () => {
  if (!notificationStore.currentNotification) return
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!notificationStore.currentNotification) return
  isDeleting.value = true
  try {
    await notificationStore.deleteNotification(notificationStore.currentNotification.notification_id)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Notificação excluída')
    router.push('/notifications')
  } catch (error) {
    console.error('Erro ao excluir:', error)
  } finally {
    isDeleting.value = false
    showDeleteDialog.value = false
  }
}

const getTypeIcon = (type: string) => {
  switch (type) {
    case 'INFO':
      return Info
    case 'WARNING':
      return AlertTriangle
    case 'SUCCESS':
      return CheckCircle2
    case 'ERROR':
      return XCircle
    default:
      return Bell
  }
}

const getTypeIconBg = (type: string) => {
  switch (type) {
    case 'INFO':
      return 'bg-blue-500'
    case 'WARNING':
      return 'bg-yellow-500'
    case 'SUCCESS':
      return 'bg-green-500'
    case 'ERROR':
      return 'bg-red-500'
    default:
      return 'bg-gray-500'
  }
}

const getTypeLabel = (type: string) => {
  switch (type) {
    case 'INFO':
      return 'Informação'
    case 'WARNING':
      return 'Aviso'
    case 'SUCCESS':
      return 'Sucesso'
    case 'ERROR':
      return 'Erro'
    default:
      return type
  }
}

const getTypeBadgeVariant = (type: string) => {
  switch (type) {
    case 'INFO':
      return 'default'
    case 'WARNING':
      return 'secondary'
    case 'SUCCESS':
      return 'default'
    case 'ERROR':
      return 'destructive'
    default:
      return 'outline'
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const formatMetadataValue = (value: any) => {
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  return String(value)
}

onMounted(() => {
  loadNotification()
})

onUnmounted(() => {
  notificationStore.clearCurrentNotification()
})
</script>

