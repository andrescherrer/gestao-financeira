<template>
  <div class="flex min-h-screen flex-col">
    <Header />
    <div class="flex flex-1">
      <Sidebar :is-open="sidebarOpen" />
      <main class="flex-1 md:ml-64">
        <div class="min-h-[calc(100vh-4rem)] p-4 sm:p-6 lg:p-8">
          <slot />
        </div>
      </main>
    </div>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Header from './Header.vue'
import Sidebar from './Sidebar.vue'
import Footer from './Footer.vue'

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

