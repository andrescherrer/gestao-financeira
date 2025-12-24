<template>
  <div class="rounded-lg border border-gray-200 bg-white overflow-hidden">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Data
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Descrição
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Tipo
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Valor
          </th>
          <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
            Conta
          </th>
          <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
            Ações
          </th>
        </tr>
      </thead>
      <tbody class="bg-white divide-y divide-gray-200">
        <tr
          v-for="transaction in transactions"
          :key="transaction.transaction_id"
          class="hover:bg-gray-50 transition-colors"
        >
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
            {{ formatDate(transaction.date) }}
          </td>
          <td class="px-6 py-4 text-sm text-gray-900">
            <div class="max-w-xs truncate" :title="transaction.description">
              {{ transaction.description }}
            </div>
          </td>
          <td class="px-6 py-4 whitespace-nowrap">
            <span
              class="px-2 py-1 text-xs font-semibold rounded-full"
              :class="
                transaction.type === 'INCOME'
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800'
              "
            >
              {{ transaction.type === 'INCOME' ? 'Receita' : 'Despesa' }}
            </span>
          </td>
          <td
            class="px-6 py-4 whitespace-nowrap text-sm font-medium"
            :class="
              transaction.type === 'INCOME' ? 'text-green-600' : 'text-red-600'
            "
          >
            {{ formatCurrency(parseFloat(transaction.amount)) }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
            {{ getAccountName(transaction.account_id) }}
          </td>
          <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
            <router-link
              :to="`/transactions/${transaction.transaction_id}`"
              class="text-blue-600 hover:text-blue-900 transition-colors"
            >
              Ver detalhes
            </router-link>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Empty State -->
    <div
      v-if="transactions.length === 0"
      class="px-6 py-12 text-center"
    >
      <i class="pi pi-inbox text-4xl text-gray-400 mb-4"></i>
      <p class="text-gray-600">Nenhuma transação encontrada</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAccountsStore } from '@/stores/accounts'
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
  
  // Tentar parsear a data
  const date = new Date(dateString)
  
  // Verificar se a data é válida
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

