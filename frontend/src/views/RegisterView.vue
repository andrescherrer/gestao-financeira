<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-md space-y-8 rounded-lg border bg-white p-8 shadow-lg">
      <div class="text-center">
        <h1 class="text-3xl font-bold">Criar Conta</h1>
        <p class="mt-2 text-gray-600">
          Preencha os dados para criar sua conta
        </p>
      </div>

      <Form @submit="handleSubmit" :validation-schema="validationSchema" v-slot="{ errors }">
        <div v-if="error" class="rounded-md bg-red-50 p-3 text-sm text-red-600">
          {{ error }}
        </div>

        <div v-if="success" class="rounded-md bg-green-50 p-3 text-sm text-green-600">
          Conta criada com sucesso! Redirecionando...
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="first_name" class="block text-sm font-medium text-gray-700">
              Nome
            </label>
            <Field
              id="first_name"
              name="first_name"
              type="text"
              class="mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
              :class="
                errors.first_name
                  ? 'border-red-300 focus:border-red-500'
                  : 'border-gray-300 focus:border-blue-500'
              "
            />
            <ErrorMessage name="first_name" class="mt-1 text-sm text-red-600" />
          </div>

          <div>
            <label for="last_name" class="block text-sm font-medium text-gray-700">
              Sobrenome
            </label>
            <Field
              id="last_name"
              name="last_name"
              type="text"
              class="mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
              :class="
                errors.last_name
                  ? 'border-red-300 focus:border-red-500'
                  : 'border-gray-300 focus:border-blue-500'
              "
            />
            <ErrorMessage name="last_name" class="mt-1 text-sm text-red-600" />
          </div>
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
          />
          <ErrorMessage name="password" class="mt-1 text-sm text-red-600" />
          <p class="mt-1 text-xs text-gray-500">
            Mínimo 8 caracteres, com letra maiúscula, minúscula e número
          </p>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
            Confirmar Senha
          </label>
          <Field
            id="confirmPassword"
            name="confirmPassword"
            type="password"
            class="mt-1 block w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-blue-500"
            :class="
              errors.confirmPassword
                ? 'border-red-300 focus:border-red-500'
                : 'border-gray-300 focus:border-blue-500'
            "
          />
          <ErrorMessage name="confirmPassword" class="mt-1 text-sm text-red-600" />
        </div>

        <button
          type="submit"
          :disabled="isLoading || success"
          class="w-full rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {{ isLoading ? 'Criando...' : success ? 'Conta criada!' : 'Criar Conta' }}
        </button>
      </Form>

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
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useAuthStore } from '@/stores/auth'
import { registerSchema } from '@/validations/auth'

const validationSchema = toTypedSchema(registerSchema)

const router = useRouter()
const authStore = useAuthStore()

const error = ref<string | null>(null)
const success = ref(false)
const isLoading = computed(() => authStore.isLoading)

async function handleSubmit(values: any) {
  error.value = null
  try {
    // Separar confirmPassword do payload
    const { confirmPassword, ...registerData } = values
    await authStore.register(registerData)
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

