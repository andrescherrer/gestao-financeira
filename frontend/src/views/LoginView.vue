<template>
  <div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-4">
    <div class="w-full max-w-md">
      <!-- Card Principal -->
      <Card class="p-8 shadow-2xl">
        <!-- Logo e Título -->
        <div class="mb-8 text-center">
          <div class="mb-4 flex justify-center">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-indigo-600 shadow-lg">
              <Wallet class="h-8 w-8 text-white" />
            </div>
          </div>
          <h1 class="text-3xl font-bold text-foreground">Bem-vindo de volta</h1>
          <p class="mt-2 text-sm text-muted-foreground">
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
          <Card v-if="error" class="mb-6 border-destructive">
            <CardContent class="p-4">
              <div class="flex items-center gap-2">
                <AlertCircle class="h-4 w-4 text-destructive" />
                <p class="text-sm font-medium text-destructive">{{ error }}</p>
              </div>
            </CardContent>
          </Card>

          <!-- Campo Email -->
          <div class="mb-5">
            <Label for="email" class="mb-2">
              Endereço de Email
            </Label>
            <div class="relative">
              <Field
                id="email"
                name="email"
                type="email"
                v-slot="{ field, meta }"
              >
                <Input
                  id="email"
                  :name="field.name"
                  :value="field.value"
                  @input="field.onInput"
                  @change="field.onChange"
                  @blur="field.onBlur"
                  :class="errors.email || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
                  placeholder="seu@email.com"
                  class="pr-10"
                />
              </Field>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
                <Mail class="h-4 w-4 text-muted-foreground" />
              </div>
            </div>
            <ErrorMessage name="email" class="mt-1 text-sm text-destructive" />
          </div>

          <!-- Campo Senha -->
          <div class="mb-5">
            <Label for="password" class="mb-2">
              Senha
            </Label>
            <div class="relative">
              <Field
                id="password"
                name="password"
                :type="showPassword ? 'text' : 'password'"
                v-slot="{ field, meta }"
              >
                <Input
                  id="password"
                  :name="field.name"
                  :type="showPassword ? 'text' : 'password'"
                  :value="field.value"
                  @input="field.onInput"
                  @change="field.onChange"
                  @blur="field.onBlur"
                  :class="errors.password || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
                  placeholder="••••••••"
                  class="pr-10"
                />
              </Field>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                class="absolute inset-y-0 right-0 h-full"
                @click="showPassword = !showPassword"
              >
                <Eye v-if="!showPassword" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </Button>
            </div>
            <ErrorMessage name="password" class="mt-1 text-sm text-destructive" />
          </div>

          <!-- Remember me e Forgot password -->
          <div class="mb-6 flex items-center justify-between">
            <label class="flex items-center gap-2 cursor-pointer">
              <input
                v-model="rememberMe"
                type="checkbox"
                class="h-4 w-4 rounded border-input text-primary focus:ring-2 focus:ring-ring"
              />
              <span class="text-sm font-medium text-foreground">Lembrar-me</span>
            </label>
            <router-link
              to="/forgot-password"
              class="text-sm font-semibold text-green-600 transition-colors hover:text-green-700"
            >
              Esqueceu sua senha?
            </router-link>
          </div>

          <!-- Botão de Login -->
          <Button
            type="submit"
            :disabled="isLoading"
            class="w-full"
            size="lg"
          >
            <User v-if="!isLoading" class="h-5 w-5" />
            <Loader2 v-else class="h-5 w-5 animate-spin" />
            {{ isLoading ? 'Entrando...' : 'Entrar' }}
          </Button>
        </Form>
      </Card>

      <!-- Footer -->
      <div class="mt-6 text-center text-sm text-muted-foreground">
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
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Wallet, Mail, Eye, EyeOff, User, Loader2, AlertCircle } from 'lucide-vue-next'
import { translateError } from '@/utils/errorTranslations'

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
    const savedToken = localStorage.getItem('auth_token')
    if (!savedToken) {
      throw new Error('Token não foi salvo corretamente')
    }
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    const rawError =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao fazer login. Verifique suas credenciais.'
    error.value = translateError(rawError)
  }
}
</script>
