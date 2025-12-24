<template>
  <div
    class="group relative cursor-pointer overflow-hidden rounded-xl border border-gray-200 bg-white p-6 shadow-sm transition-all hover:border-blue-300 hover:shadow-lg"
    @click="handleClick"
  >
    <!-- Background gradient effect -->
    <div
      class="absolute inset-0 opacity-0 transition-opacity group-hover:opacity-5"
      :class="getAccountGradient(account.type)"
    ></div>

    <div class="relative">
      <div class="mb-4 flex items-start justify-between">
        <div class="flex items-center gap-3">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-lg"
            :class="getAccountIconBg(account.type)"
          >
            <i
              class="pi text-xl text-white"
              :class="getAccountIcon(account.type)"
            ></i>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">
              {{ account.name }}
            </h3>
            <p class="text-sm text-gray-500">
              {{ getAccountTypeLabel(account.type) }}
            </p>
          </div>
        </div>
        <span
          class="rounded-full px-3 py-1 text-xs font-semibold"
          :class="
            account.is_active
              ? 'bg-green-100 text-green-700'
              : 'bg-gray-100 text-gray-700'
          "
        >
          {{ account.is_active ? 'Ativa' : 'Inativa' }}
        </span>
      </div>

      <div class="mb-4 rounded-lg bg-gray-50 p-4">
        <div class="text-xs font-medium uppercase tracking-wide text-gray-500 mb-1">Saldo</div>
        <div class="text-3xl font-bold" :class="getBalanceColor(account.balance)">
          {{ formatCurrency(account.balance, account.currency) }}
        </div>
      </div>

      <div class="flex items-center justify-between">
        <span
          class="rounded-full px-3 py-1 text-xs font-medium"
          :class="
            account.context === 'PERSONAL'
              ? 'bg-purple-100 text-purple-700'
              : 'bg-indigo-100 text-indigo-700'
          "
        >
          {{ getContextLabel(account.context) }}
        </span>
        <i
          class="pi pi-chevron-right text-gray-400 transition-all group-hover:translate-x-1 group-hover:text-blue-600"
        ></i>
      </div>
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

function getAccountIcon(type: Account['type']): string {
  const icons: Record<Account['type'], string> = {
    BANK: 'pi-building',
    WALLET: 'pi-wallet',
    INVESTMENT: 'pi-chart-line',
    CREDIT_CARD: 'pi-credit-card',
  }
  return icons[type] || 'pi-wallet'
}

function getAccountIconBg(type: Account['type']): string {
  const colors: Record<Account['type'], string> = {
    BANK: 'bg-blue-500',
    WALLET: 'bg-green-500',
    INVESTMENT: 'bg-purple-500',
    CREDIT_CARD: 'bg-orange-500',
  }
  return colors[type] || 'bg-gray-500'
}

function getAccountGradient(type: Account['type']): string {
  const gradients: Record<Account['type'], string> = {
    BANK: 'bg-gradient-to-br from-blue-500 to-blue-600',
    WALLET: 'bg-gradient-to-br from-green-500 to-green-600',
    INVESTMENT: 'bg-gradient-to-br from-purple-500 to-purple-600',
    CREDIT_CARD: 'bg-gradient-to-br from-orange-500 to-orange-600',
  }
  return gradients[type] || 'bg-gradient-to-br from-gray-500 to-gray-600'
}
</script>

