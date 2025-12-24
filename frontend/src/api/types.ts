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
  amount: string
  currency?: 'BRL' | 'USD' | 'EUR'
  description: string
  date: string
}

export interface ListTransactionsResponse {
  transactions: Transaction[]
  count: number
}

