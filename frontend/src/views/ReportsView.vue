<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Relatórios' }]" />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Relatórios</h1>
        <p class="text-muted-foreground">
          Visualize seus dados financeiros com gráficos e análises
        </p>
      </div>

      <!-- Filters -->
      <Card class="mb-6">
        <CardContent class="p-4">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
            <div>
              <label class="text-sm font-medium mb-2 block">Período</label>
              <Select v-model="filters.period" @update:model-value="handleFilterChange">
                <SelectTrigger>
                  <SelectValue placeholder="Selecione" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="monthly">Mensal</SelectItem>
                  <SelectItem value="annual">Anual</SelectItem>
                  <SelectItem value="custom">Personalizado</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div v-if="filters.period === 'monthly' || filters.period === 'annual'">
              <label class="text-sm font-medium mb-2 block">Ano</label>
              <Input
                v-model.number="filters.year"
                type="number"
                :min="2020"
                :max="2100"
                placeholder="Ano"
                @input="handleFilterChange"
              />
            </div>
            <div v-if="filters.period === 'monthly'">
              <label class="text-sm font-medium mb-2 block">Mês</label>
              <Select v-model="filters.month" @update:model-value="handleFilterChange">
                <SelectTrigger>
                  <SelectValue placeholder="Mês" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="1">Janeiro</SelectItem>
                  <SelectItem :value="2">Fevereiro</SelectItem>
                  <SelectItem :value="3">Março</SelectItem>
                  <SelectItem :value="4">Abril</SelectItem>
                  <SelectItem :value="5">Maio</SelectItem>
                  <SelectItem :value="6">Junho</SelectItem>
                  <SelectItem :value="7">Julho</SelectItem>
                  <SelectItem :value="8">Agosto</SelectItem>
                  <SelectItem :value="9">Setembro</SelectItem>
                  <SelectItem :value="10">Outubro</SelectItem>
                  <SelectItem :value="11">Novembro</SelectItem>
                  <SelectItem :value="12">Dezembro</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <label class="text-sm font-medium mb-2 block">Moeda</label>
              <Select v-model="filters.currency" @update:model-value="handleFilterChange">
                <SelectTrigger>
                  <SelectValue placeholder="Moeda" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="BRL">BRL - Real</SelectItem>
                  <SelectItem value="USD">USD - Dólar</SelectItem>
                  <SelectItem value="EUR">EUR - Euro</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Reports Content -->
      <div class="space-y-6">
        <!-- Income vs Expense Chart -->
        <Card>
          <CardHeader>
            <CardTitle>Receitas vs Despesas</CardTitle>
            <CardDescription>
              Comparação entre receitas e despesas no período selecionado
            </CardDescription>
          </CardHeader>
          <CardContent>
            <IncomeVsExpenseChart
              :year="filters.year || currentYear"
              :month="filters.month"
              :currency="filters.currency || 'BRL'"
              :period="filters.period"
            />
          </CardContent>
        </Card>

        <!-- Category Chart -->
        <Card>
          <CardHeader>
            <CardTitle>Gastos por Categoria</CardTitle>
            <CardDescription>
              Distribuição de gastos por categoria
            </CardDescription>
          </CardHeader>
          <CardContent>
            <CategoryChart
              :year="filters.year || currentYear"
              :month="filters.month"
              :currency="filters.currency || 'BRL'"
            />
          </CardContent>
        </Card>

        <!-- Trends Chart -->
        <Card>
          <CardHeader>
            <CardTitle>Tendências Temporais</CardTitle>
            <CardDescription>
              Evolução de receitas e despesas ao longo do tempo
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TrendsChart
              :year="filters.year || currentYear"
              :currency="filters.currency || 'BRL'"
              :period="filters.period"
            />
          </CardContent>
        </Card>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import IncomeVsExpenseChart from '@/components/reports/IncomeVsExpenseChart.vue'
import CategoryChart from '@/components/reports/CategoryChart.vue'
import TrendsChart from '@/components/reports/TrendsChart.vue'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const currentYear = new Date().getFullYear()
const currentMonth = new Date().getMonth() + 1

const filters = ref<{
  period: 'monthly' | 'annual' | 'custom'
  year?: number
  month?: number
  currency?: 'BRL' | 'USD' | 'EUR'
}>({
  period: 'monthly',
  year: currentYear,
  month: currentMonth,
  currency: 'BRL',
})

function handleFilterChange() {
  // Os componentes de gráfico vão reagir automaticamente às mudanças
}
</script>

