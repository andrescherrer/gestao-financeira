<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-lg">Filtros</CardTitle>
    </CardHeader>
    <CardContent>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
        <!-- Tipo -->
        <div>
          <Label for="type" class="mb-1">Tipo</Label>
          <Select 
            :model-value="localFilters.type" 
            @update:model-value="handleTypeChange"
          >
            <SelectTrigger id="type">
              <SelectValue placeholder="Todos os tipos" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todos os tipos</SelectItem>
              <SelectItem value="INCOME">Receitas</SelectItem>
              <SelectItem value="EXPENSE">Despesas</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- Conta -->
        <div>
          <Label for="account" class="mb-1">Conta</Label>
          <Select v-model="localFilters.accountId" @update:model-value="handleChange">
            <SelectTrigger id="account">
              <SelectValue placeholder="Todas as contas" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todas as contas</SelectItem>
              <SelectItem
                v-for="account in accounts"
                :key="account.account_id"
                :value="account.account_id"
              >
                {{ account.name }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <!-- Data Inicial -->
        <div>
          <Label for="startDate" class="mb-1">Data Inicial</Label>
          <Input
            id="startDate"
            type="date"
            :value="localFilters.startDate || ''"
            @update:model-value="(val: string | number) => { localFilters.startDate = String(val); handleChange() }"
          />
        </div>

        <!-- Data Final -->
        <div>
          <Label for="endDate" class="mb-1">Data Final</Label>
          <Input
            id="endDate"
            type="date"
            :value="localFilters.endDate || ''"
            @update:model-value="(val: string | number) => { localFilters.endDate = String(val); handleChange() }"
          />
        </div>
      </div>

      <div class="mt-4 flex gap-2">
        <Button
          v-if="hasActiveFilters"
          variant="outline"
          size="sm"
          @click="clearFilters"
        >
          <X class="h-4 w-4 mr-2" />
          Limpar Filtros
        </Button>
      </div>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { X } from 'lucide-vue-next'
import type { Account } from '@/api/types'
import type { AcceptableValue } from 'reka-ui'

interface Filters {
  type?: 'INCOME' | 'EXPENSE' | ''
  accountId?: string
  startDate?: string
  endDate?: string
}

interface Props {
  filters: Filters
  accounts: Account[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:filters': [filters: Filters]
  'clear': []
}>()

const localFilters = ref<Filters>({ ...props.filters })

watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })

const hasActiveFilters = computed(() => {
  return !!(
    localFilters.value.type ||
    localFilters.value.accountId ||
    localFilters.value.startDate ||
    localFilters.value.endDate
  )
})

function handleTypeChange(val: AcceptableValue) {
  if (val === null || val === undefined) {
    localFilters.value.type = undefined
  } else if (val === '' || val === 'INCOME' || val === 'EXPENSE') {
    localFilters.value.type = val as 'INCOME' | 'EXPENSE' | ''
  } else {
    // Se for outro tipo (número, objeto, etc), converter para string e validar
    const strVal = String(val)
    if (strVal === 'INCOME' || strVal === 'EXPENSE') {
      localFilters.value.type = strVal
    } else if (strVal === '') {
      localFilters.value.type = ''
    } else {
      localFilters.value.type = undefined
    }
  }
  handleChange()
}

function handleChange() {
  emit('update:filters', { ...localFilters.value })
}

function clearFilters() {
  localFilters.value = {}
  emit('clear')
}
</script>

