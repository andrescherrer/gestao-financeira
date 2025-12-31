import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import EmptyState from '../EmptyState.vue'
import { Inbox, Plus } from 'lucide-vue-next'

describe('EmptyState', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('deve renderizar o componente', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Nenhum item encontrado',
        description: 'Não há itens para exibir no momento.',
      },
      global: {
        stubs: {
          Card: true,
          CardContent: true,
          Button: true,
          RouterLink: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('deve exibir título e descrição', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Título de teste',
        description: 'Descrição de teste',
      },
      global: {
        stubs: {
          Card: true,
          CardContent: true,
          Button: true,
          RouterLink: true,
        },
      },
    })

    expect(wrapper.text()).toContain('Título de teste')
    expect(wrapper.text()).toContain('Descrição de teste')
  })

  it('deve exibir botão de ação quando fornecido', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Título',
        description: 'Descrição',
        actionLabel: 'Criar item',
        actionTo: '/create',
        actionIcon: Plus,
      },
      global: {
        stubs: {
          Card: true,
          CardContent: true,
          Button: {
            template: '<button><slot /></button>',
          },
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).toContain('Criar item')
  })

  it('não deve exibir botão de ação quando não fornecido', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Título',
        description: 'Descrição',
      },
      global: {
        stubs: {
          Card: true,
          CardContent: true,
          Button: true,
          RouterLink: true,
        },
      },
    })

    const button = wrapper.find('button')
    expect(button.exists()).toBe(false)
  })
})

