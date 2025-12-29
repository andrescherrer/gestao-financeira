import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AccountForm from '../AccountForm.vue'

describe('AccountForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('deve renderizar o formulário', () => {
    const wrapper = mount(AccountForm, {
      global: {
        stubs: {
          Form: true,
          Field: true,
          ErrorMessage: true,
          Label: true,
          Input: true,
          Button: true,
          Card: true,
          CardContent: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
  })


  it('deve aceitar valores iniciais', () => {
    const initialValues = {
      name: 'Conta Teste',
      type: 'BANK',
      context: 'PERSONAL',
      currency: 'BRL',
      initial_balance: 1000.00,
    }

    const wrapper = mount(AccountForm, {
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
          Button: true,
          Card: true,
          CardContent: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
    // Verificar se a prop foi passada
    expect(wrapper.vm.$props.initialValues).toEqual(initialValues)
  })

  it('deve usar submitLabel customizado', () => {
    const customLabel = 'Criar Conta'
    const wrapper = mount(AccountForm, {
      props: {
        submitLabel: customLabel,
      },
      global: {
        stubs: {
          Form: {
            template: '<form><slot :isLoading="false" :errors="{}" /></form>',
          },
          Field: true,
          ErrorMessage: true,
          Label: true,
          Input: true,
          Button: {
            template: '<button><slot /></button>',
          },
          Card: true,
          CardContent: true,
        },
      },
    })

    // Verificar se a prop foi passada corretamente
    expect(wrapper.exists()).toBe(true)
    expect(wrapper.props('submitLabel')).toBe(customLabel)
  })
})
