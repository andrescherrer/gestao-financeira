import { beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

/**
 * Setup para testes de stores Pinia
 * Deve ser chamado no beforeEach de cada arquivo de teste
 */
export function setupPinia() {
  beforeEach(() => {
    // Criar nova instância do Pinia para cada teste
    setActivePinia(createPinia())
    
    // Limpar localStorage
    localStorage.clear()
    
    // Limpar sessionStorage
    sessionStorage.clear()
  })
}

/**
 * Mock do localStorage
 */
export function mockLocalStorage() {
  const store: Record<string, string> = {}
  
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      Object.keys(store).forEach(key => delete store[key])
    }),
  }
}
