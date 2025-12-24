<template>
  <nav v-if="items.length > 0" class="mb-6 flex items-center gap-2 text-sm text-muted-foreground">
    <template v-for="(item, index) in items" :key="index">
      <router-link
        v-if="item.to && index < items.length - 1"
        :to="item.to"
        class="transition-colors hover:text-foreground"
      >
        {{ item.label }}
      </router-link>
      <span v-else class="font-medium text-foreground">{{ item.label }}</span>
      <ChevronRight
        v-if="index < items.length - 1"
        class="h-4 w-4 text-muted-foreground"
      />
    </template>
  </nav>
</template>

<script setup lang="ts">
import { ChevronRight } from 'lucide-vue-next'

interface BreadcrumbItem {
  label: string
  to?: string
}

defineProps<{
  items: BreadcrumbItem[]
}>()
</script>

<script setup lang="ts">
interface BreadcrumbItem {
  label: string
  to?: string
}

defineProps<{
  items: BreadcrumbItem[]
}>()
</script>

