import { ref, watch, onMounted } from 'vue'

type Theme = 'light' | 'dark' | 'system'

/**
 * Composable para gerenciar tema (dark mode)
 */
export function useTheme() {
  // Verificar localStorage primeiro
  const getStoredTheme = (): Theme => {
    const stored = localStorage.getItem('theme') as Theme | null
    if (stored && ['light', 'dark', 'system'].includes(stored)) {
      return stored
    }
    // Fallback para system
    return 'system'
  }

  const theme = ref<Theme>(getStoredTheme())

  const isDark = ref(false)

  /**
   * Aplica o tema ao documento
   */
  function applyTheme(newTheme: Theme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)

    const root = document.documentElement

    if (newTheme === 'system') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      isDark.value = prefersDark
    } else {
      isDark.value = newTheme === 'dark'
    }

    if (isDark.value) {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }

  /**
   * Alterna entre light e dark
   */
  function toggleTheme() {
    const newTheme = isDark.value ? 'light' : 'dark'
    applyTheme(newTheme)
  }

  /**
   * Define o tema explicitamente
   */
  function setTheme(newTheme: Theme) {
    applyTheme(newTheme)
  }

  // Inicializar tema ao montar
  onMounted(() => {
    applyTheme(theme.value)

    // Ouvir mudanças na preferência do sistema
    if (theme.value === 'system') {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      const handleChange = (e: MediaQueryListEvent) => {
        if (theme.value === 'system') {
          isDark.value = e.matches
          const root = document.documentElement
          if (e.matches) {
            root.classList.add('dark')
          } else {
            root.classList.remove('dark')
          }
        }
      }
      mediaQuery.addEventListener('change', handleChange)
    }
  })

  // Watch para mudanças no tema
  watch(theme, (newTheme) => {
    applyTheme(newTheme)
  })

  return {
    theme,
    isDark,
    toggleTheme,
    setTheme,
  }
}

