<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
    <div class="w-full max-w-md">
      <!-- Card Principal -->
      <div class="rounded-2xl bg-white p-8 shadow-2xl">
        <!-- Logo e Título -->
        <div class="mb-8 text-center">
          <div class="mb-4 flex justify-center">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-indigo-600 shadow-lg">
              <i class="pi pi-wallet text-3xl text-white"></i>
            </div>
          </div>
          <h1 class="text-3xl font-bold text-gray-900">Bem-vindo de volta</h1>
          <p class="mt-2 text-sm text-gray-600">
            Não tem uma conta?
            <router-link
              to="/register"
              class="font-semibold text-green-600 transition-colors hover:text-green-700"
            >
              Criar conta agora!
            </router-link>
          </p>
        </div>

        <!-- Formulário -->
        <Form @submit="handleSubmit" :validation-schema="validationSchema" v-slot="{ errors }">
          <!-- Mensagem de Erro -->
          <div
            v-if="error"
            class="mb-6 rounded-lg border border-red-200 bg-red-50 p-4"
          >
            <div class="flex items-center gap-2">
              <i class="pi pi-exclamation-circle text-red-600"></i>
              <p class="text-sm font-medium text-red-600">{{ error }}</p>
            </div>
          </div>

          <!-- Campo Email -->
          <div class="mb-5">
            <label
              for="email"
              class="mb-2 block text-sm font-semibold text-gray-700"
            >
              Endereço de Email
            </label>
            <div class="relative">
              <Field
                id="email"
                name="email"
                type="email"
                class="block w-full rounded-lg border px-4 py-3 pr-10 text-gray-900 placeholder-gray-400 shadow-sm transition-all focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                :class="
                  errors.email
                    ? 'border-red-300 focus:border-red-500 focus:ring-red-500/20'
                    : 'border-gray-300'
                "
                placeholder="seu@email.com"
              />
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
                <i class="pi pi-envelope text-gray-400"></i>
              </div>
            </div>
            <ErrorMessage name="email" class="mt-1 text-sm text-red-600" />
          </div>

          <!-- Campo Senha -->
          <div class="mb-5">
            <label
              for="password"
              class="mb-2 block text-sm font-semibold text-gray-700"
            >
              Senha
            </label>
            <div class="relative">
              <Field
                id="password"
                name="password"
                :type="showPassword ? 'text' : 'password'"
                class="block w-full rounded-lg border px-4 py-3 pr-10 text-gray-900 placeholder-gray-400 shadow-sm transition-all focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20"
                :class="
                  errors.password
                    ? 'border-red-300 focus:border-red-500 focus:ring-red-500/20'
                    : 'border-gray-300'
                "
                placeholder="••••••••"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
              >
                <i :class="showPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
              </button>
            </div>
            <ErrorMessage name="password" class="mt-1 text-sm text-red-600" />
          </div>

          <!-- Remember me e Forgot password -->
          <div class="mb-6 flex items-center justify-between">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="rememberMe"
                type="checkbox"
                class="h-4 w-4 rounded border-gray-300 text-green-600 focus:ring-2 focus:ring-green-500/20"
              />
              <span class="text-sm font-medium text-gray-700">Lembrar-me</span>
            </label>
            <router-link
              to="/forgot-password"
              class="text-sm font-semibold text-green-600 transition-colors hover:text-green-700"
            >
              Esqueceu sua senha?
            </router-link>
          </div>

          <!-- Botão de Login -->
          <button
            type="submit"
            :disabled="isLoading"
            class="flex w-full items-center justify-center gap-3 rounded-lg bg-gradient-to-r from-green-500 to-green-600 px-6 py-3.5 text-base font-semibold text-white shadow-lg transition-all hover:from-green-600 hover:to-green-700 hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <i v-if="!isLoading" class="pi pi-user text-lg"></i>
            <i v-else class="pi pi-spinner pi-spin text-lg"></i>
            <span>{{ isLoading ? 'Entrando...' : 'Entrar' }}</span>
          </button>
        </Form>
      </div>

      <!-- Footer -->
      <div class="mt-6 text-center text-sm text-gray-600">
        <p>© 2025 Gestão Financeira. Todos os direitos reservados.</p>
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
const showPassword = ref(false)
const rememberMe = ref(false)

async function handleSubmit(values: any) {
  error.value = null
  try {
    const response = await authStore.login(values)
    // Verificar se o token foi salvo
    const savedToken = localStorage.getItem('auth_token')
    if (!savedToken) {
      throw new Error('Token não foi salvo corretamente')
    }
    // Redirecionar após login bem-sucedido
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
