<template>
  <Layout>
    <div>
      <!-- Breadcrumbs -->
      <Breadcrumbs
        :items="[
          { label: 'Investimentos', to: '/investments' },
          { label: 'Novo Investimento' },
        ]"
      />

      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-4xl font-bold mb-2">Novo Investimento</h1>
        <p class="text-muted-foreground">
          Preencha os dados para adicionar um novo investimento à sua carteira
        </p>
      </div>

      <!-- Form -->
      <Card>
        <CardContent class="p-6">
          <InvestmentForm
            ref="formRef"
            :isLoading="investmentsStore.isLoading"
            submitLabel="Criar Investimento"
            @submit="handleSubmit"
          >
            <template #actions="{ isLoading: formLoading }">
              <div class="flex gap-3 pt-4">
                <Button
                  type="submit"
                  :disabled="formLoading"
                >
                  <Loader2 v-if="formLoading" class="h-4 w-4 animate-spin" />
                  <Check v-else class="h-4 w-4" />
                  {{ formLoading ? 'Criando...' : 'Criar Investimento' }}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  @click="goBack"
                >
                  <X class="h-4 w-4" />
                  Cancelar
                </Button>
              </div>
            </template>
          </InvestmentForm>
        </CardContent>
      </Card>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useInvestmentsStore } from '@/stores/investments'
import { useAccountsStore } from '@/stores/accounts'
import { useAuthStore } from '@/stores/auth'
import Layout from '@/components/layout/Layout.vue'
import InvestmentForm from '@/components/InvestmentForm.vue'
import Breadcrumbs from '@/components/Breadcrumbs.vue'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Loader2, Check, X } from 'lucide-vue-next'
import type { CreateInvestmentFormData } from '@/validations/investment'

const router = useRouter()
const investmentsStore = useInvestmentsStore()
const accountsStore = useAccountsStore()
const authStore = useAuthStore()
const formRef = ref<InstanceType<typeof InvestmentForm> | null>(null)

onMounted(async () => {
  if (!authStore.token) {
    authStore.init()
  }
  
  // Carregar contas se não tiver
  if (accountsStore.accounts.length === 0) {
    await accountsStore.listAccounts()
  }
})

async function handleSubmit(values: CreateInvestmentFormData) {
  try {
    const investment = await investmentsStore.createInvestment({
      account_id: values.account_id,
      type: values.type,
      name: values.name,
      ticker: values.ticker || undefined,
      purchase_date: values.purchase_date,
      purchase_amount: values.purchase_amount,
      currency: values.currency || 'BRL',
      quantity: values.quantity,
      context: values.context,
    })
    
    router.push(`/investments/${investment.investment_id}`)
  } catch (err: any) {
    if (formRef.value) {
      const { extractErrorMessage } = await import('@/utils/errorTranslations')
      formRef.value.setError(extractErrorMessage(err))
    }
  }
}

function goBack() {
  router.push('/investments')
}
</script>

