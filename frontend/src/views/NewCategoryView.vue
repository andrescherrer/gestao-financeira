<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[
        { label: 'Categorias', to: '/categories' },
        { label: 'Nova Categoria' }
      ]" />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Nova Categoria</h1>
        <p class="text-muted-foreground">
          Crie uma nova categoria para organizar suas transações
        </p>
      </div>

      <!-- Form -->
      <Card>
        <CardHeader>
          <CardTitle>Informações da Categoria</CardTitle>
          <CardDescription>
            Preencha os dados abaixo para criar uma nova categoria
          </CardDescription>
        </CardHeader>
        <CardContent>
          <CategoryForm
            ref="formRef"
            :is-submitting="isSubmitting"
            :error="error"
            @submit="handleSubmit"
            @cancel="goBack"
          />
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import CategoryForm from '@/components/CategoryForm.vue'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { toTypedSchema } from '@vee-validate/zod'
import { createCategorySchema } from '@/validations/category'
import type { CreateCategoryFormData, UpdateCategoryFormData } from '@/validations/category'
import type { CreateCategoryRequest } from '@/api/types'

const router = useRouter()
const categoriesStore = useCategoriesStore()
const formRef = ref<InstanceType<typeof CategoryForm> | null>(null)

const isSubmitting = ref(false)
const error = ref<string | null>(null)

async function handleSubmit(values: CreateCategoryFormData | UpdateCategoryFormData) {
  if (formRef.value) {
    // Clear previous errors
    error.value = null
  }

  isSubmitting.value = true

  try {
    // Garantir que name existe (CreateCategoryFormData sempre tem name)
    if (!values.name) {
      throw new Error('Nome é obrigatório')
    }

    const trimmedName = values.name.trim()
    if (trimmedName.length < 2) {
      throw new Error('Nome deve ter no mínimo 2 caracteres')
    }

    const categoryData: CreateCategoryRequest = {
      name: trimmedName,
      description: values.description?.trim() || undefined,
    }

    if (import.meta.env.DEV) {
      console.log('[NewCategoryView] Dados sendo enviados:', categoryData)
    }

    const category = await categoriesStore.createCategory(categoryData)
    router.push(`/categories/${category.category_id}`)
  } catch (err: any) {
    if (import.meta.env.DEV) {
      console.error('[NewCategoryView] Erro ao criar categoria:', {
        message: err.message,
        response: err.response?.data,
        status: err.response?.status,
        statusText: err.response?.statusText,
      })
    }

    const errorMessage =
      err.response?.data?.error ||
      err.response?.data?.message ||
      err.message ||
      'Erro ao criar categoria. Tente novamente.'
    
    error.value = errorMessage
  } finally {
    isSubmitting.value = false
  }
}

function goBack() {
  router.push('/categories')
}
</script>

