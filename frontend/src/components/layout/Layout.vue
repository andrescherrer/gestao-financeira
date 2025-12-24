<template>
  <div class="flex min-h-screen bg-gray-50">
    <Sidebar :is-open="sidebarOpen" />
    <div class="flex flex-1 flex-col md:ml-64">
      <Header />
      <main class="flex-1 overflow-y-auto">
        <div class="p-6">
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'

const sidebarOpen = ref(true)

function handleResize() {
  // Em telas grandes, sidebar sempre visível
  // Em telas pequenas, pode ser controlado por um botão
  if (window.innerWidth >= 768) {
    sidebarOpen.value = true
  }
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>
