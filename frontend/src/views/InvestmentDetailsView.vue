<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Investimentos', to: '/investments' },
          { label: investmentsStore.currentInvestment?.name || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="investmentsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando detalhes do investimento...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="investmentsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ investmentsStore.error }}</p>
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

      <!-- Investment Details -->
      <div v-else-if="investmentsStore.currentInvestment" class="space-y-6">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold text-foreground mb-2">
              {{ investmentsStore.currentInvestment.name }}
              <span
                v-if="investmentsStore.currentInvestment.ticker"
                class="text-xl text-muted-foreground"
              >
                ({{ investmentsStore.currentInvestment.ticker }})
              </span>
            </h1>
            <p class="text-muted-foreground">
              {{ getInvestmentTypeLabel(investmentsStore.currentInvestment.type) }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <Badge
              :variant="investmentsStore.currentInvestment.return_percentage >= 0 ? 'default' : 'destructive'"
            >
              {{ formatReturnPercentage(investmentsStore.currentInvestment.return_percentage) }}
            </Badge>
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

        <!-- Investment Card -->
        <Card>
          <CardContent class="p-6">
            <!-- Current Value -->
            <div class="mb-6 rounded-lg bg-gradient-to-br from-muted to-muted/50 p-6">
              <div class="text-sm font-medium text-muted-foreground mb-2">Valor Atual</div>
              <div class="text-4xl font-bold">
                {{ formatCurrency(investmentsStore.currentInvestment.current_value, investmentsStore.currentInvestment.currency) }}
              </div>
            </div>

            <!-- Return Card -->
            <div
              class="mb-6 rounded-lg p-6"
              :class="investmentsStore.currentInvestment.return_percentage >= 0
                ? 'bg-green-50 dark:bg-green-950 border border-green-200 dark:border-green-800'
                : 'bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800'"
            >
              <div class="text-sm font-medium text-muted-foreground mb-2">Retorno</div>
              <div
                class="text-3xl font-bold"
                :class="investmentsStore.currentInvestment.return_percentage >= 0
                  ? 'text-green-600 dark:text-green-400'
                  : 'text-red-600 dark:text-red-400'"
              >
                {{ formatCurrency(investmentsStore.currentInvestment.return_absolute, investmentsStore.currentInvestment.currency) }}
              </div>
              <div
                class="text-lg font-semibold mt-1"
                :class="investmentsStore.currentInvestment.return_percentage >= 0
                  ? 'text-green-600 dark:text-green-400'
                  : 'text-red-600 dark:text-red-400'"
              >
                {{ formatReturnPercentage(investmentsStore.currentInvestment.return_percentage) }}
              </div>
            </div>

            <!-- Investment Information Grid -->
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Tipo de Investimento
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getInvestmentTypeLabel(investmentsStore.currentInvestment.type) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Contexto
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getContextLabel(investmentsStore.currentInvestment.context) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Valor de Compra
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatCurrency(investmentsStore.currentInvestment.purchase_amount, investmentsStore.currentInvestment.currency) }}
                </div>
              </div>

              <div
                v-if="investmentsStore.currentInvestment.quantity"
                class="rounded-lg border border-border bg-muted/50 p-4"
              >
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Quantidade
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatQuantity(investmentsStore.currentInvestment.quantity) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data de Compra
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(investmentsStore.currentInvestment.purchase_date) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Moeda
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ investmentsStore.currentInvestment.currency }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data de Criação
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(investmentsStore.currentInvestment.created_at) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Última Atualização
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(investmentsStore.currentInvestment.updated_at) }}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Update Form -->
        <Card>
          <CardHeader>
            <CardTitle>Atualizar Investimento</CardTitle>
            <CardDescription>
              Atualize o valor atual ou a quantidade do investimento
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form @submit.prevent="handleUpdate" class="space-y-4">
              <div>
                <Label for="current_value" class="mb-1">
                  Valor Atual
                </Label>
                <Input
                  id="current_value"
                  v-model="updateForm.current_value"
                  type="number"
                  step="0.01"
                  min="0.01"
                  placeholder="0.00"
                />
              </div>
              <div v-if="investmentsStore.currentInvestment.quantity">
                <Label for="quantity" class="mb-1">
                  Quantidade
                </Label>
                <Input
                  id="quantity"
                  v-model="updateForm.quantity"
                  type="number"
                  step="0.0001"
                  min="0.0001"
                  placeholder="0.0000"
                />
              </div>
              <div class="flex gap-3">
                <Button
                  type="submit"
                  :disabled="isUpdating || (!updateForm.current_value && !updateForm.quantity)"
                >
                  <Loader2 v-if="isUpdating" class="h-4 w-4 mr-2 animate-spin" />
                  <Check v-else class="h-4 w-4 mr-2" />
                  {{ isUpdating ? 'Atualizando...' : 'Atualizar' }}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>

      <!-- Not Found State -->
      <Card v-else>
        <CardContent class="p-12 text-center">
          <AlertCircle class="mx-auto h-16 w-16 text-muted-foreground mb-4" />
          <h3 class="text-xl font-semibold text-foreground mb-2">
            Investimento não encontrado
          </h3>
          <p class="text-muted-foreground mb-6">
            O investimento que você está procurando não existe ou foi removido.
          </p>
          <Button
            @click="goBack"
            as-child
          >
            <router-link to="/investments">
              <ArrowLeft class="h-4 w-4 mr-2" />
              Voltar para investimentos
            </router-link>
          </Button>
        </CardContent>
      </Card>

      <!-- Delete Confirmation Dialog -->
      <ConfirmDialog
        v-model:open="showDeleteDialog"
        title="Excluir Investimento"
        description="Tem certeza que deseja excluir este investimento? Esta ação não pode ser desfeita."
        confirm-label="Excluir"
        cancel-label="Cancelar"
        variant="destructive"
        @confirm="confirmDelete"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useInvestmentsStore } from '@/stores/investments'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Loader2, AlertCircle, Trash2, Check, ArrowLeft } from 'lucide-vue-next'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import type { Investment } from '@/api/types'

const route = useRoute()
const router = useRouter()
const investmentsStore = useInvestmentsStore()

const investmentId = route.params.id as string
const isDeleting = ref(false)
const isUpdating = ref(false)
const showDeleteDialog = ref(false)
const updateForm = ref({
  current_value: '',
  quantity: '',
})

onMounted(async () => {
  await loadInvestment()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadInvestment()
    }
  }
)

async function loadInvestment() {
  if (!investmentId) return

  const existingInvestment = investmentsStore.investments.find(
    (inv) => inv.investment_id === investmentId
  )

  if (existingInvestment && !investmentsStore.currentInvestment) {
    investmentsStore.currentInvestment = existingInvestment
  } else {
    try {
      await investmentsStore.getInvestment(investmentId)
      // Preencher formulário de atualização
      if (investmentsStore.currentInvestment) {
        updateForm.value.current_value = investmentsStore.currentInvestment.current_value
        updateForm.value.quantity = investmentsStore.currentInvestment.quantity || ''
      }
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

function goBack() {
  router.push('/investments')
}

function handleRetry() {
  investmentsStore.clearError()
  loadInvestment()
}

async function handleUpdate() {
  if (!investmentsStore.currentInvestment) return

  isUpdating.value = true
  try {
    const updateData: any = {}
    if (updateForm.value.current_value) {
      updateData.current_value = parseFloat(updateForm.value.current_value)
    }
    if (updateForm.value.quantity) {
      updateData.quantity = parseFloat(updateForm.value.quantity)
    }

    await investmentsStore.updateInvestment(investmentId, updateData)
    
    // Atualizar formulário
    if (investmentsStore.currentInvestment) {
      updateForm.value.current_value = investmentsStore.currentInvestment.current_value
      updateForm.value.quantity = investmentsStore.currentInvestment.quantity || ''
    }
  } catch (error) {
    // Erro já é tratado no store
  } finally {
    isUpdating.value = false
  }
}

async function handleDelete() {
  showDeleteDialog.value = true
}

async function confirmDelete() {
  if (!investmentsStore.currentInvestment) return

  isDeleting.value = true
  try {
    await investmentsStore.deleteInvestment(investmentId)
    router.push('/investments')
  } catch (error) {
    // Erro já é tratado no store
  } finally {
    isDeleting.value = false
    showDeleteDialog.value = false
  }
}

function getInvestmentTypeLabel(type: Investment['type']): string {
  const labels: Record<Investment['type'], string> = {
    STOCK: 'Ação',
    FUND: 'Fundo',
    CDB: 'CDB',
    TREASURY: 'Tesouro Direto',
    CRYPTO: 'Criptomoeda',
    OTHER: 'Outro',
  }
  return labels[type] || type
}

function getContextLabel(context: Investment['context']): string {
  return context === 'PERSONAL' ? 'Pessoal' : 'Negócio'
}

function formatCurrency(amount: string, currency: Investment['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function formatReturnPercentage(percentage: number): string {
  const sign = percentage >= 0 ? '+' : ''
  return `${sign}${percentage.toFixed(2)}%`
}

function formatQuantity(quantity: string): string {
  const value = parseFloat(quantity)
  return value.toLocaleString('pt-BR', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 4,
  })
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}
</script>

