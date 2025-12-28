/**
 * Composable para otimização de imagens
 */

/**
 * Gera srcset para diferentes densidades de tela
 */
export function generateSrcSet(
  baseSrc: string,
  widths: number[] = [320, 640, 960, 1280, 1920]
): string {
  return widths
    .map((width) => {
      // Se a URL já tem parâmetros, adicionar width como query param
      const separator = baseSrc.includes('?') ? '&' : '?'
      return `${baseSrc}${separator}w=${width} ${width}w`
    })
    .join(', ')
}

/**
 * Gera sizes attribute para responsive images
 */
export function generateSizes(breakpoints: Record<string, string>): string {
  const defaultSize = breakpoints.default || '100vw'
  const mediaQueries = Object.entries(breakpoints)
    .filter(([key]) => key !== 'default')
    .map(([breakpoint, size]) => `(max-width: ${breakpoint}px) ${size}`)
    .join(', ')

  return mediaQueries ? `${mediaQueries}, ${defaultSize}` : defaultSize
}

/**
 * Adiciona lazy loading a uma imagem
 */
export function useLazyImage(src: string, lazy: boolean = true): {
  loading: 'lazy' | 'eager'
  decoding: 'async' | 'auto'
} {
  return {
    loading: lazy ? 'lazy' : 'eager',
    decoding: lazy ? 'async' : 'auto',
  }
}

/**
 * Gera URL otimizada com parâmetros de otimização
 */
export function optimizeImageUrl(
  src: string,
  options: {
    width?: number
    height?: number
    quality?: number
    format?: 'webp' | 'avif' | 'jpg' | 'png'
  } = {}
): string {
  const { width, height, quality = 80, format } = options

  // Se já é uma URL externa ou data URL, retornar como está
  if (src.startsWith('http') || src.startsWith('data:')) {
    return src
  }

  const params = new URLSearchParams()

  if (width) params.set('w', width.toString())
  if (height) params.set('h', height.toString())
  if (quality) params.set('q', quality.toString())
  if (format) params.set('f', format)

  const separator = src.includes('?') ? '&' : '?'
  return params.toString() ? `${src}${separator}${params.toString()}` : src
}

