<template>
  <Card>
    <CardContent class="p-12 text-center">
      <component
        :is="icon"
        class="mx-auto h-16 w-16 text-muted-foreground mb-4"
      />
      <h3 class="text-xl font-semibold text-foreground mb-2">
        {{ title }}
      </h3>
      <p class="text-muted-foreground mb-6">
        {{ description }}
      </p>
      <slot name="action">
        <Button v-if="actionLabel && actionTo" as-child>
          <router-link :to="actionTo">
            <component v-if="actionIcon" :is="actionIcon" class="h-4 w-4 mr-2" />
            {{ actionLabel }}
          </router-link>
        </Button>
      </slot>
    </CardContent>
  </Card>
</template>

<script setup lang="ts">
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import type { Component } from 'vue'

interface Props {
  icon: Component
  title: string
  description: string
  actionLabel?: string
  actionTo?: string
  actionIcon?: Component
}

withDefaults(defineProps<Props>(), {
  actionLabel: undefined,
  actionTo: undefined,
  actionIcon: undefined,
})
</script>

