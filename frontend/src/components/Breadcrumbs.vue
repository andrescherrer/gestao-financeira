<template>
  <nav v-if="items.length > 0" class="mb-6 flex items-center gap-2 text-sm text-gray-600">
    <template v-for="(item, index) in items" :key="index">
      <router-link
        v-if="item.to && index < items.length - 1"
        :to="item.to"
        class="transition-colors hover:text-gray-900"
      >
        {{ item.label }}
      </router-link>
      <span v-else class="font-medium text-gray-900">{{ item.label }}</span>
      <i
        v-if="index < items.length - 1"
        class="pi pi-chevron-right text-xs text-gray-400"
      ></i>
    </template>
  </nav>
</template>

<script setup lang="ts">
interface BreadcrumbItem {
  label: string
  to?: string
}

defineProps<{
  items: BreadcrumbItem[]
}>()
</script>

