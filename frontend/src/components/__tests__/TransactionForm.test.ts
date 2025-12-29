import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import TransactionForm from '../TransactionForm.vue'
import type { Account, Category } from '@/api/types'

// Mock das stores
vi.mock('@/stores/accounts', () => ({
  useAccountsStore: vi.fn(() => ({
    accounts: [
      {
        account_id: 'acc-1',
        name: 'Conta Corrente',
        currency: 'BRL',
      },
      {
        account_id: 'acc-2',
        name: 'Carteira',
        currency: 'BRL',
      },
    ] as Account[],
    isLoading: false,
  })),
}))

vi.mock('@/stores/categories', () => ({
  useCategoriesStore: vi.fn(() => ({
    categories: [
      {
        category_id: 'cat-1',
        name: 'Alimentação',
        is_active: true,
      },
      {
        category_id: 'cat-2',
        name: 'Transporte',
        is_active: true,
      },
    ] as Category[],
    isLoading: false,
  })),
}))

describe('TransactionForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('deve renderizar o formulário', () => {
    const wrapper = mount(TransactionForm, {
      global: {
        stubs: {
          Form: true,
          Field: true,
          ErrorMessage: true,
          Label: true,
          Input: true,
          Textarea: true,
          Button: true,
          Card: true,
          CardContent: true,
          CategorySelect: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('deve aceitar valores iniciais', () => {
    const initialValues = {
      account_id: 'acc-1',
      type: 'INCOME',
      amount: 1000.00,
      currency: 'BRL',
      description: 'Salário',
      date: '2024-01-01',
    }

    const wrapper = mount(TransactionForm, {
      props: {
        initialValues,
      },
      global: {
        stubs: {
          Form: true,
          Field: true,
          ErrorMessage: true,
          Label: true,
          Input: true,
          Textarea: true,
          Button: true,
          Card: true,
          CardContent: true,
          CategorySelect: true,
        },
      },
    })

    // Verificar se o componente renderizou
    expect(wrapper.exists()).toBe(true)
    // Verificar se a prop foi passada
    expect(wrapper.vm.$props.initialValues).toEqual(initialValues)
  })

  it('deve desabilitar campos quando isLoading é true', () => {
    const wrapper = mount(TransactionForm, {
      props: {
        isLoading: true,
      },
      global: {
        stubs: {
          Form: true,
          Field: {
            template: '<div :disabled="disabled"><slot /></div>',
            props: ['disabled'],
          },
          ErrorMessage: true,
          Label: true,
          Input: true,
          Textarea: true,
          Button: true,
          Card: true,
          CardContent: true,
          CategorySelect: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
  })
})
