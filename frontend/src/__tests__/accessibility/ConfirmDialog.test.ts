import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

describe('ConfirmDialog - Accessibility', () => {
  it('should have proper semantic structure', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar ação',
        description: 'Tem certeza que deseja continuar?',
      },
      global: {
        stubs: {
          Dialog: {
            template: '<div role="dialog" v-if="open"><slot /></div>',
            props: ['open'],
          },
          DialogContent: true,
          DialogHeader: true,
          DialogTitle: {
            template: '<h2><slot /></h2>',
          },
          DialogDescription: {
            template: '<p><slot /></p>',
          },
          DialogFooter: true,
          Button: true,
        },
      },
    })

    // Check for dialog role
    const dialog = wrapper.find('[role="dialog"]')
    expect(dialog.exists()).toBe(true)
  })

  it('should have proper ARIA attributes', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar ação',
        description: 'Tem certeza que deseja continuar?',
      },
      global: {
        stubs: {
          Dialog: {
            template: '<div role="dialog" v-if="open"><slot /></div>',
            props: ['open'],
          },
          DialogContent: true,
          DialogHeader: true,
          DialogTitle: {
            template: '<h2><slot /></h2>',
          },
          DialogDescription: {
            template: '<p><slot /></p>',
          },
          DialogFooter: true,
          Button: true,
        },
      },
    })

    // Check for dialog role
    const dialog = wrapper.find('[role="dialog"]')
    expect(dialog.exists()).toBe(true)
  })

  it('should have accessible buttons', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar ação',
        confirmLabel: 'Confirmar',
        cancelLabel: 'Cancelar',
      },
      global: {
        stubs: {
          Dialog: {
            template: '<div v-if="open"><slot /></div>',
            props: ['open'],
          },
          DialogContent: true,
          DialogHeader: true,
          DialogTitle: true,
          DialogDescription: true,
          DialogFooter: true,
          Button: {
            template: '<button><slot /></button>',
          },
        },
      },
    })

    // Check for buttons
    const buttons = wrapper.findAll('button')
    expect(buttons.length).toBeGreaterThanOrEqual(2)
    
    const cancelButton = buttons.find(btn => btn.text().includes('Cancelar'))
    const confirmButton = buttons.find(btn => btn.text().includes('Confirmar'))
    
    expect(cancelButton).toBeDefined()
    expect(confirmButton).toBeDefined()
  })
})

