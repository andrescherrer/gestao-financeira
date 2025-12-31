import { z } from 'zod'

/**
 * Schema de validação para criação de meta
 */
export const createGoalSchema = z.object({
  name: z
    .string()
    .min(1, 'Nome da meta é obrigatório')
    .min(3, 'Nome deve ter no mínimo 3 caracteres')
    .max(200, 'Nome deve ter no máximo 200 caracteres'),
  target_amount: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const num = parseFloat(val)
          return isNaN(num) ? undefined : num
        }
        return val
      },
      z
        .number({
          required_error: 'Valor alvo é obrigatório',
        })
        .positive('Valor alvo deve ser maior que zero')
    ),
  currency: z.enum(['BRL', 'USD', 'EUR']).optional().default('BRL'),
  deadline: z
    .string()
    .min(1, 'Data limite é obrigatória')
    .regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD')
    .refine(
      (date) => {
        const deadlineDate = new Date(date)
        const today = new Date()
        today.setHours(0, 0, 0, 0)
        return deadlineDate >= today
      },
      {
        message: 'Data limite deve ser no futuro',
      }
    ),
  context: z.enum(['PERSONAL', 'BUSINESS'], {
    errorMap: () => ({ message: 'Contexto é obrigatório' }),
  }),
})

export type CreateGoalFormData = z.infer<typeof createGoalSchema>

/**
 * Schema de validação para adicionar contribuição
 */
export const addContributionSchema = z.object({
  amount: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const num = parseFloat(val)
          return isNaN(num) ? undefined : num
        }
        return val
      },
      z
        .number({
          required_error: 'Valor da contribuição é obrigatório',
        })
        .positive('Valor da contribuição deve ser maior que zero')
    ),
})

export type AddContributionFormData = z.infer<typeof addContributionSchema>

/**
 * Schema de validação para atualizar progresso
 */
export const updateProgressSchema = z.object({
  amount: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const num = parseFloat(val)
          return isNaN(num) ? undefined : num
        }
        return val
      },
      z
        .number({
          required_error: 'Valor do progresso é obrigatório',
        })
        .min(0, 'Valor do progresso não pode ser negativo')
    ),
})

export type UpdateProgressFormData = z.infer<typeof updateProgressSchema>

