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
          <h1 class="text-3xl font-bold text-foreground">Criar Conta</h1>
          <p class="mt-2 text-sm text-muted-foreground">
            Já tem uma conta?
            <router-link
              to="/login"
              class="font-semibold text-green-600 transition-colors hover:text-green-700"
            >
              Fazer login
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

          <!-- Mensagem de Sucesso -->
          <Card v-if="success" class="mb-6 border-green-200 bg-green-50">
            <CardContent class="p-4">
              <div class="flex items-center gap-2">
                <CheckCircle class="h-4 w-4 text-green-600" />
                <p class="text-sm font-medium text-green-600">Conta criada com sucesso! Redirecionando...</p>
              </div>
            </CardContent>
          </Card>

          <!-- Nome e Sobrenome -->
          <div class="mb-5 grid grid-cols-2 gap-4">
            <div>
              <Label for="first_name" class="mb-2">
                Nome
              </Label>
              <Field
                id="first_name"
                name="first_name"
                type="text"
                v-slot="{ field, meta }"
              >
                <Input
                  id="first_name"
                  :name="field.name"
                  :value="field.value"
                  @input="field.onInput"
                  @change="field.onChange"
                  @blur="field.onBlur"
                  :class="errors.first_name || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
                  placeholder="Seu nome"
                />
              </Field>
              <ErrorMessage name="first_name" class="mt-1 text-sm text-destructive" />
            </div>

            <div>
              <Label for="last_name" class="mb-2">
                Sobrenome
              </Label>
              <Field
                id="last_name"
                name="last_name"
                type="text"
                v-slot="{ field, meta }"
              >
                <Input
                  id="last_name"
                  :name="field.name"
                  :value="field.value"
                  @input="field.onInput"
                  @change="field.onChange"
                  @blur="field.onBlur"
                  :class="errors.last_name || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
                  placeholder="Seu sobrenome"
                />
              </Field>
              <ErrorMessage name="last_name" class="mt-1 text-sm text-destructive" />
            </div>
          </div>

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
            <p class="mt-1 text-xs text-muted-foreground">
              Mínimo 8 caracteres, com letra maiúscula, minúscula e número
            </p>
          </div>

          <!-- Campo Confirmar Senha -->
          <div class="mb-6">
            <Label for="confirmPassword" class="mb-2">
              Confirmar Senha
            </Label>
            <div class="relative">
              <Field
                id="confirmPassword"
                name="confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                v-slot="{ field, meta }"
              >
                <Input
                  id="confirmPassword"
                  :name="field.name"
                  :type="showConfirmPassword ? 'text' : 'password'"
                  :value="field.value"
                  @input="field.onInput"
                  @change="field.onChange"
                  @blur="field.onBlur"
                  :class="errors.confirmPassword || (meta.touched && !meta.valid) ? 'border-destructive' : ''"
                  placeholder="••••••••"
                  class="pr-10"
                />
              </Field>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                class="absolute inset-y-0 right-0 h-full"
                @click="showConfirmPassword = !showConfirmPassword"
              >
                <Eye v-if="!showConfirmPassword" class="h-4 w-4" />
                <EyeOff v-else class="h-4 w-4" />
              </Button>
            </div>
            <ErrorMessage name="confirmPassword" class="mt-1 text-sm text-destructive" />
          </div>

          <!-- Botão de Registro -->
          <Button
            type="submit"
            :disabled="isLoading || success"
            class="w-full"
            size="lg"
          >
            <UserPlus v-if="!isLoading && !success" class="h-5 w-5" />
            <Loader2 v-else-if="isLoading" class="h-5 w-5 animate-spin" />
            <CheckCircle v-else class="h-5 w-5" />
            {{ isLoading ? 'Criando...' : success ? 'Conta criada!' : 'Criar Conta' }}
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
import { useRouter } from 'vue-router'
import { Form, Field, ErrorMessage } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useAuthStore } from '@/stores/auth'
import { registerSchema } from '@/validations/auth'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Wallet, Mail, Eye, EyeOff, UserPlus, Loader2, AlertCircle, CheckCircle } from 'lucide-vue-next'

const validationSchema = toTypedSchema(registerSchema)

const router = useRouter()
const authStore = useAuthStore()

const error = ref<string | null>(null)
const success = ref(false)
const isLoading = computed(() => authStore.isLoading)
const showPassword = ref(false)
const showConfirmPassword = ref(false)

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
