<template>
  <Dialog :open="open" @update:open="handleOpenChange">
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{{ title }}</DialogTitle>
        <DialogDescription v-if="description">
          {{ description }}
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button variant="outline" @click="handleCancel">
          {{ cancelLabel }}
        </Button>
        <Button :variant="confirmVariant" @click="handleConfirm" :disabled="isLoading">
          <Loader2 v-if="isLoading" class="h-4 w-4 mr-2 animate-spin" />
          {{ confirmLabel }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-vue-next'

interface Props {
  open: boolean
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  confirmVariant?: 'default' | 'destructive'
  isLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  confirmLabel: 'Confirmar',
  cancelLabel: 'Cancelar',
  confirmVariant: 'default',
  isLoading: false,
})

const emit = defineEmits<{
  'update:open': [open: boolean]
  'confirm': []
  'cancel': []
}>()

function handleOpenChange(open: boolean) {
  emit('update:open', open)
}

function handleConfirm() {
  emit('confirm')
}

function handleCancel() {
  emit('update:open', false)
  emit('cancel')
}
</script>

