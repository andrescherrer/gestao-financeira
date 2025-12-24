// Tipos para as respostas da API

export interface ApiError {
  code: number;
  error: string;
  message?: string;
}

export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  error?: string;
}

// Tipos para autenticação
export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    user_id: string;
    email: string;
    first_name: string;
    last_name: string;
    full_name: string;
  };
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
}

export interface RegisterResponse {
  message: string;
  data: {
    user_id: string;
    email: string;
    first_name: string;
    last_name: string;
    full_name: string;
  };
}

// Tipos para contas
export interface Account {
  account_id: string;
  user_id: string;
  name: string;
  type: 'BANK' | 'WALLET' | 'INVESTMENT' | 'CREDIT_CARD';
  balance: string;
  currency: string;
  context: 'PERSONAL' | 'BUSINESS';
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateAccountRequest {
  name: string;
  type: 'BANK' | 'WALLET' | 'INVESTMENT' | 'CREDIT_CARD';
  initial_balance?: number;
  currency: 'BRL' | 'USD' | 'EUR';
  context: 'PERSONAL' | 'BUSINESS';
}

export interface ListAccountsResponse {
  accounts: Account[];
  total: number;
}

// Tipos para transações
export interface Transaction {
  transaction_id: string;
  account_id: string;
  user_id: string;
  type: 'INCOME' | 'EXPENSE' | 'TRANSFER';
  amount: string;
  currency: string;
  description: string;
  date: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTransactionRequest {
  account_id: string;
  type: 'INCOME' | 'EXPENSE' | 'TRANSFER';
  amount: string;
  currency?: string;
  description: string;
  date: string;
}

export interface UpdateTransactionRequest {
  type?: 'INCOME' | 'EXPENSE' | 'TRANSFER';
  amount?: string;
  currency?: string;
  description?: string;
  date?: string;
}

export interface ListTransactionsResponse {
  transactions: Transaction[];
  total: number;
}

