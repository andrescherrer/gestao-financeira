<template>
  <Card>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="cursor-pointer select-none" @click="sortBy('date')">
            <div class="flex items-center gap-2">
              Data
              <component
                :is="sortIcon('date')"
                class="h-4 w-4"
              />
            </div>
          </TableHead>
          <TableHead class="cursor-pointer select-none" @click="sortBy('description')">
            <div class="flex items-center gap-2">
              Descrição
              <component
                :is="sortIcon('description')"
                class="h-4 w-4"
              />
            </div>
          </TableHead>
          <TableHead class="cursor-pointer select-none" @click="sortBy('type')">
            <div class="flex items-center gap-2">
              Tipo
              <component
                :is="sortIcon('type')"
                class="h-4 w-4"
              />
            </div>
          </TableHead>
          <TableHead class="cursor-pointer select-none text-right" @click="sortBy('amount')">
            <div class="flex items-center justify-end gap-2">
              Valor
              <component
                :is="sortIcon('amount')"
                class="h-4 w-4"
              />
            </div>
          </TableHead>
          <TableHead>Conta</TableHead>
          <TableHead class="text-right">Ações</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow
          v-for="transaction in sortedTransactions"
          :key="transaction.transaction_id"
          class="hover:bg-muted/50 transition-colors"
        >
          <TableCell class="whitespace-nowrap">
            {{ formatDate(transaction.date) }}
          </TableCell>
          <TableCell>
            <div class="max-w-xs truncate" :title="transaction.description">
              {{ transaction.description }}
            </div>
          </TableCell>
          <TableCell class="whitespace-nowrap">
            <Badge
              :variant="transaction.type === 'INCOME' ? 'default' : 'destructive'"
              class="bg-green-100 text-green-800 hover:bg-green-100"
              v-if="transaction.type === 'INCOME'"
            >
              Receita
            </Badge>
            <Badge
              variant="destructive"
              class="bg-red-100 text-red-800 hover:bg-red-100"
              v-else
            >
              Despesa
            </Badge>
          </TableCell>
          <TableCell
            class="whitespace-nowrap font-medium text-right"
            :class="transaction.type === 'INCOME' ? 'text-green-600' : 'text-red-600'"
          >
            {{ formatCurrency(parseFloat(transaction.amount)) }}
          </TableCell>
          <TableCell class="whitespace-nowrap text-muted-foreground">
            {{ getAccountName(transaction.account_id) }}
          </TableCell>
          <TableCell class="text-right">
            <Button variant="link" as-child>
              <router-link :to="`/transactions/${transaction.transaction_id}`">
                Ver detalhes
              </router-link>
            </Button>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <!-- Empty State -->
    <TableEmpty v-if="transactions.length === 0">
      <Inbox class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
      <p class="text-muted-foreground">Nenhuma transação encontrada</p>
    </TableEmpty>
  </Card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { Card } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableEmpty, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ArrowUpDown, ArrowUp, ArrowDown } from 'lucide-vue-next'
import type { Transaction } from '@/api/types'

interface Props {
  transactions: Transaction[]
}

const props = defineProps<Props>()

const accountsStore = useAccountsStore()

type SortField = 'date' | 'description' | 'type' | 'amount'
type SortOrder = 'asc' | 'desc' | null

const sortField = ref<SortField | null>(null)
const sortOrder = ref<SortOrder>(null)

function sortBy(field: SortField) {
  if (sortField.value === field) {
    // Toggle order: null -> asc -> desc -> null
    if (sortOrder.value === null) {
      sortOrder.value = 'asc'
    } else if (sortOrder.value === 'asc') {
      sortOrder.value = 'desc'
    } else {
      sortField.value = null
      sortOrder.value = null
    }
  } else {
    sortField.value = field
    sortOrder.value = 'asc'
  }
}

function sortIcon(field: SortField) {
  if (sortField.value !== field) {
    return ArrowUpDown
  }
  if (sortOrder.value === 'asc') {
    return ArrowUp
  }
  if (sortOrder.value === 'desc') {
    return ArrowDown
  }
  return ArrowUpDown
}

const sortedTransactions = computed(() => {
  if (!sortField.value || !sortOrder.value) {
    return props.transactions
  }

  const sorted = [...props.transactions]

  sorted.sort((a, b) => {
    let aVal: any
    let bVal: any

    switch (sortField.value) {
      case 'date':
        aVal = new Date(a.date).getTime()
        bVal = new Date(b.date).getTime()
        break
      case 'description':
        aVal = a.description.toLowerCase()
        bVal = b.description.toLowerCase()
        break
      case 'type':
        aVal = a.type
        bVal = b.type
        break
      case 'amount':
        aVal = parseFloat(a.amount)
        bVal = parseFloat(b.amount)
        break
      default:
        return 0
    }

    if (aVal < bVal) {
      return sortOrder.value === 'asc' ? -1 : 1
    }
    if (aVal > bVal) {
      return sortOrder.value === 'asc' ? 1 : -1
    }
    return 0
  })

  return sorted
})

function getAccountName(accountId: string): string {
  const account = accountsStore.accounts.find((acc) => acc.account_id === accountId)
  return account?.name || 'Conta não encontrada'
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value)
}
</script>
