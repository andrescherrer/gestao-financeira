import apiClient from './client'

/**
 * Tipos de relatórios
 */
export interface MonthlyReport {
  user_id: string
  year: number
  month: number
  currency: 'BRL' | 'USD' | 'EUR'
  total_income: number
  total_expense: number
  balance: number
  income_count: number
  expense_count: number
  total_count: number
  category_breakdown?: MonthlyCategorySummary[]
}

export interface MonthlyCategorySummary {
  category_id: string
  category_name?: string
  total_amount: number
  count: number
  type: 'INCOME' | 'EXPENSE'
}

export interface AnnualReport {
  user_id: string
  year: number
  currency: 'BRL' | 'USD' | 'EUR'
  total_income: number
  total_expense: number
  balance: number
  income_count: number
  expense_count: number
  total_count: number
  monthly_breakdown: MonthlySummary[]
}

export interface MonthlyBreakdown {
  month: number
  total_income: number
  total_expense: number
  balance: number
  income_count: number
  expense_count: number
}

export interface MonthlySummary {
  month: number
  total_income: number
  total_expense: number
  balance: number
  income_count: number
  expense_count: number
}

export interface CategoryReport {
  user_id: string
  currency: 'BRL' | 'USD' | 'EUR'
  category_breakdown: CategorySummary[]
  total_income: number
  total_expense: number
  balance: number
  total_count: number
}

export interface CategorySummary {
  category_id?: string
  category_name?: string
  type: 'INCOME' | 'EXPENSE'
  total_amount: number
  count: number
  percentage: number
}

export interface IncomeVsExpenseReport {
  user_id: string
  currency: 'BRL' | 'USD' | 'EUR'
  total_income: number
  total_expense: number
  balance: number
  difference: number
  income_count: number
  expense_count: number
  total_count: number
  period_breakdown?: PeriodSummary[]
}

export interface PeriodSummary {
  period: string
  total_income: number
  total_expense: number
  balance: number
  income_count: number
  expense_count: number
}

/**
 * Serviço de API para relatórios
 */
export const reportService = {
  /**
   * Obtém relatório mensal
   */
  async getMonthlyReport(params: {
    year: number
    month: number
    currency?: 'BRL' | 'USD' | 'EUR'
  }): Promise<MonthlyReport> {
    const response = await apiClient.get<{
      message: string
      data: MonthlyReport
    }>('/reports/monthly', {
      params,
    })

    return response.data.data
  },

  /**
   * Obtém relatório anual
   */
  async getAnnualReport(params: {
    year: number
    currency?: 'BRL' | 'USD' | 'EUR'
  }): Promise<AnnualReport> {
    const response = await apiClient.get<{
      message: string
      data: AnnualReport
    }>('/reports/annual', {
      params,
    })

    return response.data.data
  },

  /**
   * Obtém relatório por categoria
   */
  async getCategoryReport(params?: {
    category_id?: string
    start_date?: string
    end_date?: string
    currency?: 'BRL' | 'USD' | 'EUR'
  }): Promise<CategoryReport> {
    const response = await apiClient.get<{
      message: string
      data: CategoryReport
    }>('/reports/category', {
      params,
    })

    return response.data.data
  },

  /**
   * Obtém relatório de receitas vs despesas
   */
  async getIncomeVsExpenseReport(params?: {
    start_date?: string
    end_date?: string
    currency?: 'BRL' | 'USD' | 'EUR'
    group_by?: 'day' | 'week' | 'month' | 'year'
  }): Promise<IncomeVsExpenseReport> {
    const response = await apiClient.get<{
      message: string
      data: IncomeVsExpenseReport
    }>('/reports/income-vs-expense', {
      params,
    })

    return response.data.data
  },
}

