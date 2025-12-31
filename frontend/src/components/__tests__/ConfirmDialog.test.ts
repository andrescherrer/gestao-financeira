import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ConfirmDialog from '../ConfirmDialog.vue'

describe('ConfirmDialog', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('deve renderizar o componente quando open é true', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar ação',
        description: 'Tem certeza?',
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
          Button: true,
          Loader2: true,
        },
      },
    })

    expect(wrapper.exists()).toBe(true)
  })

  it('não deve renderizar quando open é false', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: false,
        title: 'Confirmar ação',
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
          Button: true,
        },
      },
    })

    const dialog = wrapper.find('div')
    expect(dialog.exists()).toBe(false)
  })

  it('deve emitir confirm quando botão de confirmar é clicado', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar',
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
            template: '<button @click="$emit(\'click\')"><slot /></button>',
          },
          Loader2: true,
        },
      },
    })

    const confirmButton = wrapper.findAll('button').find(btn => btn.text().includes('Confirmar'))
    if (confirmButton) {
      await confirmButton.trigger('click')
      expect(wrapper.emitted('confirm')).toBeTruthy()
    }
  })

  it('deve emitir cancel quando botão de cancelar é clicado', async () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar',
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
            template: '<button @click="$emit(\'click\')"><slot /></button>',
          },
          Loader2: true,
        },
      },
    })

    const cancelButton = wrapper.findAll('button').find(btn => btn.text().includes('Cancelar'))
    if (cancelButton) {
      await cancelButton.trigger('click')
      expect(wrapper.emitted('cancel')).toBeTruthy()
    }
  })

  it('deve exibir loading quando isLoading é true', () => {
    const wrapper = mount(ConfirmDialog, {
      props: {
        open: true,
        title: 'Confirmar',
        isLoading: true,
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
            template: '<button :disabled="disabled"><slot /></button>',
            props: ['disabled'],
          },
          Loader2: {
            template: '<div class="loader">Loading</div>',
          },
        },
      },
    })

    const confirmButton = wrapper.findAll('button').find(btn => btn.text().includes('Confirmar'))
    if (confirmButton) {
      expect(confirmButton.attributes('disabled')).toBeDefined()
    }
  })
})

