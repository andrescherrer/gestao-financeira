<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-md space-y-8 rounded-lg border bg-white p-8 shadow-lg">
      <div class="text-center">
        <h1 class="text-3xl font-bold">Gestão Financeira</h1>
        <p class="mt-2 text-gray-600">
          Faça login para acessar sua conta
        </p>
      </div>

      <Form @submit="handleSubmit" :validation-schema="validationSchema" v-slot="{ errors }">
        <div v-if="error" class="rounded-md bg-red-50 p-3 text-sm text-red-600">
          {{ error }}
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700">
            Email
          </label>
          <Field
            id="email"
            name="email"
            type="email"
            class="mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
            :class="
              errors.email
                ? 'border-red-300 focus:border-red-500'
                : 'border-gray-300 focus:border-blue-500'
            "
            placeholder="seu@email.com"
          />
          <ErrorMessage name="email" class="mt-1 text-sm text-red-600" />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700">
            Senha
          </label>
          <Field
            id="password"
            name="password"
            type="password"
            class="mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
            :class="
              errors.password
                ? 'border-red-300 focus:border-red-500'
                : 'border-gray-300 focus:border-blue-500'
            "
            placeholder="••••••••"
          />
          <ErrorMessage name="password" class="mt-1 text-sm text-red-600" />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {{ isLoading ? 'Entrando...' : 'Entrar' }}
        </button>
      </Form>

      <div class="text-center text-sm">
        <span class="text-gray-600">Não tem uma conta? </span>
        <router-link to="/register" class="text-blue-600 hover:underline">
          Criar conta
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useAuthStore } from '@/stores/auth'
import { loginSchema } from '@/validations/auth'

const validationSchema = toTypedSchema(loginSchema)

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const error = ref<string | null>(null)
const isLoading = computed(() => authStore.isLoading)

async function handleSubmit(values: any) {
  error.value = null
  try {
    await authStore.login(values)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    error.value =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao fazer login. Verifique suas credenciais.'
  }
}
</script>

