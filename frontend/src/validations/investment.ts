import { z } from 'zod'

/**
 * Schema de validação para criação de investimento
 */
export const createInvestmentSchema = z.object({
  account_id: z
    .string()
    .min(1, 'Conta é obrigatória')
    .uuid('ID da conta inválido'),
  type: z.enum(['STOCK', 'FUND', 'CDB', 'TREASURY', 'CRYPTO', 'OTHER'], {
    errorMap: () => ({ message: 'Tipo de investimento é obrigatório' }),
  }),
  name: z
    .string()
    .min(1, 'Nome do investimento é obrigatório')
    .min(2, 'Nome deve ter no mínimo 2 caracteres')
    .max(200, 'Nome deve ter no máximo 200 caracteres'),
  ticker: z
    .string()
    .max(20, 'Ticker deve ter no máximo 20 caracteres')
    .optional()
    .or(z.literal('')),
  purchase_date: z
    .string()
    .min(1, 'Data de compra é obrigatória')
    .regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
  purchase_amount: z
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
          required_error: 'Valor de compra é obrigatório',
        })
        .positive('Valor de compra deve ser maior que zero')
    ),
  currency: z.enum(['BRL', 'USD', 'EUR']).optional().default('BRL'),
  quantity: z
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
        .number()
        .positive('Quantidade deve ser maior que zero')
        .optional()
    ),
  context: z.enum(['PERSONAL', 'BUSINESS'], {
    errorMap: () => ({ message: 'Contexto é obrigatório' }),
  }),
}).refine(
  (data) => {
    // STOCK, FUND e CRYPTO requerem quantity
    if (['STOCK', 'FUND', 'CRYPTO'].includes(data.type)) {
      return data.quantity !== undefined && data.quantity > 0
    }
    return true
  },
  {
    message: 'Quantidade é obrigatória para este tipo de investimento',
    path: ['quantity'],
  }
)

export type CreateInvestmentFormData = z.infer<typeof createInvestmentSchema>

/**
 * Schema de validação para atualização de investimento
 */
export const updateInvestmentSchema = z.object({
  current_value: z
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
        .number()
        .positive('Valor atual deve ser maior que zero')
        .optional()
    ),
  quantity: z
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
        .number()
        .positive('Quantidade deve ser maior que zero')
        .optional()
    ),
}).refine(
  (data) => {
    // Pelo menos um campo deve ser fornecido
    return data.current_value !== undefined || data.quantity !== undefined
  },
  {
    message: 'Pelo menos um campo (valor atual ou quantidade) deve ser fornecido',
  }
)

export type UpdateInvestmentFormData = z.infer<typeof updateInvestmentSchema>

