/**
 * Configurações de ambiente
 * Centraliza o acesso às variáveis de ambiente do Next.js
 */

export const env = {
  // API Configuration
  apiUrl: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
  
  // Environment
  environment: process.env.NEXT_PUBLIC_ENV || 'development',
  isDevelopment: process.env.NEXT_PUBLIC_ENV === 'development',
  isProduction: process.env.NEXT_PUBLIC_ENV === 'production',
  isTest: process.env.NEXT_PUBLIC_ENV === 'test',
  
  // App Configuration
  appName: process.env.NEXT_PUBLIC_APP_NAME || 'Gestão Financeira',
  appVersion: process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0',
} as const;

/**
 * Valida se todas as variáveis de ambiente obrigatórias estão definidas
 * Apenas em runtime, não durante o build
 */
export function validateEnv(): void {
  // Não validar durante o build (SSR)
  if (process.env.NODE_ENV === 'production' && typeof window === 'undefined') {
    return;
  }

  const required = ['NEXT_PUBLIC_API_URL'];
  const missing: string[] = [];

  required.forEach((key) => {
    if (!process.env[key]) {
      missing.push(key);
    }
  });

  if (missing.length > 0 && typeof window !== 'undefined') {
    console.warn(
      `Missing required environment variables: ${missing.join(', ')}`
    );
  }
}

