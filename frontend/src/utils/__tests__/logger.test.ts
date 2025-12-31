import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { logger, log, LogLevel } from '../logger'

describe('Logger', () => {
  beforeEach(() => {
    // Limpar mocks
    vi.clearAllMocks()
    // Mock console methods
    global.console = {
      ...console,
      debug: vi.fn(),
      info: vi.fn(),
      warn: vi.fn(),
      error: vi.fn(),
    }
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Log Levels', () => {
    it('deve logar mensagens de debug', () => {
      logger.debug('Test debug message', { test: 'data' })
      expect(console.debug).toHaveBeenCalled()
    })

    it('deve logar mensagens de info', () => {
      logger.info('Test info message', { test: 'data' })
      expect(console.info).toHaveBeenCalled()
    })

    it('deve logar mensagens de warn', () => {
      logger.warn('Test warn message', { test: 'data' })
      expect(console.warn).toHaveBeenCalled()
    })

    it('deve logar mensagens de error', () => {
      const error = new Error('Test error')
      logger.error('Test error message', error, { test: 'data' })
      expect(console.error).toHaveBeenCalled()
    })
  })

  describe('User Context', () => {
    it('deve definir contexto do usuário', () => {
      logger.setUserContext('user-123', 'user@example.com')
      logger.info('Test message')
      
      const call = (console.info as any).mock.calls[0]
      expect(call[0]).toContain('INFO')
      expect(JSON.stringify(call[1])).toContain('user-123')
      expect(JSON.stringify(call[1])).toContain('user@example.com')
    })

    it('deve limpar contexto do usuário', () => {
      logger.setUserContext('user-123', 'user@example.com')
      logger.setUserContext(null, null)
      logger.info('Test message')
      
      const call = (console.info as any).mock.calls[0]
      expect(JSON.stringify(call[1])).not.toContain('user-123')
    })
  })

  describe('Correlation IDs', () => {
    it('deve definir correlation IDs', () => {
      logger.setCorrelationIds('req-123', 'trace-456', 'span-789')
      logger.info('Test message')
      
      const call = (console.info as any).mock.calls[0]
      expect(JSON.stringify(call[1])).toContain('req-123')
      expect(JSON.stringify(call[1])).toContain('trace-456')
      expect(JSON.stringify(call[1])).toContain('span-789')
    })

    it('deve gerar novo request ID', () => {
      const requestId = logger.newRequestId()
      expect(requestId).toBeTruthy()
      expect(typeof requestId).toBe('string')
    })
  })

  describe('Error Handling', () => {
    it('deve logar erros com stack trace', () => {
      const error = new Error('Test error')
      error.stack = 'Error: Test error\n    at test.js:1:1'
      
      logger.error('Error occurred', error)
      
      const call = (console.error as any).mock.calls[0]
      expect(JSON.stringify(call[1])).toContain('Test error')
      expect(JSON.stringify(call[1])).toContain('stack')
    })

    it('deve logar objetos de erro não Error', () => {
      const errorObj = { code: 'ERR001', message: 'Custom error' }
      logger.error('Error occurred', errorObj)
      
      expect(console.error).toHaveBeenCalled()
    })
  })

  describe('Convenience Functions', () => {
    it('deve exportar funções de conveniência', () => {
      expect(log.debug).toBeDefined()
      expect(log.info).toBeDefined()
      expect(log.warn).toBeDefined()
      expect(log.error).toBeDefined()
      expect(log.setUserContext).toBeDefined()
      expect(log.setCorrelationIds).toBeDefined()
      expect(log.newRequestId).toBeDefined()
    })

    it('deve usar funções de conveniência corretamente', () => {
      log.info('Test message', { test: 'data' })
      expect(console.info).toHaveBeenCalled()
    })
  })

  describe('Context Data', () => {
    it('deve incluir contexto nos logs', () => {
      const context = { userId: '123', action: 'login' }
      logger.info('User action', context)
      
      const call = (console.info as any).mock.calls[0]
      expect(JSON.stringify(call[1])).toContain('123')
      expect(JSON.stringify(call[1])).toContain('login')
    })

    it('deve incluir timestamp nos logs', () => {
      logger.info('Test message')
      
      const call = (console.info as any).mock.calls[0]
      expect(JSON.stringify(call[1])).toContain('timestamp')
    })
  })
})

