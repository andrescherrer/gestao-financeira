<template>
  <Card
    class="group relative cursor-pointer overflow-hidden transition-all hover:border-primary hover:shadow-lg"
    :class="{ 'border-primary/50 bg-primary/5': notification.status === 'UNREAD' }"
    @click="handleClick"
    role="button"
    :aria-label="`Notificação ${notification.title}`"
    tabindex="0"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <CardContent class="relative p-4">
      <div class="flex items-start gap-4">
        <!-- Icon -->
        <div
          class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg"
          :class="getTypeIconBg(notification.type)"
        >
          <component :is="getTypeIcon(notification.type)" class="h-5 w-5 text-white" />
        </div>

        <!-- Content -->
        <div class="flex-1 min-w-0">
          <div class="mb-2 flex items-start justify-between gap-2">
            <div class="flex-1 min-w-0">
              <h3
                class="text-base font-semibold text-foreground mb-1"
                :class="{ 'font-bold': notification.status === 'UNREAD' }"
              >
                {{ notification.title }}
              </h3>
              <p class="text-sm text-muted-foreground line-clamp-2">
                {{ notification.message }}
              </p>
            </div>
            <Badge
              :variant="getTypeBadgeVariant(notification.type)"
              class="shrink-0"
            >
              {{ getTypeLabel(notification.type) }}
            </Badge>
          </div>

          <div class="flex items-center justify-between text-xs text-muted-foreground">
            <span>{{ formatDate(notification.created_at) }}</span>
            <div class="flex items-center gap-2">
              <Badge
                v-if="notification.status === 'UNREAD'"
                variant="secondary"
                class="text-xs"
              >
                Não lida
              </Badge>
              <Badge
                v-else-if="notification.status === 'ARCHIVED'"
                variant="outline"
                class="text-xs"
              >
                Arquivada
              </Badge>
            </div>
          </div>
        </div>
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type { Notification } from '@/api/types'
import { Info, AlertTriangle, CheckCircle2, XCircle, Bell } from 'lucide-vue-next'

interface Props {
  notification: Notification
}

const props = defineProps<Props>()
const router = useRouter()

const handleClick = () => {
  router.push(`/notifications/${props.notification.notification_id}`)
}

const getTypeIcon = (type: Notification['type']) => {
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

const getTypeIconBg = (type: Notification['type']) => {
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

const getTypeLabel = (type: Notification['type']) => {
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

const getTypeBadgeVariant = (type: Notification['type']) => {
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
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) {
    return 'Agora'
  } else if (diffMins < 60) {
    return `${diffMins} min atrás`
  } else if (diffHours < 24) {
    return `${diffHours}h atrás`
  } else if (diffDays < 7) {
    return `${diffDays}d atrás`
  } else {
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    })
  }
}
</script>

