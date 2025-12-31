/**
 * Tipos da API
 */

export interface User {
  user_id: string
  email: string
  first_name: string
  last_name: string
  full_name: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
}

export interface RegisterResponse {
  message: string
  data: {
    user_id: string
    email: string
    first_name: string
    last_name: string
    full_name: string
  }
}

export interface Account {
  account_id: string
  user_id: string
  name: string
  type: 'BANK' | 'WALLET' | 'INVESTMENT' | 'CREDIT_CARD'
  balance: string
  currency: 'BRL' | 'USD' | 'EUR'
  context: 'PERSONAL' | 'BUSINESS'
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateAccountRequest {
  name: string
  type: 'BANK' | 'WALLET' | 'INVESTMENT' | 'CREDIT_CARD'
  initial_balance: number // Backend espera float64, não string
  currency: 'BRL' | 'USD' | 'EUR' // Backend espera obrigatório
  context: 'PERSONAL' | 'BUSINESS'
}

export interface ListAccountsResponse {
  accounts: Account[]
  count: number
}

export interface Transaction {
  transaction_id: string
  account_id: string
  user_id: string
  type: 'INCOME' | 'EXPENSE'
  amount: string
  currency: 'BRL' | 'USD' | 'EUR'
  description: string
  date: string
  created_at: string
  updated_at: string
}

export interface CreateTransactionRequest {
  account_id: string
  type: 'INCOME' | 'EXPENSE'
  amount: number // Backend espera float64
  currency: 'BRL' | 'USD' | 'EUR' // Backend exige obrigatório
  description: string
  date: string // ISO 8601 format: YYYY-MM-DD
}

export interface ListTransactionsResponse {
  transactions: Transaction[]
  count: number
}

export interface Category {
  category_id: string
  user_id: string
  name: string
  slug: string
  description: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateCategoryRequest {
  name: string
  description?: string
}

export interface UpdateCategoryRequest {
  name?: string
  description?: string
}

export interface ListCategoriesResponse {
  categories: Category[]
  count: number
}

export interface Budget {
  budget_id: string
  user_id: string
  category_id: string
  amount: number
  currency: 'BRL' | 'USD' | 'EUR'
  period_type: 'MONTHLY' | 'YEARLY'
  year: number
  month?: number
  context: 'PERSONAL' | 'BUSINESS'
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CreateBudgetRequest {
  category_id: string
  amount: number
  currency: 'BRL' | 'USD' | 'EUR'
  period_type: 'MONTHLY' | 'YEARLY'
  year: number
  month?: number
  context: 'PERSONAL' | 'BUSINESS'
}

export interface UpdateBudgetRequest {
  amount?: number
  currency?: 'BRL' | 'USD' | 'EUR'
  period_type?: 'MONTHLY' | 'YEARLY'
  year?: number
  month?: number
  context?: 'PERSONAL' | 'BUSINESS'
  is_active?: boolean
}

export interface ListBudgetsResponse {
  budgets: Budget[]
  total: number
}

export interface BudgetProgress {
  budget_id: string
  category_id: string
  budgeted: number
  spent: number
  remaining: number
  percentage_used: number
  currency: 'BRL' | 'USD' | 'EUR'
  is_exceeded: boolean
  period_type: 'MONTHLY' | 'YEARLY'
  year: number
  month?: number
}

export interface Investment {
  investment_id: string
  user_id: string
  account_id: string
  type: 'STOCK' | 'FUND' | 'CDB' | 'TREASURY' | 'CRYPTO' | 'OTHER'
  name: string
  ticker?: string
  purchase_date: string
  purchase_amount: string
  current_value: string
  currency: 'BRL' | 'USD' | 'EUR'
  quantity?: string
  context: 'PERSONAL' | 'BUSINESS'
  return_absolute: string
  return_percentage: number
  created_at: string
  updated_at: string
}

export interface CreateInvestmentRequest {
  account_id: string
  type: 'STOCK' | 'FUND' | 'CDB' | 'TREASURY' | 'CRYPTO' | 'OTHER'
  name: string
  ticker?: string
  purchase_date: string // YYYY-MM-DD
  purchase_amount: number
  currency: 'BRL' | 'USD' | 'EUR'
  quantity?: number
  context: 'PERSONAL' | 'BUSINESS'
}

export interface UpdateInvestmentRequest {
  current_value?: number
  quantity?: number
}

export interface ListInvestmentsResponse {
  investments: Investment[]
  count: number
  pagination?: {
    page: number
    limit: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
}

export interface Goal {
  goal_id: string
  user_id: string
  name: string
  target_amount: string
  current_amount: string
  currency: 'BRL' | 'USD' | 'EUR'
  deadline: string
  context: 'PERSONAL' | 'BUSINESS'
  status: 'IN_PROGRESS' | 'COMPLETED' | 'OVERDUE' | 'CANCELLED'
  progress: number
  remaining_days: number
  created_at: string
  updated_at: string
}

export interface CreateGoalRequest {
  name: string
  target_amount: number
  currency: 'BRL' | 'USD' | 'EUR'
  deadline: string // YYYY-MM-DD
  context: 'PERSONAL' | 'BUSINESS'
}

export interface AddContributionRequest {
  amount: number
}

export interface UpdateProgressRequest {
  amount: number
}

export interface ListGoalsResponse {
  goals: Goal[]
  count: number
  pagination?: {
    page: number
    limit: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
}

// Re-export report types from reports.ts
export type {
  MonthlyReport,
  MonthlyCategorySummary,
  AnnualReport,
  MonthlyBreakdown,
  MonthlySummary,
  CategoryReport,
  CategorySummary,
  IncomeVsExpenseReport,
  PeriodSummary,
} from './reports'

