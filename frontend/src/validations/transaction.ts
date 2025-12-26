import { z } from 'zod'

/**
 * Schema de validação para criação de transação
 */
export const createTransactionSchema = z.object({
  account_id: z
    .string()
    .min(1, 'Conta é obrigatória')
    .uuid('ID da conta inválido'),
  type: z.enum(['INCOME', 'EXPENSE'], {
    errorMap: () => ({ message: 'Tipo de transação é obrigatório' }),
  }),
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
            if (!val || val === '') return false
            const num = parseFloat(val)
            return !isNaN(num) && isFinite(num) && num > 0
          },
          { message: 'Valor deve ser um número maior que zero' }
        )
    ),
  currency: z.enum(['BRL', 'USD', 'EUR'], {
    errorMap: () => ({ message: 'Moeda é obrigatória' }),
  }).default('BRL'),
  description: z
    .string()
    .min(1, 'Descrição é obrigatória')
    .min(3, 'Descrição deve ter no mínimo 3 caracteres') // Backend exige mínimo 3
    .max(500, 'Descrição deve ter no máximo 500 caracteres'),
  date: z
    .string()
    .min(1, 'Data é obrigatória')
    .regex(/^\d{4}-\d{2}-\d{2}$/, 'Data deve estar no formato YYYY-MM-DD'),
  category_id: z
    .preprocess(
      (val) => {
        // Se for objeto, retornar string vazia
        if (typeof val === 'object' && val !== null) {
          return ''
        }
        // Se for null ou undefined, retornar string vazia
        if (val === null || val === undefined) {
          return ''
        }
        // Converter para string
        return String(val)
      },
      z
        .string()
        .uuid('ID da categoria inválido')
        .optional()
        .or(z.literal(''))
    ),
})

export type CreateTransactionFormData = z.infer<typeof createTransactionSchema>

