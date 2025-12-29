import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { sanitizeHtml, escapeHtml, sanitizeUrl, sanitizeText } from './sanitize'

describe('sanitize', () => {
  beforeEach(() => {
    // DOM já está disponível no ambiente de teste do Vitest
  })

  describe('escapeHtml', () => {
    it('should escape HTML special characters', () => {
      expect(escapeHtml('<script>alert("XSS")</script>')).toBe('&lt;script&gt;alert(&quot;XSS&quot;)&lt;/script&gt;')
      expect(escapeHtml('Hello & World')).toBe('Hello &amp; World')
      expect(escapeHtml("It's a test")).toBe('It&#039;s a test')
    })

    it('should handle empty strings', () => {
      expect(escapeHtml('')).toBe('')
      expect(escapeHtml(null as any)).toBe('')
      expect(escapeHtml(undefined as any)).toBe('')
    })
  })

  describe('sanitizeUrl', () => {
    it('should allow safe URLs', () => {
      expect(sanitizeUrl('https://example.com')).toBe('https://example.com')
      expect(sanitizeUrl('http://example.com')).toBe('http://example.com')
      expect(sanitizeUrl('mailto:test@example.com')).toBe('mailto:test@example.com')
      expect(sanitizeUrl('tel:+1234567890')).toBe('tel:+1234567890')
      expect(sanitizeUrl('/relative/path')).toBe('/relative/path')
    })

    it('should block dangerous URLs', () => {
      expect(sanitizeUrl('javascript:alert("XSS")')).toBe('')
      expect(sanitizeUrl('data:text/html,<script>alert("XSS")</script>')).toBe('')
      expect(sanitizeUrl('vbscript:msgbox("XSS")')).toBe('')
    })

    it('should handle empty strings', () => {
      expect(sanitizeUrl('')).toBe('')
      expect(sanitizeUrl(null as any)).toBe('')
    })
  })

  describe('sanitizeText', () => {
    it('should remove HTML tags', () => {
      expect(sanitizeText('<p>Hello</p>')).toBe('Hello')
      expect(sanitizeText('<script>alert("XSS")</script>')).toBe('alert("XSS")')
    })

    it('should handle empty strings', () => {
      expect(sanitizeText('')).toBe('')
      expect(sanitizeText(null as any)).toBe('')
    })
  })

  describe('sanitizeHtml', () => {
    it('should remove dangerous tags', () => {
      const result = sanitizeHtml('<script>alert("XSS")</script><p>Safe</p>')
      expect(result).not.toContain('<script>')
      expect(result).toContain('<p>Safe</p>')
    })

    it('should preserve safe tags', () => {
      const result = sanitizeHtml('<p>Hello <strong>World</strong></p>')
      expect(result).toContain('<p>')
      expect(result).toContain('<strong>')
      expect(result).toContain('World')
    })

    it('should escape text content', () => {
      const result = sanitizeHtml('<p>Hello & World</p>')
      expect(result).toContain('&amp;')
    })

    it('should handle empty strings', () => {
      expect(sanitizeHtml('')).toBe('')
      expect(sanitizeHtml(null as any)).toBe('')
    })
  })
})
