<template>
  <Card
    class="group relative cursor-pointer overflow-hidden transition-all hover:border-primary hover:shadow-lg"
    @click="handleClick"
    role="button"
    :aria-label="`Investimento ${investment.name}, valor ${formatCurrency(investment.current_value, investment.currency)}`"
    tabindex="0"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <!-- Background gradient effect -->
    <div
      class="absolute inset-0 opacity-0 transition-opacity group-hover:opacity-5"
      :class="getInvestmentGradient(investment.type)"
    ></div>

    <CardContent class="relative p-6">
      <div class="mb-4 flex items-start justify-between">
        <div class="flex items-center gap-3">
          <div
            class="flex h-12 w-12 items-center justify-center rounded-lg"
            :class="getInvestmentIconBg(investment.type)"
          >
            <component :is="getInvestmentIcon(investment.type)" class="h-6 w-6 text-white" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-foreground">
              {{ investment.name }}
              <span v-if="investment.ticker" class="text-sm text-muted-foreground">
                ({{ investment.ticker }})
              </span>
            </h3>
            <p class="text-sm text-muted-foreground">
              {{ getInvestmentTypeLabel(investment.type) }}
            </p>
          </div>
        </div>
        <Badge
          :variant="investment.return_percentage >= 0 ? 'default' : 'destructive'"
        >
          {{ formatReturnPercentage(investment.return_percentage) }}
        </Badge>
      </div>

      <div class="mb-4 space-y-2">
        <div class="rounded-lg bg-muted p-4">
          <div class="text-xs font-medium uppercase tracking-wide text-muted-foreground mb-1">
            Valor Atual
          </div>
          <div class="text-2xl font-bold">
            {{ formatCurrency(investment.current_value, investment.currency) }}
          </div>
        </div>
        <div class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Valor de Compra:</span>
          <span class="font-medium">
            {{ formatCurrency(investment.purchase_amount, investment.currency) }}
          </span>
        </div>
        <div v-if="investment.quantity" class="flex items-center justify-between text-sm">
          <span class="text-muted-foreground">Quantidade:</span>
          <span class="font-medium">{{ formatQuantity(investment.quantity) }}</span>
        </div>
      </div>

      <div class="flex items-center justify-between">
        <Badge
          :variant="investment.context === 'PERSONAL' ? 'secondary' : 'outline'"
          class="bg-purple-100 text-purple-700"
        >
          {{ getContextLabel(investment.context) }}
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
import {
  TrendingUp,
  BarChart3,
  Coins,
  Shield,
  Bitcoin,
  Package,
  ChevronRight,
} from 'lucide-vue-next'
import type { Investment } from '@/api/types'

interface Props {
  investment: Investment
}

const props = defineProps<Props>()
const router = useRouter()

function handleClick() {
  router.push(`/investments/${props.investment.investment_id}`)
}

function getInvestmentTypeLabel(type: Investment['type']): string {
  const labels: Record<Investment['type'], string> = {
    STOCK: 'Ação',
    FUND: 'Fundo',
    CDB: 'CDB',
    TREASURY: 'Tesouro',
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

function getInvestmentIcon(type: Investment['type']) {
  const icons: Record<Investment['type'], any> = {
    STOCK: TrendingUp,
    FUND: BarChart3,
    CDB: Shield,
    TREASURY: Shield,
    CRYPTO: Bitcoin,
    OTHER: Package,
  }
  return icons[type] || Package
}

function getInvestmentIconBg(type: Investment['type']): string {
  const colors: Record<Investment['type'], string> = {
    STOCK: 'bg-blue-500',
    FUND: 'bg-purple-500',
    CDB: 'bg-green-500',
    TREASURY: 'bg-yellow-500',
    CRYPTO: 'bg-orange-500',
    OTHER: 'bg-gray-500',
  }
  return colors[type] || 'bg-gray-500'
}

function getInvestmentGradient(type: Investment['type']): string {
  const gradients: Record<Investment['type'], string> = {
    STOCK: 'bg-gradient-to-br from-blue-500 to-blue-600',
    FUND: 'bg-gradient-to-br from-purple-500 to-purple-600',
    CDB: 'bg-gradient-to-br from-green-500 to-green-600',
    TREASURY: 'bg-gradient-to-br from-yellow-500 to-yellow-600',
    CRYPTO: 'bg-gradient-to-br from-orange-500 to-orange-600',
    OTHER: 'bg-gradient-to-br from-gray-500 to-gray-600',
  }
  return gradients[type] || 'bg-gradient-to-br from-gray-500 to-gray-600'
}
</script>

