<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[{ label: 'Categorias' }]" />

      <!-- Header -->
      <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-4xl font-bold mb-2">Categorias</h1>
          <p class="text-muted-foreground">
            Gerencie suas categorias de transações
          </p>
        </div>
        <Button as-child>
          <router-link to="/categories/new">
            <Plus class="h-4 w-4 mr-2" />
            Nova Categoria
          </router-link>
        </Button>
      </div>

      <!-- Loading State -->
      <div v-if="categoriesStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando categorias...</p>
        </div>
      </div>

      <!-- Error State -->
      <Card
        v-else-if="categoriesStore.error"
        class="mb-6 border-destructive"
      >
        <CardContent class="p-4">
          <div class="flex items-center gap-2 mb-3">
            <AlertCircle class="h-4 w-4 text-destructive" />
            <p class="text-destructive">{{ categoriesStore.error }}</p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="link"
              @click="handleRetry"
              class="text-destructive"
            >
              Tentar novamente
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Empty State -->
      <Card
        v-else-if="categoriesStore.categories.length === 0"
      >
        <CardContent class="p-12 text-center">
          <Tag class="mx-auto h-16 w-16 text-muted-foreground mb-4" />
          <h3 class="text-xl font-semibold text-foreground mb-2">
            Nenhuma categoria encontrada
          </h3>
          <p class="text-muted-foreground mb-6">
            Comece criando sua primeira categoria
          </p>
          <Button as-child>
            <router-link to="/categories/new">
              <Plus class="h-4 w-4 mr-2" />
              Criar Primeira Categoria
            </router-link>
          </Button>
        </CardContent>
      </Card>

      <!-- Categories List -->
      <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card
          v-for="category in categoriesStore.categories"
          :key="category.category_id"
          class="hover:shadow-md transition-shadow cursor-pointer"
          @click="goToCategoryDetails(category.category_id)"
        >
          <CardHeader>
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <CardTitle class="text-lg">{{ category.name }}</CardTitle>
                <CardDescription v-if="category.description" class="mt-1">
                  {{ category.description }}
                </CardDescription>
              </div>
              <Badge :variant="category.is_active ? 'default' : 'secondary'">
                {{ category.is_active ? 'Ativa' : 'Inativa' }}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <div class="flex items-center justify-between text-sm text-muted-foreground">
              <span>
                Criada em {{ formatDate(category.created_at) }}
              </span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Loader2, Plus, Tag, AlertCircle } from 'lucide-vue-next'

const router = useRouter()
const categoriesStore = useCategoriesStore()

onMounted(async () => {
  if (categoriesStore.categories.length === 0) {
    await categoriesStore.listCategories()
  }
})

function handleRetry() {
  categoriesStore.clearError()
  categoriesStore.listCategories()
}

function goToCategoryDetails(categoryId: string) {
  router.push(`/categories/${categoryId}`)
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
}
</script>

