<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs :items="[
        { label: 'Categorias', to: '/categories' },
        { label: category?.name || 'Carregando...' }
      ]" />

      <!-- Loading State -->
      <div v-if="categoriesStore.isLoading" class="flex items-center justify-center py-12">
        <div class="text-center">
          <Loader2 class="mx-auto h-12 w-12 text-primary mb-4 animate-spin" />
          <p class="text-muted-foreground">Carregando categoria...</p>
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
            <Button
              variant="outline"
              @click="goBack"
            >
              Voltar
            </Button>
          </div>
        </CardContent>
      </Card>

      <!-- Category Details -->
      <div v-else-if="category">
        <!-- Header -->
        <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-4xl font-bold mb-2">{{ category.name }}</h1>
            <p class="text-muted-foreground">
              Detalhes da categoria
            </p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="outline"
              @click="goBack"
            >
              Voltar
            </Button>
            <Button
              variant="destructive"
              @click="handleDelete"
              :disabled="isDeleting"
            >
              <Loader2 v-if="isDeleting" class="h-4 w-4 mr-2 animate-spin" />
              <Trash2 v-else class="h-4 w-4 mr-2" />
              Deletar
            </Button>
          </div>
        </div>

        <!-- Category Info -->
        <div class="grid gap-6 md:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Informações</CardTitle>
            </CardHeader>
            <CardContent class="space-y-4">
              <div>
                <Label class="text-muted-foreground">Nome</Label>
                <p class="text-lg font-semibold">{{ category.name }}</p>
              </div>
              <div v-if="category.description">
                <Label class="text-muted-foreground">Descrição</Label>
                <p class="text-base">{{ category.description }}</p>
              </div>
              <div>
                <Label class="text-muted-foreground">Status</Label>
                <Badge :variant="category.is_active ? 'default' : 'secondary'">
                  {{ category.is_active ? 'Ativa' : 'Inativa' }}
                </Badge>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Datas</CardTitle>
            </CardHeader>
            <CardContent class="space-y-4">
              <div>
                <Label class="text-muted-foreground">Criada em</Label>
                <p class="text-base">{{ formatDate(category.created_at) }}</p>
              </div>
              <div>
                <Label class="text-muted-foreground">Atualizada em</Label>
                <p class="text-base">{{ formatDate(category.updated_at) }}</p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useCategoriesStore } from '@/stores/categories'
import Layout from '@/components/layout/Layout.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Loader2, AlertCircle, Trash2 } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const categoriesStore = useCategoriesStore()

const isDeleting = ref(false)

const category = computed(() => categoriesStore.currentCategory)

onMounted(async () => {
  const categoryId = route.params.id as string
  if (categoryId) {
    await categoriesStore.getCategory(categoryId)
  }
})

function handleRetry() {
  categoriesStore.clearError()
  const categoryId = route.params.id as string
  if (categoryId) {
    categoriesStore.getCategory(categoryId)
  }
}

function goBack() {
  router.push('/categories')
}

async function handleDelete() {
  if (!category.value) return

  // Por enquanto, usar confirm nativo
  if (!confirm(`Tem certeza que deseja deletar a categoria "${category.value.name}"?`)) {
    return
  }

  isDeleting.value = true
  try {
    await categoriesStore.deleteCategory(category.value.category_id)
    
    // Mostrar toast de sucesso
    const { toast } = await import('@/components/ui/toast')
    toast.success('Categoria deletada com sucesso!')
    
    router.push('/categories')
  } catch (err: any) {
    console.error('Erro ao deletar categoria:', err)
    
    // Mostrar toast de erro
    const { toast } = await import('@/components/ui/toast')
    toast.error('Erro ao deletar categoria', {
      description: err.message || 'Tente novamente mais tarde.',
    })
  } finally {
    isDeleting.value = false
  }
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}
</script>

