<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Metas', to: '/goals' },
          { label: 'Nova Meta' },
        ]"
      />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Nova Meta</h1>
        <p class="text-muted-foreground">
          Defina uma nova meta financeira e acompanhe seu progresso
        </p>
      </div>

      <!-- Form -->
      <Card>
        <CardContent class="p-6">
          <GoalForm
            ref="formRef"
            :isLoading="goalsStore.isLoading"
            submitLabel="Criar Meta"
            @submit="handleSubmit"
            @cancel="goBack"
          />
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGoalsStore } from '@/stores/goals'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import GoalForm from '@/components/GoalForm.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import type { CreateGoalFormData } from '@/validations/goal'

const router = useRouter()
const goalsStore = useGoalsStore()
const authStore = useAuthStore()
const formRef = ref<InstanceType<typeof GoalForm> | null>(null)

onMounted(async () => {
  if (!authStore.token) {
    authStore.init()
  }
})

async function handleSubmit(values: CreateGoalFormData) {
  try {
    const goal = await goalsStore.createGoal({
      name: values.name,
      target_amount: values.target_amount,
      currency: values.currency || 'BRL',
      deadline: values.deadline,
      context: values.context,
    })
    
    router.push(`/goals/${goal.goal_id}`)
  } catch (err: any) {
    if (formRef.value) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      formRef.value.setError(extractErrorMessage(err))
    }
  }
}

function goBack() {
  router.push('/goals')
}
</script>

