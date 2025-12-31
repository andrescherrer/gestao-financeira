/**
 * Sistema de Logging Estruturado para Frontend
 * 
 * Fornece logging estruturado com níveis, contexto e integração com backend
 */

import { env } from '@/config/env'

export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
}

export interface LogContext {
  [key: string]: unknown
}

export interface LogEntry {
  level: LogLevel
  message: string
  timestamp: string
  context?: LogContext
  user?: {
    id?: string
    email?: string
  }
  requestId?: string
  traceId?: string
  spanId?: string
  url?: string
  userAgent?: string
  environment?: string
}

class Logger {
  private minLevel: LogLevel
  private enableConsole: boolean
  private enableBackend: boolean
  private backendUrl?: string
  private requestId: string | null = null
  private traceId: string | null = null
  private spanId: string | null = null
  private userId: string | null = null
  private userEmail: string | null = null
  private logQueue: LogEntry[] = []
  private isSending = false

  constructor() {
    // Configurar nível mínimo baseado no ambiente
    const isDev = typeof window !== 'undefined' && window.location.hostname === 'localhost'
    this.enableConsole = isDev || (typeof window !== 'undefined' && (window as any).__VITE_ENABLE_CONSOLE_LOGS__ === 'true')
    this.enableBackend = !isDev || (typeof window !== 'undefined' && (window as any).__VITE_ENABLE_BACKEND_LOGS__ === 'true')
    this.backendUrl = env.apiUrl

    // Determinar nível mínimo de log
    const minLevelEnv = (typeof window !== 'undefined' && (window as any).__VITE_LOG_LEVEL__) || (isDev ? 'debug' : 'info')
    this.minLevel = this.parseLogLevel(minLevelEnv)

    // Gerar request ID inicial
    this.requestId = this.generateRequestId()

    // Capturar correlation IDs de headers de resposta
    this.setupCorrelationIdCapture()

    // Enviar logs pendentes periodicamente
    if (this.enableBackend) {
      this.startLogSender()
    }
  }

  private parseLogLevel(level: string): LogLevel {
    switch (level.toLowerCase()) {
      case 'debug':
        return LogLevel.DEBUG
      case 'info':
        return LogLevel.INFO
      case 'warn':
        return LogLevel.WARN
      case 'error':
        return LogLevel.ERROR
      default:
        return LogLevel.INFO
    }
  }

  private generateRequestId(): string {
    return `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
  }

  private setupCorrelationIdCapture(): void {
    // Interceptar respostas para capturar correlation IDs
    if (typeof window !== 'undefined' && 'fetch' in window) {
      const originalFetch = window.fetch
      const self = this

      window.fetch = async function (...args) {
        const response = await originalFetch.apply(this, args)
        
        // Capturar correlation IDs dos headers
        const requestId = response.headers.get('X-Request-ID')
        const traceId = response.headers.get('X-Trace-ID')
        const spanId = response.headers.get('X-Span-ID')

        if (requestId) {
          self.requestId = requestId
        }
        if (traceId) {
          self.traceId = traceId
        }
        if (spanId) {
          self.spanId = spanId
        }

        return response
      }
    }
  }

  private shouldLog(level: LogLevel): boolean {
    const levels = [LogLevel.DEBUG, LogLevel.INFO, LogLevel.WARN, LogLevel.ERROR]
    const currentIndex = levels.indexOf(level)
    const minIndex = levels.indexOf(this.minLevel)
    return currentIndex >= minIndex
  }

  private createLogEntry(
    level: LogLevel,
    message: string,
    context?: LogContext
  ): LogEntry {
    return {
      level,
      message,
      timestamp: new Date().toISOString(),
      context,
      user: this.userId || this.userEmail ? {
        id: this.userId || undefined,
        email: this.userEmail || undefined,
      } : undefined,
      requestId: this.requestId || undefined,
      traceId: this.traceId || undefined,
      spanId: this.spanId || undefined,
      url: typeof window !== 'undefined' ? window.location.href : undefined,
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : undefined,
      environment: env.environment,
    }
  }

  private logToConsole(entry: LogEntry): void {
    if (!this.enableConsole) return

    const { level, message, timestamp, context, ...meta } = entry
    const logData = {
      message,
      timestamp,
      ...meta,
      ...(context && { context }),
    }

    switch (level) {
      case LogLevel.DEBUG:
        console.debug(`[${level.toUpperCase()}]`, logData)
        break
      case LogLevel.INFO:
        console.info(`[${level.toUpperCase()}]`, logData)
        break
      case LogLevel.WARN:
        console.warn(`[${level.toUpperCase()}]`, logData)
        break
      case LogLevel.ERROR:
        console.error(`[${level.toUpperCase()}]`, logData)
        break
    }
  }

  private async sendToBackend(entry: LogEntry): Promise<void> {
    if (!this.enableBackend || !this.backendUrl) return

    // Adicionar à fila
    this.logQueue.push(entry)

    // Enviar em lote (não bloquear)
    if (this.logQueue.length >= 10 || entry.level === LogLevel.ERROR) {
      this.flushLogs()
    }
  }

  private async flushLogs(): Promise<void> {
    if (this.isSending || this.logQueue.length === 0) return

    this.isSending = true
    const logsToSend = [...this.logQueue]
    this.logQueue = []

    try {
      const token = localStorage.getItem('auth_token')
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      }

      if (token) {
        headers.Authorization = `Bearer ${token}`
      }

      // Enviar logs em lote
      await fetch(`${this.backendUrl}/logs`, {
        method: 'POST',
        headers,
        body: JSON.stringify({ logs: logsToSend }),
        keepalive: true, // Não bloquear navegação
      }).catch(() => {
        // Silenciosamente falhar - não queremos que logging quebre a aplicação
        // Re-adicionar logs à fila se falhar (limitado para evitar crescimento infinito)
        if (this.logQueue.length < 100) {
          this.logQueue.unshift(...logsToSend)
        }
      })
    } catch (error) {
      // Silenciosamente falhar
      if (this.logQueue.length < 100) {
        this.logQueue.unshift(...logsToSend)
      }
    } finally {
      this.isSending = false
    }
  }

  private startLogSender(): void {
    // Enviar logs pendentes a cada 5 segundos
    setInterval(() => {
      if (this.logQueue.length > 0) {
        this.flushLogs()
      }
    }, 5000)

    // Enviar logs pendentes antes de fechar a página
    if (typeof window !== 'undefined') {
      window.addEventListener('beforeunload', () => {
        if (this.logQueue.length > 0) {
          // Usar sendBeacon para garantir envio
          const logsToSend = JSON.stringify({ logs: this.logQueue })
          navigator.sendBeacon(`${this.backendUrl}/logs`, logsToSend)
        }
      })
    }
  }

  /**
   * Define o contexto do usuário
   */
  setUserContext(userId: string | null, userEmail: string | null): void {
    this.userId = userId
    this.userEmail = userEmail
  }

  /**
   * Define correlation IDs manualmente
   */
  setCorrelationIds(requestId?: string, traceId?: string, spanId?: string): void {
    if (requestId) this.requestId = requestId
    if (traceId) this.traceId = traceId
    if (spanId) this.spanId = spanId
  }

  /**
   * Gera um novo request ID
   */
  newRequestId(): string {
    this.requestId = this.generateRequestId()
    return this.requestId
  }

  /**
   * Log de debug
   */
  debug(message: string, context?: LogContext): void {
    if (!this.shouldLog(LogLevel.DEBUG)) return

    const entry = this.createLogEntry(LogLevel.DEBUG, message, context)
    this.logToConsole(entry)
    this.sendToBackend(entry)
  }

  /**
   * Log de informação
   */
  info(message: string, context?: LogContext): void {
    if (!this.shouldLog(LogLevel.INFO)) return

    const entry = this.createLogEntry(LogLevel.INFO, message, context)
    this.logToConsole(entry)
    this.sendToBackend(entry)
  }

  /**
   * Log de aviso
   */
  warn(message: string, context?: LogContext): void {
    if (!this.shouldLog(LogLevel.WARN)) return

    const entry = this.createLogEntry(LogLevel.WARN, message, context)
    this.logToConsole(entry)
    this.sendToBackend(entry)
  }

  /**
   * Log de erro
   */
  error(message: string, error?: Error | unknown, context?: LogContext): void {
    if (!this.shouldLog(LogLevel.ERROR)) return

    const errorContext: LogContext = {
      ...context,
    }

    if (error instanceof Error) {
      errorContext.error = {
        name: error.name,
        message: error.message,
        stack: error.stack,
      }
    } else if (error) {
      errorContext.error = error
    }

    const entry = this.createLogEntry(LogLevel.ERROR, message, errorContext)
    this.logToConsole(entry)
    this.sendToBackend(entry)
  }
}

// Exportar instância singleton
export const logger = new Logger()

// Exportar funções de conveniência
export const log = {
  debug: (message: string, context?: LogContext) => logger.debug(message, context),
  info: (message: string, context?: LogContext) => logger.info(message, context),
  warn: (message: string, context?: LogContext) => logger.warn(message, context),
  error: (message: string, error?: Error | unknown, context?: LogContext) =>
    logger.error(message, error, context),
  setUserContext: (userId: string | null, userEmail: string | null) =>
    logger.setUserContext(userId, userEmail),
  setCorrelationIds: (requestId?: string, traceId?: string, spanId?: string) =>
    logger.setCorrelationIds(requestId, traceId, spanId),
  newRequestId: () => logger.newRequestId(),
}

