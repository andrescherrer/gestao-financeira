import { watch, onUnmounted } from 'vue'
import { useBudgets } from '@/hooks/useBudgets'
import { toast } from '@/components/ui/toast'

/**
 * Composable para monitorar orçamentos e exibir alertas
 */
export function useBudgetAlerts() {
  const { budgets, useBudgetProgress } = useBudgets()
  
  // IDs de orçamentos que já exibiram alerta (para evitar spam)
  const alertedBudgets = new Set<string>()
  
  /**
   * Monitora um orçamento específico e exibe alertas
   */
  function monitorBudget(budgetId: string) {
    const { data: progress } = useBudgetProgress(budgetId)
    
    const unwatch = watch(
      progress,
      (newProgress) => {
        if (!newProgress) return
        
        const alertKey = `${budgetId}-${newProgress.percentage_used.toFixed(0)}`
        
        // Se já exibiu alerta para este nível, não exibir novamente
        if (alertedBudgets.has(alertKey)) {
          return
        }
        
        // Orçamento excedido
        if (newProgress.is_exceeded) {
          toast.error('Orçamento Excedido!', {
            description: `O orçamento de ${formatCurrency(newProgress.budgeted, newProgress.currency)} foi excedido. Gasto atual: ${formatCurrency(newProgress.spent, newProgress.currency)}`,
            duration: 10000, // 10 segundos
          })
          alertedBudgets.add(alertKey)
          return
        }
        
        // Próximo do limite (>= 90%)
        if (newProgress.percentage_used >= 90 && newProgress.percentage_used < 100) {
          toast.warning('Orçamento Próximo do Limite', {
            description: `Você já gastou ${newProgress.percentage_used.toFixed(1)}% do orçamento. Restam ${formatCurrency(newProgress.remaining, newProgress.currency)}`,
            duration: 8000, // 8 segundos
          })
          alertedBudgets.add(alertKey)
          return
        }
        
        // Aviso de 80%
        if (newProgress.percentage_used >= 80 && newProgress.percentage_used < 90) {
          toast.info('Atenção ao Orçamento', {
            description: `Você já gastou ${newProgress.percentage_used.toFixed(1)}% do orçamento. Restam ${formatCurrency(newProgress.remaining, newProgress.currency)}`,
            duration: 6000, // 6 segundos
          })
          alertedBudgets.add(alertKey)
          return
        }
      },
      { immediate: true }
    )
    
    return unwatch
  }
  
  /**
   * Monitora todos os orçamentos ativos
   */
  function monitorAllBudgets() {
    const unwatchers: Array<() => void> = []
    
    // Monitorar cada orçamento ativo
    budgets.forEach((budget) => {
      if (budget.is_active) {
        const unwatch = monitorBudget(budget.budget_id)
        unwatchers.push(unwatch)
      }
    })
    
    // Cleanup function
    return () => {
      unwatchers.forEach((unwatch) => unwatch())
    }
  }
  
  /**
   * Limpa os alertas exibidos (útil para resetar após período)
   */
  function clearAlerts() {
    alertedBudgets.clear()
  }
  
  /**
   * Formata moeda
   */
  function formatCurrency(amount: number, currency: string): string {
    return new Intl.NumberFormat('pt-BR', {
      style: 'currency',
      currency: currency || 'BRL',
    }).format(amount)
  }
  
  return {
    monitorBudget,
    monitorAllBudgets,
    clearAlerts,
  }
}

