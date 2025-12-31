<template>
  <div
    v-if="needRefresh"
    class="fixed bottom-4 left-4 right-4 z-50 rounded-lg border bg-background p-4 shadow-lg sm:left-auto sm:max-w-md"
  >
    <div class="flex items-start gap-3">
      <div class="flex-1">
        <h3 class="font-semibold">Nova versão disponível</h3>
        <p class="text-sm text-muted-foreground">
          Uma nova versão da aplicação está disponível. A página será recarregada automaticamente.
        </p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" size="sm" @click="dismiss">
          Depois
        </Button>
        <Button size="sm" @click="update">
          Atualizar Agora
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Button } from '@/components/ui/button'

const needRefresh = ref(false)

const update = () => {
  // Com registerType: 'autoUpdate', o service worker já atualiza automaticamente
  // Apenas recarregar a página
  window.location.reload()
}

const dismiss = () => {
  needRefresh.value = false
  // Verificar novamente após 5 minutos
  setTimeout(() => {
    checkForUpdate()
  }, 5 * 60 * 1000)
}

const checkForUpdate = () => {
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.ready.then((registration) => {
      registration.update()
      
      // Escutar atualizações disponíveis
      registration.addEventListener('updatefound', () => {
        const newWorker = registration.installing
        if (newWorker) {
          newWorker.addEventListener('statechange', () => {
            if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
              // Nova versão disponível
              needRefresh.value = true
            }
          })
        }
      })
    })
  }
}

onMounted(() => {
  // Verificar atualizações a cada hora
  setInterval(() => {
    checkForUpdate()
  }, 60 * 60 * 1000)
  
  // Verificar imediatamente
  checkForUpdate()
  
  // Escutar mensagens do service worker
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      window.location.reload()
    })
  }
})
</script>

