<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Contas', to: '/accounts' },
          { label: accountsStore.currentAccount?.name || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="accountsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando detalhes da conta...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="accountsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ accountsStore.error }}</p>
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

      <!-- Account Details -->
      <div v-else-if="accountsStore.currentAccount" class="space-y-6">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold text-foreground mb-2">
              {{ accountsStore.currentAccount.name }}
            </h1>
            <p class="text-muted-foreground">
              {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
            </p>
          </div>
          <div class="flex items-center gap-3">
            <Badge
              :variant="accountsStore.currentAccount.is_active ? 'default' : 'secondary'"
            >
              {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
            </Badge>
            <Button
              as-child
              variant="default"
            >
              <router-link
                :to="`/accounts/${accountsStore.currentAccount.account_id}/edit`"
              >
                <Pencil class="h-4 w-4 mr-2" />
                Editar
              </router-link>
            </Button>
          </div>
        </div>

        <!-- Account Card -->
        <Card>
          <CardContent class="p-6">
            <!-- Balance -->
            <div class="mb-6 rounded-lg bg-gradient-to-br from-muted to-muted/50 p-6">
              <div class="text-sm font-medium text-muted-foreground mb-2">Saldo Atual</div>
              <div
                class="text-4xl font-bold"
                :class="getBalanceColor(accountsStore.currentAccount.balance)"
              >
                {{ formatCurrency(accountsStore.currentAccount.balance, accountsStore.currentAccount.currency) }}
              </div>
            </div>

            <!-- Account Information Grid -->
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Tipo de Conta
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getAccountTypeLabel(accountsStore.currentAccount.type) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Contexto
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getContextLabel(accountsStore.currentAccount.context) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Moeda
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ accountsStore.currentAccount.currency }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Status
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ accountsStore.currentAccount.is_active ? 'Ativa' : 'Inativa' }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data de Criação
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(accountsStore.currentAccount.created_at) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Última Atualização
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(accountsStore.currentAccount.updated_at) }}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Not Found State -->
      <Card v-else>
        <CardContent class="p-12 text-center">
          <AlertCircle class="mx-auto h-16 w-16 text-muted-foreground mb-4" />
          <h3 class="text-xl font-semibold text-foreground mb-2">
            Conta não encontrada
          </h3>
          <p class="text-muted-foreground mb-6">
            A conta que você está procurando não existe ou foi removida.
          </p>
          <Button
            @click="goBack"
            as-child
          >
            <router-link to="/accounts">
              <ArrowLeft class="h-4 w-4 mr-2" />
              Voltar para contas
            </router-link>
          </Button>
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, AlertCircle, Pencil, ArrowLeft } from 'lucide-vue-next'
import type { Account } from '@/api/types'

const route = useRoute()
const router = useRouter()
const accountsStore = useAccountsStore()

const accountId = route.params.id as string

onMounted(async () => {
  await loadAccount()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadAccount()
    }
  }
)

async function loadAccount() {
  if (!accountId) return

  const existingAccount = accountsStore.accounts.find(
    (acc) => acc.account_id === accountId
  )

  if (existingAccount && !accountsStore.currentAccount) {
    accountsStore.currentAccount = existingAccount
  } else {
    try {
      await accountsStore.getAccount(accountId)
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

function goBack() {
  router.push('/accounts')
}

function handleRetry() {
  accountsStore.clearError()
  loadAccount()
}

function getAccountTypeLabel(type: Account['type']): string {
  const labels: Record<Account['type'], string> = {
    BANK: 'Banco',
    WALLET: 'Carteira',
    INVESTMENT: 'Investimento',
    CREDIT_CARD: 'Cartão de Crédito',
  }
  return labels[type] || type
}

function getContextLabel(context: Account['context']): string {
  return context === 'PERSONAL' ? 'Pessoal' : 'Negócio'
}

function formatCurrency(amount: string, currency: Account['currency']): string {
  const value = parseFloat(amount)
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function getBalanceColor(balance: string): string {
  const value = parseFloat(balance)
  if (value > 0) return 'text-green-600'
  if (value < 0) return 'text-red-600'
  return 'text-foreground'
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
</script>
