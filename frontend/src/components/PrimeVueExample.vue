<template>
  <div class="p-4">
    <h2 class="text-2xl font-bold mb-4">PrimeVue Configurado</h2>
    <p class="mb-4">PrimeVue está instalado e configurado corretamente!</p>
    
    <!-- Exemplo de componente PrimeVue -->
    <div class="space-y-4">
      <Button label="Botão PrimeVue" icon="pi pi-check" />
      <InputText v-model="value" placeholder="Input PrimeVue" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'

const value = ref('')
</script>

