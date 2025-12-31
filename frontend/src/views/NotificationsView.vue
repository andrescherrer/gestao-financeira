<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Notificações' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl sm:text-4xl font-bold mb-2">Notificações</h1>
          <p class="text-sm sm:text-base text-muted-foreground">
            Acompanhe suas notificações e atualizações
          </p>
        </div>
        <div class="flex gap-2">
          <Button
            variant="outline"
            @click="handleFilterChange('UNREAD')"
            :class="{ 'bg-primary text-primary-foreground': currentFilter === 'UNREAD' }"
          >
            Não lidas ({{ notificationsStore.unreadCount }})
          </Button>
          <Button
            variant="outline"
            @click="handleFilterChange(null)"
            :class="{ 'bg-primary text-primary-foreground': currentFilter === null }"
          >
            Todas
          </Button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="notificationsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando notificações...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="notificationsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ notificationsStore.error }}</p>
          </div>
          <Button
            variant="link"
            @click="handleRetry"
            class="text-destructive"
          >
            Tentar novamente
          </Button>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <EmptyState
        v-else-if="notificationsStore.notifications.length === 0"
        :icon="Bell"
        title="Nenhuma notificação encontrada"
        description="Você não possui notificações no momento."
      />

      <!-- Notifications List -->
      <div v-else class="space-y-4">
        <NotificationCard
          v-for="notification in notificationsStore.notifications"
          :key="notification.notification_id"
          :notification="notification"
        />
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useNotificationsStore } from '@/stores/notifications'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import EmptyState from '@/components/EmptyState.vue'
import NotificationCard from '@/components/NotificationCard.vue'
import { Bell, Loader2, AlertCircle } from 'lucide-vue-next'

const router = useRouter()
const notificationsStore = useNotificationsStore()

const currentFilter = ref<'UNREAD' | null>(null)

const handleFilterChange = async (filter: 'UNREAD' | null) => {
  currentFilter.value = filter
  await loadNotifications()
}

const loadNotifications = async () => {
  try {
    await notificationsStore.listNotifications({
      status: currentFilter.value || undefined,
      page: 1,
      page_size: 50,
    })
  } catch (error) {
    console.error('Erro ao carregar notificações:', error)
  }
}

const handleRetry = () => {
  notificationsStore.clearError()
  loadNotifications()
}

onMounted(() => {
  loadNotifications()
})
</script>

