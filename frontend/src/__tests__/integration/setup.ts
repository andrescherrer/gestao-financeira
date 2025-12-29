import { beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

/**
 * Setup para testes de integração
 * Limpa todas as stores e mocks antes de cada teste
 */
export function setupIntegrationTests() {
  beforeEach(() => {
    // Criar nova instância do Pinia para cada teste
    setActivePinia(createPinia())
    
    // Limpar localStorage e sessionStorage
    localStorage.clear()
    sessionStorage.clear()
    
    // Limpar todos os mocks
    vi.clearAllMocks()
  })
}

/**
 * Helper para criar dados mock de usuário
 */
export function createMockUser() {
  return {
    user_id: 'user-123',
    email: 'test@example.com',
    first_name: 'Test',
    last_name: 'User',
    full_name: 'Test User',
  }
}

/**
 * Helper para criar dados mock de conta
 */
export function createMockAccount() {
  return {
    account_id: 'acc-1',
    user_id: 'user-123',
    name: 'Conta Corrente',
    type: 'BANK' as const,
    balance: '1000.00',
    currency: 'BRL' as const,
    context: 'PERSONAL' as const,
    is_active: true,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

/**
 * Helper para criar dados mock de transação
 */
export function createMockTransaction() {
  return {
    transaction_id: 'tx-1',
    account_id: 'acc-1',
    user_id: 'user-123',
    type: 'INCOME' as const,
    amount: '1000.00',
    currency: 'BRL' as const,
    description: 'Salário',
    date: '2024-01-01',
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}

/**
 * Helper para criar dados mock de categoria
 */
export function createMockCategory() {
  return {
    category_id: 'cat-1',
    user_id: 'user-123',
    name: 'Alimentação',
    slug: 'alimentacao',
    description: 'Gastos com alimentação',
    is_active: true,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  }
}
