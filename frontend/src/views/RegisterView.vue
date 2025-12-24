<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-md space-y-8 rounded-lg border bg-white p-8 shadow-lg">
      <div class="text-center">
        <h1 class="text-3xl font-bold">Criar Conta</h1>
        <p class="mt-2 text-gray-600">
          Preencha os dados para criar sua conta
        </p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div v-if="error" class="rounded-md bg-red-50 p-3 text-sm text-red-600">
          {{ error }}
        </div>

        <div v-if="success" class="rounded-md bg-green-50 p-3 text-sm text-green-600">
          Conta criada com sucesso! Redirecionando...
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="firstName" class="block text-sm font-medium text-gray-700">
              Nome
            </label>
            <input
              id="firstName"
              v-model="form.first_name"
              type="text"
              required
              minlength="2"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
            />
          </div>

          <div>
            <label for="lastName" class="block text-sm font-medium text-gray-700">
              Sobrenome
            </label>
            <input
              id="lastName"
              v-model="form.last_name"
              type="text"
              required
              minlength="2"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
            />
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">
            Email
          </label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">
            Senha
          </label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            minlength="8"
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
          />
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
            Confirmar Senha
          </label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            required
            minlength="8"
            class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500"
          />
          <p v-if="form.password && confirmPassword && form.password !== confirmPassword" class="mt-1 text-sm text-red-600">
            As senhas não coincidem
          </p>
        </div>

        <button
          type="submit"
          :disabled="isLoading || success || form.password !== confirmPassword"
          class="w-full rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {{ isLoading ? 'Criando...' : success ? 'Conta criada!' : 'Criar Conta' }}
        </button>
      </form>

      <div class="text-center text-sm">
        <span class="text-gray-600">Já tem uma conta? </span>
        <router-link to="/login" class="text-blue-600 hover:underline">
          Fazer login
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
})

const confirmPassword = ref('')
const error = ref<string | null>(null)
const success = ref(false)
const isLoading = computed(() => authStore.isLoading)

async function handleSubmit() {
  if (form.value.password !== confirmPassword.value) {
    error.value = 'As senhas não coincidem'
    return
  }

  error.value = null
  try {
    await authStore.register(form.value)
    success.value = true
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (err: any) {
    error.value =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao criar conta. Tente novamente.'
  }
}
</script>

