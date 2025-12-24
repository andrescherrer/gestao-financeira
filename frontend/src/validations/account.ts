import { z } from 'zod'

/**
 * Schema de validação para criação de conta
 */
export const createAccountSchema = z.object({
  name: z
    .string()
    .min(1, 'Nome da conta é obrigatório')
    .min(2, 'Nome deve ter no mínimo 2 caracteres')
    .max(100, 'Nome deve ter no máximo 100 caracteres'),
  type: z.enum(['BANK', 'WALLET', 'INVESTMENT', 'CREDIT_CARD'], {
    errorMap: () => ({ message: 'Tipo de conta é obrigatório' }),
  }),
  initial_balance: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val.toString()
        if (typeof val === 'string' && val.trim() === '') return undefined
        return val
      },
      z
        .string()
        .optional()
        .refine(
          (val) => {
            if (!val || val === '') return true
            const num = parseFloat(val)
            return !isNaN(num) && isFinite(num)
          },
          { message: 'Saldo inicial deve ser um número válido' }
        )
    ),
  currency: z.enum(['BRL', 'USD', 'EUR']).optional().default('BRL'),
  context: z.enum(['PERSONAL', 'BUSINESS'], {
    errorMap: () => ({ message: 'Contexto é obrigatório' }),
  }),
})

export type CreateAccountFormData = z.infer<typeof createAccountSchema>

