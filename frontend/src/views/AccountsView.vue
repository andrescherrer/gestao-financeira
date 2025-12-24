<template>
  <Layout>
    <div>
      <h1 class="text-4xl font-bold mb-2">Contas</h1>
      <p class="text-gray-600 mb-6">
        Gerencie suas contas bancárias e financeiras
      </p>
      <router-link
        to="/accounts/new"
        class="inline-block rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
      >
        Nova Conta
      </router-link>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Layout from '@/components/layout/Layout.vue'
</script>

