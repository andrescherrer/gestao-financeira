import { useQuery } from '@tanstack/vue-query'
import { reportService } from '@/api/reports'
import type {
  MonthlyReport,
  AnnualReport,
  CategoryReport,
  IncomeVsExpenseReport,
} from '@/api/reports'

/**
 * Hook para gerenciar relatórios usando TanStack Query
 */
export function useReports() {
  /**
   * Obtém relatório mensal
   */
  const useMonthlyReport = (params: {
    year: number
    month: number
    currency?: 'BRL' | 'USD' | 'EUR'
  }) => {
    return useQuery<MonthlyReport>({
      queryKey: ['reports', 'monthly', params],
      queryFn: () => reportService.getMonthlyReport(params),
      enabled: !!params.year && !!params.month,
    })
  }

  /**
   * Obtém relatório anual
   */
  const useAnnualReport = (params: {
    year: number
    currency?: 'BRL' | 'USD' | 'EUR'
  }) => {
    return useQuery<AnnualReport>({
      queryKey: ['reports', 'annual', params],
      queryFn: () => reportService.getAnnualReport(params),
      enabled: !!params.year,
    })
  }

  /**
   * Obtém relatório por categoria
   */
  const useCategoryReport = (params?: {
    category_id?: string
    start_date?: string
    end_date?: string
    currency?: 'BRL' | 'USD' | 'EUR'
  }) => {
    return useQuery<CategoryReport>({
      queryKey: ['reports', 'category', params],
      queryFn: () => reportService.getCategoryReport(params),
    })
  }

  /**
   * Obtém relatório de receitas vs despesas
   */
  const useIncomeVsExpenseReport = (params?: {
    start_date?: string
    end_date?: string
    currency?: 'BRL' | 'USD' | 'EUR'
    group_by?: 'day' | 'week' | 'month' | 'year'
  }) => {
    return useQuery<IncomeVsExpenseReport>({
      queryKey: ['reports', 'income-vs-expense', params],
      queryFn: () => reportService.getIncomeVsExpenseReport(params),
    })
  }

  return {
    useMonthlyReport,
    useAnnualReport,
    useCategoryReport,
    useIncomeVsExpenseReport,
  }
}

