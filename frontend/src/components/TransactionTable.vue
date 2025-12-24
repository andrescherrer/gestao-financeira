<template>
  <Card>
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Data</TableHead>
          <TableHead>Descrição</TableHead>
          <TableHead>Tipo</TableHead>
          <TableHead>Valor</TableHead>
          <TableHead>Conta</TableHead>
          <TableHead class="text-right">Ações</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow
          v-for="transaction in transactions"
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
            class="whitespace-nowrap font-medium"
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
    <div
      v-if="transactions.length === 0"
      class="px-6 py-12 text-center"
    >
      <Inbox class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
      <p class="text-muted-foreground">Nenhuma transação encontrada</p>
    </div>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
import { Card } from '@/components/ui/card'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Inbox } from 'lucide-vue-next'
import type { Transaction } from '@/api/types'

interface Props {
  transactions: Transaction[]
}

const props = defineProps<Props>()

const accountsStore = useAccountsStore()

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(value)
}

function formatDate(dateString: string): string {
  if (!dateString) return 'Data inválida'
  
  const date = new Date(dateString)
  
  if (isNaN(date.getTime())) {
    console.warn('[TransactionTable] Data inválida:', dateString)
    return 'Data inválida'
  }
  
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}

function getAccountName(accountId: string): string {
  const account = accountsStore.accounts.find((acc) => acc.account_id === accountId)
  return account?.name || 'Conta não encontrada'
}
</script>
