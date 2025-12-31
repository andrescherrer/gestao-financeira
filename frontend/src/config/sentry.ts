/**
 * Configuração do Sentry para Error Tracking
 */

import * as Sentry from '@sentry/vue'
import type { App } from 'vue'
import { env } from './env'

/**
 * Inicializa o Sentry
 */
export function initSentry(app: App, router: any): void {
  const dsn = import.meta.env.VITE_SENTRY_DSN
  const environment = env.environment

  // Só inicializar se DSN estiver configurado
  if (!dsn) {
    if (environment === 'development') {
      console.info('[Sentry] DSN não configurado. Error tracking desabilitado.')
    }
    return
  }

  Sentry.init({
    app,
    dsn,
    environment,
    integrations: [
      Sentry.browserTracingIntegration({
        // Rastrear navegação do Vue Router
        router,
      }),
      Sentry.replayIntegration({
        // Session Replay para debug
        maskAllText: true,
        blockAllMedia: true,
      }),
    ],
    // Performance Monitoring
    tracesSampleRate: environment === 'production' ? 0.1 : 1.0, // 10% em produção, 100% em dev
    // Session Replay
    replaysSessionSampleRate: environment === 'production' ? 0.1 : 1.0,
    replaysOnErrorSampleRate: 1.0, // Sempre gravar replays de erros
    // Configurações adicionais
    beforeSend(event, hint) {
      // Filtrar erros conhecidos que não precisam ser reportados
      if (event.exception) {
        const error = hint.originalException
        // Ignorar erros de CORS em desenvolvimento
        if (error instanceof TypeError && error.message.includes('Failed to fetch')) {
          if (environment === 'development') {
            return null
          }
        }
        // Ignorar erros de rede conhecidos
        if (error instanceof Error && error.message.includes('Network Error')) {
          if (environment === 'development') {
            return null
          }
        }
      }
      return event
    },
    // Tags padrão
    initialScope: {
      tags: {
        component: 'frontend',
      },
    },
  })
}

/**
 * Define contexto do usuário no Sentry
 */
export function setSentryUser(userId: string | null, email: string | null): void {
  if (userId || email) {
    Sentry.setUser({
      id: userId || undefined,
      email: email || undefined,
    })
  } else {
    Sentry.setUser(null)
  }
}

/**
 * Adiciona contexto adicional ao Sentry
 */
export function setSentryContext(key: string, context: Record<string, unknown>): void {
  Sentry.setContext(key, context)
}

/**
 * Captura exceção manualmente
 */
export function captureException(error: Error, context?: Record<string, unknown>): void {
  if (context) {
    Sentry.withScope((scope) => {
      Object.entries(context).forEach(([key, value]) => {
        scope.setContext(key, value as Record<string, unknown>)
      })
      Sentry.captureException(error)
    })
  } else {
    Sentry.captureException(error)
  }
}

/**
 * Captura mensagem manualmente
 */
export function captureMessage(message: string, level: Sentry.SeverityLevel = 'info'): void {
  Sentry.captureMessage(message, level)
}

