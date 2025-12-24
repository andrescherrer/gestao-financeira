<template>
  <div
    class="group cursor-pointer rounded-lg border border-gray-200 bg-white p-6 transition-all hover:border-blue-300 hover:shadow-md"
    @click="handleClick"
  >
    <div class="mb-4 flex items-start justify-between">
      <div>
        <h3 class="text-lg font-semibold text-gray-900">
          {{ account.name }}
        </h3>
        <p class="text-sm text-gray-500">
          {{ getAccountTypeLabel(account.type) }}
        </p>
      </div>
      <span
        class="rounded-full px-2 py-1 text-xs font-medium"
        :class="
          account.is_active
            ? 'bg-green-100 text-green-700'
            : 'bg-gray-100 text-gray-700'
        "
      >
        {{ account.is_active ? 'Ativa' : 'Inativa' }}
      </span>
    </div>

    <div class="mb-4">
      <div class="text-sm text-gray-600">Saldo</div>
      <div class="text-2xl font-bold" :class="getBalanceColor(account.balance)">
        {{ formatCurrency(account.balance, account.currency) }}
      </div>
    </div>

    <div class="flex items-center justify-between text-sm text-gray-500">
      <span>{{ getContextLabel(account.context) }}</span>
      <i
        class="pi pi-chevron-right text-gray-400 transition-colors group-hover:text-blue-600"
      ></i>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import type { Account } from '@/api/types'

interface Props {
  account: Account
}

const props = defineProps<Props>()
const router = useRouter()

function handleClick() {
  router.push(`/accounts/${props.account.account_id}`)
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
  return 'text-gray-900'
}
</script>

