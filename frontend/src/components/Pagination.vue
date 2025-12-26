<template>
  <div class="flex items-center justify-between px-2 py-4">
    <div class="text-sm text-muted-foreground">
      Mostrando {{ startItem }} a {{ endItem }} de {{ total }} itens
    </div>
    <div class="flex items-center gap-2">
      <Button
        variant="outline"
        size="sm"
        :disabled="currentPage === 1"
        @click="goToPage(currentPage - 1)"
      >
        <ChevronLeft class="h-4 w-4" />
        Anterior
      </Button>
      <div class="flex items-center gap-1">
        <Button
          v-for="page in visiblePages"
          :key="page"
          variant="outline"
          size="sm"
          :class="page === currentPage ? 'bg-primary text-primary-foreground' : ''"
          @click="goToPage(page)"
        >
          {{ page }}
        </Button>
      </div>
      <Button
        variant="outline"
        size="sm"
        :disabled="currentPage === totalPages"
        @click="goToPage(currentPage + 1)"
      >
        Próxima
        <ChevronRight class="h-4 w-4" />
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

interface Props {
  currentPage: number
  totalPages: number
  total: number
  itemsPerPage: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:currentPage': [page: number]
}>()

const startItem = computed(() => {
  return props.total === 0 ? 0 : (props.currentPage - 1) * props.itemsPerPage + 1
})

const endItem = computed(() => {
  const end = props.currentPage * props.itemsPerPage
  return end > props.total ? props.total : end
})

const visiblePages = computed(() => {
  const pages: number[] = []
  const maxVisible = 5
  let start = Math.max(1, props.currentPage - Math.floor(maxVisible / 2))
  let end = Math.min(props.totalPages, start + maxVisible - 1)

  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

function goToPage(page: number) {
  if (page >= 1 && page <= props.totalPages) {
    emit('update:currentPage', page)
  }
}
</script>

