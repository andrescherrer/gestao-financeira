<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Metas', to: '/goals' },
          { label: goalsStore.currentGoal?.name || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="goalsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando detalhes da meta...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="goalsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ goalsStore.error }}</p>
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

      <!-- Goal Details -->
      <div v-else-if="goalsStore.currentGoal" class="space-y-6">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold text-foreground mb-2">
              {{ goalsStore.currentGoal.name }}
            </h1>
            <p class="text-muted-foreground">
              {{ getStatusLabel(goalsStore.currentGoal.status) }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <Badge
              :variant="getStatusBadgeVariant(goalsStore.currentGoal.status)"
            >
              {{ getStatusLabel(goalsStore.currentGoal.status) }}
            </Badge>
            <Button
              v-if="canCancel"
              variant="outline"
              @click="handleCancel"
              :disabled="isCancelling"
            >
              <XCircle v-if="!isCancelling" class="h-4 w-4 mr-2" />
              <Loader2 v-else class="h-4 w-4 mr-2 animate-spin" />
              Cancelar Meta
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
        </div>

        <!-- Progress Card -->
        <Card>
          <CardContent class="p-6">
            <GoalProgress :goal="goalsStore.currentGoal" />
          </CardContent>
        </Card>

        <!-- Actions Card -->
        <Card v-if="canAddContribution || canUpdateProgress">
          <CardContent class="p-6">
            <h2 class="text-xl font-semibold mb-4">Ações</h2>
            <div class="grid gap-4 md:grid-cols-2">
              <!-- Add Contribution -->
              <div v-if="canAddContribution" class="space-y-4">
                <h3 class="font-medium">Adicionar Contribuição</h3>
                <Form
                  @submit="handleAddContribution"
                  :validation-schema="addContributionSchema"
                  v-slot="{ errors }"
                >
                  <div class="space-y-4">
                    <div>
                      <Label for="contribution_amount">Valor da Contribuição</Label>
                      <Field
                        id="contribution_amount"
                        name="amount"
                        type="number"
                        step="0.01"
                        min="0.01"
                        v-slot="{ field }"
                      >
                        <Input
                          :name="field.name"
                          type="number"
                          step="0.01"
                          min="0.01"
                          :value="field.value"
                          @input="field.onInput"
                          @change="field.onChange"
                          @blur="field.onBlur"
                          placeholder="0.00"
                        />
                      </Field>
                      <ErrorMessage name="amount" class="mt-1 text-sm text-destructive" />
                    </div>
                    <Button type="submit" :disabled="isAddingContribution">
                      <Loader2 v-if="isAddingContribution" class="h-4 w-4 mr-2 animate-spin" />
                      <Plus v-else class="h-4 w-4 mr-2" />
                      Adicionar Contribuição
                    </Button>
                  </div>
                </Form>
              </div>

              <!-- Update Progress -->
              <div v-if="canUpdateProgress" class="space-y-4">
                <h3 class="font-medium">Atualizar Progresso</h3>
                <Form
                  @submit="handleUpdateProgress"
                  :validation-schema="updateProgressSchema"
                  v-slot="{ errors }"
                >
                  <div class="space-y-4">
                    <div>
                      <Label for="progress_amount">Valor Atual</Label>
                      <Field
                        id="progress_amount"
                        name="amount"
                        type="number"
                        step="0.01"
                        min="0"
                        v-slot="{ field }"
                      >
                        <Input
                          :name="field.name"
                          type="number"
                          step="0.01"
                          min="0"
                          :value="field.value"
                          @input="field.onInput"
                          @change="field.onChange"
                          @blur="field.onBlur"
                          :placeholder="goalsStore.currentGoal.current_amount"
                        />
                      </Field>
                      <ErrorMessage name="amount" class="mt-1 text-sm text-destructive" />
                    </div>
                    <Button type="submit" :disabled="isUpdatingProgress">
                      <Loader2 v-if="isUpdatingProgress" class="h-4 w-4 mr-2 animate-spin" />
                      <Check v-else class="h-4 w-4 mr-2" />
                      Atualizar Progresso
                    </Button>
                  </div>
                </Form>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Goal Information -->
        <Card>
          <CardContent class="p-6">
            <h2 class="text-xl font-semibold mb-4">Informações</h2>
            <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Contexto
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getContextLabel(goalsStore.currentGoal.context) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Moeda
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ goalsStore.currentGoal.currency }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data Limite
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(goalsStore.currentGoal.deadline) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Criada em
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDateTime(goalsStore.currentGoal.created_at) }}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Cancel Dialog -->
        <ConfirmDialog
          v-model:open="showCancelDialog"
          title="Cancelar Meta"
          description="Tem certeza que deseja cancelar esta meta? Esta ação não pode ser desfeita."
          confirm-label="Cancelar Meta"
          cancel-label="Manter Meta"
          variant="destructive"
          :isLoading="isCancelling"
          @confirm="confirmCancel"
        />

        <!-- Delete Confirmation Dialog -->
        <ConfirmDialog
          v-model:open="showDeleteDialog"
          title="Excluir Meta"
          description="Tem certeza que deseja excluir esta meta? Esta ação não pode ser desfeita."
          confirm-label="Excluir"
          cancel-label="Cancelar"
          variant="destructive"
          :isLoading="isDeleting"
          @confirm="confirmDelete"
        />
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { addContributionSchema, updateProgressSchema } from '@/validations/goal'
import { useGoalsStore } from '@/stores/goals'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import GoalProgress from '@/components/GoalProgress.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  Loader2,
  AlertCircle,
  Trash2,
  XCircle,
  Plus,
  Check,
} from 'lucide-vue-next'
import type { Goal } from '@/api/types'

const route = useRoute()
const router = useRouter()
const goalsStore = useGoalsStore()
const authStore = useAuthStore()

const isDeleting = ref(false)
const isCancelling = ref(false)
const isAddingContribution = ref(false)
const isUpdatingProgress = ref(false)
const showCancelDialog = ref(false)
const showDeleteDialog = ref(false)

const goalId = computed(() => route.params.id as string)

const canAddContribution = computed(() => {
  const goal = goalsStore.currentGoal
  if (!goal) return false
  return goal.status === 'IN_PROGRESS' || goal.status === 'OVERDUE'
})

const canUpdateProgress = computed(() => {
  const goal = goalsStore.currentGoal
  if (!goal) return false
  return goal.status === 'IN_PROGRESS' || goal.status === 'OVERDUE'
})

const canCancel = computed(() => {
  const goal = goalsStore.currentGoal
  if (!goal) return false
  return goal.status === 'IN_PROGRESS' || goal.status === 'OVERDUE'
})

onMounted(async () => {
  if (!authStore.token) {
    authStore.init()
  }
  
  if (authStore.isValidating) {
    await new Promise(resolve => setTimeout(resolve, 100))
  }
  
  if (!goalsStore.currentGoal || goalsStore.currentGoal.goal_id !== goalId.value) {
    await goalsStore.getGoal(goalId.value)
  }
})

async function handleRetry() {
  goalsStore.clearError()
  await goalsStore.getGoal(goalId.value)
}

function goBack() {
  router.push('/goals')
}

async function handleAddContribution(values: any) {
  isAddingContribution.value = true
  try {
    await goalsStore.addContribution(goalId.value, {
      amount: values.amount,
    })
    const { toast } = await import('@/components/ui/toast')
    toast.success('Contribuição adicionada com sucesso!')
  } catch (err: any) {
    const { extractErrorMessage } = await import('@/utils/errorTranslations')
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao adicionar contribuição', {
      description: extractErrorMessage(err),
    })
  } finally {
    isAddingContribution.value = false
  }
}

async function handleUpdateProgress(values: any) {
  isUpdatingProgress.value = true
  try {
    await goalsStore.updateProgress(goalId.value, {
      amount: values.amount,
    })
    const { toast } = await import('@/components/ui/toast')
    toast.success('Progresso atualizado com sucesso!')
  } catch (err: any) {
    const { extractErrorMessage } = await import('@/utils/errorTranslations')
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao atualizar progresso', {
      description: extractErrorMessage(err),
    })
  } finally {
    isUpdatingProgress.value = false
  }
}

function handleCancel() {
  showCancelDialog.value = true
}

async function confirmCancel() {
  isCancelling.value = true
  try {
    await goalsStore.cancelGoal(goalId.value)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Meta cancelada com sucesso!')
    showCancelDialog.value = false
  } catch (err: any) {
    const { extractErrorMessage } = await import('@/utils/errorTranslations')
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao cancelar meta', {
      description: extractErrorMessage(err),
  })
  } finally {
    isCancelling.value = false
  }
}

function handleDelete() {
  showDeleteDialog.value = true
}

async function confirmDelete() {
  isDeleting.value = true
  try {
    await goalsStore.deleteGoal(goalId.value)
    const { toast } = await import('@/components/ui/toast')
    toast.success('Meta excluída com sucesso!')
    router.push('/goals')
  } catch (err: any) {
    const { extractErrorMessage } = await import('@/utils/errorTranslations')
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao excluir meta', {
      description: extractErrorMessage(err),
    })
  } finally {
    isDeleting.value = false
    showDeleteDialog.value = false
  }
}

function getStatusLabel(status: Goal['status']): string {
  const labels: Record<Goal['status'], string> = {
    IN_PROGRESS: 'Em Progresso',
    COMPLETED: 'Concluída',
    OVERDUE: 'Vencida',
    CANCELLED: 'Cancelada',
  }
  return labels[status] || status
}

function getStatusBadgeVariant(status: Goal['status']): 'default' | 'secondary' | 'destructive' | 'outline' {
  const variants: Record<Goal['status'], 'default' | 'secondary' | 'destructive' | 'outline'> = {
    IN_PROGRESS: 'default',
    COMPLETED: 'default',
    OVERDUE: 'destructive',
    CANCELLED: 'outline',
  }
  return variants[status] || 'default'
}

function getContextLabel(context: Goal['context']): string {
  return context === 'PERSONAL' ? 'Pessoal' : 'Negócio'
}

function formatDate(date: string): string {
  return new Date(date).toLocaleDateString('pt-BR')
}

function formatDateTime(dateTime: string): string {
  return new Date(dateTime).toLocaleString('pt-BR')
}
</script>

