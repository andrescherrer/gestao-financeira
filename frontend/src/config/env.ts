/**
 * Configurações de ambiente
 */
export const env = {
  apiUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  environment: import.meta.env.VITE_ENV || 'development',
  appName: import.meta.env.VITE_APP_NAME || 'Gestão Financeira',
  appVersion: import.meta.env.VITE_APP_VERSION || '1.0.0',
  sentryDsn: import.meta.env.VITE_SENTRY_DSN,
} as const

