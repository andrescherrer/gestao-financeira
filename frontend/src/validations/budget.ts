import { z } from 'zod'

/**
 * Schema de validação para criação de orçamento
 */
export const createBudgetSchema = z.object({
  category_id: z
    .string()
    .min(1, 'Categoria é obrigatória')
    .uuid('Categoria inválida'),
  amount: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val.toString()
        if (typeof val === 'string' && val.trim() === '') return undefined
        return val
      },
      z
        .string()
        .min(1, 'Valor é obrigatório')
        .refine(
          (val) => {
            const num = parseFloat(val)
            return !isNaN(num) && isFinite(num) && num > 0
          },
          { message: 'Valor deve ser um número maior que zero' }
        )
    )
    .transform((val) => parseFloat(val)),
  currency: z.enum(['BRL', 'USD', 'EUR']).optional().default('BRL'),
  period_type: z.enum(['MONTHLY', 'YEARLY'], {
    errorMap: () => ({ message: 'Tipo de período é obrigatório' }),
  }),
  year: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const trimmed = val.trim()
          if (trimmed === '') return undefined
          const parsed = parseInt(trimmed, 10)
          return isNaN(parsed) ? undefined : parsed
        }
        return undefined
      },
      z
        .number()
        .min(1900, 'Ano deve ser maior ou igual a 1900')
        .max(3000, 'Ano deve ser menor ou igual a 3000')
    ),
  month: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const trimmed = val.trim()
          if (trimmed === '') return undefined
          const parsed = parseInt(trimmed, 10)
          return isNaN(parsed) ? undefined : parsed
        }
        return undefined
      },
      z
        .number()
        .min(1, 'Mês deve ser entre 1 e 12')
        .max(12, 'Mês deve ser entre 1 e 12')
        .optional()
    ),
  context: z.enum(['PERSONAL', 'BUSINESS'], {
    errorMap: () => ({ message: 'Contexto é obrigatório' }),
  }),
}).superRefine((data, ctx) => {
  // Se for MONTHLY, mês é obrigatório
  if (data.period_type === 'MONTHLY' && !data.month) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Mês é obrigatório para orçamentos mensais',
      path: ['month'],
    })
  }
  // Se for YEARLY, mês deve ser undefined
  if (data.period_type === 'YEARLY' && data.month !== undefined) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Mês não deve ser informado para orçamentos anuais',
      path: ['month'],
    })
  }
})

export type CreateBudgetFormData = z.infer<typeof createBudgetSchema>

/**
 * Schema de validação para atualização de orçamento
 */
export const updateBudgetSchema = z.object({
  amount: z
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
            return !isNaN(num) && isFinite(num) && num > 0
          },
          { message: 'Valor deve ser um número maior que zero' }
        )
    )
    .transform((val) => val ? parseFloat(val) : undefined)
    .optional(),
  currency: z.enum(['BRL', 'USD', 'EUR']).optional(),
  period_type: z.enum(['MONTHLY', 'YEARLY']).optional(),
  year: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const trimmed = val.trim()
          if (trimmed === '') return undefined
          const parsed = parseInt(trimmed, 10)
          return isNaN(parsed) ? undefined : parsed
        }
        return undefined
      },
      z
        .number()
        .min(1900)
        .max(3000)
        .optional()
    ),
  month: z
    .preprocess(
      (val) => {
        if (val === undefined || val === null || val === '') return undefined
        if (typeof val === 'number') return val
        if (typeof val === 'string') {
          const trimmed = val.trim()
          if (trimmed === '') return undefined
          const parsed = parseInt(trimmed, 10)
          return isNaN(parsed) ? undefined : parsed
        }
        return undefined
      },
      z
        .number()
        .min(1)
        .max(12)
        .optional()
    ),
  context: z.enum(['PERSONAL', 'BUSINESS']).optional(),
  is_active: z.boolean().optional(),
})

export type UpdateBudgetFormData = z.infer<typeof updateBudgetSchema>

