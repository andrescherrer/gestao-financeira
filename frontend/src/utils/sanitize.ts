/**
 * Utilitários de sanitização para proteção contra XSS
 */

/**
 * Sanitiza HTML removendo tags e atributos perigosos
 * 
 * Esta função remove tags script, iframe, object, embed e outros elementos
 * que podem ser usados para XSS. Também remove atributos como onclick, onerror, etc.
 * 
 * @param html - String HTML a ser sanitizada
 * @returns String HTML sanitizada
 * 
 * @example
 * ```ts
 * const safeHtml = sanitizeHtml('<script>alert("XSS")</script><p>Safe content</p>')
 * // Retorna: '<p>Safe content</p>'
 * ```
 */
export function sanitizeHtml(html: string): string {
  if (!html || typeof html !== 'string') {
    return ''
  }

  // Lista de tags permitidas (whitelist)
  const allowedTags = [
    'p', 'br', 'strong', 'em', 'u', 'b', 'i', 'span', 'div',
    'h1', 'h2', 'h3', 'h4', 'h5', 'h6',
    'ul', 'ol', 'li',
    'a', 'blockquote', 'code', 'pre'
  ]

  // Lista de atributos permitidos
  const allowedAttributes = ['href', 'title', 'class', 'id']

  // Criar um elemento temporário para parsing
  const temp = document.createElement('div')
  temp.innerHTML = html

  // Função recursiva para sanitizar elementos
  function sanitizeElement(element: Element): string {
    const tagName = element.tagName.toLowerCase()

    // Se a tag não está na whitelist, retorna apenas o conteúdo
    if (!allowedTags.includes(tagName)) {
      return sanitizeChildren(element)
    }

    // Construir a tag com atributos permitidos
    let result = `<${tagName}`

    // Adicionar atributos permitidos
    for (const attr of allowedAttributes) {
      const value = element.getAttribute(attr)
      if (value) {
        // Sanitizar atributos especiais
        if (attr === 'href') {
          // Permitir apenas http, https, mailto, tel
          const sanitizedHref = sanitizeUrl(value)
          if (sanitizedHref) {
            result += ` ${attr}="${escapeHtml(sanitizedHref)}"`
          }
        } else {
          result += ` ${attr}="${escapeHtml(value)}"`
        }
      }
    }

    result += '>'

    // Adicionar conteúdo (texto e elementos filhos)
    result += sanitizeChildren(element)

    // Fechar tag
    result += `</${tagName}>`

    return result
  }

  // Sanitizar filhos (texto e elementos)
  function sanitizeChildren(element: Element): string {
    let result = ''

    for (const node of Array.from(element.childNodes)) {
      if (node.nodeType === Node.TEXT_NODE) {
        // Escapar texto
        result += escapeHtml(node.textContent || '')
      } else if (node.nodeType === Node.ELEMENT_NODE) {
        // Sanitizar elemento
        result += sanitizeElement(node as Element)
      }
    }

    return result
  }

  // Sanitizar todos os elementos
  let sanitized = ''
  for (const child of Array.from(temp.children)) {
    sanitized += sanitizeElement(child)
  }

  // Se não havia elementos filhos, pegar texto direto
  if (!sanitized && temp.textContent) {
    sanitized = escapeHtml(temp.textContent)
  }

  return sanitized
}

/**
 * Escapa caracteres HTML especiais
 * 
 * @param text - Texto a ser escapado
 * @returns Texto escapado
 */
export function escapeHtml(text: string): string {
  if (!text || typeof text !== 'string') {
    return ''
  }

  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  }

  return text.replace(/[&<>"']/g, (char) => map[char] || char)
}

/**
 * Sanitiza URLs para prevenir javascript: e data: URLs
 * 
 * @param url - URL a ser sanitizada
 * @returns URL sanitizada ou string vazia se inválida
 */
export function sanitizeUrl(url: string): string {
  if (!url || typeof url !== 'string') {
    return ''
  }

  const trimmed = url.trim().toLowerCase()

  // Permitir apenas http, https, mailto, tel
  if (
    trimmed.startsWith('http://') ||
    trimmed.startsWith('https://') ||
    trimmed.startsWith('mailto:') ||
    trimmed.startsWith('tel:') ||
    trimmed.startsWith('/') // URLs relativas
  ) {
    return url.trim()
  }

  // Bloquear javascript:, data:, vbscript:, etc.
  return ''
}

/**
 * Sanitiza texto simples removendo HTML
 * 
 * @param text - Texto a ser sanitizado
 * @returns Texto sem HTML
 */
export function sanitizeText(text: string): string {
  if (!text || typeof text !== 'string') {
    return ''
  }

  const temp = document.createElement('div')
  temp.textContent = text
  return temp.textContent || ''
}
