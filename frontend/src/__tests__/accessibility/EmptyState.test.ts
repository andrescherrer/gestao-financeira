import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmptyState from '@/components/EmptyState.vue'
import { Inbox } from 'lucide-vue-next'

describe('EmptyState - Accessibility', () => {
  it('should have proper semantic structure', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Nenhum item encontrado',
        description: 'Não há itens para exibir no momento.',
      },
      global: {
        stubs: {
          Card: {
            template: '<div><slot /></div>',
          },
          CardContent: {
            template: '<div><slot /></div>',
          },
          Button: true,
          RouterLink: true,
        },
      },
    })

    // Check for proper heading structure
    const heading = wrapper.find('h3')
    expect(heading.exists()).toBe(true)
  })

  it('should have proper heading structure', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Nenhum item encontrado',
        description: 'Não há itens para exibir no momento.',
      },
      global: {
        stubs: {
          Card: {
            template: '<div><slot /></div>',
          },
          CardContent: {
            template: '<div><slot /></div>',
          },
          Button: true,
          RouterLink: true,
        },
      },
    })

    // Check for heading
    const heading = wrapper.find('h3')
    expect(heading.exists()).toBe(true)
    expect(heading.text()).toBe('Nenhum item encontrado')
  })

  it('should have accessible button when action is provided', () => {
    const wrapper = mount(EmptyState, {
      props: {
        icon: Inbox,
        title: 'Nenhum item encontrado',
        description: 'Não há itens para exibir no momento.',
        actionLabel: 'Criar item',
        actionTo: '/create',
      },
      global: {
        stubs: {
          Card: {
            template: '<div><slot /></div>',
          },
          CardContent: {
            template: '<div><slot /></div>',
          },
          Button: {
            template: '<button><slot /></button>',
          },
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    // Check for action button
    const button = wrapper.find('button')
    expect(button.exists()).toBe(true)
    expect(button.text()).toContain('Criar item')
  })
})

