<template>
  <Card
    class="group relative cursor-pointer overflow-hidden transition-all hover:border-primary hover:shadow-lg"
    @click="handleClick"
  >
    <!-- Background gradient effect -->
    <div
      class="absolute inset-0 opacity-0 transition-opacity group-hover:opacity-5"
      :class="getAccountGradient(account.type)"
    ></div>

    <CardContent class="relative p-6">
      <div class="mb-4 flex items-start justify-between">
        <div class="flex items-center gap-3">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-lg"
            :class="getAccountIconBg(account.type)"
          >
            <component :is="getAccountIcon(account.type)" class="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-foreground">
              {{ account.name }}
            </h3>
            <p class="text-sm text-muted-foreground">
              {{ getAccountTypeLabel(account.type) }}
            </p>
          </div>
        </div>
        <Badge
          :variant="account.is_active ? 'default' : 'secondary'"
        >
          {{ account.is_active ? 'Ativa' : 'Inativa' }}
        </Badge>
      </div>

      <div class="mb-4 rounded-lg bg-muted p-4">
        <div class="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-1">Saldo</div>
        <div class="text-3xl font-bold" :class="getBalanceColor(account.balance)">
          {{ formatCurrency(account.balance, account.currency) }}
        </div>
      </div>

      <div class="flex items-center justify-between">
        <Badge
          :variant="account.context === 'PERSONAL' ? 'secondary' : 'outline'"
          class="bg-purple-100 text-purple-700"
        >
          {{ getContextLabel(account.context) }}
        </Badge>
        <ChevronRight
          class="h-5 w-5 text-muted-foreground transition-all group-hover:translate-x-1 group-hover:text-primary"
        />
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Building2, Wallet, TrendingUp, CreditCard, ChevronRight } from 'lucide-vue-next'
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
  return 'text-foreground'
}

function getAccountIcon(type: Account['type']) {
  const icons: Record<Account['type'], any> = {
    BANK: Building2,
    WALLET: Wallet,
    INVESTMENT: TrendingUp,
    CREDIT_CARD: CreditCard,
  }
  return icons[type] || Wallet
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
