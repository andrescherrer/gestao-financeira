<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Metas' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl sm:text-4xl font-bold mb-2">Metas</h1>
          <p class="text-sm sm:text-base text-muted-foreground">
            Defina e acompanhe suas metas financeiras
          </p>
        </div>
        <Button as-child class="w-full sm:w-auto">
          <router-link to="/goals/new">
            <Plus class="h-4 w-4 mr-2" />
            Nova Meta
          </router-link>
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="goalsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando metas...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="goalsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ goalsStore.error }}</p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="link"
              @click="handleRetry"
              class="text-destructive"
            >
              Tentar novamente
            </Button>
            <Button
              v-if="isAuthError"
              variant="outline"
              @click="handleLogin"
              class="text-destructive border-destructive"
            >
              Fazer login novamente
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <EmptyState
        v-else-if="goalsStore.goals.length === 0"
        :icon="Target"
        title="Nenhuma meta encontrada"
        description="Comece criando sua primeira meta financeira para acompanhar seu progresso."
      >
        <Button as-child>
          <router-link to="/goals/new">
            <Plus class="h-4 w-4 mr-2" />
            Nova Meta
          </router-link>
        </Button>
      </EmptyState>

      <!-- Goals Grid -->
      <div v-else class="space-y-6">
        <!-- Summary Cards -->
        <div class="grid gap-4 md:grid-cols-4">
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Total de Metas
              </div>
              <div class="text-2xl font-bold">
                {{ goalsStore.totalGoals }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Em Progresso
              </div>
              <div class="text-2xl font-bold text-blue-600">
                {{ goalsStore.inProgressGoals.length }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Concluídas
              </div>
              <div class="text-2xl font-bold text-green-600">
                {{ goalsStore.completedGoals.length }}
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardContent class="p-4">
              <div class="text-sm font-medium text-muted-foreground mb-1">
                Valor Total Alvo
              </div>
              <div class="text-2xl font-bold">
                {{ formatCurrency(goalsStore.totalTargetAmount.toString(), 'BRL') }}
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Goals List -->
        <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <GoalCard
            v-for="goal in goalsStore.goals"
            :key="goal.goal_id"
            :goal="goal"
          />
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGoalsStore } from '@/stores/goals'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import GoalCard from '@/components/GoalCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus, Loader2, AlertCircle, Target } from 'lucide-vue-next'

const router = useRouter()
const goalsStore = useGoalsStore()
const authStore = useAuthStore()

const isAuthError = computed(() => {
  const error = goalsStore.error?.toLowerCase() || ''
  return error.includes('token') || 
         error.includes('autenticação') || 
         error.includes('unauthorized') ||
         error.includes('invalid') ||
         error.includes('expired')
})

function formatCurrency(amount: string, currency: string): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

onMounted(async () => {
  if (!authStore.token) {
    authStore.init()
  }
  
  if (authStore.isValidating) {
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  
  if (goalsStore.goals.length === 0 && !goalsStore.isLoading) {
    await goalsStore.listGoals()
  }
})

function handleRetry() {
  goalsStore.clearError()
  goalsStore.listGoals()
}

function handleLogin() {
  router.push('/login')
}
</script>

