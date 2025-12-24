<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Transações', to: '/transactions' },
          { label: transactionsStore.currentTransaction?.description || 'Detalhes' },
        ]"
      />

      <!-- Loading State -->
      <div v-if="transactionsStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando detalhes da transação...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="transactionsStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-4">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ transactionsStore.error }}</p>
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

      <!-- Transaction Details -->
      <div v-else-if="transactionsStore.currentTransaction" class="space-y-6">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold text-foreground mb-2">
              {{ transactionsStore.currentTransaction.description }}
            </h1>
            <p class="text-muted-foreground">
              {{ getTransactionTypeLabel(transactionsStore.currentTransaction.type) }}
            </p>
          </div>
          <Badge
            :variant="transactionsStore.currentTransaction.type === 'INCOME' ? 'default' : 'destructive'"
            class="bg-green-100 text-green-700 hover:bg-green-100"
            v-if="transactionsStore.currentTransaction.type === 'INCOME'"
          >
            Receita
          </Badge>
          <Badge
            variant="destructive"
            class="bg-red-100 text-red-700 hover:bg-red-100"
            v-else
          >
            Despesa
          </Badge>
        </div>

        <!-- Transaction Card -->
        <Card>
          <CardContent class="p-6">
            <!-- Amount -->
            <div class="mb-6 rounded-lg bg-gradient-to-br p-6"
              :class="
                transactionsStore.currentTransaction.type === 'INCOME'
                  ? 'from-green-50 to-green-100'
                  : 'from-red-50 to-red-100'
              ">
              <div class="text-sm font-medium mb-2"
                :class="
                  transactionsStore.currentTransaction.type === 'INCOME'
                    ? 'text-green-700'
                    : 'text-red-700'
                ">
                Valor
              </div>
              <div
                class="text-4xl font-bold"
                :class="
                  transactionsStore.currentTransaction.type === 'INCOME'
                    ? 'text-green-600'
                    : 'text-red-600'
                "
              >
                {{ formatCurrency(
                  parseFloat(transactionsStore.currentTransaction.amount),
                  transactionsStore.currentTransaction.currency
                ) }}
              </div>
            </div>

            <!-- Transaction Information Grid -->
            <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Tipo
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getTransactionTypeLabel(transactionsStore.currentTransaction.type) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Conta
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ getAccountName(transactionsStore.currentTransaction.account_id) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDate(transactionsStore.currentTransaction.date) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Moeda
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ transactionsStore.currentTransaction.currency }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Data de Criação
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDateTime(transactionsStore.currentTransaction.created_at) }}
                </div>
              </div>

              <div class="rounded-lg border border-border bg-muted/50 p-4">
                <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
                  Última Atualização
                </div>
                <div class="text-lg font-semibold text-foreground">
                  {{ formatDateTime(transactionsStore.currentTransaction.updated_at) }}
                </div>
              </div>
            </div>

            <!-- Description -->
            <div class="mt-6 rounded-lg border border-border bg-muted/50 p-4">
              <div class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-2">
                Descrição
              </div>
              <div class="text-lg font-medium text-foreground">
                {{ transactionsStore.currentTransaction.description }}
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
            Transação não encontrada
          </h3>
          <p class="text-muted-foreground mb-6">
            A transação que você está procurando não existe ou foi removida.
          </p>
          <Button
            @click="goBack"
            as-child
          >
            <router-link to="/transactions">
              <ArrowLeft class="h-4 w-4 mr-2" />
              Voltar para transações
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
import { useTransactionsStore } from '@/stores/transactions'
import { useAccountsStore } from '@/stores/accounts'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, AlertCircle, ArrowLeft } from 'lucide-vue-next'
import type { Transaction } from '@/api/types'

const route = useRoute()
const router = useRouter()
const transactionsStore = useTransactionsStore()
const accountsStore = useAccountsStore()

const transactionId = route.params.id as string

onMounted(async () => {
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
  await loadTransaction()
})

watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadTransaction()
    }
  }
)

async function loadTransaction() {
  if (!transactionId) return

  const existingTransaction = transactionsStore.transactions.find(
    (tx) => tx.transaction_id === transactionId
  )

  if (existingTransaction && !transactionsStore.currentTransaction) {
    transactionsStore.currentTransaction = existingTransaction
  } else {
    try {
      await transactionsStore.getTransaction(transactionId)
    } catch (error) {
      // Erro já é tratado no store
    }
  }
}

function goBack() {
  router.push('/transactions')
}

function handleRetry() {
  transactionsStore.clearError()
  loadTransaction()
}

function getTransactionTypeLabel(type: Transaction['type']): string {
  return type === 'INCOME' ? 'Receita' : 'Despesa'
}

function getAccountName(accountId: string): string {
  const account = accountsStore.accounts.find((acc) => acc.account_id === accountId)
  return account?.name || 'Conta não encontrada'
}

function formatCurrency(value: number, currency: Transaction['currency']): string {
  const formatter = new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency || 'BRL',
  })
  return formatter.format(value)
}

function formatDate(dateString: string): string {
  if (!dateString) return 'Data inválida'
  
  const date = new Date(dateString)
  
  if (isNaN(date.getTime())) {
    console.warn('[TransactionDetailsView] Data inválida:', dateString)
    return 'Data inválida'
  }
  
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

function formatDateTime(dateString: string): string {
  if (!dateString) return 'Data inválida'
  
  const date = new Date(dateString)
  
  if (isNaN(date.getTime())) {
    console.warn('[TransactionDetailsView] Data/hora inválida:', dateString)
    return 'Data inválida'
  }
  
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
</script>
